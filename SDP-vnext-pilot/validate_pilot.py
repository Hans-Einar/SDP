#!/usr/bin/env python3
"""Deterministic standard-library validation for the provisional vNext pilot."""

from __future__ import annotations

import json
import hashlib
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
    "issue-assignment", "work-domain-registry",
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
        "slices", "acceptedEvidence", "releaseInclusions",
    },
    "refactor": {
        "schemaVersion", "experimental", "kind", "domainUid", "id",
        "revision", "declaredState", "title", "behaviorBaseline",
        "targetStructure", "compatibility", "temporaryAdapters", "exitEvidence",
        "relations", "issueAuthorities", "slices",
    },
    "fix": {
        "schemaVersion", "experimental", "kind", "domainUid", "id",
        "revision", "declaredState", "title", "defectEvidence", "correction",
        "risk", "invariants", "affectedWork", "standaloneReviewedUnit",
        "slices", "verificationCriteria", "reviewCriteria",
    },
    "study": {
        "schemaVersion", "experimental", "kind", "domainUid", "id",
        "revision", "declaredState", "title", "question", "evidenceBoundary",
        "ownerRef", "independentStudyRefs", "convergenceGate", "findings",
        "limitations", "informs",
    },
    "slice": {
        "schemaVersion", "experimental", "kind", "domainUid", "id",
        "revision", "declaredState", "title", "ownerRef", "assignmentIssue",
        "outcome", "whySmallestCoherent", "decisionRefs", "ownedPaths",
        "sharedTouchpoints", "invariants", "nonGoals", "verificationCriteria",
        "reviewCriteria", "discoveryRule", "completionSignal", "hardStop",
    },
    "issue-assignment": {
        "schemaVersion", "experimental", "kind", "revision", "authority",
        "workRef", "activeSlices", "baseline", "delivery", "coordination",
        "boundaries", "requiredEvidence", "stopCondition",
    },
    "work-domain-registry": {
        "schemaVersion", "experimental", "kind", "repository",
        "defaultDomainUid", "domains",
    },
}
WORK_OWNER_KINDS = {"feature", "refactor", "fix"}


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
    path = value.replace("\\", "/")
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
        if any(ord(char) < 32 for char in part) or ":" in part:
            return None
        normalized = portable_text(part)
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
                if normalized is None:
                    add(errors, "INVALID_PATH")
                else:
                    root_entries.append((uid, normalized))
            by_uid[uid].append((registry.get("repository", ""), domain))
        for index, (left_uid, left_path) in enumerate(root_entries):
            for right_uid, right_path in root_entries[index + 1:]:
                if left_uid != right_uid and paths_overlap(left_path, right_path):
                    add(errors, "DOMAIN_ROOT_COLLISION")

    for entries in by_uid.values():
        active_entries = [item for item in entries if item[1].get("state") == "active"]
        if len(active_entries) > 1:
            add(errors, "MULTIPLE_ACTIVE_DOMAIN_HOSTS")
        keys = {item[1].get("key") for item in entries}
        if len(keys) > 1:
            add(errors, "MOVE_KEY_REWRITE")
        active_inventory = set(active_entries[0][1].get("issuedIds", [])) if len(active_entries) == 1 else set()
        for _, domain in entries:
            if domain.get("state") == "moved" and not set(domain.get("issuedIds", [])).issubset(active_inventory):
                add(errors, "MOVE_ID_INVENTORY_REWRITE")
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
        for record_id in domain.get("issuedIds", [])
        if isinstance(record_id, str)
    }


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
    elif kind == "work-domain-registry":
        if not isinstance(record.get("domains"), list) or not record.get("domains"):
            add(errors, code)


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
        if kind in WORK_OWNER_KINDS and "schemaVersion" in record:
            for slice_ref in record.get("slices", []):
                slice_key = reference_key(slice_ref)
                if slice_key is None:
                    add(errors, "INVALID_REFERENCE")
                    continue
                slice_record = records.get(slice_key)
                if slice_record is None:
                    continue
                if slice_record.get("kind") != "slice" or reference_key(slice_record.get("ownerRef")) != key:
                    add(errors, "NONRECIPROCAL_SLICE_OWNER")


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


def validate_assignments(document: dict[str, Any], errors: list[str]) -> None:
    assignments = document.get("assignments", [])
    domain_by_uid = active_domain_map(document)
    records = {
        reference_key(record): record
        for record in document.get("records", [])
        if reference_key(record) is not None
    }
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
        authority = assignment.get("authority", {}).get("issue", "")
        coordination = assignment.get("coordination", {})
        base = coordination.get("integrationBase")
        if base is not None and not SHA_RE.fullmatch(str(base)):
            add(errors, "INVALID_INTEGRATION_BASE")
        digest = coordination.get("reservationSet", {}).get("digest")
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


def validate_document(document: dict[str, Any]) -> list[str]:
    errors: list[str] = []
    if document.get("experimental") is not True:
        add(errors, "EXPERIMENTAL_MARKER_REQUIRED")
    validate_registry_set(document, errors)
    validate_records(document, errors)
    validate_record_graph(document, errors)
    validate_relations(document, errors)
    validate_assignments(document, errors)
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
    elif kind == "issue-assignment":
        baseline = record.get("baseline", {}).get("commit", "")
        integration = record.get("coordination", {}).get("integrationBase", "")
        if not SHA_RE.fullmatch(str(baseline)) or baseline != integration:
            add(errors, "INVALID_INTEGRATION_BASE")
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
    if len(scoped.get("assignments", [])) < 2:
        failures.append("coverage: concurrent assignment example missing")
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
    }
    missing = sorted(required_negative - negative_codes)
    if missing:
        failures.append("coverage: missing negative diagnostics " + ", ".join(missing))
    return failures


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
        actual = validate_document(document)
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
        }
        errors = validate_document(dogfood)
        if errors:
            failures.append(f"Steering dogfood: {', '.join(errors)}")
        declared_reservation = assignment.get("coordination", {}).get("reservationSet", {})
        if reservation.get("experimental") is not True or reservation.get("kind") != "reservation-set":
            failures.append("Steering reservation set: invalid pilot marker/kind")
        if declared_reservation.get("id") != reservation.get("id"):
            failures.append("Steering reservation set: assignment ID mismatch")
        expected_reservation_path = "SDP-vnext-pilot/Steering/Reservations/RSV-ISSUE-007-001.json"
        if declared_reservation.get("path") != expected_reservation_path:
            failures.append("Steering reservation set: declared path mismatch")
        resolved_reservation_path = (ROOT.parent / str(declared_reservation.get("path", ""))).resolve()
        if resolved_reservation_path != reservation_path.resolve() or not resolved_reservation_path.is_file():
            failures.append("Steering reservation set: declared path does not resolve to the validated file")
        if declared_reservation.get("digest") != canonical_digest(reservation):
            failures.append("Steering reservation set: canonical SHA-256 mismatch")
        if assignment.get("coordination", {}).get("integrationBase") != reservation.get("integrationBase"):
            failures.append("Steering reservation set: integration-base mismatch")
        reservation_assignment = reservation.get("assignments", [{}])[0]
        if assignment.get("authority", {}).get("issue") != reservation_assignment.get("issue"):
            failures.append("Steering reservation set: Issue mismatch")
        assignment_work_identity = {
            "domainUid": assignment.get("workRef", {}).get("domainUid"),
            "id": assignment.get("workRef", {}).get("id"),
        }
        if assignment_work_identity != reservation_assignment.get("workRef"):
            failures.append("Steering reservation set: work-reference mismatch")
        if assignment.get("workRef", {}).get("domainUid") != reservation_assignment.get("domainUid"):
            failures.append("Steering reservation set: domain mismatch")
        if assignment.get("coordination", {}).get("ownedPaths") != reservation_assignment.get("ownedPaths"):
            failures.append("Steering reservation set: owned-path mismatch")
        if assignment.get("coordination", {}).get("sharedTouchpoints") != reservation_assignment.get("sharedTouchpoints"):
            failures.append("Steering reservation set: shared-touchpoint mismatch")
        assignment_reserved_ids = [
            reference.get("id")
            for reference in assignment.get("coordination", {}).get("reservedIds", [])
        ]
        if assignment_reserved_ids != reservation_assignment.get("reservedIds"):
            failures.append("Steering reservation set: reserved-ID mismatch")
        reservation_validation_errors: list[str] = []
        dogfood_domains = active_domain_map(dogfood)
        for record_id in reservation_assignment.get("reservedIds", []):
            validate_id_against_domain(
                reservation_assignment.get("domainUid"), record_id,
                dogfood_domains, reservation_validation_errors,
                invalid_code="INVALID_RESERVED_ID",
                style_code="RESERVED_ID_STYLE_MISMATCH",
            )
        if reservation_validation_errors:
            failures.append("Steering reservation set: " + ", ".join(sorted(reservation_validation_errors)))
        expected_record_path = "SDP-vnext-pilot/Studies/STU-007.json"
        if reservation_assignment.get("recordPaths") != [expected_record_path]:
            failures.append("Steering reservation set: dogfood record-path mismatch")
        else:
            resolved_record_path = (ROOT.parent / expected_record_path).resolve()
            if resolved_record_path != study_path.resolve() or not resolved_record_path.is_file():
                failures.append("Steering reservation set: dogfood record path does not resolve")
        dependency_rows = {
            row.get("issue"): row.get("dependsOn", [])
            for row in reservation.get("dependencies", [])
        }
        if dependency_rows.get(assignment.get("authority", {}).get("issue")) != assignment.get("coordination", {}).get("dependsOnIssues"):
            failures.append("Steering reservation set: dependency mismatch")
        if reservation.get("conflicts") != assignment.get("coordination", {}).get("conflictsWithIssues"):
            failures.append("Steering reservation set: conflict mismatch")
        if reservation.get("mergeOrder") != assignment.get("coordination", {}).get("mergeOrder"):
            failures.append("Steering reservation set: merge-order mismatch")
        if reservation.get("convergence") != assignment.get("coordination", {}).get("convergence"):
            failures.append("Steering reservation set: convergence mismatch")
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
