#!/usr/bin/env python3
"""Deterministically validate the Issue #5 repository-analysis corpus.

The validator intentionally uses only the Python standard library.  It treats
RepositoryInventory.md as the report manifest and checks report structure and
evidence-labelled prose; the presence of SDP-looking path names alone is never
accepted as evidence that a required topic was analysed.
"""

from __future__ import annotations

import argparse
import hashlib
import re
import sys
import unicodedata
from dataclasses import dataclass
from pathlib import Path
from typing import Iterable, Sequence
from urllib.parse import unquote, urlsplit


EXPECTED_CONSIDERED = 35
EXPECTED_IN_SCOPE = 17
EXPECTED_EXCLUDED = 18
OWNER = "Hans-Einar"
LABELS = ("OBSERVED", "OWNER DIRECTION", "INFERENCE", "RECOMMENDATION")
HEX40_RE = re.compile(r"(?<![0-9a-fA-F])[0-9a-fA-F]{40}(?![0-9a-fA-F])")
REPOSITORY_RE = re.compile(r"`?(Hans-Einar/[A-Za-z0-9_.-]+)`?")
HEADING_RE = re.compile(r"^(#{1,6})[ \t]+(.+?)[ \t]*#*[ \t]*$")
LABELLED_ITEM_RE = re.compile(
    r"^\s*[-+*]\s+\*\*(OBSERVED|OWNER DIRECTION|INFERENCE|RECOMMENDATION)"
    r"(?:\s*(?:[-—:]|/)[^*]*)?\*\*",
    re.IGNORECASE,
)
INLINE_LINK_RE = re.compile(
    r"!?\[[^\]\n]*\]\(\s*(?:<([^>]+)>|([^\s)]+))(?:\s+[^)]*)?\s*\)"
)
REFERENCE_LINK_RE = re.compile(
    r"^\s{0,3}\[[^\]]+\]:\s*(?:<([^>]+)>|(\S+))", re.MULTILINE
)


@dataclass(frozen=True)
class Diagnostic:
    path: str
    code: str
    message: str


@dataclass(frozen=True)
class Heading:
    level: int
    line: int
    text: str
    key: str


@dataclass(frozen=True)
class LabelledBlock:
    label: str
    start: int
    end: int
    text: str


@dataclass(frozen=True)
class InventoryRow:
    repository: str
    name: str
    default_branch: str
    study_commit: str
    report_path: str | None


@dataclass
class Stats:
    considered: int = 0
    in_scope: int = 0
    excluded: int = 0
    reports: int = 0


def display_path(path: Path, root: Path) -> str:
    try:
        return path.resolve().relative_to(root.resolve()).as_posix()
    except ValueError:
        return path.as_posix()


def add_error(
    errors: list[Diagnostic], root: Path, path: Path, code: str, message: str
) -> None:
    errors.append(Diagnostic(display_path(path, root), code, message))


def read_utf8(path: Path, root: Path, errors: list[Diagnostic]) -> str | None:
    try:
        data = path.read_bytes()
    except OSError as exc:
        add_error(errors, root, path, "FILE_READ", str(exc))
        return None
    try:
        return data.decode("utf-8")
    except UnicodeDecodeError as exc:
        add_error(errors, root, path, "UTF8", f"not valid UTF-8: {exc}")
        return None


def lines_outside_fences(text: str) -> list[tuple[int, str]]:
    result: list[tuple[int, str]] = []
    fence: str | None = None
    for number, line in enumerate(text.splitlines(), 1):
        marker = re.match(r"^\s{0,3}(`{3,}|~{3,})", line)
        if marker:
            run = marker.group(1)
            if fence is None:
                fence = run[0]
            elif run[0] == fence:
                fence = None
            continue
        if fence is None:
            result.append((number, line))
    return result


def plain_inline(text: str) -> str:
    text = re.sub(r"!\[([^]]*)\]\([^)]*\)", r"\1", text)
    text = re.sub(r"\[([^]]+)\]\([^)]*\)", r"\1", text)
    text = re.sub(r"<[^>]+>", " ", text)
    text = text.replace("`", "").replace("*", "")
    return text.strip()


def normalized_key(text: str) -> str:
    value = unicodedata.normalize("NFKC", plain_inline(text)).casefold()
    value = re.sub(r"[^\w]+", " ", value, flags=re.UNICODE)
    return " ".join(value.split())


def parse_headings(text: str) -> list[Heading]:
    headings: list[Heading] = []
    for number, line in lines_outside_fences(text):
        match = HEADING_RE.match(line)
        if match:
            heading_text = match.group(2).strip()
            headings.append(
                Heading(len(match.group(1)), number, heading_text, normalized_key(heading_text))
            )
    return headings


def validate_headings(
    text: str, path: Path, root: Path, errors: list[Diagnostic]
) -> list[Heading]:
    headings = parse_headings(text)
    seen: dict[str, Heading] = {}
    for heading in headings:
        previous = seen.get(heading.key)
        if previous is not None:
            add_error(
                errors,
                root,
                path,
                "DUPLICATE_HEADING",
                f"heading {heading.text!r} on line {heading.line} duplicates line "
                f"{previous.line}",
            )
        else:
            seen[heading.key] = heading
    return headings


def validate_whitespace(
    text: str, path: Path, root: Path, errors: list[Diagnostic]
) -> None:
    tab_lines = [str(index) for index, line in enumerate(text.splitlines(), 1) if "\t" in line]
    if tab_lines:
        add_error(
            errors,
            root,
            path,
            "TAB",
            f"tab character on line(s) {', '.join(tab_lines)}",
        )
    trailing_lines = [
        str(index)
        for index, line in enumerate(text.splitlines(), 1)
        if line.endswith((" ", "\t"))
    ]
    if trailing_lines:
        add_error(
            errors,
            root,
            path,
            "TRAILING_WHITESPACE",
            f"trailing whitespace on line(s) {', '.join(trailing_lines)}",
        )


def h2_sections(text: str, headings: Sequence[Heading]) -> list[tuple[Heading, int, int]]:
    lines = text.splitlines()
    h2s = [heading for heading in headings if heading.level == 2]
    return [
        (
            heading,
            heading.line,
            (h2s[index + 1].line - 1 if index + 1 < len(h2s) else len(lines)),
        )
        for index, heading in enumerate(h2s)
    ]


def section_text(text: str, start: int, end: int) -> str:
    return "\n".join(text.splitlines()[start:end])


def word_count(text: str) -> int:
    visible = plain_inline(text)
    return len(re.findall(r"\b[^\W_][\w'-]*\b", visible, flags=re.UNICODE))


def parse_labelled_blocks(text: str) -> list[LabelledBlock]:
    lines = text.splitlines()
    blocks: list[LabelledBlock] = []
    fence: str | None = None
    index = 0
    while index < len(lines):
        line = lines[index]
        marker = re.match(r"^\s{0,3}(`{3,}|~{3,})", line)
        if marker:
            run = marker.group(1)
            if fence is None:
                fence = run[0]
            elif run[0] == fence:
                fence = None
            index += 1
            continue
        match = LABELLED_ITEM_RE.match(line) if fence is None else None
        if not match:
            index += 1
            continue
        start = index + 1
        block_lines = [line]
        index += 1
        while index < len(lines):
            candidate = lines[index]
            if HEADING_RE.match(candidate) or re.match(r"^\s*[-+*]\s+", candidate):
                break
            block_lines.append(candidate)
            index += 1
        blocks.append(
            LabelledBlock(
                match.group(1).upper(), start, start + len(block_lines) - 1, "\n".join(block_lines)
            )
        )
    return blocks


def blocks_in_range(
    blocks: Sequence[LabelledBlock], start: int, end: int
) -> list[LabelledBlock]:
    return [block for block in blocks if block.start >= start and block.start <= end]


def split_table_row(line: str) -> list[str]:
    stripped = line.strip()
    if stripped.startswith("|"):
        stripped = stripped[1:]
    if stripped.endswith("|"):
        stripped = stripped[:-1]
    return [cell.strip() for cell in stripped.split("|")]


def extract_h2_lines(text: str, wanted: str) -> list[str] | None:
    lines = text.splitlines()
    headings = parse_headings(text)
    matches = [heading for heading in headings if heading.level == 2 and heading.key == wanted]
    if len(matches) != 1:
        return None
    start = matches[0].line
    later = [heading.line for heading in headings if heading.level == 2 and heading.line > start]
    end = min(later) - 1 if later else len(lines)
    return lines[start:end]


def parse_inventory_table(
    text: str,
    section_key: str,
    required_columns: Sequence[str],
    inventory_path: Path,
    root: Path,
    errors: list[Diagnostic],
) -> list[dict[str, str]]:
    section = extract_h2_lines(text, section_key)
    if section is None:
        add_error(
            errors,
            root,
            inventory_path,
            "INVENTORY_SECTION",
            f"expected exactly one H2 section {section_key!r}",
        )
        return []
    wanted = [normalized_key(column) for column in required_columns]
    for index, line in enumerate(section):
        if not line.lstrip().startswith("|"):
            continue
        header = split_table_row(line)
        header_keys = [normalized_key(cell) for cell in header]
        if header_keys != wanted:
            continue
        if index + 1 >= len(section):
            break
        separator = split_table_row(section[index + 1])
        if len(separator) != len(header) or not all(
            re.fullmatch(r":?-{3,}:?", cell) for cell in separator
        ):
            add_error(
                errors,
                root,
                inventory_path,
                "INVENTORY_TABLE",
                f"invalid separator row in {section_key!r}",
            )
            return []
        rows: list[dict[str, str]] = []
        for row_line in section[index + 2 :]:
            if not row_line.lstrip().startswith("|"):
                break
            cells = split_table_row(row_line)
            if len(cells) != len(header):
                add_error(
                    errors,
                    root,
                    inventory_path,
                    "INVENTORY_TABLE",
                    f"table row has {len(cells)} cells; expected {len(header)}",
                )
                continue
            rows.append(dict(zip(header_keys, cells)))
        return rows
    add_error(
        errors,
        root,
        inventory_path,
        "INVENTORY_TABLE",
        f"required table not found in {section_key!r}",
    )
    return []


def parse_inventory_row(
    row: dict[str, str],
    include_report: bool,
    inventory_path: Path,
    root: Path,
    errors: list[Diagnostic],
) -> InventoryRow | None:
    repository_cell = row.get("repository", "")
    repository_matches = REPOSITORY_RE.findall(repository_cell)
    if len(repository_matches) != 1 or repository_cell.strip().strip("`") != repository_matches[0]:
        add_error(
            errors,
            root,
            inventory_path,
            "INVENTORY_REPOSITORY",
            f"invalid repository identity cell {repository_cell!r}",
        )
        return None
    repository = repository_matches[0]
    branch_cell = plain_inline(row.get("default branch", "")).strip("`")
    if not re.fullmatch(r"[A-Za-z0-9._/-]+", branch_cell):
        add_error(
            errors,
            root,
            inventory_path,
            "INVENTORY_BRANCH",
            f"{repository}: invalid default branch {branch_cell!r}",
        )
        return None
    commit_key = "exact default branch study commit committer utc"
    commits = HEX40_RE.findall(row.get(commit_key, ""))
    if len(commits) != 1:
        add_error(
            errors,
            root,
            inventory_path,
            "INVENTORY_COMMIT",
            f"{repository}: expected one exact 40-character study commit, found {len(commits)}",
        )
        return None
    report_path: str | None = None
    if include_report:
        report_path = plain_inline(row.get("planned report", "")).strip("`")
        candidate = Path(report_path)
        if (
            not report_path
            or candidate.is_absolute()
            or ".." in candidate.parts
            or candidate.suffix.casefold() != ".md"
            or not report_path.replace("\\", "/").startswith("repositories/")
        ):
            add_error(
                errors,
                root,
                inventory_path,
                "INVENTORY_REPORT_PATH",
                f"{repository}: invalid planned report path {report_path!r}",
            )
            return None
        report_path = report_path.replace("\\", "/")
    return InventoryRow(
        repository=repository,
        name=repository.split("/", 1)[1],
        default_branch=branch_cell,
        study_commit=commits[0].lower(),
        report_path=report_path,
    )


def validate_inventory(
    text: str, inventory_path: Path, root: Path, errors: list[Diagnostic]
) -> tuple[list[InventoryRow], list[InventoryRow], Stats]:
    in_columns = (
        "Repository",
        "Visibility",
        "Default branch",
        "Exact default-branch study commit (committer UTC)",
        "Refs/commits",
        "Window commits",
        "First qualifying commit (UTC; containing ref)",
        "Last qualifying commit (UTC; containing ref)",
        "Planned report",
    )
    out_columns = (
        "Repository",
        "Visibility / flags",
        "Default branch",
        "Exact default-branch study commit (committer UTC)",
        "Refs/commits",
        "Window commits",
        "Nearest earlier reachable commit (UTC; containing ref)",
        "Evidence-backed exclusion",
    )
    in_raw = parse_inventory_table(
        text,
        "in scope repositories and report manifest",
        in_columns,
        inventory_path,
        root,
        errors,
    )
    out_raw = parse_inventory_table(
        text, "excluded repositories", out_columns, inventory_path, root, errors
    )
    included = [
        parsed
        for row in in_raw
        if (parsed := parse_inventory_row(row, True, inventory_path, root, errors))
        is not None
    ]
    excluded = [
        parsed
        for row in out_raw
        if (parsed := parse_inventory_row(row, False, inventory_path, root, errors))
        is not None
    ]
    stats = Stats(
        considered=len(in_raw) + len(out_raw),
        in_scope=len(in_raw),
        excluded=len(out_raw),
        reports=0,
    )
    expected_counts = (
        ("considered", stats.considered, EXPECTED_CONSIDERED),
        ("in-scope", stats.in_scope, EXPECTED_IN_SCOPE),
        ("excluded", stats.excluded, EXPECTED_EXCLUDED),
    )
    for label, actual, expected in expected_counts:
        if actual != expected:
            add_error(
                errors,
                root,
                inventory_path,
                "INVENTORY_COUNT",
                f"{label} row count is {actual}; expected {expected}",
            )
    identities: dict[str, str] = {}
    for row in included + excluded:
        key = row.repository.casefold()
        if key in identities:
            add_error(
                errors,
                root,
                inventory_path,
                "DUPLICATE_IDENTITY",
                f"repository {row.repository!r} duplicates {identities[key]!r}",
            )
        else:
            identities[key] = row.repository
    paths: dict[str, str] = {}
    for row in included:
        assert row.report_path is not None
        key = row.report_path.casefold()
        if key in paths:
            add_error(
                errors,
                root,
                inventory_path,
                "DUPLICATE_REPORT_PATH",
                f"planned path {row.report_path!r} duplicates {paths[key]!r}",
            )
        else:
            paths[key] = row.report_path
        if Path(row.report_path).stem != row.name:
            add_error(
                errors,
                root,
                inventory_path,
                "REPORT_FILENAME",
                f"{row.repository}: report filename must be {row.name}.md, got "
                f"{Path(row.report_path).name}",
            )
    if len(paths) != EXPECTED_IN_SCOPE:
        add_error(
            errors,
            root,
            inventory_path,
            "REPORT_MANIFEST_COUNT",
            f"unique planned report count is {len(paths)}; expected {EXPECTED_IN_SCOPE}",
        )
    return included, excluded, stats


def validate_required_section(
    report_text: str,
    report_path: Path,
    root: Path,
    errors: list[Diagnostic],
    heading: Heading,
    end: int,
    blocks: Sequence[LabelledBlock],
    minimum_words: int,
    required_label: str,
) -> None:
    body = section_text(report_text, heading.line, end)
    count = word_count(body)
    if count < minimum_words:
        add_error(
            errors,
            root,
            report_path,
            "SUBSTANTIVE_SECTION",
            f"section {heading.text!r} has {count} words; expected at least {minimum_words}",
        )
    section_blocks = blocks_in_range(blocks, heading.line, end)
    if not any(block.label == required_label for block in section_blocks):
        add_error(
            errors,
            root,
            report_path,
            "SECTION_EVIDENCE",
            f"section {heading.text!r} requires a substantive {required_label} block",
        )


def validate_report(
    row: InventoryRow,
    report_path: Path,
    text: str,
    root: Path,
    errors: list[Diagnostic],
) -> str | None:
    headings = validate_headings(text, report_path, root, errors)
    h1s = [heading for heading in headings if heading.level == 1]
    if len(h1s) != 1:
        add_error(
            errors,
            root,
            report_path,
            "REPORT_H1",
            f"expected exactly one H1 heading, found {len(h1s)}",
        )
        h1_key = None
    else:
        h1_key = h1s[0].key
        if row.name.casefold() not in plain_inline(h1s[0].text).casefold():
            add_error(
                errors,
                root,
                report_path,
                "REPORT_IDENTITY",
                f"H1 {h1s[0].text!r} does not identify {row.name!r}",
            )
    sections = h2_sections(text, headings)
    if len(sections) < 6:
        add_error(
            errors,
            root,
            report_path,
            "REPORT_SECTIONS",
            f"expected at least six H2 sections, found {len(sections)}",
        )
        return h1_key
    h2s = [section[0] for section in sections]
    if [heading.key for heading in h2s[-2:]] != [
        "carry forward",
        "legacy do not carry forward",
    ]:
        add_error(
            errors,
            root,
            report_path,
            "FINAL_SECTIONS",
            "Carry forward and Legacy / do not carry forward must be the final two H2 sections",
        )

    facts_matches = [
        section
        for section in sections
        if "repository facts" in section[0].key
        or section[0].key in {"repository facts and evidence boundary", "evidence frame"}
    ]
    if not facts_matches:
        add_error(
            errors,
            root,
            report_path,
            "FACTS_SECTION",
            "missing an H2 repository facts/evidence section",
        )
    work_matches = [
        section
        for section in sections
        if section[0].key.startswith("how ")
        and (" work" in f" {section[0].key}" or "sdp" in section[0].key)
    ]
    if not work_matches:
        add_error(
            errors,
            root,
            report_path,
            "WORK_CONTROL_SECTION",
            "missing an H2 actual SDP use/work-control section",
        )
    worked_matches = [section for section in sections if section[0].key == "what worked well"]
    pain_matches = [section for section in sections if section[0].key.startswith("pain points")]
    carry_matches = [section for section in sections if section[0].key == "carry forward"]
    legacy_matches = [
        section for section in sections if section[0].key == "legacy do not carry forward"
    ]
    for name, matches in (
        ("What worked well", worked_matches),
        ("Pain points", pain_matches),
        ("Carry forward", carry_matches),
        ("Legacy / do not carry forward", legacy_matches),
    ):
        if len(matches) != 1:
            add_error(
                errors,
                root,
                report_path,
                "REQUIRED_SECTION",
                f"expected exactly one {name!r} H2 section, found {len(matches)}",
            )

    blocks = parse_labelled_blocks(text)
    represented = {block.label for block in blocks}
    for label in LABELS:
        if label not in represented:
            add_error(
                errors,
                root,
                report_path,
                "EVIDENCE_LABEL",
                f"no substantive list item uses evidence label {label}",
            )

    if facts_matches:
        facts_start = facts_matches[0][0].line
        work_start = work_matches[0][0].line if work_matches else facts_matches[0][2]
        facts_end = max(facts_matches[0][2], work_start - 1)
        facts_body = section_text(text, facts_start, facts_end)
        facts_blocks = blocks_in_range(blocks, facts_start, facts_end)
        if word_count(facts_body) < 80 or sum(
            block.label == "OBSERVED" for block in facts_blocks
        ) < 2:
            add_error(
                errors,
                root,
                report_path,
                "FACTS_EVIDENCE",
                "repository facts/evidence must contain at least 80 words and two OBSERVED blocks",
            )
        metadata = "\n".join(text.splitlines()[: max(0, work_start - 1)])
        if row.repository not in metadata:
            add_error(
                errors,
                root,
                report_path,
                "REPORT_REPOSITORY",
                f"facts/evidence does not identify exact repository {row.repository}",
            )
        metadata_commits = [commit.lower() for commit in HEX40_RE.findall(metadata)]
        if not metadata_commits or metadata_commits[0] != row.study_commit:
            add_error(
                errors,
                root,
                report_path,
                "REPORT_COMMIT",
                "the first exact commit in facts/evidence must be the inventory study "
                f"commit {row.study_commit}",
            )
        branch_pattern = re.compile(
            rf"(?:`{re.escape(row.default_branch)}`|\b{re.escape(row.default_branch)}@{row.study_commit}\b)"
        )
        if not branch_pattern.search(metadata):
            add_error(
                errors,
                root,
                report_path,
                "REPORT_BRANCH",
                f"facts/evidence does not identify inventory default branch {row.default_branch!r}",
            )

    if work_matches and worked_matches:
        work_start = work_matches[0][0].line
        work_end = worked_matches[0][0].line - 1
        work_body = section_text(text, work_start, work_end)
        work_blocks = blocks_in_range(blocks, work_start, work_end)
        labelled_work = "\n".join(block.text for block in work_blocks)
        if word_count(work_body) < 140 or sum(
            block.label == "OBSERVED" for block in work_blocks
        ) < 3:
            add_error(
                errors,
                root,
                report_path,
                "WORK_CONTROL_EVIDENCE",
                "actual-use/work-control analysis must contain at least 140 words and three OBSERVED blocks",
            )
        topics = (
            ("GitHub Issue", r"\bIssues?\b"),
            ("Master", r"\bMaster\b"),
            ("Worker", r"\bWorkers?\b"),
            ("Architect", r"\bArchitect\b"),
            ("Verifier", r"\bVerifier\b"),
            ("Reviewer", r"\bReviewers?\b"),
            ("branch", r"\bbranch(?:es)?\b"),
            ("pull request / PR", r"\b(?:PRs?|pull requests?)\b"),
            ("Feature", r"\bFeatures?\b"),
            ("Refactor", r"\bRefactors?\b"),
            ("Fix", r"\bFix(?:es)?\b"),
            ("Sprint", r"\bSprints?\b"),
            ("Iteration", r"\bIterations?\b"),
            ("Slice", r"\bSlices?\b"),
            ("Study / Studies", r"\bStud(?:y|ies)\b"),
            ("review", r"\breviews?\b"),
            ("handoff", r"\bHandoff\b"),
            ("CurrentIndex", r"\bCurrentIndex\b"),
            ("Relations", r"\bRelations\b"),
            ("Ledger", r"\bLedger\b"),
            ("Steering", r"\bSteering\b"),
            ("CurrentAssignment", r"\bCurrentAssignment\b"),
            ("automation", r"\bautomat\w*\b"),
        )
        for topic, pattern in topics:
            if not re.search(pattern, labelled_work, re.IGNORECASE):
                add_error(
                    errors,
                    root,
                    report_path,
                    "MISSING_TOPIC",
                    f"evidence-labelled actual-use/work-control prose does not cover {topic}",
                )

    for matches, minimum, label in (
        (worked_matches, 30, "OBSERVED"),
        (pain_matches, 30, "OBSERVED"),
        (carry_matches, 25, "RECOMMENDATION"),
        (legacy_matches, 25, "RECOMMENDATION"),
    ):
        if len(matches) == 1:
            heading, _start, end = matches[0]
            validate_required_section(
                text, report_path, root, errors, heading, end, blocks, minimum, label
            )
    return h1_key


def is_external_target(target: str) -> bool:
    if target.startswith(("#", "/", "\\", "//")):
        return True
    return re.match(r"^[A-Za-z][A-Za-z0-9+.-]*:", target) is not None


def validate_links(
    text: str, path: Path, root: Path, errors: list[Diagnostic]
) -> None:
    visible = "\n".join(line for _number, line in lines_outside_fences(text))
    targets: list[str] = []
    for match in INLINE_LINK_RE.finditer(visible):
        targets.append(match.group(1) or match.group(2))
    for match in REFERENCE_LINK_RE.finditer(visible):
        targets.append(match.group(1) or match.group(2))
    for raw_target in targets:
        target = raw_target.strip()
        if is_external_target(target):
            continue
        parsed = urlsplit(target)
        local_part = unquote(parsed.path)
        if not local_part:
            continue
        candidate = (path.parent / local_part).resolve()
        if not candidate.exists():
            add_error(
                errors,
                root,
                path,
                "BROKEN_LINK",
                f"relative Markdown link {raw_target!r} does not resolve",
            )


def corpus_digest(paths: Iterable[Path], root: Path) -> str:
    digest = hashlib.sha256()
    for path in sorted(paths, key=lambda item: display_path(item, root)):
        relative = display_path(path, root).encode("utf-8")
        data = path.read_bytes()
        digest.update(relative)
        digest.update(b"\0")
        digest.update(len(data).to_bytes(8, "big"))
        digest.update(data)
    return digest.hexdigest()


def validate(analysis_root: Path) -> tuple[list[Diagnostic], Stats, str | None]:
    analysis_root = analysis_root.resolve()
    repository_root = analysis_root.parent
    errors: list[Diagnostic] = []
    readme_path = analysis_root / "README.md"
    inventory_path = analysis_root / "RepositoryInventory.md"
    corpus_paths: list[Path] = []

    texts: dict[Path, str] = {}
    for path in (readme_path, inventory_path):
        if not path.is_file():
            add_error(errors, repository_root, path, "MISSING_FILE", "required file is missing")
            continue
        text = read_utf8(path, repository_root, errors)
        if text is not None:
            texts[path] = text
            corpus_paths.append(path)
            validate_whitespace(text, path, repository_root, errors)
            validate_headings(text, path, repository_root, errors)
            validate_links(text, path, repository_root, errors)

    stats = Stats()
    included: list[InventoryRow] = []
    if inventory_path in texts:
        included, _excluded, stats = validate_inventory(
            texts[inventory_path], inventory_path, repository_root, errors
        )

    planned = {
        row.report_path: row for row in included if row.report_path is not None
    }
    reports_dir = analysis_root / "repositories"
    actual_reports = {
        path.relative_to(analysis_root).as_posix(): path
        for path in reports_dir.rglob("*.md")
    } if reports_dir.is_dir() else {}
    stats.reports = len(actual_reports)
    for relative in sorted(set(planned) - set(actual_reports)):
        add_error(
            errors,
            repository_root,
            analysis_root / relative,
            "MISSING_REPORT",
            "planned inventory report is missing",
        )
    for relative in sorted(set(actual_reports) - set(planned)):
        add_error(
            errors,
            repository_root,
            actual_reports[relative],
            "EXTRA_REPORT",
            "repository report is not present in the inventory manifest",
        )

    h1_identities: dict[str, str] = {}
    for relative in sorted(set(planned) & set(actual_reports)):
        path = actual_reports[relative]
        text = read_utf8(path, repository_root, errors)
        if text is None:
            continue
        corpus_paths.append(path)
        validate_whitespace(text, path, repository_root, errors)
        validate_links(text, path, repository_root, errors)
        h1_key = validate_report(planned[relative], path, text, repository_root, errors)
        if h1_key is not None:
            if h1_key in h1_identities:
                add_error(
                    errors,
                    repository_root,
                    path,
                    "DUPLICATE_REPORT_IDENTITY",
                    f"H1 identity duplicates {h1_identities[h1_key]}",
                )
            else:
                h1_identities[h1_key] = display_path(path, repository_root)

    digest = None
    if not errors and len(corpus_paths) == 2 + EXPECTED_IN_SCOPE:
        digest = corpus_digest(corpus_paths, repository_root)
    errors.sort(key=lambda item: (item.path.casefold(), item.code, item.message))
    return errors, stats, digest


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        description="Validate the deterministic Hans-Einar/SDP Issue #5 analysis corpus."
    )
    parser.add_argument(
        "--analysis-root",
        type=Path,
        default=Path(__file__).resolve().parent,
        help="analysis directory (default: directory containing this script)",
    )
    return parser


def main(argv: Sequence[str] | None = None) -> int:
    args = build_parser().parse_args(argv)
    errors, stats, digest = validate(args.analysis_root)
    if errors:
        for error in errors:
            print(f"ERROR [{error.code}] {error.path}: {error.message}")
        print(
            "FAIL analysis corpus: "
            f"errors={len(errors)} considered={stats.considered} "
            f"in_scope={stats.in_scope} excluded={stats.excluded} reports={stats.reports}"
        )
        return 1
    assert digest is not None
    print(
        "PASS analysis corpus: "
        f"considered={stats.considered} in_scope={stats.in_scope} "
        f"excluded={stats.excluded} reports={stats.reports} sha256={digest}"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
