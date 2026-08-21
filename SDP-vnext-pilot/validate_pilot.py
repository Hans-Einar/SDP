#!/usr/bin/env python3
"""Deterministic standard-library validation for the provisional vNext pilot."""

from __future__ import annotations

import json
import hashlib
import copy
import re
import sys
import unicodedata
import uuid
from collections import Counter, defaultdict
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parent
UUID_URN_RE = re.compile(
    r"^urn:uuid:[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
)
KEY_RE = re.compile(r"^[A-Z][A-Z0-9]{1,11}$")
SCOPED_ID_RE = re.compile(r"^(?P<key>[A-Z][A-Z0-9]{1,11})-(?P<type>FEAT|REF|FIX|STU|SLC|SPR|ITR)-[0-9]{3}(?:-[0-9]{3})?$")
UNSCOPED_ID_RE = re.compile(r"^(?P<type>FEAT|REF|FIX|STU|SLC|SPR|ITR)-[0-9]{3}(?:-[0-9]{3})?$")
SHA_RE = re.compile(r"^[0-9a-f]{40}$")
DIGEST_RE = re.compile(r"^sha256:[0-9a-f]{64}$")
WINDOWS_RESERVED = {
    "con", "prn", "aux", "nul",
    *(f"com{i}" for i in range(1, 10)),
    *(f"lpt{i}" for i in range(1, 10)),
}
REQUIRED_TEMPLATE_KINDS = {
    "feature", "refactor", "fix", "study", "slice",
    "issue-assignment", "reservation-set", "work-domain-registry",
}
REQUIRED_TOP_LEVEL_DOCS = {
    "README.md", "WorkflowContract.md", "WorkDomains.md", "IssueContract.md",
    "ConcurrentAssignments.md", "Profiles.md", "Examples.md",
    "HSX-Pilot-Handoff.md", "OpenQuestions.md",
}
KIND_TO_ID_TOKEN = {
    "feature": "FEAT", "refactor": "REF", "fix": "FIX",
    "study": "STU", "slice": "SLC",
}
TYPE_REQUIRED_FIELDS = {
    "feature": {
        "schemaVersion", "experimental", "kind", "domainUid", "id",
        "revision", "declaredState", "title", "intent", "scope", "nonGoals",
        "constraints", "acceptanceCriteria", "relations", "issueAuthorities",
        "slices", "acceptedEvidence",
    },
    "refactor": {
        "schemaVersion", "experimental", "kind", "domainUid", "id",
        "revision", "declaredState", "title", "behaviorBaseline",
        "targetStructure", "compatibility", "temporaryAdapters", "exitEvidence",
        "relations", "issueAuthorities", "slices", "acceptedEvidence",
    },
    "fix": {
        "schemaVersion", "experimental", "kind", "domainUid", "id",
        "revision", "declaredState", "title", "defectEvidence", "correction",
        "risk", "invariants", "affectedWork", "standaloneReviewedUnit",
        "issueAuthorities", "slices", "verificationCriteria", "reviewCriteria",
        "acceptedEvidence",
    },
    "study": {
        "schemaVersion", "experimental", "kind", "domainUid", "id",
        "revision", "declaredState", "title", "question", "evidenceBoundary",
        "ownerRef", "issueAuthorities", "independentStudyRefs",
        "convergenceGate", "findings", "limitations", "informs",
        "acceptedEvidence",
    },
    "slice": {
        "schemaVersion", "experimental", "kind", "domainUid", "id",
        "revision", "declaredState", "title", "ownerRef", "assignmentIssue",
        "outcome", "whySmallestCoherent", "decisionRefs", "ownedPaths",
        "sharedTouchpoints", "invariants", "nonGoals", "verificationCriteria",
        "reviewCriteria", "discoveryRule", "completionSignal", "hardStop",
        "acceptedEvidence",
    },
    "issue-assignment": {
        "schemaVersion", "experimental", "kind", "revision", "declaredState", "authority",
        "workRef", "activeSlices", "baseline", "delivery", "coordination",
        "boundaries", "requiredEvidence", "acceptedEvidence", "stopCondition",
    },
    "reservation-set": {
        "schemaVersion", "experimental", "kind", "id", "integrationBase",
        "assignments", "mergeOrder", "convergence",
    },
    "work-domain-registry": {
        "schemaVersion", "experimental", "kind", "repository",
        "defaultDomainUid", "domains",
    },
}
WORK_OWNER_KINDS = {"feature", "refactor", "fix"}
STATEFUL_KINDS = WORK_OWNER_KINDS | {"study", "slice", "issue-assignment"}
TERMINAL_ACCEPTED_STATES = {"accepted", "delivered", "released"}
EVIDENCE_FIELDS = {
    "candidate", "verificationRefs", "currentReviewRef",
    "steeringDisposition", "releaseRefs",
}
RESERVATION_ROW_FIELDS = {
    "issue", "workRef", "activeSlices", "reservedIds", "ownedPaths",
    "sharedTouchpoints", "dependsOnIssues", "conflictsWithIssues",
}


def load_json(path: Path) -> Any:
    with path.open("r", encoding="utf-8") as handle:
        return json.load(handle)


def canonical_digest(document: Any) -> str:
    encoded = json.dumps(
        document, sort_keys=True, separators=(",", ":"), ensure_ascii=False
    ).encode("utf-8")
    return "sha256:" + hashlib.sha256(encoded).hexdigest()


def portable_text(value: str) -> str:
    return unicodedata.normalize("NFKC", value).casefold()


def valid_uid(value: Any) -> bool:
    if not isinstance(value, str) or not UUID_URN_RE.fullmatch(value):
        return False
    try:
        return str(uuid.UUID(value.removeprefix("urn:uuid:"))) == value.removeprefix("urn:uuid:")
    except ValueError:
        return False


def normalize_path(value: Any) -> tuple[str, bool] | None:
    if not isinstance(value, str) or not value:
        return None
    semantic_chars = ":/\\?*[]{}"
    for character in value:
        normalized_character = unicodedata.normalize("NFKC", character)
        if character not in semantic_chars and any(
            token in normalized_character for token in semantic_chars
        ):
            return None
    path = unicodedata.normalize("NFKC", value)
    if any(ord(character) < 32 for character in path):
        return None
    path = path.replace("\\", "/")
    recursive = path.endswith("/**")
    if recursive:
        path = path[:-3]
    if not path or path.startswith("/") or re.match(r"^[A-Za-z]:", path):
        return None
    if any(char in path for char in "?*[]{}"):
        return None
    parts = path.split("/")
    if any(not part or part in {".", ".."} for part in parts):
        return None
    normalized_parts: list[str] = []
    for part in parts:
        normalized = portable_text(part)
        if any(ord(char) < 32 for char in normalized) or ":" in normalized:
            return None
        if any(char in normalized for char in "?/\\*[]{}"):
            return None
        if normalized in {".", ".."}:
            return None
        if normalized.rstrip(" .") != normalized:
            return None
        if normalized.split(".", 1)[0] in WINDOWS_RESERVED:
            return None
        normalized_parts.append(normalized)
    return "/".join(normalized_parts), recursive


def paths_overlap(left: tuple[str, bool], right: tuple[str, bool]) -> bool:
    left_path, left_recursive = left
    right_path, right_recursive = right
    if left_path == right_path:
        return True
    if left_recursive and right_path.startswith(left_path + "/"):
        return True
    if right_recursive and left_path.startswith(right_path + "/"):
        return True
    return False


def tree_roots_overlap(left: tuple[str, bool], right: tuple[str, bool]) -> bool:
    left_path = left[0]
    right_path = right[0]
    return (
        left_path == right_path
        or left_path.startswith(right_path + "/")
        or right_path.startswith(left_path + "/")
    )


def path_is_contained(
    reservation: tuple[str, bool], requested: tuple[str, bool]
) -> bool:
    reservation_path, reservation_recursive = reservation
    requested_path, requested_recursive = requested
    if reservation_path == requested_path:
        return not requested_recursive or reservation_recursive
    return reservation_recursive and requested_path.startswith(reservation_path + "/")


def add(errors: list[str], code: str) -> None:
    if code not in errors:
        errors.append(code)


def declarations(document: dict[str, Any]) -> list[tuple[str, dict[str, Any]]]:
    result: list[tuple[str, dict[str, Any]]] = []
    for registry in document.get("registries", []):
        repository = registry.get("repository", "")
        for domain in registry.get("domains", []):
            result.append((repository, domain))
    return result


def inventory_member_id(member: Any) -> str | None:
    if isinstance(member, str):
        return member
    if isinstance(member, dict) and isinstance(member.get("id"), str):
        return member["id"]
    return None


def inventory_map(domain: dict[str, Any]) -> dict[str, Any]:
    result: dict[str, Any] = {}
    for member in domain.get("issuedIds", []):
        record_id = inventory_member_id(member)
        if isinstance(record_id, str):
            result[portable_text(record_id)] = member
    return result


def validate_inventory_member(
    domain: dict[str, Any], member: Any, errors: list[str]
) -> None:
    if not isinstance(member, dict):
        add(errors, "INVALID_ID_INVENTORY_MEMBER")
        return
    record_id = member.get("id")
    status = member.get("status")
    authority = member.get("authorityIssue")
    source = member.get("source")
    if not isinstance(record_id, str) or status not in {
        "prospective", "legacy-preserved"
    }:
        add(errors, "INVALID_ID_INVENTORY_MEMBER")
        return
    if not isinstance(authority, str) or not authority or not isinstance(source, str) or not source:
        add(errors, "INVALID_ID_INVENTORY_MEMBER")
    if status == "prospective":
        scoped = SCOPED_ID_RE.fullmatch(record_id)
        unscoped = UNSCOPED_ID_RE.fullmatch(record_id)
        style = domain.get("newRecordIdStyle")
        valid = bool(scoped or unscoped)
        if style == "scoped":
            valid = bool(scoped and scoped.group("key") == domain.get("key"))
        elif style == "unscoped":
            valid = bool(unscoped)
        if not valid:
            add(errors, "INVALID_PROSPECTIVE_ID")
    else:
        provenance = member.get("provenance")
        if not isinstance(provenance, dict):
            add(errors, "LEGACY_PROVENANCE_REQUIRED")
            return
        repository = provenance.get("repository")
        commit = provenance.get("commit")
        path = provenance.get("path")
        if (
            not isinstance(repository, str)
            or not repository
            or not isinstance(commit, str)
            or not SHA_RE.fullmatch(commit)
            or normalize_path(path) is None
            or source != path
        ):
            add(errors, "LEGACY_PROVENANCE_REQUIRED")


def validate_registry_set(document: dict[str, Any], errors: list[str]) -> None:
    by_uid: dict[str, list[tuple[str, dict[str, Any]]]] = defaultdict(list)
    for registry in document.get("registries", []):
        strict_registry = "schemaVersion" in registry
        if strict_registry:
            validate_typed_shape(registry, errors, "REGISTRY_REQUIRED_FIELD_MISSING")
        per_repo_keys: dict[str, str] = {}
        active = [domain for domain in registry.get("domains", []) if domain.get("state") == "active"]
        if len(active) > 1 and any(domain.get("newRecordIdStyle") != "scoped" for domain in active):
            add(errors, "MULTIDOMAIN_REQUIRES_SCOPED_STYLE")
        default_uid = registry.get("defaultDomainUid")
        if default_uid is not None and default_uid not in {domain.get("domainUid") for domain in active}:
            add(errors, "INVALID_DEFAULT_DOMAIN")
        root_entries: list[tuple[str, tuple[str, bool]]] = []
        for domain in registry.get("domains", []):
            if strict_registry:
                required_domain_fields = {
                    "domainUid", "key", "name", "state", "owners",
                    "newRecordIdStyle", "roots", "issuedIds",
                }
                if not required_domain_fields.issubset(domain):
                    add(errors, "DOMAIN_REQUIRED_FIELD_MISSING")
            uid = domain.get("domainUid")
            key = domain.get("key")
            if not valid_uid(uid):
                add(errors, "INVALID_DOMAIN_UID")
            if not isinstance(key, str) or not KEY_RE.fullmatch(key):
                add(errors, "INVALID_DOMAIN_KEY")
            if isinstance(key, str):
                normalized_key = portable_text(key)
                other_uid = per_repo_keys.get(normalized_key)
                if other_uid is not None and other_uid != uid:
                    add(errors, "DOMAIN_KEY_COLLISION")
                per_repo_keys[normalized_key] = uid
            if not domain.get("owners"):
                add(errors, "DOMAIN_OWNER_REQUIRED")
            if domain.get("newRecordIdStyle") not in {"scoped", "unscoped"}:
                add(errors, "INVALID_ID_STYLE")
            state = domain.get("state")
            roots = domain.get("roots", [])
            if state == "active" and not roots:
                add(errors, "ACTIVE_DOMAIN_ROOT_REQUIRED")
            if state == "moved":
                if roots or not domain.get("successorRepository") or not domain.get("successorRegistry"):
                    add(errors, "INVALID_MOVE_TOMBSTONE")
            elif state != "active":
                add(errors, "INVALID_DOMAIN_STATE")
            for root in roots:
                normalized = normalize_path(root)
                if normalized is None or normalized[1]:
                    add(errors, "INVALID_PATH")
                else:
                    if state == "active":
                        root_entries.append((uid, normalized))
            inventory = domain.get("issuedIds", [])
            if not isinstance(inventory, list):
                add(errors, "INVALID_ID_INVENTORY_MEMBER")
                inventory = []
            normalized_inventory_ids: set[str] = set()
            for member in inventory:
                record_id = inventory_member_id(member)
                if not isinstance(record_id, str):
                    add(errors, "INVALID_ID_INVENTORY_MEMBER")
                else:
                    normalized_id = portable_text(record_id)
                    if normalized_id in normalized_inventory_ids:
                        add(errors, "ISSUED_ID_COLLISION")
                    normalized_inventory_ids.add(normalized_id)
                if strict_registry:
                    validate_inventory_member(domain, member, errors)
            by_uid[uid].append((registry.get("repository", ""), domain))
        for index, (left_uid, left_path) in enumerate(root_entries):
            for right_uid, right_path in root_entries[index + 1:]:
                if left_uid != right_uid and tree_roots_overlap(left_path, right_path):
                    add(errors, "DOMAIN_ROOT_COLLISION")

    for entries in by_uid.values():
        active_entries = [item for item in entries if item[1].get("state") == "active"]
        if len(active_entries) > 1:
            add(errors, "MULTIPLE_ACTIVE_DOMAIN_HOSTS")
        keys = {item[1].get("key") for item in entries}
        if len(keys) > 1:
            add(errors, "MOVE_KEY_REWRITE")
        active_inventory = inventory_map(active_entries[0][1]) if len(active_entries) == 1 else {}
        for _, domain in entries:
            if domain.get("state") != "moved":
                continue
            moved_inventory = inventory_map(domain)
            for normalized_id, moved_member in moved_inventory.items():
                if active_inventory.get(normalized_id) != moved_member:
                    add(errors, "MOVE_ID_INVENTORY_REWRITE")
                    break
        for _, domain in entries:
            if domain.get("state") != "moved":
                continue
            active_repositories = {item[0] for item in active_entries}
            if domain.get("successorRepository") not in active_repositories:
                add(errors, "MOVE_SUCCESSOR_MISMATCH")
        if len(entries) > 1 and active_entries:
            _, active_domain = active_entries[0]
            predecessors = set(active_domain.get("predecessorRepositories", []))
            moved_repositories = {item[0] for item in entries if item[1].get("state") == "moved"}
            if not moved_repositories.issubset(predecessors):
                add(errors, "MOVE_PREDECESSOR_MISSING")


def reference_key(reference: Any) -> tuple[str, str] | None:
    if not isinstance(reference, dict):
        return None
    uid = reference.get("domainUid")
    record_id = reference.get("id")
    if not isinstance(uid, str) or not isinstance(record_id, str):
        return None
    return uid, portable_text(record_id)


def active_domain_map(document: dict[str, Any]) -> dict[str, dict[str, Any]]:
    return {
        domain.get("domainUid"): domain
        for _, domain in declarations(document)
        if domain.get("state") == "active"
    }


def issued_reference_keys(document: dict[str, Any]) -> set[tuple[str, str]]:
    return {
        (domain.get("domainUid"), portable_text(record_id))
        for _, domain in declarations(document)
        for member in domain.get("issuedIds", [])
        for record_id in [inventory_member_id(member)]
        if isinstance(record_id, str)
    }


def issued_authorities(document: dict[str, Any]) -> dict[tuple[str, str], str]:
    result: dict[tuple[str, str], str] = {}
    for _, domain in declarations(document):
        uid = domain.get("domainUid")
        for member in domain.get("issuedIds", []):
            if not isinstance(member, dict):
                continue
            record_id = inventory_member_id(member)
            authority = member.get("authorityIssue")
            if isinstance(record_id, str) and isinstance(authority, str):
                result[(uid, portable_text(record_id))] = authority
    return result


def validate_typed_shape(record: dict[str, Any], errors: list[str], code: str) -> None:
    kind = record.get("kind")
    fields = TYPE_REQUIRED_FIELDS.get(kind)
    if fields is None:
        add(errors, "UNKNOWN_RECORD_KIND")
        return
    if any(field not in record for field in fields):
        add(errors, code)
    if record.get("experimental") is not True:
        add(errors, "EXPERIMENTAL_MARKER_REQUIRED")
    if kind == "feature":
        intent = record.get("intent")
        if not isinstance(intent, dict) or not {"problem", "outcome", "usersOrValue"}.issubset(intent):
            add(errors, code)
    elif kind == "refactor":
        compatibility = record.get("compatibility")
        if not isinstance(compatibility, dict) or not {"preserved", "changed", "migrationRange"}.issubset(compatibility):
            add(errors, code)
    elif kind == "study":
        boundary = record.get("evidenceBoundary")
        if not isinstance(boundary, dict) or not {"inScope", "outOfScope"}.issubset(boundary):
            add(errors, code)
    elif kind == "issue-assignment":
        authority = record.get("authority")
        work_ref = record.get("workRef")
        baseline = record.get("baseline")
        delivery = record.get("delivery")
        coordination = record.get("coordination")
        if not isinstance(authority, dict) or not {"issue", "amendments"}.issubset(authority):
            add(errors, code)
        if not isinstance(work_ref, dict) or not {"type", "domainUid", "id"}.issubset(work_ref):
            add(errors, code)
        if not isinstance(baseline, dict) or not {"branch", "commit"}.issubset(baseline):
            add(errors, code)
        if not isinstance(delivery, dict) or not {"branch", "pullRequest"}.issubset(delivery):
            add(errors, code)
        required_coordination = {
            "integrationBase", "reservationSet", "dependsOnIssues",
            "conflictsWithIssues", "ownedPaths", "sharedTouchpoints",
            "reservedIds", "mergeOrder", "convergence", "staleBasePolicy",
        }
        if not isinstance(coordination, dict) or not required_coordination.issubset(coordination):
            add(errors, code)
    elif kind == "reservation-set":
        reservation_assignments = record.get("assignments")
        if not isinstance(reservation_assignments, list) or any(
            not isinstance(row, dict) or not RESERVATION_ROW_FIELDS.issubset(row)
            for row in reservation_assignments or []
        ):
            add(errors, code)
        convergence = record.get("convergence")
        if not isinstance(convergence, dict) or not {
            "owner", "command", "terminal"
        }.issubset(convergence):
            add(errors, code)
    elif kind == "work-domain-registry":
        if not isinstance(record.get("domains"), list) or not record.get("domains"):
            add(errors, code)


def validate_accepted_evidence(record: dict[str, Any], errors: list[str]) -> None:
    evidence = record.get("acceptedEvidence")
    if not isinstance(evidence, dict) or not EVIDENCE_FIELDS.issubset(evidence):
        add(errors, "ACCEPTED_EVIDENCE_SHAPE_INVALID")
        return
    verification_refs = evidence.get("verificationRefs")
    release_refs = evidence.get("releaseRefs")
    if not isinstance(verification_refs, list) or not isinstance(release_refs, list):
        add(errors, "ACCEPTED_EVIDENCE_SHAPE_INVALID")
        return
    state = record.get("declaredState")
    if state in TERMINAL_ACCEPTED_STATES:
        if (
            not isinstance(evidence.get("candidate"), str)
            or not SHA_RE.fullmatch(evidence["candidate"])
            or not verification_refs
            or not all(isinstance(item, str) and item for item in verification_refs)
            or not isinstance(evidence.get("currentReviewRef"), str)
            or not evidence.get("currentReviewRef")
            or evidence.get("steeringDisposition") != "accepted"
        ):
            add(errors, "ACCEPTED_EVIDENCE_REQUIRED")
        if state == "released" and (
            not release_refs
            or not all(isinstance(item, str) and item for item in release_refs)
        ):
            add(errors, "RELEASE_EVIDENCE_REQUIRED")


def validate_id_against_domain(
    uid: Any,
    record_id: Any,
    domain_by_uid: dict[str, dict[str, Any]],
    errors: list[str],
    *,
    invalid_code: str,
    style_code: str,
) -> None:
    if not valid_uid(uid) or not isinstance(record_id, str):
        add(errors, invalid_code)
        return
    scoped_match = SCOPED_ID_RE.fullmatch(record_id)
    unscoped_match = UNSCOPED_ID_RE.fullmatch(record_id)
    if not scoped_match and not unscoped_match:
        add(errors, invalid_code)
        return
    domain = domain_by_uid.get(uid)
    if domain is None:
        return
    style = domain.get("newRecordIdStyle")
    if style == "scoped" and (not scoped_match or scoped_match.group("key") != domain.get("key")):
        add(errors, style_code)
    if style == "unscoped" and not unscoped_match:
        add(errors, style_code)


def validate_records(document: dict[str, Any], errors: list[str]) -> None:
    decls = declarations(document)
    domain_by_uid = active_domain_map(document)
    issued = issued_reference_keys(document)
    active_by_repo = {
        registry.get("repository", ""): [domain for domain in registry.get("domains", []) if domain.get("state") == "active"]
        for registry in document.get("registries", [])
    }
    seen: set[tuple[str, str]] = set()
    for record in document.get("records", []):
        if "schemaVersion" in record:
            validate_typed_shape(record, errors, "RECORD_REQUIRED_FIELD_MISSING")
            if record.get("kind") in STATEFUL_KINDS:
                validate_accepted_evidence(record, errors)
        uid = record.get("domainUid")
        record_id = record.get("id")
        if not isinstance(record_id, str):
            add(errors, "INVALID_RECORD_ID")
            continue
        normalized_pair = (uid, portable_text(record_id))
        if normalized_pair in seen:
            add(errors, "RECORD_ID_COLLISION")
        seen.add(normalized_pair)
        scoped_match = SCOPED_ID_RE.fullmatch(record_id)
        unscoped_match = UNSCOPED_ID_RE.fullmatch(record_id)
        if not scoped_match and not unscoped_match:
            add(errors, "INVALID_RECORD_ID")
            continue
        domain = domain_by_uid.get(uid)
        if domain is None:
            if decls:
                add(errors, "UNKNOWN_RECORD_DOMAIN")
            continue
        expected_token = KIND_TO_ID_TOKEN.get(record.get("kind"))
        actual_token = (scoped_match or unscoped_match).group("type")
        if expected_token and expected_token != actual_token:
            add(errors, "RECORD_KIND_ID_MISMATCH")
        if scoped_match and scoped_match.group("key") != domain.get("key"):
            add(errors, "RECORD_DOMAIN_PREFIX_MISMATCH")
        if unscoped_match and record.get("identityStatus") != "legacy-preserved":
            containing_registries = [
                (registry, active_by_repo.get(registry.get("repository", ""), []))
                for registry in document.get("registries", [])
                if any(item.get("domainUid") == uid for item in registry.get("domains", []))
            ]
            if any(len(active) != 1 or registry.get("defaultDomainUid") != uid or domain.get("newRecordIdStyle") != "unscoped"
                   for registry, active in containing_registries):
                add(errors, "UNSCOPED_ID_MULTIDOMAIN")
        if "schemaVersion" in record and normalized_pair not in issued:
            add(errors, "RECORD_NOT_IN_ISSUED_INVENTORY")


def validate_record_graph(document: dict[str, Any], errors: list[str]) -> None:
    records = {
        reference_key(record): record
        for record in document.get("records", [])
        if reference_key(record) is not None
    }
    assignments = document.get("assignments", [])
    assignments_by_issue: dict[str, list[dict[str, Any]]] = defaultdict(list)
    assignments_by_slice: dict[tuple[str, str], list[dict[str, Any]]] = defaultdict(list)
    for assignment in assignments:
        issue = assignment.get("authority", {}).get("issue")
        if isinstance(issue, str):
            assignments_by_issue[issue].append(assignment)
        for slice_ref in assignment.get("activeSlices", []):
            slice_key = reference_key(slice_ref)
            if slice_key is not None:
                assignments_by_slice[slice_key].append(assignment)

    for key, record in records.items():
        kind = record.get("kind")
        if kind == "slice" and "schemaVersion" in record:
            owner_key = reference_key(record.get("ownerRef"))
            owner = records.get(owner_key)
            if owner is None:
                add(errors, "UNRESOLVED_SLICE_OWNER")
            elif owner.get("kind") not in WORK_OWNER_KINDS:
                add(errors, "INVALID_SLICE_OWNER")
            else:
                owner_slice_keys = {reference_key(item) for item in owner.get("slices", [])}
                if key not in owner_slice_keys:
                    add(errors, "NONRECIPROCAL_SLICE_OWNER")
            issue = record.get("assignmentIssue")
            matching_assignments = [
                assignment
                for assignment in assignments_by_slice.get(key, [])
                if assignment.get("authority", {}).get("issue") == issue
            ]
            if len(matching_assignments) != 1:
                add(errors, "SLICE_ASSIGNMENT_CARDINALITY")
            else:
                assignment = matching_assignments[0]
                if reference_key(assignment.get("workRef")) != owner_key:
                    add(errors, "SLICE_ASSIGNMENT_OWNER_MISMATCH")
                if owner is not None:
                    authorities = owner.get("issueAuthorities", [])
                    if authorities.count(issue) != 1:
                        add(errors, "OWNER_ISSUE_AUTHORITY_MISMATCH")
                coordination = assignment.get("coordination", {})
                reserved_owned = [
                    normalize_path(path)
                    for path in coordination.get("ownedPaths", [])
                ]
                reserved_owned = [path for path in reserved_owned if path is not None]
                reserved_shared = {
                    normalize_path(shared_path(item))
                    for item in coordination.get("sharedTouchpoints", [])
                }
                reserved_shared.discard(None)
                for path in record.get("ownedPaths", []):
                    normalized = normalize_path(path)
                    if normalized is None:
                        add(errors, "INVALID_SLICE_PATH")
                    elif not any(
                        path_is_contained(reservation, normalized)
                        for reservation in reserved_owned
                    ):
                        add(errors, "UNRESERVED_SLICE_PATH")
                for path in record.get("sharedTouchpoints", []):
                    normalized = normalize_path(shared_path(path))
                    if normalized is None:
                        add(errors, "INVALID_SLICE_PATH")
                    elif normalized not in reserved_shared:
                        add(errors, "UNDECLARED_SLICE_SHARED_PATH")
            for decision_ref in record.get("decisionRefs", []):
                decision_key = reference_key(decision_ref)
                if decision_key is None:
                    add(errors, "UNQUALIFIED_SLICE_DECISION_REF")
                elif decision_key not in records:
                    add(errors, "UNRESOLVED_SLICE_DECISION_REF")
        if kind in WORK_OWNER_KINDS and "schemaVersion" in record:
            slice_keys: list[tuple[str, str]] = []
            for slice_ref in record.get("slices", []):
                slice_key = reference_key(slice_ref)
                if slice_key is None:
                    add(errors, "INVALID_REFERENCE")
                    continue
                slice_keys.append(slice_key)
                slice_record = records.get(slice_key)
                if slice_record is None:
                    add(errors, "MISSING_OWNER_REFERENCED_SLICE")
                    continue
                if slice_record.get("kind") != "slice" or reference_key(slice_record.get("ownerRef")) != key:
                    add(errors, "NONRECIPROCAL_SLICE_OWNER")
            if len(slice_keys) != len(set(slice_keys)):
                add(errors, "NONRECIPROCAL_SLICE_OWNER")
            if record.get("declaredState") in {"delivered", "released"}:
                if not record.get("issueAuthorities"):
                    add(errors, "DELIVERED_WORK_ISSUE_REQUIRED")
                if kind != "fix" and not record.get("slices"):
                    add(errors, "DELIVERED_WORK_SLICE_REQUIRED")


def validate_relations(document: dict[str, Any], errors: list[str]) -> None:
    local_domains = {domain.get("domainUid") for _, domain in declarations(document)}
    active_domains = {domain.get("domainUid") for _, domain in declarations(document) if domain.get("state") == "active"}
    known = {
        reference_key(record)
        for record in document.get("records", [])
        if reference_key(record) is not None
    } | issued_reference_keys(document)
    for relation in document.get("relations", []):
        for name in ("sourceRef", "targetRef"):
            reference = relation.get(name)
            key = reference_key(reference)
            if key is None:
                if len(active_domains) > 1:
                    add(errors, "UNQUALIFIED_CROSS_DOMAIN_REF")
                else:
                    add(errors, "INVALID_REFERENCE")
                continue
            if not valid_uid(key[0]):
                add(errors, "INVALID_REFERENCE")
                continue
            if key[0] in local_domains and key not in known:
                add(errors, "UNRESOLVED_RELATION_REF")


def shared_path(item: Any) -> Any:
    return item.get("path") if isinstance(item, dict) else item


def assignment_projection(assignment: dict[str, Any]) -> dict[str, Any]:
    coordination = assignment.get("coordination", {})
    return {
        "issue": assignment.get("authority", {}).get("issue"),
        "workRef": assignment.get("workRef"),
        "activeSlices": assignment.get("activeSlices", []),
        "reservedIds": coordination.get("reservedIds", []),
        "ownedPaths": coordination.get("ownedPaths", []),
        "sharedTouchpoints": coordination.get("sharedTouchpoints", []),
        "dependsOnIssues": coordination.get("dependsOnIssues", []),
        "conflictsWithIssues": coordination.get("conflictsWithIssues", []),
    }


def validate_reservation_sets(document: dict[str, Any], errors: list[str]) -> None:
    assignments = document.get("assignments", [])
    reservation_sets = document.get("reservationSets", [])
    by_id: dict[str, list[dict[str, Any]]] = defaultdict(list)
    for reservation in reservation_sets:
        if not isinstance(reservation, dict):
            add(errors, "RESERVATION_SET_INCOMPLETE")
            continue
        if reservation.get("kind") == "reservation-set":
            validate_typed_shape(
                reservation, errors, "RESERVATION_SET_INCOMPLETE"
            )
        reservation_id = reservation.get("id")
        if not isinstance(reservation_id, str) or not reservation_id:
            add(errors, "RESERVATION_SET_INCOMPLETE")
            continue
        by_id[reservation_id].append(reservation)
    if any(len(values) != 1 for values in by_id.values()):
        add(errors, "RESERVATION_SET_INCOMPLETE")

    grouped_assignments: dict[str, list[dict[str, Any]]] = defaultdict(list)
    for assignment in assignments:
        reference = assignment.get("coordination", {}).get("reservationSet")
        if not isinstance(reference, dict):
            add(errors, "RESERVATION_SET_REQUIRED")
            continue
        reservation_id = reference.get("id")
        digest = reference.get("digest")
        if not isinstance(reservation_id, str) or not reservation_id:
            add(errors, "RESERVATION_SET_REQUIRED")
            continue
        grouped_assignments[reservation_id].append(assignment)
        candidates = by_id.get(reservation_id, [])
        if len(candidates) != 1:
            add(errors, "UNRESOLVED_RESERVATION_SET")
            continue
        if not isinstance(digest, str) or not DIGEST_RE.fullmatch(digest):
            add(errors, "INVALID_RESERVATION_DIGEST")
        elif digest != canonical_digest(candidates[0]):
            add(errors, "RESERVATION_DIGEST_MISMATCH")

    for reservation_id, group in grouped_assignments.items():
        candidates = by_id.get(reservation_id, [])
        if len(candidates) != 1:
            continue
        reservation = candidates[0]
        if not SHA_RE.fullmatch(str(reservation.get("integrationBase", ""))):
            add(errors, "INVALID_INTEGRATION_BASE")
        rows = reservation.get("assignments")
        if not isinstance(rows, list) or not rows:
            add(errors, "RESERVATION_SET_INCOMPLETE")
            continue
        if any(
            not isinstance(row, dict) or not RESERVATION_ROW_FIELDS.issubset(row)
            for row in rows
        ):
            add(errors, "RESERVATION_SET_INCOMPLETE")
        group_issues = [
            assignment.get("authority", {}).get("issue") for assignment in group
        ]
        row_issues = [row.get("issue") for row in rows if isinstance(row, dict)]
        if Counter(row_issues) != Counter(group_issues):
            add(errors, "RESERVATION_SET_INCOMPLETE")
        rows_by_issue = {
            row.get("issue"): row for row in rows if isinstance(row, dict)
        }
        for assignment in group:
            issue = assignment.get("authority", {}).get("issue")
            if rows_by_issue.get(issue) != assignment_projection(assignment):
                add(errors, "RESERVATION_ASSIGNMENT_MISMATCH")
        bases = {
            assignment.get("coordination", {}).get("integrationBase")
            for assignment in group
        }
        if bases != {reservation.get("integrationBase")}:
            add(errors, "RESERVATION_ASSIGNMENT_MISMATCH")
        shared_order = {
            canonical_digest({
                "mergeOrder": assignment.get("coordination", {}).get("mergeOrder"),
                "convergence": assignment.get("coordination", {}).get("convergence"),
            })
            for assignment in group
        }
        if len(shared_order) != 1:
            add(errors, "CONVERGENCE_CONTRACT_MISMATCH")
        for assignment in group:
            coordination = assignment.get("coordination", {})
            if (
                coordination.get("mergeOrder") != reservation.get("mergeOrder")
                or coordination.get("convergence") != reservation.get("convergence")
            ):
                add(errors, "RESERVATION_ASSIGNMENT_MISMATCH")

        issue_set = set(group_issues)
        edge_map: dict[str, list[str]] = {}
        conflict_map: dict[str, list[str]] = {}
        for row in rows:
            if not isinstance(row, dict) or not isinstance(row.get("issue"), str):
                continue
            issue = row["issue"]
            dependencies = row.get("dependsOnIssues", [])
            conflicts = row.get("conflictsWithIssues", [])
            if not isinstance(dependencies, list) or not isinstance(conflicts, list):
                add(errors, "RESERVATION_SET_INCOMPLETE")
                continue
            edge_map[issue] = dependencies
            conflict_map[issue] = conflicts
            if any(dependency not in issue_set for dependency in dependencies):
                add(errors, "UNRESOLVED_DEPENDENCY_ISSUE")
            if any(conflict not in issue_set for conflict in conflicts):
                add(errors, "UNRESOLVED_CONFLICT_ISSUE")
            if issue in dependencies:
                add(errors, "DEPENDENCY_CYCLE")

        visiting: set[str] = set()
        visited: set[str] = set()

        def visit(issue: str) -> bool:
            if issue in visiting:
                return True
            if issue in visited:
                return False
            visiting.add(issue)
            if any(
                dependency in issue_set and visit(dependency)
                for dependency in edge_map.get(issue, [])
            ):
                return True
            visiting.remove(issue)
            visited.add(issue)
            return False

        if any(visit(issue) for issue in issue_set if issue not in visited):
            add(errors, "DEPENDENCY_CYCLE")
        for issue, conflicts in conflict_map.items():
            for conflict in conflicts:
                if conflict in issue_set and issue not in conflict_map.get(conflict, []):
                    add(errors, "ASYMMETRIC_CONFLICT")

        merge_order = reservation.get("mergeOrder")
        if not isinstance(merge_order, list) or not merge_order:
            add(errors, "MERGE_ORDER_MISSING_MEMBER")
            merge_order = []
        if len(merge_order) != len(set(merge_order)):
            add(errors, "MERGE_ORDER_DUPLICATE")
        if any(issue not in issue_set for issue in merge_order):
            add(errors, "MERGE_ORDER_UNKNOWN_MEMBER")
        if issue_set - set(merge_order):
            add(errors, "MERGE_ORDER_MISSING_MEMBER")
        positions = {issue: index for index, issue in enumerate(merge_order)}
        for issue, dependencies in edge_map.items():
            for dependency in dependencies:
                if (
                    issue in positions
                    and dependency in positions
                    and positions[dependency] >= positions[issue]
                ):
                    add(errors, "MERGE_ORDER_DEPENDENCY_VIOLATION")
        convergence = reservation.get("convergence")
        if (
            not isinstance(convergence, dict)
            or not isinstance(convergence.get("owner"), str)
            or not convergence.get("owner").strip()
            or not isinstance(convergence.get("command"), str)
            or not convergence.get("command").strip()
        ):
            add(errors, "EMPTY_CONVERGENCE")
        if not isinstance(convergence, dict) or convergence.get("terminal") is not True:
            add(errors, "INVALID_CONVERGENCE_TERMINAL")


def validate_assignments(document: dict[str, Any], errors: list[str]) -> None:
    assignments = document.get("assignments", [])
    domain_by_uid = active_domain_map(document)
    records = {
        reference_key(record): record
        for record in document.get("records", [])
        if reference_key(record) is not None
    }
    issues = [
        assignment.get("authority", {}).get("issue") for assignment in assignments
    ]
    if any(
        issue is not None and count > 1
        for issue, count in Counter(issues).items()
    ):
        add(errors, "DUPLICATE_ISSUE_AUTHORITY")
    if len(assignments) > 1:
        bases = {item.get("coordination", {}).get("integrationBase") for item in assignments}
        reservation_sets = {
            (
                item.get("coordination", {}).get("reservationSet", {}).get("id"),
                item.get("coordination", {}).get("reservationSet", {}).get("digest"),
            )
            for item in assignments
        }
        if len(bases) != 1 or None in bases or len(reservation_sets) != 1:
            add(errors, "STALE_COMMON_BASE")
        convergence_contracts = {
            canonical_digest({
                "mergeOrder": item.get("coordination", {}).get("mergeOrder"),
                "convergence": item.get("coordination", {}).get("convergence"),
            })
            for item in assignments
        }
        if len(convergence_contracts) != 1:
            add(errors, "CONVERGENCE_CONTRACT_MISMATCH")
        issue_map = {
            item.get("authority", {}).get("issue"): item
            for item in assignments
            if item.get("authority", {}).get("issue")
        }
        for issue, assignment in issue_map.items():
            for conflict in assignment.get("coordination", {}).get("conflictsWithIssues", []):
                peer = issue_map.get(conflict)
                if peer is not None and issue not in peer.get("coordination", {}).get("conflictsWithIssues", []):
                    add(errors, "ASYMMETRIC_CONFLICT")

    owned: list[tuple[str, tuple[str, bool]]] = []
    shared: list[tuple[str, tuple[str, bool], Any]] = []
    reserved: dict[tuple[str, str], str] = {}
    for assignment in assignments:
        if assignment.get("kind") == "issue-assignment":
            validate_typed_shape(assignment, errors, "ASSIGNMENT_REQUIRED_FIELD_MISSING")
            validate_accepted_evidence(assignment, errors)
        authority = assignment.get("authority", {}).get("issue", "")
        coordination = assignment.get("coordination", {})
        base = coordination.get("integrationBase")
        if base is not None and not SHA_RE.fullmatch(str(base)):
            add(errors, "INVALID_INTEGRATION_BASE")
        reservation_reference = coordination.get("reservationSet")
        digest = (
            reservation_reference.get("digest")
            if isinstance(reservation_reference, dict)
            else None
        )
        if digest is not None and not DIGEST_RE.fullmatch(str(digest)):
            add(errors, "INVALID_RESERVATION_DIGEST")
        if assignment.get("kind") == "issue-assignment":
            baseline_commit = assignment.get("baseline", {}).get("commit")
            if not SHA_RE.fullmatch(str(baseline_commit)) or baseline_commit != base:
                add(errors, "INVALID_INTEGRATION_BASE")
            work_ref = assignment.get("workRef")
            validate_id_against_domain(
                work_ref.get("domainUid") if isinstance(work_ref, dict) else None,
                work_ref.get("id") if isinstance(work_ref, dict) else None,
                domain_by_uid, errors, invalid_code="INVALID_WORK_REF",
                style_code="WORK_REF_STYLE_MISMATCH",
            )
            work_key = reference_key(work_ref)
            work_record = records.get(work_key)
            if work_record is None:
                add(errors, "UNRESOLVED_WORK_REF")
            elif work_record.get("kind") != work_ref.get("type"):
                add(errors, "WORK_REF_TYPE_MISMATCH")
            elif work_record.get("kind") in WORK_OWNER_KINDS:
                authorities = work_record.get("issueAuthorities", [])
                if authorities.count(authority) != 1:
                    add(errors, "OWNER_ISSUE_AUTHORITY_MISMATCH")
            for active_slice in assignment.get("activeSlices", []):
                validate_id_against_domain(
                    active_slice.get("domainUid") if isinstance(active_slice, dict) else None,
                    active_slice.get("id") if isinstance(active_slice, dict) else None,
                    domain_by_uid, errors, invalid_code="INVALID_REFERENCE",
                    style_code="ACTIVE_SLICE_STYLE_MISMATCH",
                )
                slice_record = records.get(reference_key(active_slice))
                if slice_record is None or slice_record.get("kind") != "slice":
                    add(errors, "UNRESOLVED_ACTIVE_SLICE")
                elif slice_record.get("assignmentIssue") != authority:
                    add(errors, "SLICE_ASSIGNMENT_MISMATCH")
                elif reference_key(slice_record.get("ownerRef")) != work_key:
                    add(errors, "ACTIVE_SLICE_OWNER_MISMATCH")
        for value in coordination.get("ownedPaths", []):
            normalized = normalize_path(value)
            if normalized is None:
                add(errors, "INVALID_PATH")
            else:
                owned.append((authority, normalized))
        for item in coordination.get("sharedTouchpoints", []):
            normalized = normalize_path(shared_path(item))
            if normalized is None:
                add(errors, "INVALID_PATH")
            else:
                owner_uid = item.get("ownerDomainUid") if isinstance(item, dict) else None
                if not isinstance(item, dict) or not item.get("allowedMutation"):
                    add(errors, "INVALID_SHARED_TOUCHPOINT")
                if not valid_uid(owner_uid):
                    add(errors, "INVALID_SHARED_OWNER")
                elif domain_by_uid and owner_uid not in domain_by_uid:
                    add(errors, "UNKNOWN_SHARED_OWNER_DOMAIN")
                shared.append((authority, normalized, owner_uid))
        for reference in coordination.get("reservedIds", []):
            if not isinstance(reference, dict) or not reference.get("domainUid") or not reference.get("id"):
                add(errors, "INVALID_REFERENCE")
                continue
            validate_id_against_domain(
                reference.get("domainUid"), reference.get("id"), domain_by_uid,
                errors, invalid_code="INVALID_RESERVED_ID",
                style_code="RESERVED_ID_STYLE_MISMATCH",
            )
            if domain_by_uid and reference.get("domainUid") not in domain_by_uid:
                add(errors, "UNKNOWN_RESERVED_DOMAIN")
            key = (reference["domainUid"], portable_text(reference["id"]))
            if key in reserved and reserved[key] != authority:
                add(errors, "RESERVED_ID_COLLISION")
            reserved[key] = authority
            issued_by = issued_authorities(document).get(key)
            if issued_by is not None and issued_by != authority:
                add(errors, "ISSUED_ID_RECLAIM")

    for index, (left_issue, left_path) in enumerate(owned):
        for right_issue, right_path in owned[index + 1:]:
            if left_issue != right_issue and paths_overlap(left_path, right_path):
                add(errors, "PATH_COLLISION")
        for _, shared_value, _ in shared:
            if paths_overlap(left_path, shared_value):
                add(errors, "PATH_COLLISION")

    if len(assignments) > 1:
        shared_by_path: dict[tuple[str, bool], list[tuple[str, Any]]] = defaultdict(list)
        for issue, path, owner_uid in shared:
            shared_by_path[path].append((issue, owner_uid))
        if any(len({issue for issue, _ in values}) < 2 for values in shared_by_path.values()):
            add(errors, "ASYMMETRIC_SHARED_TOUCHPOINT")
        if any(len({owner for _, owner in values}) != 1 or None in {owner for _, owner in values}
               for values in shared_by_path.values()):
            add(errors, "SHARED_OWNER_MISMATCH")


def validate_document(
    document: dict[str, Any], *, isolated_fixture: bool = False
) -> list[str]:
    errors: list[str] = []
    if document.get("experimental") is not True:
        add(errors, "EXPERIMENTAL_MARKER_REQUIRED")
    validate_registry_set(document, errors)
    validate_records(document, errors)
    validate_record_graph(document, errors)
    validate_relations(document, errors)
    validate_assignments(document, errors)
    if not isolated_fixture:
        validate_reservation_sets(document, errors)
    return sorted(errors)


def validate_template(record: dict[str, Any]) -> list[str]:
    errors: list[str] = []
    kind = record.get("kind")
    validate_typed_shape(record, errors, "TEMPLATE_REQUIRED_FIELD_MISSING")
    if kind in KIND_TO_ID_TOKEN:
        if not valid_uid(record.get("domainUid")):
            add(errors, "INVALID_DOMAIN_UID")
        match = SCOPED_ID_RE.fullmatch(str(record.get("id", ""))) or UNSCOPED_ID_RE.fullmatch(str(record.get("id", "")))
        if not match or match.group("type") != KIND_TO_ID_TOKEN[kind]:
            add(errors, "RECORD_KIND_ID_MISMATCH")
        validate_accepted_evidence(record, errors)
    elif kind == "issue-assignment":
        baseline = record.get("baseline", {}).get("commit", "")
        integration = record.get("coordination", {}).get("integrationBase", "")
        if not SHA_RE.fullmatch(str(baseline)) or baseline != integration:
            add(errors, "INVALID_INTEGRATION_BASE")
        validate_accepted_evidence(record, errors)
    elif kind == "reservation-set":
        if not SHA_RE.fullmatch(str(record.get("integrationBase", ""))):
            add(errors, "INVALID_INTEGRATION_BASE")
        wrapper = {
            "experimental": True,
            "registries": [],
            "records": [],
            "relations": [],
            "assignments": [],
            "reservationSets": [record],
        }
        errors.extend(validate_document(wrapper))
    elif kind == "work-domain-registry":
        wrapper = {"experimental": True, "registries": [record], "records": [], "relations": [], "assignments": []}
        errors.extend(validate_document(wrapper))
    else:
        add(errors, "UNKNOWN_TEMPLATE_KIND")
    return sorted(set(errors))


def check_coverage(example_documents: dict[str, dict[str, Any]], negative_codes: set[str]) -> list[str]:
    failures: list[str] = []
    simple = example_documents.get("simple-single-domain.json", {})
    scoped = example_documents.get("scoped-monorepo.json", {})
    moved = example_documents.get("domain-move.json", {})
    if not any(UNSCOPED_ID_RE.fullmatch(str(record.get("id", ""))) for record in simple.get("records", [])):
        failures.append("coverage: simple unscoped example missing")
    active_scoped_domains = [domain for _, domain in declarations(scoped) if domain.get("state") == "active"]
    if len(active_scoped_domains) < 4:
        failures.append("coverage: four-domain scoped example missing")
    if len(scoped.get("assignments", [])) < 3:
        failures.append("coverage: three-Issue concurrent assignment example missing")
    if len(scoped.get("reservationSets", [])) != 1:
        failures.append("coverage: canonical embedded reservation set missing")
    if not any(
        assignment.get("authority", {}).get("issue", "").endswith("/203")
        and assignment.get("activeSlices")
        for assignment in scoped.get("assignments", [])
    ):
        failures.append("coverage: shared-domain Issue/Slice assignment missing")
    if not any(
        isinstance(relation.get("sourceRef"), dict)
        and isinstance(relation.get("targetRef"), dict)
        and relation["sourceRef"].get("domainUid") != relation["targetRef"].get("domainUid")
        for relation in scoped.get("relations", [])
    ):
        failures.append("coverage: qualified cross-domain relation missing")
    moved_uids = Counter(domain.get("domainUid") for _, domain in declarations(moved))
    if not any(count > 1 for count in moved_uids.values()):
        failures.append("coverage: move/tombstone example missing")
    if not any(
        isinstance(member, dict)
        and member.get("id") == "DBG-RF-001"
        and member.get("status") == "legacy-preserved"
        and isinstance(member.get("provenance"), dict)
        for _, domain in declarations(moved)
        for member in domain.get("issuedIds", [])
    ):
        failures.append("coverage: HSX-shaped preserved legacy inventory missing")
    required_negative = {
        "DOMAIN_KEY_COLLISION", "RECORD_ID_COLLISION", "PATH_COLLISION",
        "UNQUALIFIED_CROSS_DOMAIN_REF", "UNSCOPED_ID_MULTIDOMAIN",
        "STALE_COMMON_BASE", "RESERVED_ID_COLLISION",
        "ASYMMETRIC_SHARED_TOUCHPOINT", "MOVE_KEY_REWRITE",
        "MOVE_ID_INVENTORY_REWRITE",
        "CONVERGENCE_CONTRACT_MISMATCH", "ASYMMETRIC_CONFLICT",
        "UNRESOLVED_WORK_REF", "UNRESOLVED_RELATION_REF",
        "RESERVED_ID_STYLE_MISMATCH",
        "NONRECIPROCAL_SLICE_OWNER",
        "SLICE_ASSIGNMENT_CARDINALITY", "MISSING_OWNER_REFERENCED_SLICE",
        "INVALID_SLICE_PATH", "UNRESERVED_SLICE_PATH",
        "UNDECLARED_SLICE_SHARED_PATH", "UNQUALIFIED_SLICE_DECISION_REF",
        "UNRESOLVED_SLICE_DECISION_REF", "DELIVERED_WORK_ISSUE_REQUIRED",
        "DELIVERED_WORK_SLICE_REQUIRED", "DUPLICATE_ISSUE_AUTHORITY",
        "RESERVATION_SET_REQUIRED", "UNRESOLVED_RESERVATION_SET",
        "RESERVATION_SET_INCOMPLETE", "RESERVATION_DIGEST_MISMATCH",
        "RESERVATION_ASSIGNMENT_MISMATCH", "DEPENDENCY_CYCLE",
        "UNRESOLVED_DEPENDENCY_ISSUE", "MERGE_ORDER_MISSING_MEMBER",
        "MERGE_ORDER_DUPLICATE", "MERGE_ORDER_DEPENDENCY_VIOLATION",
        "EMPTY_CONVERGENCE", "INVALID_CONVERGENCE_TERMINAL",
        "ACCEPTED_EVIDENCE_REQUIRED",
        "RELEASE_EVIDENCE_REQUIRED", "ISSUED_ID_COLLISION",
        "INVALID_PROSPECTIVE_ID", "LEGACY_PROVENANCE_REQUIRED",
        "ISSUED_ID_RECLAIM", "DOMAIN_ROOT_COLLISION", "INVALID_PATH",
    }
    missing = sorted(required_negative - negative_codes)
    if missing:
        failures.append("coverage: missing negative diagnostics " + ", ".join(missing))
    return failures


def pointer_parts(pointer: str) -> list[str]:
    if not pointer.startswith("/"):
        raise ValueError("JSON pointer must start with '/'")
    return [part.replace("~1", "/").replace("~0", "~") for part in pointer[1:].split("/")]


def pointer_parent(document: Any, pointer: str) -> tuple[Any, str]:
    parts = pointer_parts(pointer)
    if not parts:
        raise ValueError("root mutation is not supported")
    current = document
    for part in parts[:-1]:
        current = current[int(part)] if isinstance(current, list) else current[part]
    return current, parts[-1]


def pointer_value(document: Any, pointer: str) -> Any:
    current = document
    for part in pointer_parts(pointer):
        current = current[int(part)] if isinstance(current, list) else current[part]
    return current


def apply_mutations(document: dict[str, Any], mutations: list[dict[str, Any]]) -> None:
    for mutation in mutations:
        operation = mutation.get("op")
        pointer = mutation.get("path")
        if not isinstance(pointer, str):
            raise ValueError("mutation path must be a JSON pointer")
        parent, key = pointer_parent(document, pointer)
        if operation == "remove":
            if isinstance(parent, list):
                parent.pop(int(key))
            else:
                del parent[key]
        elif operation in {"add", "replace"}:
            value = copy.deepcopy(mutation.get("value"))
            if isinstance(parent, list):
                if key == "-":
                    parent.append(value)
                elif operation == "add":
                    parent.insert(int(key), value)
                else:
                    parent[int(key)] = value
            else:
                parent[key] = value
        elif operation == "copy":
            source = mutation.get("from")
            if not isinstance(source, str):
                raise ValueError("copy mutation requires 'from'")
            value = copy.deepcopy(pointer_value(document, source))
            if isinstance(parent, list):
                if key == "-":
                    parent.append(value)
                else:
                    parent.insert(int(key), value)
            else:
                parent[key] = value
        else:
            raise ValueError(f"unsupported mutation operation {operation!r}")


def rebind_reservation_digests(document: dict[str, Any]) -> None:
    digests = {
        reservation.get("id"): canonical_digest(reservation)
        for reservation in document.get("reservationSets", [])
        if isinstance(reservation, dict) and reservation.get("id")
    }
    for assignment in document.get("assignments", []):
        reference = assignment.get("coordination", {}).get("reservationSet")
        if isinstance(reference, dict) and reference.get("id") in digests:
            reference["digest"] = digests[reference["id"]]


def validate_text_corpus() -> list[str]:
    failures: list[str] = []
    markdown_files = sorted(ROOT.rglob("*.md"))
    top_level_names = {path.name for path in ROOT.glob("*.md")}
    missing_docs = sorted(REQUIRED_TOP_LEVEL_DOCS - top_level_names)
    if missing_docs:
        failures.append("corpus: missing top-level documents " + ", ".join(missing_docs))

    status_re = re.compile(
        r"^Status:\s.*(?:pilot|provisional|experimental)",
        flags=re.IGNORECASE | re.MULTILINE,
    )
    link_re = re.compile(r"\[[^\]]+\]\(([^)]+)\)")
    for path in markdown_files:
        content = path.read_text(encoding="utf-8")
        if path.parent == ROOT and not status_re.search("\n".join(content.splitlines()[:8])):
            failures.append(f"corpus: {path.name} lacks a clear top-level pilot/provisional Status marker")
        for target in link_re.findall(content):
            target = target.strip().strip("<>")
            if not target or "://" in target or target.startswith(("#", "mailto:")):
                continue
            local_target = target.split("#", 1)[0]
            if local_target and not (path.parent / local_target).resolve().exists():
                failures.append(f"corpus: broken local link {path.relative_to(ROOT)} -> {target}")

    text_suffixes = {".md", ".json", ".py", ".yaml"}
    for path in sorted(item for item in ROOT.rglob("*") if item.is_file() and item.suffix in text_suffixes):
        for line_number, line in enumerate(path.read_text(encoding="utf-8").splitlines(), start=1):
            if line.endswith((" ", "\t")):
                failures.append(f"corpus: trailing whitespace {path.relative_to(ROOT)}:{line_number}")
    return failures


def main() -> int:
    failures: list[str] = []

    templates: dict[str, dict[str, Any]] = {}
    template_kinds: set[str] = set()
    for path in sorted((ROOT / "templates").glob("*.json")):
        try:
            record = load_json(path)
        except (OSError, json.JSONDecodeError) as exc:
            failures.append(f"{path.relative_to(ROOT)}: parse error: {exc}")
            continue
        templates[path.name] = record
        template_kinds.add(record.get("kind"))
        errors = validate_template(record)
        if errors:
            failures.append(f"{path.relative_to(ROOT)}: {', '.join(errors)}")
    if template_kinds != REQUIRED_TEMPLATE_KINDS:
        failures.append(
            "templates: required kinds mismatch; missing="
            + repr(sorted(REQUIRED_TEMPLATE_KINDS - template_kinds))
            + " extra=" + repr(sorted(template_kinds - REQUIRED_TEMPLATE_KINDS))
        )

    example_documents: dict[str, dict[str, Any]] = {}
    for path in sorted((ROOT / "examples").glob("*.json")):
        try:
            document = load_json(path)
        except (OSError, json.JSONDecodeError) as exc:
            failures.append(f"{path.relative_to(ROOT)}: parse error: {exc}")
            continue
        example_documents[path.name] = document
        errors = validate_document(document)
        if errors:
            failures.append(f"{path.relative_to(ROOT)}: unexpected {', '.join(errors)}")

    negative_codes: set[str] = set()
    negative_count = 0
    for path in sorted((ROOT / "fixtures" / "negative").glob("*.json")):
        negative_count += 1
        try:
            document = load_json(path)
        except (OSError, json.JSONDecodeError) as exc:
            failures.append(f"{path.relative_to(ROOT)}: parse error: {exc}")
            continue
        expected = sorted(set(document.get("expectedErrors", [])))
        try:
            source_example = document.get("sourceExample")
            if source_example:
                source_path = ROOT / "examples" / source_example
                subject = copy.deepcopy(example_documents[source_path.name])
                apply_mutations(subject, document.get("mutations", []))
                if document.get("rebindReservationDigests") is True:
                    rebind_reservation_digests(subject)
                actual = validate_document(subject)
            else:
                if document.get("autoBindReservationDigests") is True:
                    rebind_reservation_digests(document)
                actual = validate_document(document, isolated_fixture=True)
        except (KeyError, IndexError, TypeError, ValueError) as exc:
            failures.append(f"{path.relative_to(ROOT)}: invalid mutation fixture: {exc}")
            continue
        negative_codes.update(expected)
        if actual != expected:
            failures.append(f"{path.relative_to(ROOT)}: expected {expected!r}, got {actual!r}")

    failures.extend(check_coverage(example_documents, negative_codes))

    assignment_path = ROOT / "Steering" / "Assignments" / "ISSUE-007.yaml"
    registry_path = ROOT / "Steering" / "WorkDomains.json"
    reservation_path = ROOT / "Steering" / "Reservations" / "RSV-ISSUE-007-001.json"
    study_path = ROOT / "Studies" / "STU-007.json"
    try:
        assignment = load_json(assignment_path)
        registry = load_json(registry_path)
        reservation = load_json(reservation_path)
        study = load_json(study_path)
        dogfood = {
            "experimental": True,
            "registries": [registry],
            "records": [study],
            "relations": [],
            "assignments": [assignment],
            "reservationSets": [reservation],
        }
        errors = validate_document(dogfood)
        if errors:
            failures.append(f"Steering dogfood: {', '.join(errors)}")
        declared_reservation = assignment.get("coordination", {}).get("reservationSet", {})
        expected_reservation_path = "SDP-vnext-pilot/Steering/Reservations/RSV-ISSUE-007-001.json"
        if declared_reservation.get("path") != expected_reservation_path:
            failures.append("Steering reservation set: declared path mismatch")
        resolved_reservation_path = (ROOT.parent / str(declared_reservation.get("path", ""))).resolve()
        if resolved_reservation_path != reservation_path.resolve() or not resolved_reservation_path.is_file():
            failures.append("Steering reservation set: declared path does not resolve to the validated file")
        issued_members = registry.get("domains", [{}])[0].get("issuedIds", [])
        study_sources = [
            member.get("source")
            for member in issued_members
            if isinstance(member, dict) and member.get("id") == study.get("id")
        ]
        if study_sources != ["SDP-vnext-pilot/Studies/STU-007.json"]:
            failures.append("Steering dogfood: Study inventory source mismatch")
        else:
            resolved_study_path = (ROOT.parent / study_sources[0]).resolve()
            if resolved_study_path != study_path.resolve() or not resolved_study_path.is_file():
                failures.append("Steering dogfood: Study inventory source does not resolve")
    except (OSError, json.JSONDecodeError) as exc:
        failures.append(f"Steering dogfood parse error: {exc}")

    failures.extend(validate_text_corpus())

    if failures:
        print("SDP vNext pilot validation: FAIL")
        for failure in failures:
            print(f"- {failure}")
        return 1

    print("SDP vNext pilot validation: PASS")
    print(f"- templates: {len(templates)} ({', '.join(sorted(template_kinds))})")
    print(f"- positive examples: {len(example_documents)}")
    print(f"- negative fixtures: {negative_count}")
    print(f"- negative diagnostic coverage: {len(negative_codes)} codes")
    print("- Steering Issue #7 Study/assignment/domain/reservation binding: valid")
    print("- local Markdown links/status markers/trailing whitespace: valid")
    return 0


if __name__ == "__main__":
    sys.exit(main())
