#!/usr/bin/env python3
"""Deterministic standard-library validation for the provisional vNext pilot."""

from __future__ import annotations

import json
import hashlib
import copy
import re
import subprocess
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
SCOPED_ID_RE = re.compile(r"^(?P<key>[A-Z][A-Z0-9]{1,11})-(?P<type>FEAT|REF|FIX|STU|SLC|SPR|ITR)-[0-9]{3}$")
UNSCOPED_ID_RE = re.compile(r"^(?P<type>FEAT|REF|FIX|STU|SLC|SPR|ITR)-[0-9]{3}$")
SHA_RE = re.compile(r"^[0-9a-f]{40}$")
DIGEST_RE = re.compile(r"^sha256:[0-9a-f]{64}$")
REPOSITORY_RE = re.compile(
    r"^https://github\.com/(?P<owner>[A-Za-z0-9](?:[A-Za-z0-9-]{0,37}[A-Za-z0-9])?)/"
    r"(?P<repo>[A-Za-z0-9_.-]+)$"
)
ISSUE_RE = re.compile(
    r"^https://github\.com/(?P<owner>[A-Za-z0-9](?:[A-Za-z0-9-]{0,37}[A-Za-z0-9])?)/"
    r"(?P<repo>[A-Za-z0-9_.-]+)/issues/(?P<number>[1-9][0-9]*)$"
)
PULL_REQUEST_RE = re.compile(
    r"^https://github\.com/(?P<owner>[A-Za-z0-9](?:[A-Za-z0-9-]{0,37}[A-Za-z0-9])?)/"
    r"(?P<repo>[A-Za-z0-9_.-]+)/pull/(?P<number>[1-9][0-9]*)$"
)
ISSUE_COMMENT_RE = re.compile(
    r"^https://github\.com/(?P<owner>[A-Za-z0-9](?:[A-Za-z0-9-]{0,37}[A-Za-z0-9])?)/"
    r"(?P<repo>[A-Za-z0-9_.-]+)/issues/(?P<number>[1-9][0-9]*)"
    r"#issuecomment-(?P<comment>[1-9][0-9]*)$"
)
EVIDENCE_ID_RE = re.compile(r"^[A-Za-z0-9][A-Za-z0-9._-]*$")
FORBIDDEN_PATH_CATEGORIES = {"Cc", "Cf", "Cs", "Co"}
WINDOWS_FORBIDDEN_COMPONENT_CHARS = set('<>:"|?*')
WINDOWS_RESERVED = {
    "con", "prn", "aux", "nul",
    *(f"com{i}" for i in range(1, 10)),
    *(f"lpt{i}" for i in range(1, 10)),
}
SCHEMA_VERSIONS = {
    "feature": "sdp-vnext-pilot-feature-v0",
    "refactor": "sdp-vnext-pilot-refactor-v0",
    "fix": "sdp-vnext-pilot-fix-v0",
    "study": "sdp-vnext-pilot-study-v0",
    "slice": "sdp-vnext-pilot-slice-v0",
    "issue-assignment": "sdp-vnext-pilot-assignment-v0",
    "reservation-set": "sdp-vnext-pilot-reservation-set-v0",
    "work-domain-registry": "sdp-vnext-pilot-work-domains-v0",
}
SCHEMA_KIND_BY_VERSION = {version: kind for kind, version in SCHEMA_VERSIONS.items()}
ASSIGNMENT_HISTORY_SCHEMA_VERSION = "sdp-vnext-pilot-assignment-history-v0"
SEMANTIC_RELATION_TYPES = {
    "depends_on", "informs", "refines", "supersedes", "preserves",
    "requires_revision", "owned_by", "independent_of", "corrects",
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
DECLARED_STATES_BY_KIND = {
    "feature": {
        "proposed", "studying", "ready", "active", "blocked", "rejected",
        "superseded", "delivered", "released",
    },
    "refactor": {
        "proposed", "studying", "ready", "active", "blocked", "rejected",
        "superseded", "delivered", "released",
    },
    "fix": {
        "proposed", "ready", "active", "blocked", "rejected", "superseded",
        "delivered", "released",
    },
    "study": {"proposed", "active", "blocked", "rejected", "superseded", "accepted"},
    "slice": {"proposed", "active", "blocked", "rejected", "superseded", "accepted"},
    "issue-assignment": {
        "proposed", "active", "blocked", "accepted", "rejected", "cancelled",
        "superseded",
    },
}
MAXIMUM_SEVERITIES = {"none", "low", "medium", "high", "blocking"}
TYPE_REQUIRED_FIELDS = {
    "feature": {
        "schemaVersion", "experimental", "kind", "domainUid", "id",
        "source", "revision", "declaredState", "title", "intent", "scope", "nonGoals",
        "constraints", "acceptanceCriteria", "relations", "issueAuthorities",
        "slices", "acceptedEvidence",
    },
    "refactor": {
        "schemaVersion", "experimental", "kind", "domainUid", "id",
        "source", "revision", "declaredState", "title", "behaviorBaseline",
        "targetStructure", "compatibility", "temporaryAdapters", "exitEvidence",
        "relations", "issueAuthorities", "slices", "acceptedEvidence",
    },
    "fix": {
        "schemaVersion", "experimental", "kind", "domainUid", "id",
        "source", "revision", "declaredState", "title", "defectEvidence", "correction",
        "risk", "invariants", "affectedWork", "standaloneReviewedUnit",
        "issueAuthorities", "slices", "verificationCriteria", "reviewCriteria",
        "acceptedEvidence",
    },
    "study": {
        "schemaVersion", "experimental", "kind", "domainUid", "id",
        "source", "revision", "declaredState", "title", "question", "evidenceBoundary",
        "ownerRef", "issueAuthorities", "independentStudyRefs",
        "convergenceGate", "findings", "limitations", "informs",
        "acceptedEvidence",
    },
    "slice": {
        "schemaVersion", "experimental", "kind", "domainUid", "id",
        "source", "revision", "declaredState", "title", "ownerRef", "assignmentIssue",
        "outcome", "whySmallestCoherent", "decisionRefs", "ownedPaths",
        "sharedTouchpoints", "invariants", "nonGoals", "verificationCriteria",
        "reviewCriteria", "discoveryRule", "completionSignal", "hardStop",
        "acceptedEvidence",
    },
    "issue-assignment": {
        "schemaVersion", "experimental", "kind", "source", "revision", "declaredState", "authority",
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


def contains_invalid_unicode_scalar(value: Any) -> bool:
    """Reject Python surrogate code points before canonical UTF-8 encoding."""
    if isinstance(value, str):
        return any(unicodedata.category(character) == "Cs" for character in value)
    if isinstance(value, dict):
        return any(
            contains_invalid_unicode_scalar(key)
            or contains_invalid_unicode_scalar(item)
            for key, item in value.items()
        )
    if isinstance(value, (list, tuple)):
        return any(contains_invalid_unicode_scalar(item) for item in value)
    return False


def canonical_digest(document: Any) -> str | None:
    if contains_invalid_unicode_scalar(document):
        return None
    encoded = json.dumps(
        document, sort_keys=True, separators=(",", ":"), ensure_ascii=False
    ).encode("utf-8")
    return "sha256:" + hashlib.sha256(encoded).hexdigest()


def portable_text(value: str) -> str:
    return unicodedata.normalize("NFKC", value).casefold()


def contains_nonportable_unicode(value: str) -> bool:
    return any(
        unicodedata.category(character) in FORBIDDEN_PATH_CATEGORIES
        for character in value
    )


def valid_uid(value: Any) -> bool:
    if not isinstance(value, str) or not UUID_URN_RE.fullmatch(value):
        return False
    try:
        return str(uuid.UUID(value.removeprefix("urn:uuid:"))) == value.removeprefix("urn:uuid:")
    except ValueError:
        return False


def nonblank(value: Any) -> bool:
    return isinstance(value, str) and bool(value.strip())


def nonblank_string_list(value: Any, *, nonempty: bool = False) -> bool:
    return (
        isinstance(value, list)
        and (not nonempty or bool(value))
        and all(nonblank(item) for item in value)
        and len(value) == len({portable_text(item.strip()) for item in value})
    )


def repository_key(value: Any) -> tuple[str, str] | None:
    if not isinstance(value, str):
        return None
    match = REPOSITORY_RE.fullmatch(value)
    if match is None or match.group("repo").casefold().endswith(".git"):
        return None
    return match.group("owner").casefold(), match.group("repo").casefold()


def issue_key(value: Any) -> tuple[str, str, int] | None:
    if not isinstance(value, str):
        return None
    match = ISSUE_RE.fullmatch(value)
    if match is None:
        return None
    return (
        match.group("owner").casefold(),
        match.group("repo").casefold(),
        int(match.group("number")),
    )


def pull_request_key(value: Any) -> tuple[str, str, int] | None:
    if not isinstance(value, str):
        return None
    match = PULL_REQUEST_RE.fullmatch(value)
    if match is None:
        return None
    return (
        match.group("owner").casefold(),
        match.group("repo").casefold(),
        int(match.group("number")),
    )


def issue_comment_key(value: Any) -> tuple[str, str, int, int] | None:
    if not isinstance(value, str):
        return None
    match = ISSUE_COMMENT_RE.fullmatch(value)
    if match is None:
        return None
    return (
        match.group("owner").casefold(),
        match.group("repo").casefold(),
        int(match.group("number")),
        int(match.group("comment")),
    )


def repository_for_domain(
    document: dict[str, Any], uid: Any
) -> tuple[str, str] | None:
    active_hosts = {
        repository_key(registry.get("repository"))
        for registry in document.get("registries", [])
        for domain in registry.get("domains", [])
        if domain.get("domainUid") == uid and domain.get("state") == "active"
    }
    active_hosts.discard(None)
    return next(iter(active_hosts)) if len(active_hosts) == 1 else None


def valid_branch(value: Any) -> bool:
    if not nonblank(value) or len(value.encode("utf-8")) > 255:
        return False
    if re.fullmatch(r"[A-Za-z0-9._/-]+", value) is None:
        return False
    if value in {"@", "HEAD"} or value.startswith(("-", "/", ".")):
        return False
    if value.endswith(("/", ".")) or "//" in value or ".." in value or "@{" in value:
        return False
    if any(character.isspace() or ord(character) < 32 for character in value):
        return False
    if any(character in value for character in "~^:?*[\\"):
        return False
    return all(not part.startswith(".") and not part.endswith(".lock") for part in value.split("/"))


def normalize_path(value: Any) -> tuple[str, bool] | None:
    if not isinstance(value, str) or not value or contains_nonportable_unicode(value):
        return None
    semantic_chars = '<>:"/\\|?*[]{}'
    for character in value:
        normalized_character = unicodedata.normalize("NFKC", character)
        if character not in semantic_chars and any(
            token in normalized_character for token in semantic_chars
        ):
            return None
    path = unicodedata.normalize("NFKC", value)
    if contains_nonportable_unicode(path):
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
        if contains_nonportable_unicode(normalized):
            return None
        if any(char in normalized for char in WINDOWS_FORBIDDEN_COMPONENT_CHARS | set("/\\[]{}")):
            return None
        if normalized in {".", ".."}:
            return None
        if normalized.rstrip(" .") != normalized:
            return None
        if normalized.split(".", 1)[0] in WINDOWS_RESERVED:
            return None
        normalized_parts.append(normalized)
    return "/".join(normalized_parts), recursive


def portable_path_identity(value: Any) -> str | None:
    """Return a collision identity even for slash/case aliases.

    This is deliberately separate from validity: a backslash-authored record
    path is invalid, but it must still collide with its forward-slash alias.
    """
    if not isinstance(value, str) or not value:
        return None
    normalized = unicodedata.normalize("NFKC", value).replace("\\", "/")
    return normalized.casefold()


def valid_record_path(value: Any) -> bool:
    normalized = normalize_path(value)
    return (
        isinstance(value, str)
        and normalized is not None
        and not normalized[1]
        and "\\" not in value
        and unicodedata.normalize("NFKC", value) == value
    )


def resolve_local_path(repository_root: Path, value: Any) -> Path | None:
    if not valid_record_path(value):
        return None
    candidate = (repository_root / value).resolve()
    try:
        candidate.relative_to(repository_root.resolve())
    except ValueError:
        return None
    return candidate


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


def mapping(value: Any) -> dict[str, Any]:
    return value if isinstance(value, dict) else {}


def effective_kind(record: Any) -> str | None:
    if not isinstance(record, dict):
        return None
    kind = record.get("kind")
    if kind in SCHEMA_VERSIONS:
        return kind
    return SCHEMA_KIND_BY_VERSION.get(record.get("schemaVersion"))


def assignment_issue(assignment: Any) -> Any:
    return mapping(mapping(assignment).get("authority")).get("issue")


def assignment_coordination(assignment: Any) -> dict[str, Any]:
    return mapping(mapping(assignment).get("coordination"))


def declarations(document: dict[str, Any]) -> list[tuple[str, dict[str, Any]]]:
    result: list[tuple[str, dict[str, Any]]] = []
    registries = document.get("registries", [])
    for registry in registries if isinstance(registries, list) else []:
        if not isinstance(registry, dict):
            continue
        repository = registry.get("repository", "")
        domains = registry.get("domains", [])
        for domain in domains if isinstance(domains, list) else []:
            if isinstance(domain, dict):
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
    if not nonblank(record_id) or status not in {
        "prospective", "legacy-preserved"
    }:
        add(errors, "INVALID_ID_INVENTORY_MEMBER")
        return
    if not valid_record_path(source):
        add(errors, "INVALID_ID_INVENTORY_MEMBER")
    if status == "prospective":
        if issue_key(authority) is None:
            add(errors, "INVALID_INVENTORY_AUTHORITY")
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
            repository_key(repository) is None
            or not isinstance(commit, str)
            or not SHA_RE.fullmatch(commit)
            or not valid_record_path(path)
            or source != path
        ):
            add(errors, "LEGACY_PROVENANCE_REQUIRED")
        if authority is None:
            if not nonblank(member.get("authorityMissingReason")):
                add(errors, "LEGACY_AUTHORITY_REASON_REQUIRED")
        elif issue_key(authority) is None:
            add(errors, "INVALID_INVENTORY_AUTHORITY")


def validate_registry_set(document: dict[str, Any], errors: list[str]) -> None:
    registries = document.get("registries", [])
    if not isinstance(registries, list):
        add(errors, "INVALID_REGISTRY_COLLECTION")
        return
    repository_keys: list[tuple[str, str] | None] = []
    for registry in registries:
        if not isinstance(registry, dict):
            add(errors, "REGISTRY_REQUIRED_FIELD_MISSING")
            repository_keys.append(None)
            continue
        key = repository_key(registry.get("repository"))
        repository_keys.append(key)
        if key is None:
            add(errors, "INVALID_GITHUB_REPOSITORY")
    valid_repository_keys = [key for key in repository_keys if key is not None]
    if len(valid_repository_keys) != len(set(valid_repository_keys)):
        # Pilot v0 deliberately chooses one registry object per canonical repo.
        # Partitioned registries are rejected before domain/default/root checks.
        add(errors, "DUPLICATE_REPOSITORY_REGISTRY")
        return

    by_uid: dict[str, list[tuple[str, dict[str, Any]]]] = defaultdict(list)
    for registry in registries:
        validate_typed_shape(
            registry, errors, "REGISTRY_REQUIRED_FIELD_MISSING",
            expected_kind="work-domain-registry",
        )
        domains = registry.get("domains")
        if not isinstance(domains, list) or not domains or any(
            not isinstance(domain, dict) for domain in domains
        ):
            add(errors, "REGISTRY_REQUIRED_FIELD_MISSING")
            continue
        per_repo_keys: dict[str, str] = {}
        active = [domain for domain in domains if domain.get("state") == "active"]
        if len(active) > 1 and any(domain.get("newRecordIdStyle") != "scoped" for domain in active):
            add(errors, "MULTIDOMAIN_REQUIRES_SCOPED_STYLE")
        default_uid = registry.get("defaultDomainUid")
        if default_uid is not None and default_uid not in {domain.get("domainUid") for domain in active}:
            add(errors, "INVALID_DEFAULT_DOMAIN")
        root_entries: list[tuple[str, tuple[str, bool]]] = []
        for domain in domains:
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
            owners = domain.get("owners")
            if not nonblank_string_list(owners, nonempty=True):
                add(errors, "DOMAIN_OWNER_REQUIRED")
            if domain.get("newRecordIdStyle") not in {"scoped", "unscoped"}:
                add(errors, "INVALID_ID_STYLE")
            state = domain.get("state")
            roots = domain.get("roots", [])
            if not isinstance(roots, list):
                add(errors, "INVALID_PATH")
                roots = []
            if state == "active" and not roots:
                add(errors, "ACTIVE_DOMAIN_ROOT_REQUIRED")
            if state == "moved":
                if (
                    roots
                    or repository_key(domain.get("successorRepository")) is None
                    or not valid_record_path(domain.get("successorRegistry"))
                ):
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
            active_repositories = {repository_key(item[0]) for item in active_entries}
            if repository_key(domain.get("successorRepository")) not in active_repositories:
                add(errors, "MOVE_SUCCESSOR_MISMATCH")
        if len(entries) > 1 and active_entries:
            _, active_domain = active_entries[0]
            predecessor_values = active_domain.get("predecessorRepositories", [])
            if not isinstance(predecessor_values, list) or any(
                repository_key(item) is None for item in predecessor_values
            ):
                add(errors, "INVALID_GITHUB_REPOSITORY")
                predecessor_values = []
            predecessors = {repository_key(item) for item in predecessor_values}
            moved_repositories = {
                repository_key(item[0])
                for item in entries if item[1].get("state") == "moved"
            }
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


def domain_key_for_uid(document: dict[str, Any], uid: Any) -> str | None:
    keys = {
        domain.get("key")
        for _, domain in declarations(document)
        if domain.get("domainUid") == uid and isinstance(domain.get("key"), str)
    }
    return next(iter(keys)) if len(keys) == 1 else None


def validate_key_hint(
    document: dict[str, Any], reference: Any, errors: list[str]
) -> None:
    if not isinstance(reference, dict) or "keyHint" not in reference:
        return
    expected = domain_key_for_uid(document, reference.get("domainUid"))
    if not nonblank(reference.get("keyHint")) or expected is None or reference.get("keyHint") != expected:
        add(errors, "KEY_HINT_MISMATCH")


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


def issued_authorities(
    document: dict[str, Any]
) -> dict[tuple[str, str], tuple[str, str, int]]:
    result: dict[tuple[str, str], tuple[str, str, int]] = {}
    for _, domain in declarations(document):
        uid = domain.get("domainUid")
        for member in domain.get("issuedIds", []):
            if not isinstance(member, dict):
                continue
            record_id = inventory_member_id(member)
            authority = member.get("authorityIssue")
            authority_identity = issue_key(authority)
            if isinstance(record_id, str) and authority_identity is not None:
                result[(uid, portable_text(record_id))] = authority_identity
    return result


def issued_statuses(document: dict[str, Any]) -> dict[tuple[str, str], Any]:
    result: dict[tuple[str, str], Any] = {}
    for _, domain in declarations(document):
        uid = domain.get("domainUid")
        for member in domain.get("issuedIds", []):
            if not isinstance(member, dict):
                continue
            record_id = inventory_member_id(member)
            if isinstance(record_id, str):
                result[(uid, portable_text(record_id))] = member.get("status")
    return result


def validate_typed_shape(
    record: dict[str, Any], errors: list[str], code: str,
    *, expected_kind: str | None = None,
) -> None:
    authored_kind = record.get("kind")
    kind = expected_kind or effective_kind(record)
    if "kind" not in record:
        add(errors, "KIND_MARKER_REQUIRED")
    elif expected_kind is not None and authored_kind != expected_kind:
        add(errors, "KIND_MARKER_MISMATCH")
    if "schemaVersion" not in record:
        add(errors, "SCHEMA_VERSION_REQUIRED")
    elif kind is not None and record.get("schemaVersion") != SCHEMA_VERSIONS.get(kind):
        add(errors, "UNSUPPORTED_SCHEMA_VERSION")
    if "experimental" not in record:
        add(errors, "EXPERIMENTAL_MARKER_REQUIRED")
    elif record.get("experimental") is not True:
        add(errors, "EXPERIMENTAL_MARKER_REQUIRED")
    fields = TYPE_REQUIRED_FIELDS.get(kind)
    if fields is None:
        if authored_kind is not None:
            add(errors, "UNKNOWN_RECORD_KIND")
        return
    if any(
        field not in record
        for field in fields - {"schemaVersion", "experimental", "kind"}
    ):
        add(errors, code)
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


def validate_issue_authorities(value: Any, errors: list[str]) -> None:
    if not isinstance(value, list):
        add(errors, "INVALID_REQUIRED_COLLECTION")
        return
    keys: list[tuple[str, str, int]] = []
    for authority in value:
        key = issue_key(authority)
        if key is None:
            add(errors, "INVALID_GITHUB_ISSUE")
        else:
            keys.append(key)
    if len(keys) != len(set(keys)):
        add(errors, "DUPLICATE_ISSUE_AUTHORITY")


def validate_stateful_values(
    record: dict[str, Any], errors: list[str], *, kind_override: str | None = None
) -> None:
    kind = kind_override or effective_kind(record)
    revision = record.get("revision")
    if isinstance(revision, bool) or not isinstance(revision, int) or revision < 1:
        add(errors, "INVALID_REVISION")
    if record.get("declaredState") not in DECLARED_STATES_BY_KIND.get(kind, set()):
        add(errors, "UNSUPPORTED_DECLARED_STATE")
    if not valid_record_path(record.get("source")):
        add(errors, "INVALID_RECORD_SOURCE")

    if kind in {"feature", "refactor", "fix", "study", "slice"}:
        if not nonblank(record.get("title")):
            add(errors, "BLANK_REQUIRED_VALUE")
    if kind in {"feature", "refactor", "fix", "study"}:
        validate_issue_authorities(record.get("issueAuthorities"), errors)

    if kind == "feature":
        intent = record.get("intent")
        if not isinstance(intent, dict) or any(
            not nonblank(intent.get(field))
            for field in ("problem", "outcome", "usersOrValue")
        ):
            add(errors, "BLANK_REQUIRED_VALUE")
        for field, required in (
            ("scope", True), ("nonGoals", False), ("constraints", False),
            ("acceptanceCriteria", True),
        ):
            if not nonblank_string_list(record.get(field), nonempty=required):
                add(errors, "INVALID_REQUIRED_COLLECTION")
        for field in ("relations", "slices"):
            if not isinstance(record.get(field), list):
                add(errors, "INVALID_REQUIRED_COLLECTION")
    elif kind == "refactor":
        for field, required in (
            ("behaviorBaseline", True), ("targetStructure", True),
            ("temporaryAdapters", False), ("exitEvidence", True),
        ):
            if not nonblank_string_list(record.get(field), nonempty=required):
                add(errors, "INVALID_REQUIRED_COLLECTION")
        compatibility = record.get("compatibility")
        if not isinstance(compatibility, dict) or not nonblank(compatibility.get("migrationRange")):
            add(errors, "BLANK_REQUIRED_VALUE")
        elif (
            not nonblank_string_list(compatibility.get("preserved"))
            or not nonblank_string_list(compatibility.get("changed"))
            or not (compatibility.get("preserved") or compatibility.get("changed"))
        ):
            add(errors, "INVALID_REQUIRED_COLLECTION")
        for field in ("relations", "slices"):
            if not isinstance(record.get(field), list):
                add(errors, "INVALID_REQUIRED_COLLECTION")
    elif kind == "fix":
        if not nonblank(record.get("correction")):
            add(errors, "BLANK_REQUIRED_VALUE")
        for field, required in (
            ("defectEvidence", True), ("invariants", True),
            ("affectedWork", False), ("verificationCriteria", True),
            ("reviewCriteria", True),
        ):
            value = record.get(field)
            valid = (
                isinstance(value, list)
                and (not required or bool(value))
                and all(nonblank(item) if isinstance(item, str) else isinstance(item, dict) for item in value)
            )
            if not valid:
                add(errors, "INVALID_REQUIRED_COLLECTION")
        if record.get("risk") not in {"low", "medium", "high", "safety-critical"}:
            add(errors, "INVALID_FIX_PROFILE")
        if not isinstance(record.get("standaloneReviewedUnit"), bool):
            add(errors, "INVALID_FIX_PROFILE")
        if not isinstance(record.get("slices"), list):
            add(errors, "INVALID_REQUIRED_COLLECTION")
    elif kind == "study":
        if not nonblank(record.get("question")) or not nonblank(record.get("convergenceGate")):
            add(errors, "BLANK_REQUIRED_VALUE")
        boundary = record.get("evidenceBoundary")
        if (
            not isinstance(boundary, dict)
            or not nonblank_string_list(boundary.get("inScope"), nonempty=True)
            or not nonblank_string_list(boundary.get("outOfScope"))
        ):
            add(errors, "INVALID_REQUIRED_COLLECTION")
        for field in ("independentStudyRefs", "informs"):
            if not isinstance(record.get(field), list):
                add(errors, "INVALID_REQUIRED_COLLECTION")
        if record.get("ownerRef") is not None and reference_key(record.get("ownerRef")) is None:
            add(errors, "INVALID_REFERENCE")
        for field in ("findings", "limitations"):
            if not nonblank_string_list(record.get(field)):
                add(errors, "INVALID_REQUIRED_COLLECTION")
    elif kind == "slice":
        for field in (
            "outcome", "whySmallestCoherent", "discoveryRule",
            "completionSignal", "hardStop",
        ):
            if not nonblank(record.get(field)):
                add(errors, "BLANK_REQUIRED_VALUE")
        if issue_key(record.get("assignmentIssue")) is None:
            add(errors, "INVALID_GITHUB_ISSUE")
        for field, required in (
            ("decisionRefs", False), ("ownedPaths", False),
            ("sharedTouchpoints", False),
        ):
            value = record.get(field)
            if not isinstance(value, list) or (required and not value):
                add(errors, "INVALID_REQUIRED_COLLECTION")
        for field, required in (
            ("invariants", True), ("nonGoals", False),
            ("verificationCriteria", True), ("reviewCriteria", True),
        ):
            if not nonblank_string_list(record.get(field), nonempty=required):
                add(errors, "INVALID_REQUIRED_COLLECTION")
    elif kind == "issue-assignment":
        authority = record.get("authority")
        issue = authority.get("issue") if isinstance(authority, dict) else None
        issue_identity = issue_key(issue)
        if issue_identity is None:
            add(errors, "INVALID_GITHUB_ISSUE")
        amendments = authority.get("amendments") if isinstance(authority, dict) else None
        if not isinstance(amendments, list):
            add(errors, "INVALID_REQUIRED_COLLECTION")
        else:
            for amendment in amendments:
                amendment_identity = issue_comment_key(amendment)
                if amendment_identity is None:
                    add(errors, "INVALID_GITHUB_COMMENT")
                elif issue_identity is not None and amendment_identity[:3] != issue_identity:
                    add(errors, "REPOSITORY_AUTHORITY_MISMATCH")
        baseline = record.get("baseline")
        delivery = record.get("delivery")
        if (
            not isinstance(baseline, dict)
            or not valid_branch(baseline.get("branch"))
            or not isinstance(delivery, dict)
            or not valid_branch(delivery.get("branch"))
        ):
            add(errors, "INVALID_BRANCH")
        if not isinstance(delivery, dict) or pull_request_key(delivery.get("pullRequest")) is None:
            add(errors, "INVALID_GITHUB_PULL_REQUEST")
        if not isinstance(record.get("activeSlices"), list):
            add(errors, "INVALID_REQUIRED_COLLECTION")
        if not nonblank(record.get("stopCondition")):
            add(errors, "INVALID_STOP_CONDITION")
        boundaries = record.get("boundaries")
        if (
            not isinstance(boundaries, dict)
            or any(
                not nonblank_string_list(boundaries.get(field))
                for field in ("prohibitedRepositories", "prohibitedPaths", "prohibitedOperations")
            )
        ):
            add(errors, "INVALID_ASSIGNMENT_BOUNDARIES")
        elif any(
            repository_key(value) is None
            for value in boundaries.get("prohibitedRepositories", [])
        ):
            add(errors, "INVALID_GITHUB_REPOSITORY")
        if isinstance(boundaries, dict):
            for value in boundaries.get("prohibitedPaths", []):
                if normalize_path(value) is None:
                    add(errors, "INVALID_PATH")
        required_evidence = record.get("requiredEvidence")
        if (
            not isinstance(required_evidence, dict)
            or not nonblank_string_list(required_evidence.get("verification"), nonempty=True)
            or not isinstance(required_evidence.get("independentReview"), bool)
            or required_evidence.get("maximumUnresolvedSeverity") not in MAXIMUM_SEVERITIES
        ):
            add(errors, "INVALID_EVIDENCE_POLICY")
        coordination = record.get("coordination")
        if isinstance(coordination, dict):
            for field in (
                "dependsOnIssues", "conflictsWithIssues", "ownedPaths",
                "sharedTouchpoints", "reservedIds", "mergeOrder",
            ):
                if not isinstance(coordination.get(field), list):
                    add(errors, "INVALID_REQUIRED_COLLECTION")
            if not nonblank(coordination.get("staleBasePolicy")):
                add(errors, "BLANK_REQUIRED_VALUE")
            reservation_reference = coordination.get("reservationSet")
            if (
                isinstance(reservation_reference, dict)
                and "path" in reservation_reference
                and not valid_record_path(reservation_reference.get("path"))
            ):
                add(errors, "INVALID_PATH")


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
    candidate = evidence.get("candidate")
    review = evidence.get("currentReviewRef")
    disposition = evidence.get("steeringDisposition")
    if candidate is None:
        if verification_refs or review is not None or disposition is not None or release_refs:
            add(errors, "ACCEPTED_EVIDENCE_SHAPE_INVALID")
    elif not isinstance(candidate, str) or not SHA_RE.fullmatch(candidate):
        add(errors, "ACCEPTED_EVIDENCE_SHAPE_INVALID")
    else:
        for reference in verification_refs:
            if (
                not isinstance(reference, dict)
                or set(reference) != {"id", "candidate"}
                or not isinstance(reference.get("id"), str)
                or EVIDENCE_ID_RE.fullmatch(reference["id"]) is None
            ):
                add(errors, "ACCEPTED_EVIDENCE_SHAPE_INVALID")
            elif reference.get("candidate") != candidate:
                add(errors, "ACCEPTED_EVIDENCE_CANDIDATE_MISMATCH")
        verification_ids = [
            reference.get("id") for reference in verification_refs
            if isinstance(reference, dict) and isinstance(reference.get("id"), str)
        ]
        if len(verification_ids) != len(set(verification_ids)):
            add(errors, "ACCEPTED_EVIDENCE_SHAPE_INVALID")
        if review is not None:
            if (
                not isinstance(review, dict)
                or set(review) != {"id", "candidate", "disposition"}
                or not isinstance(review.get("id"), str)
                or EVIDENCE_ID_RE.fullmatch(review["id"]) is None
                or review.get("disposition") != "approved"
            ):
                add(errors, "ACCEPTED_EVIDENCE_SHAPE_INVALID")
            elif review.get("candidate") != candidate:
                add(errors, "ACCEPTED_EVIDENCE_CANDIDATE_MISMATCH")
        if disposition is not None:
            disposition_identity = issue_comment_key(
                disposition.get("authority") if isinstance(disposition, dict) else None
            )
            if (
                not isinstance(disposition, dict)
                or set(disposition) != {"authority", "candidate", "decision"}
                or disposition_identity is None
                or disposition.get("decision") != "accepted"
            ):
                add(errors, "ACCEPTED_EVIDENCE_SHAPE_INVALID")
            elif disposition.get("candidate") != candidate:
                add(errors, "ACCEPTED_EVIDENCE_CANDIDATE_MISMATCH")
            else:
                if effective_kind(record) == "issue-assignment":
                    authorities = [assignment_issue(record)]
                elif effective_kind(record) == "slice":
                    authorities = [record.get("assignmentIssue")]
                else:
                    authorities = record.get("issueAuthorities", [])
                authority_keys = {
                    issue_key(authority) for authority in authorities
                } if isinstance(authorities, list) else set()
                authority_keys.discard(None)
                if authority_keys and disposition_identity[:3] not in authority_keys:
                    add(errors, "STEERING_AUTHORITY_MISMATCH")
        for reference in release_refs:
            if (
                not isinstance(reference, dict)
                or set(reference) != {"id", "candidate"}
                or not isinstance(reference.get("id"), str)
                or EVIDENCE_ID_RE.fullmatch(reference["id"]) is None
            ):
                add(errors, "ACCEPTED_EVIDENCE_SHAPE_INVALID")
            elif reference.get("candidate") != candidate:
                add(errors, "ACCEPTED_EVIDENCE_CANDIDATE_MISMATCH")
        release_ids = [
            reference.get("id") for reference in release_refs
            if isinstance(reference, dict) and isinstance(reference.get("id"), str)
        ]
        if len(release_ids) != len(set(release_ids)):
            add(errors, "ACCEPTED_EVIDENCE_SHAPE_INVALID")
    if state in TERMINAL_ACCEPTED_STATES:
        if (
            not isinstance(candidate, str)
            or not SHA_RE.fullmatch(candidate)
            or not verification_refs
            or not isinstance(review, dict)
            or review.get("disposition") != "approved"
            or not isinstance(disposition, dict)
            or disposition.get("decision") != "accepted"
        ):
            add(errors, "ACCEPTED_EVIDENCE_REQUIRED")
        if state == "released" and not release_refs:
            add(errors, "RELEASE_EVIDENCE_REQUIRED")


def has_qualified_accepted_evidence(record: Any) -> bool:
    if not isinstance(record, dict):
        return False
    evidence = record.get("acceptedEvidence")
    if not isinstance(evidence, dict) or not EVIDENCE_FIELDS.issubset(evidence):
        return False
    candidate = evidence.get("candidate")
    verification_refs = evidence.get("verificationRefs")
    review = evidence.get("currentReviewRef")
    disposition = evidence.get("steeringDisposition")
    if not isinstance(candidate, str) or SHA_RE.fullmatch(candidate) is None:
        return False
    if not isinstance(verification_refs, list) or not verification_refs:
        return False
    if any(
        not isinstance(reference, dict)
        or set(reference) != {"id", "candidate"}
        or EVIDENCE_ID_RE.fullmatch(str(reference.get("id", ""))) is None
        or reference.get("candidate") != candidate
        for reference in verification_refs
    ):
        return False
    if (
        not isinstance(review, dict)
        or set(review) != {"id", "candidate", "disposition"}
        or EVIDENCE_ID_RE.fullmatch(str(review.get("id", ""))) is None
        or review.get("candidate") != candidate
        or review.get("disposition") != "approved"
    ):
        return False
    return (
        isinstance(disposition, dict)
        and set(disposition) == {"authority", "candidate", "decision"}
        and issue_comment_key(disposition.get("authority")) is not None
        and disposition.get("candidate") == candidate
        and disposition.get("decision") == "accepted"
    )


def is_evidence_qualified_terminal(record: Any) -> bool:
    return (
        isinstance(record, dict)
        and record.get("declaredState") in TERMINAL_ACCEPTED_STATES
        and has_qualified_accepted_evidence(record)
    )


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
    records_value = document.get("records", [])
    record_list = records_value if isinstance(records_value, list) else []
    if not isinstance(records_value, list):
        add(errors, "INVALID_REQUIRED_COLLECTION")
    active_by_repo = {
        registry.get("repository", ""): [domain for domain in registry.get("domains", []) if domain.get("state") == "active"]
        for registry in document.get("registries", [])
    }
    active_repository_by_uid = {
        domain.get("domainUid"): repository_key(repository)
        for repository, domain in decls
        if domain.get("state") == "active"
    }
    inventory_source_claims: dict[
        tuple[tuple[str, str], str], list[tuple[str, str]]
    ] = defaultdict(list)
    active_prospective: dict[tuple[str, str], list[dict[str, Any]]] = defaultdict(list)
    for repository, domain in decls:
        if domain.get("state") != "active":
            continue
        host = repository_key(repository)
        uid = domain.get("domainUid")
        for member in domain.get("issuedIds", []):
            if not isinstance(member, dict) or member.get("status") != "prospective":
                continue
            record_id = inventory_member_id(member)
            source_identity = portable_path_identity(member.get("source"))
            if isinstance(record_id, str):
                active_prospective[(uid, portable_text(record_id))].append(member)
                if host is not None and source_identity is not None:
                    inventory_source_claims[(host, source_identity)].append(
                        (uid, portable_text(record_id))
                    )
    if any(
        len(set(claims)) > 1
        for claims in inventory_source_claims.values()
    ):
        add(errors, "DUPLICATE_INVENTORY_SOURCE")

    record_source_claims: dict[
        tuple[tuple[str, str], str], list[tuple[str, str]]
    ] = defaultdict(list)
    record_keys: list[tuple[str, str]] = []
    for record in record_list:
        if not isinstance(record, dict):
            continue
        uid = record.get("domainUid")
        record_id = record.get("id")
        host = active_repository_by_uid.get(uid)
        source_identity = portable_path_identity(record.get("source"))
        if isinstance(record_id, str):
            key = (uid, portable_text(record_id))
            record_keys.append(key)
            if host is not None and source_identity is not None:
                record_source_claims[(host, source_identity)].append(key)
    if any(len(set(claims)) > 1 for claims in record_source_claims.values()):
        add(errors, "DUPLICATE_RECORD_SOURCE")
    if record_list:
        if any(len(active_prospective.get(key, [])) != 1 for key in record_keys):
            add(errors, "RECORD_INVENTORY_BINDING_MISMATCH")
        if any(key not in set(record_keys) for key in active_prospective):
            add(errors, "RECORD_INVENTORY_BINDING_MISMATCH")
    seen: set[tuple[str, str]] = set()
    for record in record_list:
        if not isinstance(record, dict):
            add(errors, "RECORD_REQUIRED_FIELD_MISSING")
            continue
        kind = effective_kind(record)
        validate_typed_shape(record, errors, "RECORD_REQUIRED_FIELD_MISSING")
        if kind in STATEFUL_KINDS:
            validate_stateful_values(record, errors, kind_override=kind)
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
        expected_token = KIND_TO_ID_TOKEN.get(kind)
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
        if normalized_pair not in issued:
            add(errors, "RECORD_NOT_IN_ISSUED_INVENTORY")
        if normalized_pair in issued:
            member = inventory_map(domain).get(portable_text(record_id))
            if (
                not isinstance(member, dict)
                or member.get("status") != "prospective"
                or record.get("source") != member.get("source")
            ):
                add(errors, "RECORD_SOURCE_MISMATCH")


def validate_semantic_reference(
    document: dict[str, Any],
    reference: Any,
    records: dict[tuple[str, str], dict[str, Any]],
    errors: list[str],
    *,
    allowed_kinds: set[str],
) -> tuple[str, str] | None:
    if not isinstance(reference, dict) or not {
        "domainUid", "id"
    }.issubset(reference):
        add(errors, "UNQUALIFIED_SEMANTIC_REFERENCE")
        return None
    if not set(reference).issubset({"domainUid", "id", "keyHint"}):
        add(errors, "INVALID_SEMANTIC_EDGE")
        return None
    key = reference_key(reference)
    if key is None:
        add(errors, "UNQUALIFIED_SEMANTIC_REFERENCE")
        return None
    validate_key_hint(document, reference, errors)
    target = records.get(key)
    if target is None:
        add(errors, "UNRESOLVED_SEMANTIC_REFERENCE")
        return None
    elif effective_kind(target) not in allowed_kinds:
        add(errors, "SEMANTIC_TARGET_KIND_MISMATCH")
        return None
    return key


def validate_semantic_graph(document: dict[str, Any], errors: list[str]) -> None:
    records = {
        reference_key(record): record
        for record in document.get("records", [])
        if isinstance(record, dict) and reference_key(record) is not None
    }
    local_domains = {domain.get("domainUid") for _, domain in declarations(document)}
    active_domains = {
        domain.get("domainUid")
        for _, domain in declarations(document)
        if domain.get("state") == "active"
    }
    known = set(records) | issued_reference_keys(document)
    seen: set[tuple[str, tuple[str, str], tuple[str, str]]] = set()

    def register(
        relation_type: str, source: tuple[str, str] | None,
        target: tuple[str, str] | None,
    ) -> None:
        if source is None or target is None:
            return
        identity = (relation_type, source, target)
        if source == target:
            add(errors, "SEMANTIC_SELF_EDGE")
            return
        if identity in seen:
            add(errors, "DUPLICATE_SEMANTIC_EDGE")
        seen.add(identity)
        source_record = records.get(source)
        target_record = records.get(target)
        if (
            relation_type == "depends_on"
            and isinstance(source_record, dict)
            and effective_kind(source_record) in WORK_OWNER_KINDS | {"slice"}
            and source_record.get("declaredState")
            in {"active", "accepted", "delivered", "released"}
            and not is_evidence_qualified_terminal(target_record)
        ):
            add(errors, "UNSATISFIED_SEMANTIC_DEPENDENCY")

    for source, record in records.items():
        kind = effective_kind(record)
        if kind in {"feature", "refactor"}:
            for edge in record.get("relations", []):
                if (
                    not isinstance(edge, dict)
                    or edge.get("type") not in SEMANTIC_RELATION_TYPES
                    or set(edge) != {"type", "targetRef"}
                ):
                    add(errors, "INVALID_SEMANTIC_EDGE")
                    continue
                target = validate_semantic_reference(
                    document, edge.get("targetRef"), records, errors,
                    allowed_kinds=WORK_OWNER_KINDS | {"study"},
                )
                register(edge["type"], source, target)
        elif kind == "study":
            owner = record.get("ownerRef")
            if owner is not None:
                target = validate_semantic_reference(
                    document, owner, records, errors,
                    allowed_kinds=WORK_OWNER_KINDS,
                )
                register("owned_by", source, target)
            for field, relation_type, allowed in (
                ("independentStudyRefs", "independent_of", {"study"}),
                ("informs", "informs", WORK_OWNER_KINDS | {"study"}),
            ):
                for reference in record.get(field, []):
                    target = validate_semantic_reference(
                        document, reference, records, errors,
                        allowed_kinds=allowed,
                    )
                    register(relation_type, source, target)
        elif kind == "fix":
            for reference in record.get("affectedWork", []):
                target = validate_semantic_reference(
                    document, reference, records, errors,
                    allowed_kinds=WORK_OWNER_KINDS,
                )
                register("corrects", source, target)

    for relation in document.get("relations", []):
        if (
            not isinstance(relation, dict)
            or relation.get("type") not in SEMANTIC_RELATION_TYPES
            or set(relation) != {"type", "sourceRef", "targetRef"}
        ):
            add(errors, "INVALID_SEMANTIC_EDGE")
            continue
        relation_keys: dict[str, tuple[str, str]] = {}
        for name in ("sourceRef", "targetRef"):
            reference = relation.get(name)
            if (
                not isinstance(reference, dict)
                or not {"domainUid", "id"}.issubset(reference)
            ):
                if len(active_domains) > 1:
                    add(errors, "UNQUALIFIED_CROSS_DOMAIN_REF")
                else:
                    add(errors, "INVALID_REFERENCE")
                continue
            if not set(reference).issubset({"domainUid", "id", "keyHint"}):
                add(errors, "INVALID_SEMANTIC_EDGE")
                continue
            validate_key_hint(document, reference, errors)
            key = reference_key(reference)
            if key is None:
                add(errors, "INVALID_REFERENCE")
                continue
            if not valid_uid(key[0]):
                add(errors, "INVALID_REFERENCE")
                continue
            if key[0] in local_domains and key not in known:
                add(errors, "UNRESOLVED_RELATION_REF")
            target_record = records.get(key)
            if (
                target_record is not None
                and effective_kind(target_record)
                not in STATEFUL_KINDS - {"issue-assignment"}
            ):
                add(errors, "SEMANTIC_TARGET_KIND_MISMATCH")
            relation_keys[name] = key
        if {"sourceRef", "targetRef"}.issubset(relation_keys):
            register(
                relation["type"], relation_keys["sourceRef"],
                relation_keys["targetRef"],
            )


def validate_record_graph(document: dict[str, Any], errors: list[str]) -> None:
    records = {
        reference_key(record): record
        for record in document.get("records", [])
        if reference_key(record) is not None
    }
    assignments = document.get("assignments", [])
    assignments_by_issue: dict[tuple[str, str, int], list[dict[str, Any]]] = defaultdict(list)
    assignments_by_slice: dict[tuple[str, str], list[dict[str, Any]]] = defaultdict(list)
    for assignment in assignments:
        if not isinstance(assignment, dict):
            continue
        issue = assignment_issue(assignment)
        issue_identity = issue_key(issue)
        if issue_identity is not None:
            assignments_by_issue[issue_identity].append(assignment)
        active_slices = assignment.get("activeSlices", [])
        for slice_ref in active_slices if isinstance(active_slices, list) else []:
            slice_key = reference_key(slice_ref)
            if slice_key is not None:
                assignments_by_slice[slice_key].append(assignment)

    for key, record in records.items():
        kind = effective_kind(record)
        if kind == "slice":
            owner_key = reference_key(record.get("ownerRef"))
            validate_key_hint(document, record.get("ownerRef"), errors)
            owner = records.get(owner_key)
            if owner is None:
                add(errors, "UNRESOLVED_SLICE_OWNER")
            elif effective_kind(owner) not in WORK_OWNER_KINDS:
                add(errors, "INVALID_SLICE_OWNER")
            else:
                owner_slice_keys = {reference_key(item) for item in owner.get("slices", [])}
                if key not in owner_slice_keys:
                    add(errors, "NONRECIPROCAL_SLICE_OWNER")
            issue = record.get("assignmentIssue")
            issue_identity = issue_key(issue)
            matching_assignments = [
                assignment
                for assignment in assignments_by_slice.get(key, [])
                if issue_key(assignment_issue(assignment)) == issue_identity
            ]
            if len(matching_assignments) != 1:
                add(errors, "SLICE_ASSIGNMENT_CARDINALITY")
            else:
                assignment = matching_assignments[0]
                if reference_key(assignment.get("workRef")) != owner_key:
                    add(errors, "SLICE_ASSIGNMENT_OWNER_MISMATCH")
                if owner is not None:
                    authorities = owner.get("issueAuthorities", [])
                    if sum(issue_key(authority) == issue_identity for authority in authorities) != 1:
                        add(errors, "OWNER_ISSUE_AUTHORITY_MISMATCH")
                coordination = assignment_coordination(assignment)
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
                validate_key_hint(document, decision_ref, errors)
                decision_key = reference_key(decision_ref)
                if decision_key is None:
                    add(errors, "UNQUALIFIED_SLICE_DECISION_REF")
                elif decision_key not in records:
                    add(errors, "UNRESOLVED_SLICE_DECISION_REF")
                else:
                    decision = records[decision_key]
                    if effective_kind(decision) != "study":
                        add(errors, "SLICE_DECISION_AUTHORITY_KIND_MISMATCH")
                    elif not is_evidence_qualified_terminal(decision):
                        add(errors, "UNACCEPTED_SLICE_DECISION_AUTHORITY")
        if kind in WORK_OWNER_KINDS:
            slice_keys: list[tuple[str, str]] = []
            for slice_ref in record.get("slices", []):
                validate_key_hint(document, slice_ref, errors)
                slice_key = reference_key(slice_ref)
                if slice_key is None:
                    add(errors, "INVALID_REFERENCE")
                    continue
                slice_keys.append(slice_key)
                slice_record = records.get(slice_key)
                if slice_record is None:
                    add(errors, "MISSING_OWNER_REFERENCED_SLICE")
                    continue
                if effective_kind(slice_record) != "slice" or reference_key(slice_record.get("ownerRef")) != key:
                    add(errors, "NONRECIPROCAL_SLICE_OWNER")
            if len(slice_keys) != len(set(slice_keys)):
                add(errors, "NONRECIPROCAL_SLICE_OWNER")
            if record.get("declaredState") in {"delivered", "released"}:
                authorities = record.get("issueAuthorities")
                if not isinstance(authorities, list) or not authorities:
                    add(errors, "DELIVERED_WORK_ISSUE_REQUIRED")
                if kind != "fix" and not record.get("slices"):
                    add(errors, "DELIVERED_WORK_SLICE_REQUIRED")
                for slice_key in slice_keys:
                    slice_record = records.get(slice_key)
                    if not (
                        isinstance(slice_record, dict)
                        and slice_record.get("declaredState") == "accepted"
                        and has_qualified_accepted_evidence(slice_record)
                    ):
                        add(errors, "TERMINAL_WORK_SLICE_NOT_ACCEPTED")
                    issue_identity = issue_key(
                        slice_record.get("assignmentIssue")
                        if isinstance(slice_record, dict) else None
                    )
                    authorizing = [
                        assignment
                        for assignment in assignments_by_slice.get(slice_key, [])
                        if issue_key(assignment_issue(assignment)) == issue_identity
                    ]
                    if not (
                        len(authorizing) == 1
                        and authorizing[0].get("declaredState") == "accepted"
                        and has_qualified_accepted_evidence(authorizing[0])
                    ):
                        add(errors, "TERMINAL_WORK_ASSIGNMENT_NOT_ACCEPTED")
                if kind == "fix" and not record.get("slices"):
                    authority_keys = {
                        issue_key(authority) for authority in authorities or []
                    }
                    authority_keys.discard(None)
                    matching = [
                        assignment
                        for assignment in assignments
                        if issue_key(assignment_issue(assignment)) in authority_keys
                        and reference_key(assignment.get("workRef")) == key
                    ]
                    if (
                        record.get("risk") != "low"
                        or record.get("standaloneReviewedUnit") is not True
                        or len(authority_keys) != 1
                        or len(matching) != 1
                        or matching[0].get("activeSlices") != []
                    ):
                        add(errors, "ZERO_SLICE_FIX_INVALID")
                    affected = [
                        records.get(reference_key(reference))
                        for reference in record.get("affectedWork", [])
                        if reference_key(reference) is not None
                    ]
                    if (
                        not affected
                        or any(
                            target is None
                            or effective_kind(target) not in WORK_OWNER_KINDS
                            or target.get("declaredState") not in TERMINAL_ACCEPTED_STATES
                            or mapping(target.get("acceptedEvidence")).get("candidate") is None
                            for target in affected
                        )
                    ):
                        add(errors, "ZERO_SLICE_FIX_AFFECTED_WORK_REQUIRED")


def shared_path(item: Any) -> Any:
    return item.get("path") if isinstance(item, dict) else item


def assignment_projection(assignment: dict[str, Any]) -> dict[str, Any]:
    coordination = assignment_coordination(assignment)
    return {
        "issue": assignment_issue(assignment),
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
        validate_typed_shape(
            reservation, errors, "RESERVATION_SET_INCOMPLETE",
            expected_kind="reservation-set",
        )
        for row in reservation.get("assignments", []) if isinstance(reservation.get("assignments"), list) else []:
            if not isinstance(row, dict):
                continue
            for value in row.get("ownedPaths", []):
                if normalize_path(value) is None:
                    add(errors, "INVALID_PATH")
            for item in row.get("sharedTouchpoints", []):
                if normalize_path(shared_path(item)) is None:
                    add(errors, "INVALID_PATH")
                if not isinstance(item, dict) or not nonblank(item.get("allowedMutation")):
                    add(errors, "INVALID_SHARED_TOUCHPOINT")
            seen_reserved: set[tuple[str, str]] = set()
            for reference in row.get("reservedIds", []):
                key = reference_key(reference)
                if key is not None:
                    if key in seen_reserved:
                        add(errors, "DUPLICATE_RESERVED_ID")
                    seen_reserved.add(key)
        reservation_id = reservation.get("id")
        if not isinstance(reservation_id, str) or not reservation_id:
            add(errors, "RESERVATION_SET_INCOMPLETE")
            continue
        by_id[reservation_id].append(reservation)
    if any(len(values) != 1 for values in by_id.values()):
        add(errors, "RESERVATION_SET_INCOMPLETE")

    grouped_assignments: dict[str, list[dict[str, Any]]] = defaultdict(list)
    for assignment in assignments:
        if not isinstance(assignment, dict):
            add(errors, "ASSIGNMENT_REQUIRED_FIELD_MISSING")
            continue
        reference = assignment_coordination(assignment).get("reservationSet")
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
        else:
            actual_digest = canonical_digest(candidates[0])
            if actual_digest is None:
                add(errors, "INVALID_UNICODE_SCALAR")
            elif digest != actual_digest:
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
            issue_key(assignment_issue(assignment)) for assignment in group
        ]
        row_issues = [issue_key(row.get("issue")) for row in rows if isinstance(row, dict)]
        if None in group_issues or None in row_issues:
            add(errors, "INVALID_GITHUB_ISSUE")
        if Counter(row_issues) != Counter(group_issues):
            add(errors, "RESERVATION_SET_INCOMPLETE")
        rows_by_issue = {
            issue_key(row.get("issue")): row for row in rows if isinstance(row, dict)
        }
        for assignment in group:
            issue = issue_key(assignment_issue(assignment))
            if rows_by_issue.get(issue) != assignment_projection(assignment):
                add(errors, "RESERVATION_ASSIGNMENT_MISMATCH")
        bases = {
            assignment_coordination(assignment).get("integrationBase")
            for assignment in group
        }
        if bases != {reservation.get("integrationBase")}:
            add(errors, "RESERVATION_ASSIGNMENT_MISMATCH")
        shared_order = {
            canonical_digest({
                "mergeOrder": assignment_coordination(assignment).get("mergeOrder"),
                "convergence": assignment_coordination(assignment).get("convergence"),
            })
            for assignment in group
        }
        if len(shared_order) != 1:
            add(errors, "CONVERGENCE_CONTRACT_MISMATCH")
        for assignment in group:
            coordination = assignment_coordination(assignment)
            if (
                coordination.get("mergeOrder") != reservation.get("mergeOrder")
                or coordination.get("convergence") != reservation.get("convergence")
            ):
                add(errors, "RESERVATION_ASSIGNMENT_MISMATCH")

        issue_set = {issue for issue in group_issues if issue is not None}
        edge_map: dict[tuple[str, str, int], list[tuple[str, str, int]]] = {}
        conflict_map: dict[tuple[str, str, int], list[tuple[str, str, int]]] = {}
        for row in rows:
            if not isinstance(row, dict) or not isinstance(row.get("issue"), str):
                continue
            issue = issue_key(row["issue"])
            if issue is None:
                add(errors, "INVALID_GITHUB_ISSUE")
                continue
            dependencies = row.get("dependsOnIssues", [])
            conflicts = row.get("conflictsWithIssues", [])
            if not isinstance(dependencies, list) or not isinstance(conflicts, list):
                add(errors, "RESERVATION_SET_INCOMPLETE")
                continue
            dependency_keys = [issue_key(dependency) for dependency in dependencies]
            conflict_keys = [issue_key(conflict) for conflict in conflicts]
            if None in dependency_keys or None in conflict_keys:
                add(errors, "INVALID_GITHUB_ISSUE")
            normalized_dependencies = [item for item in dependency_keys if item is not None]
            normalized_conflicts = [item for item in conflict_keys if item is not None]
            edge_map[issue] = normalized_dependencies
            conflict_map[issue] = normalized_conflicts
            if len(normalized_dependencies) != len(set(normalized_dependencies)):
                add(errors, "DUPLICATE_DEPENDENCY_EDGE")
            if len(normalized_conflicts) != len(set(normalized_conflicts)):
                add(errors, "DUPLICATE_CONFLICT_EDGE")
            if set(normalized_dependencies) & set(normalized_conflicts):
                add(errors, "DEPENDENCY_CONFLICT_OVERLAP")
            if any(dependency not in issue_set for dependency in normalized_dependencies):
                add(errors, "UNRESOLVED_DEPENDENCY_ISSUE")
            if any(conflict not in issue_set for conflict in normalized_conflicts):
                add(errors, "UNRESOLVED_CONFLICT_ISSUE")
            if issue in normalized_dependencies:
                add(errors, "DEPENDENCY_CYCLE")
            if issue in normalized_conflicts:
                add(errors, "SELF_CONFLICT")

        visiting: set[tuple[str, str, int]] = set()
        visited: set[tuple[str, str, int]] = set()

        def visit(issue: tuple[str, str, int]) -> bool:
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
        merge_keys = [issue_key(issue) for issue in merge_order]
        if None in merge_keys:
            add(errors, "INVALID_GITHUB_ISSUE")
        normalized_merge = [item for item in merge_keys if item is not None]
        if len(normalized_merge) != len(set(normalized_merge)):
            add(errors, "MERGE_ORDER_DUPLICATE")
        if any(issue not in issue_set for issue in normalized_merge):
            add(errors, "MERGE_ORDER_UNKNOWN_MEMBER")
        if issue_set - set(normalized_merge):
            add(errors, "MERGE_ORDER_MISSING_MEMBER")
        positions = {issue: index for index, issue in enumerate(normalized_merge)}
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
        issue_key(assignment_issue(assignment)) for assignment in assignments
    ]
    issue_map = {
        issue_key(assignment_issue(item)): item
        for item in assignments
        if issue_key(assignment_issue(item)) is not None
    }
    if any(
        issue is not None and count > 1
        for issue, count in Counter(issues).items()
    ):
        add(errors, "DUPLICATE_ISSUE_AUTHORITY")
    if len(assignments) > 1:
        bases = {assignment_coordination(item).get("integrationBase") for item in assignments}
        reservation_sets = {
            (
                mapping(assignment_coordination(item).get("reservationSet")).get("id"),
                mapping(assignment_coordination(item).get("reservationSet")).get("digest"),
            )
            for item in assignments
        }
        if len(bases) != 1 or None in bases or len(reservation_sets) != 1:
            add(errors, "STALE_COMMON_BASE")
        convergence_contracts = {
            canonical_digest({
                "mergeOrder": assignment_coordination(item).get("mergeOrder"),
                "convergence": assignment_coordination(item).get("convergence"),
            })
            for item in assignments
        }
        if len(convergence_contracts) != 1:
            add(errors, "CONVERGENCE_CONTRACT_MISMATCH")
        for issue, assignment in issue_map.items():
            for conflict in assignment_coordination(assignment).get("conflictsWithIssues", []):
                peer = issue_map.get(issue_key(conflict))
                peer_conflicts = {
                    issue_key(item)
                    for item in assignment_coordination(peer).get("conflictsWithIssues", [])
                } if peer is not None else set()
                if peer is not None and issue not in peer_conflicts:
                    add(errors, "ASYMMETRIC_CONFLICT")

    owned: list[tuple[Any, tuple[str, bool]]] = []
    shared: list[tuple[Any, tuple[str, bool], Any]] = []
    reserved: dict[tuple[str, str], Any] = {}
    inventory_status_by_id = issued_statuses(document)
    for assignment in assignments:
        if not isinstance(assignment, dict):
            add(errors, "ASSIGNMENT_REQUIRED_FIELD_MISSING")
            continue
        validate_typed_shape(
            assignment, errors, "ASSIGNMENT_REQUIRED_FIELD_MISSING",
            expected_kind="issue-assignment",
        )
        validate_stateful_values(
            assignment, errors, kind_override="issue-assignment"
        )
        validate_accepted_evidence(assignment, errors)
        authority = assignment_issue(assignment)
        authority_identity = issue_key(authority)
        if authority_identity is None:
            add(errors, "INVALID_GITHUB_ISSUE")
        coordination = assignment_coordination(assignment)
        dependency_keys = [
            issue_key(item) for item in coordination.get("dependsOnIssues", [])
        ] if isinstance(coordination.get("dependsOnIssues", []), list) else []
        conflict_keys = [
            issue_key(item) for item in coordination.get("conflictsWithIssues", [])
        ] if isinstance(coordination.get("conflictsWithIssues", []), list) else []
        merge_keys = [
            issue_key(item) for item in coordination.get("mergeOrder", [])
        ] if isinstance(coordination.get("mergeOrder", []), list) else []
        if None in dependency_keys or None in conflict_keys or None in merge_keys:
            add(errors, "INVALID_GITHUB_ISSUE")
        valid_dependencies = [item for item in dependency_keys if item is not None]
        valid_conflicts = [item for item in conflict_keys if item is not None]
        if len(valid_dependencies) != len(set(valid_dependencies)):
            add(errors, "DUPLICATE_DEPENDENCY_EDGE")
        if len(valid_conflicts) != len(set(valid_conflicts)):
            add(errors, "DUPLICATE_CONFLICT_EDGE")
        if set(valid_dependencies) & set(valid_conflicts):
            add(errors, "DEPENDENCY_CONFLICT_OVERLAP")
        if authority_identity is not None and authority_identity in valid_dependencies:
            add(errors, "DEPENDENCY_CYCLE")
        if authority_identity is not None and authority_identity in valid_conflicts:
            add(errors, "SELF_CONFLICT")
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
        if effective_kind(assignment) == "issue-assignment":
            baseline_commit = mapping(assignment.get("baseline")).get("commit")
            if not SHA_RE.fullmatch(str(baseline_commit)) or baseline_commit != base:
                add(errors, "INVALID_INTEGRATION_BASE")
            work_ref = assignment.get("workRef")
            validate_key_hint(document, work_ref, errors)
            validate_id_against_domain(
                work_ref.get("domainUid") if isinstance(work_ref, dict) else None,
                work_ref.get("id") if isinstance(work_ref, dict) else None,
                domain_by_uid, errors, invalid_code="INVALID_WORK_REF",
                style_code="WORK_REF_STYLE_MISMATCH",
            )
            work_key = reference_key(work_ref)
            work_record = records.get(work_key)
            host_repository = repository_for_domain(document, work_ref.get("domainUid")) if isinstance(work_ref, dict) else None
            pull_identity = pull_request_key(mapping(assignment.get("delivery")).get("pullRequest"))
            if host_repository is not None and (
                authority_identity is None
                or authority_identity[:2] != host_repository
                or pull_identity is None
                or pull_identity[:2] != host_repository
            ):
                add(errors, "REPOSITORY_AUTHORITY_MISMATCH")
            if work_record is None:
                add(errors, "UNRESOLVED_WORK_REF")
            elif effective_kind(work_record) != work_ref.get("type"):
                add(errors, "WORK_REF_TYPE_MISMATCH")
            elif effective_kind(work_record) in WORK_OWNER_KINDS | {"study"}:
                authorities = work_record.get("issueAuthorities", [])
                if not isinstance(authorities, list) or sum(
                    issue_key(item) == authority_identity for item in authorities
                ) != 1:
                    add(errors, "OWNER_ISSUE_AUTHORITY_MISMATCH")
            primary_issued_by = issued_authorities(document).get(work_key)
            if primary_issued_by is not None and primary_issued_by != authority_identity:
                add(errors, "PRIMARY_WORK_AUTHORITY_MISMATCH")
            for active_slice in assignment.get("activeSlices", []):
                validate_key_hint(document, active_slice, errors)
                validate_id_against_domain(
                    active_slice.get("domainUid") if isinstance(active_slice, dict) else None,
                    active_slice.get("id") if isinstance(active_slice, dict) else None,
                    domain_by_uid, errors, invalid_code="INVALID_REFERENCE",
                    style_code="ACTIVE_SLICE_STYLE_MISMATCH",
                )
                slice_record = records.get(reference_key(active_slice))
                if slice_record is None or effective_kind(slice_record) != "slice":
                    add(errors, "UNRESOLVED_ACTIVE_SLICE")
                elif issue_key(slice_record.get("assignmentIssue")) != authority_identity:
                    add(errors, "SLICE_ASSIGNMENT_MISMATCH")
                elif reference_key(slice_record.get("ownerRef")) != work_key:
                    add(errors, "ACTIVE_SLICE_OWNER_MISMATCH")
        for value in coordination.get("ownedPaths", []):
            normalized = normalize_path(value)
            if normalized is None:
                add(errors, "INVALID_PATH")
            else:
                owned.append((authority_identity or authority, normalized))
        for item in coordination.get("sharedTouchpoints", []):
            normalized = normalize_path(shared_path(item))
            if normalized is None:
                add(errors, "INVALID_PATH")
            else:
                owner_uid = item.get("ownerDomainUid") if isinstance(item, dict) else None
                if not isinstance(item, dict) or not nonblank(item.get("allowedMutation")):
                    add(errors, "INVALID_SHARED_TOUCHPOINT")
                if not valid_uid(owner_uid):
                    add(errors, "INVALID_SHARED_OWNER")
                elif domain_by_uid and owner_uid not in domain_by_uid:
                    add(errors, "UNKNOWN_SHARED_OWNER_DOMAIN")
                shared.append((authority_identity or authority, normalized, owner_uid))
        prohibited_paths = [
            normalize_path(value)
            for value in mapping(assignment.get("boundaries")).get(
                "prohibitedPaths", []
            )
        ]
        prohibited_paths = [path for path in prohibited_paths if path is not None]
        owned_paths = [
            normalize_path(value) for value in coordination.get("ownedPaths", [])
        ]
        owned_paths = [path for path in owned_paths if path is not None]
        shared_paths = [
            normalize_path(shared_path(item))
            for item in coordination.get("sharedTouchpoints", [])
        ]
        shared_paths = [path for path in shared_paths if path is not None]
        if any(
            paths_overlap(authorized, prohibited)
            for authorized in owned_paths + shared_paths
            for prohibited in prohibited_paths
        ):
            add(errors, "PROHIBITED_PATH_OVERLAP")
        row_reserved: set[tuple[str, str]] = set()
        for reference in coordination.get("reservedIds", []):
            if not isinstance(reference, dict) or not reference.get("domainUid") or not reference.get("id"):
                add(errors, "INVALID_REFERENCE")
                continue
            validate_key_hint(document, reference, errors)
            validate_id_against_domain(
                reference.get("domainUid"), reference.get("id"), domain_by_uid,
                errors, invalid_code="INVALID_RESERVED_ID",
                style_code="RESERVED_ID_STYLE_MISMATCH",
            )
            if domain_by_uid and reference.get("domainUid") not in domain_by_uid:
                add(errors, "UNKNOWN_RESERVED_DOMAIN")
            key = (reference["domainUid"], portable_text(reference["id"]))
            if key in row_reserved:
                add(errors, "DUPLICATE_RESERVED_ID")
            row_reserved.add(key)
            if key in reserved and reserved[key] != (authority_identity or authority):
                add(errors, "RESERVED_ID_COLLISION")
            reserved[key] = authority_identity or authority
            issued_by = issued_authorities(document).get(key)
            if issued_by is not None and issued_by != authority_identity:
                add(errors, "ISSUED_ID_RECLAIM")
            if inventory_status_by_id.get(key) == "legacy-preserved":
                add(errors, "LEGACY_PRESERVED_CURRENT_GATE")

    accepted_execution_states = {"active", "accepted"}
    for issue, assignment in issue_map.items():
        state = assignment.get("declaredState")
        coordination = assignment_coordination(assignment)
        if state == "accepted":
            candidate = mapping(assignment.get("acceptedEvidence")).get("candidate")
            work_ref = assignment.get("workRef")
            work_record = records.get(reference_key(work_ref))
            work_type = work_ref.get("type") if isinstance(work_ref, dict) else None
            slices = assignment.get("activeSlices", [])
            if work_type == "study":
                if (
                    work_record is None
                    or effective_kind(work_record) != "study"
                    or work_record.get("declaredState") != "accepted"
                ):
                    add(errors, "ASSIGNMENT_STUDY_STATE_MISMATCH")
                elif mapping(work_record.get("acceptedEvidence")).get("candidate") != candidate:
                    add(errors, "ASSIGNMENT_CANDIDATE_MISMATCH")
            elif work_type == "fix" and slices == []:
                if (
                    work_record is None
                    or work_record.get("declaredState") not in {"delivered", "released"}
                ):
                    add(errors, "ASSIGNMENT_FIX_STATE_MISMATCH")
                elif mapping(work_record.get("acceptedEvidence")).get("candidate") != candidate:
                    add(errors, "ASSIGNMENT_CANDIDATE_MISMATCH")
            else:
                if not isinstance(slices, list) or not slices:
                    add(errors, "ASSIGNMENT_SLICE_STATE_MISMATCH")
                for slice_ref in slices if isinstance(slices, list) else []:
                    slice_record = records.get(reference_key(slice_ref))
                    if slice_record is None or slice_record.get("declaredState") != "accepted":
                        add(errors, "ASSIGNMENT_SLICE_STATE_MISMATCH")
                    elif mapping(slice_record.get("acceptedEvidence")).get("candidate") != candidate:
                        add(errors, "ASSIGNMENT_CANDIDATE_MISMATCH")
            for dependency in coordination.get("dependsOnIssues", []):
                dependency_assignment = issue_map.get(issue_key(dependency))
                if (
                    dependency_assignment is not None
                    and dependency_assignment.get("declaredState") != "accepted"
                ):
                    add(errors, "UNSATISFIED_ASSIGNMENT_PREREQUISITE")
        if state in accepted_execution_states:
            for conflict in coordination.get("conflictsWithIssues", []):
                peer = issue_map.get(issue_key(conflict))
                if (
                    issue_key(conflict) != issue
                    and peer is not None
                    and peer.get("declaredState") in accepted_execution_states
                ):
                    add(errors, "ACTIVE_ASSIGNMENT_CONFLICT")

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


def validate_assignment_history(
    document: dict[str, Any], errors: list[str], *, repository_driver: bool
) -> None:
    histories = document.get("assignmentHistory", [])
    if not isinstance(histories, list):
        add(errors, "INVALID_ASSIGNMENT_REVISION_HISTORY")
        return
    histories_by_assignment: dict[str, list[dict[str, Any]]] = defaultdict(list)
    history_sources: set[str] = set()
    for history in histories:
        if not isinstance(history, dict):
            add(errors, "INVALID_ASSIGNMENT_REVISION_HISTORY")
            continue
        source = history.get("source")
        assignment_source = history.get("assignmentSource")
        if isinstance(assignment_source, str):
            histories_by_assignment[assignment_source].append(history)
        if isinstance(source, str):
            if source in history_sources:
                add(errors, "INVALID_ASSIGNMENT_REVISION_HISTORY")
            history_sources.add(source)
        reservation = history.get("reservationSet")
        if "schemaVersion" not in history:
            add(errors, "SCHEMA_VERSION_REQUIRED")
        elif history.get("schemaVersion") != ASSIGNMENT_HISTORY_SCHEMA_VERSION:
            add(errors, "UNSUPPORTED_SCHEMA_VERSION")
        if "kind" not in history:
            add(errors, "KIND_MARKER_REQUIRED")
        elif history.get("kind") != "issue-assignment-revision-snapshot":
            add(errors, "KIND_MARKER_MISMATCH")
        if "experimental" not in history:
            add(errors, "EXPERIMENTAL_MARKER_REQUIRED")
        elif history.get("experimental") is not True:
            add(errors, "EXPERIMENTAL_MARKER_REQUIRED")
        valid = (
            valid_record_path(source)
            and valid_record_path(assignment_source)
            and issue_key(history.get("authorityIssue")) is not None
            and isinstance(history.get("revision"), int)
            and not isinstance(history.get("revision"), bool)
            and history.get("revision") >= 1
            and SHA_RE.fullmatch(str(history.get("sourceCandidate", ""))) is not None
            and isinstance(reservation, dict)
            and nonblank(reservation.get("id"))
            and DIGEST_RE.fullmatch(str(reservation.get("digest", ""))) is not None
            and isinstance(history.get("supersededByRevision"), int)
            and history.get("supersededByRevision") == history.get("revision") + 1
            and isinstance(history.get("sourceFieldPresent"), bool)
            and nonblank(history.get("supersessionReason"))
            and isinstance(history.get("recordedContext"), dict)
            and bool(history.get("recordedContext"))
            and "previousSnapshot" in history
        )
        if not valid:
            add(errors, "INVALID_ASSIGNMENT_REVISION_HISTORY")

    matched_history_ids: set[int] = set()
    for assignment in document.get("assignments", []):
        if not isinstance(assignment, dict) or assignment.get("kind") != "issue-assignment":
            continue
        revision = assignment.get("revision")
        previous = assignment.get("previousRevision")
        refreeze = assignment.get("refreeze")
        source = assignment.get("source")
        candidate_histories = histories_by_assignment.get(source, []) if isinstance(source, str) else []
        if revision == 1:
            if (
                "previousRevision" in assignment
                or "refreeze" in assignment
                or candidate_histories
            ):
                add(errors, "INVALID_ASSIGNMENT_REVISION_HISTORY")
            continue
        if not isinstance(revision, int) or isinstance(revision, bool) or revision < 1:
            continue
        if not isinstance(previous, dict) or not isinstance(refreeze, dict):
            add(errors, "INVALID_ASSIGNMENT_REVISION_HISTORY")
            continue
        if not repository_driver:
            add(errors, "REPOSITORY_DRIVER_REQUIRED")
        authority = issue_key(assignment_issue(assignment))
        chain = [
            history for history in candidate_histories
            if issue_key(history.get("authorityIssue")) == authority
        ]
        if (
            len(chain) != revision - 1
            or {history.get("revision") for history in chain} != set(range(1, revision))
        ):
            add(errors, "INVALID_ASSIGNMENT_REVISION_HISTORY")
            continue
        chain.sort(key=lambda history: history["revision"])
        matched_history_ids.update(id(history) for history in chain)
        for index, history in enumerate(chain):
            prior = chain[index - 1] if index else None
            link = history.get("previousSnapshot")
            if prior is None:
                if link is not None:
                    add(errors, "INVALID_ASSIGNMENT_REVISION_HISTORY")
            elif (
                not isinstance(link, dict)
                or link.get("revision") != prior.get("revision")
                or link.get("snapshotPath") != prior.get("source")
                or link.get("sourceCandidate") != prior.get("sourceCandidate")
                or link.get("reservationDigest")
                != mapping(prior.get("reservationSet")).get("digest")
            ):
                add(errors, "INVALID_ASSIGNMENT_REVISION_HISTORY")
            if index and (
                mapping(prior.get("reservationSet")).get("id")
                != mapping(history.get("reservationSet")).get("id")
                or mapping(prior.get("reservationSet")).get("digest")
                == mapping(history.get("reservationSet")).get("digest")
            ):
                add(errors, "INVALID_ASSIGNMENT_REVISION_HISTORY")
        history = chain[-1]
        current_reservation = mapping(assignment_coordination(assignment).get("reservationSet"))
        history_reservation = history.get("reservationSet", {})
        if (
            previous.get("revision") != history.get("revision")
            or previous.get("snapshotPath") != history.get("source")
            or previous.get("sourceCandidate") != history.get("sourceCandidate")
            or previous.get("reservationDigest") != history_reservation.get("digest")
            or assignment.get("revision") != history.get("revision") + 1
            or history.get("supersededByRevision") != assignment.get("revision")
            or current_reservation.get("id") != history_reservation.get("id")
            or current_reservation.get("digest") == history_reservation.get("digest")
            or refreeze.get("sourceCandidate") != history.get("sourceCandidate")
            or not nonblank(refreeze.get("reason"))
        ):
            add(errors, "INVALID_ASSIGNMENT_REVISION_HISTORY")
    if any(id(history) not in matched_history_ids for history in histories if isinstance(history, dict)):
        add(errors, "INVALID_ASSIGNMENT_REVISION_HISTORY")


def read_local_json(repository_root: Path, value: Any) -> dict[str, Any] | None:
    path = resolve_local_path(repository_root, value)
    if path is None or not path.is_file():
        return None
    try:
        loaded = json.loads(path.read_text(encoding="utf-8"))
    except (OSError, UnicodeDecodeError, json.JSONDecodeError):
        return None
    return loaded if isinstance(loaded, dict) else None


def validate_repository_materialization(
    document: dict[str, Any], errors: list[str], repository_root: Path
) -> None:
    records = {
        reference_key(record): record
        for record in document.get("records", [])
        if isinstance(record, dict) and reference_key(record) is not None
    }
    for record in records.values():
        materialized = read_local_json(repository_root, record.get("source"))
        if materialized is None:
            add(errors, "UNRESOLVED_MATERIALIZED_SOURCE")
        elif materialized != record:
            add(errors, "MATERIALIZED_SOURCE_MISMATCH")
    for _, domain in declarations(document):
        if domain.get("state") != "active":
            continue
        uid = domain.get("domainUid")
        for member in domain.get("issuedIds", []):
            if not isinstance(member, dict) or member.get("status") != "prospective":
                continue
            record_id = inventory_member_id(member)
            key = (
                (uid, portable_text(record_id))
                if isinstance(record_id, str) else None
            )
            record = records.get(key)
            if record is None or record.get("source") != member.get("source"):
                add(errors, "RECORD_INVENTORY_BINDING_MISMATCH")

    reservation_sets = {
        reservation.get("id"): reservation
        for reservation in document.get("reservationSets", [])
        if isinstance(reservation, dict) and nonblank(reservation.get("id"))
    }
    for assignment in document.get("assignments", []):
        if not isinstance(assignment, dict) or assignment.get("kind") != "issue-assignment":
            continue
        materialized_assignment = read_local_json(repository_root, assignment.get("source"))
        if materialized_assignment is None:
            add(errors, "UNRESOLVED_MATERIALIZED_SOURCE")
        elif materialized_assignment != assignment:
            add(errors, "MATERIALIZED_SOURCE_MISMATCH")
        reservation_reference = mapping(assignment_coordination(assignment).get("reservationSet"))
        if "path" in reservation_reference:
            materialized_reservation = read_local_json(
                repository_root, reservation_reference.get("path")
            )
            expected_reservation = reservation_sets.get(reservation_reference.get("id"))
            if materialized_reservation is None:
                add(errors, "UNRESOLVED_MATERIALIZED_SOURCE")
            elif expected_reservation is None or materialized_reservation != expected_reservation:
                add(errors, "MATERIALIZED_SOURCE_MISMATCH")

    for history in document.get("assignmentHistory", []):
        if not isinstance(history, dict):
            continue
        materialized_history = read_local_json(repository_root, history.get("source"))
        if materialized_history is None:
            add(errors, "UNRESOLVED_MATERIALIZED_SOURCE")
        elif materialized_history != history:
            add(errors, "MATERIALIZED_SOURCE_MISMATCH")
        candidate = history.get("sourceCandidate")
        assignment_source = history.get("assignmentSource")
        try:
            subprocess.run(
                ["git", "cat-file", "-e", f"{candidate}^{{commit}}"],
                cwd=repository_root, check=True, capture_output=True,
            )
            historical_blob = subprocess.run(
                ["git", "show", f"{candidate}:{assignment_source}"],
                cwd=repository_root, check=True, capture_output=True,
            ).stdout.decode("utf-8")
            historical_assignment = json.loads(historical_blob)
        except (
            OSError, subprocess.CalledProcessError, UnicodeDecodeError,
            json.JSONDecodeError,
        ):
            add(errors, "UNAVAILABLE_HISTORY_GIT_OBJECT")
            continue
        historical_reservation = mapping(
            assignment_coordination(historical_assignment).get("reservationSet")
        )
        source_present = "source" in historical_assignment
        context = history.get("recordedContext")
        context_matches = isinstance(context, dict)
        if context_matches:
            for pointer, expected in context.items():
                try:
                    actual = pointer_value(historical_assignment, pointer)
                except (KeyError, IndexError, TypeError, ValueError):
                    context_matches = False
                    break
                if actual != expected:
                    context_matches = False
                    break
        if (
            issue_key(assignment_issue(historical_assignment))
            != issue_key(history.get("authorityIssue"))
            or historical_assignment.get("revision") != history.get("revision")
            or historical_reservation.get("id")
            != mapping(history.get("reservationSet")).get("id")
            or historical_reservation.get("digest")
            != mapping(history.get("reservationSet")).get("digest")
            or source_present != history.get("sourceFieldPresent")
            or (source_present and historical_assignment.get("source") != assignment_source)
            or not context_matches
        ):
            add(errors, "HISTORY_GIT_CONTENT_MISMATCH")


def validate_document(
    document: dict[str, Any], *, isolated_fixture: bool = False,
    repository_root: Path | None = None,
) -> list[str]:
    errors: list[str] = []
    if contains_invalid_unicode_scalar(document):
        add(errors, "INVALID_UNICODE_SCALAR")
    if document.get("experimental") is not True:
        add(errors, "EXPERIMENTAL_MARKER_REQUIRED")
    validate_registry_set(document, errors)
    validate_records(document, errors)
    validate_record_graph(document, errors)
    validate_semantic_graph(document, errors)
    validate_assignments(document, errors)
    validate_assignment_history(
        document, errors, repository_driver=repository_root is not None
    )
    if not isolated_fixture:
        validate_reservation_sets(document, errors)
    if repository_root is not None:
        validate_repository_materialization(document, errors, repository_root)
    return sorted(errors)


def validate_template(record: dict[str, Any]) -> list[str]:
    errors: list[str] = []
    kind = record.get("kind")
    validate_typed_shape(record, errors, "TEMPLATE_REQUIRED_FIELD_MISSING")
    if kind in KIND_TO_ID_TOKEN:
        validate_stateful_values(record, errors)
        if not valid_uid(record.get("domainUid")):
            add(errors, "INVALID_DOMAIN_UID")
        match = SCOPED_ID_RE.fullmatch(str(record.get("id", ""))) or UNSCOPED_ID_RE.fullmatch(str(record.get("id", "")))
        if not match or match.group("type") != KIND_TO_ID_TOKEN[kind]:
            add(errors, "RECORD_KIND_ID_MISMATCH")
        validate_accepted_evidence(record, errors)
    elif kind == "issue-assignment":
        validate_stateful_values(record, errors)
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
    accepted_study = example_documents.get("accepted-study-only.json", {})
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
        isinstance(edge, dict)
        and isinstance(edge.get("targetRef"), dict)
        and record.get("domainUid") != edge["targetRef"].get("domainUid")
        for record in scoped.get("records", [])
        for edge in record.get("relations", [])
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
    if not any(
        isinstance(member, dict)
        and member.get("status") == "legacy-preserved"
        and member.get("authorityIssue") is None
        and nonblank(member.get("authorityMissingReason"))
        for _, domain in declarations(moved)
        for member in domain.get("issuedIds", [])
    ):
        failures.append("coverage: truthful no-Issue legacy inventory missing")
    if not all(
        isinstance(record, dict)
        and valid_record_path(record.get("source"))
        for document in example_documents.values()
        for record in document.get("records", [])
    ):
        failures.append("coverage: stateful record source binding missing")
    simple_assignment = next(iter(simple.get("assignments", [])), {})
    simple_records = {record.get("kind"): record for record in simple.get("records", [])}
    if not (
        simple_assignment.get("declaredState") == "accepted"
        and simple_records.get("slice", {}).get("declaredState") == "accepted"
        and simple_records.get("feature", {}).get("declaredState") == "active"
    ):
        failures.append("coverage: accepted assignment/Slice with durable active Feature missing")
    if not any(
        assignment.get("declaredState") == "accepted"
        and mapping(assignment.get("workRef")).get("type") == "study"
        for assignment in accepted_study.get("assignments", [])
    ):
        failures.append("coverage: accepted Study-only assignment missing")
    scoped_issue_states = {
        issue_key(assignment_issue(assignment)): assignment.get("declaredState")
        for assignment in scoped.get("assignments", [])
    }
    if not any(
        state == "active"
        and any(
            scoped_issue_states.get(issue_key(conflict)) == "blocked"
            for conflict in assignment_coordination(assignment).get("conflictsWithIssues", [])
        )
        for assignment in scoped.get("assignments", [])
        for state in [assignment.get("declaredState")]
    ):
        failures.append("coverage: valid active/blocked symmetric conflict pair missing")
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
        "INVALID_GITHUB_REPOSITORY", "INVALID_GITHUB_ISSUE",
        "INVALID_GITHUB_PULL_REQUEST", "INVALID_GITHUB_COMMENT",
        "REPOSITORY_AUTHORITY_MISMATCH", "INVALID_INVENTORY_AUTHORITY",
        "DUPLICATE_REPOSITORY_REGISTRY", "DOMAIN_OWNER_REQUIRED",
        "INVALID_REVISION", "UNSUPPORTED_DECLARED_STATE",
        "BLANK_REQUIRED_VALUE", "INVALID_REQUIRED_COLLECTION",
        "INVALID_BRANCH", "INVALID_STOP_CONDITION", "INVALID_EVIDENCE_POLICY",
        "ZERO_SLICE_FIX_INVALID", "INVALID_RECORD_SOURCE",
        "RECORD_SOURCE_MISMATCH", "LEGACY_AUTHORITY_REASON_REQUIRED",
        "ACCEPTED_EVIDENCE_SHAPE_INVALID",
        "ACCEPTED_EVIDENCE_CANDIDATE_MISMATCH", "STEERING_AUTHORITY_MISMATCH",
        "PRIMARY_WORK_AUTHORITY_MISMATCH",
        "SELF_CONFLICT",
        "DUPLICATE_DEPENDENCY_EDGE", "DUPLICATE_CONFLICT_EDGE",
        "DEPENDENCY_CONFLICT_OVERLAP", "INVALID_ASSIGNMENT_REVISION_HISTORY",
        "ASSIGNMENT_SLICE_STATE_MISMATCH", "ASSIGNMENT_CANDIDATE_MISMATCH",
        "ASSIGNMENT_STUDY_STATE_MISMATCH", "UNSATISFIED_ASSIGNMENT_PREREQUISITE",
        "ACTIVE_ASSIGNMENT_CONFLICT", "UNQUALIFIED_SEMANTIC_REFERENCE",
        "UNRESOLVED_SEMANTIC_REFERENCE", "SEMANTIC_TARGET_KIND_MISMATCH",
        "DUPLICATE_SEMANTIC_EDGE", "ZERO_SLICE_FIX_AFFECTED_WORK_REQUIRED",
        "REPOSITORY_DRIVER_REQUIRED", "UNRESOLVED_MATERIALIZED_SOURCE",
        "MATERIALIZED_SOURCE_MISMATCH", "UNAVAILABLE_HISTORY_GIT_OBJECT",
        "HISTORY_GIT_CONTENT_MISMATCH", "DUPLICATE_RESERVED_ID",
        "DUPLICATE_INVENTORY_SOURCE", "DUPLICATE_RECORD_SOURCE",
        "RECORD_INVENTORY_BINDING_MISMATCH", "KEY_HINT_MISMATCH",
        "UNSUPPORTED_SCHEMA_VERSION", "SCHEMA_VERSION_REQUIRED",
        "KIND_MARKER_REQUIRED", "EXPERIMENTAL_MARKER_REQUIRED",
        "INVALID_SEMANTIC_EDGE", "SEMANTIC_SELF_EDGE",
        "UNSATISFIED_SEMANTIC_DEPENDENCY",
        "SLICE_DECISION_AUTHORITY_KIND_MISMATCH",
        "UNACCEPTED_SLICE_DECISION_AUTHORITY",
        "TERMINAL_WORK_SLICE_NOT_ACCEPTED",
        "TERMINAL_WORK_ASSIGNMENT_NOT_ACCEPTED",
        "PROHIBITED_PATH_OVERLAP", "INVALID_UNICODE_SCALAR",
        "LEGACY_PRESERVED_CURRENT_GATE",
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
        reservation.get("id"): digest
        for reservation in document.get("reservationSets", [])
        if isinstance(reservation, dict)
        and reservation.get("id")
        and (digest := canonical_digest(reservation)) is not None
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


def load_dogfood_document() -> dict[str, Any]:
    assignment_paths = sorted((ROOT / "Steering" / "Assignments").glob("ISSUE-*.*"))
    history_paths = sorted((ROOT / "Steering" / "Assignments" / "History").glob("*.json"))
    reservation_paths = sorted((ROOT / "Steering" / "Reservations").glob("*.json"))
    record_paths = sorted((ROOT / "Studies").glob("*.json"))
    registry_paths = sorted((ROOT / "Steering").glob("WorkDomains*.json"))
    return {
        "experimental": True,
        "registries": [load_json(path) for path in registry_paths],
        "records": [load_json(path) for path in record_paths],
        "relations": [],
        "assignments": [load_json(path) for path in assignment_paths],
        "reservationSets": [load_json(path) for path in reservation_paths],
        "assignmentHistory": [load_json(path) for path in history_paths],
    }


def main() -> int:
    failures: list[str] = []
    try:
        dogfood = load_dogfood_document()
    except (OSError, json.JSONDecodeError) as exc:
        dogfood = {}
        failures.append(f"Steering dogfood parse error: {exc}")

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

    positive_fixture_count = 0
    for path in sorted((ROOT / "fixtures" / "positive").glob("*.json")):
        positive_fixture_count += 1
        try:
            fixture = load_json(path)
            source_example = fixture.get("sourceExample")
            subject = copy.deepcopy(example_documents[source_example])
            apply_mutations(subject, fixture.get("mutations", []))
            if fixture.get("rebindReservationDigests") is True:
                rebind_reservation_digests(subject)
            actual = validate_document(subject)
        except (OSError, json.JSONDecodeError, KeyError, IndexError, TypeError, ValueError) as exc:
            failures.append(f"{path.relative_to(ROOT)}: invalid positive fixture: {exc}")
            continue
        if actual:
            failures.append(f"{path.relative_to(ROOT)}: unexpected {actual!r}")

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
            if document.get("sourceDogfood") is True:
                subject = copy.deepcopy(dogfood)
                apply_mutations(subject, document.get("mutations", []))
                if document.get("rebindReservationDigests") is True:
                    rebind_reservation_digests(subject)
                actual = validate_document(subject, repository_root=ROOT.parent)
            elif source_example:
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

    try:
        errors = validate_document(dogfood, repository_root=ROOT.parent)
        if errors:
            failures.append(f"Steering dogfood: {', '.join(errors)}")
    except (OSError, json.JSONDecodeError) as exc:
        failures.append(f"Steering dogfood repository driver error: {exc}")

    failures.extend(validate_text_corpus())

    if failures:
        print("SDP vNext pilot validation: FAIL")
        for failure in failures:
            print(f"- {failure}")
        return 1

    print("SDP vNext pilot validation: PASS")
    print(f"- templates: {len(templates)} ({', '.join(sorted(template_kinds))})")
    print(f"- positive examples: {len(example_documents)}")
    print(f"- positive mutation fixtures: {positive_fixture_count}")
    print(f"- negative fixtures: {negative_count}")
    print(f"- negative diagnostic coverage: {len(negative_codes)} codes")
    print("- Steering Issue #7 Study/assignment/domain/reservation/history binding: valid")
    print("- local Markdown links/status markers/trailing whitespace: valid")
    return 0


if __name__ == "__main__":
    sys.exit(main())
