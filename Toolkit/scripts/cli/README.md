# Small SDP shell commands

K5 defines a read-only `kanban` command. The local CardState convention is in
[the KanBan workflow](../../../SDP/KanBan/README.md). This directory contains
CLI entry scripts; the installer copies every `*.sh` file into `~/bin`, stripping
`.sh` (including `install-cli`). No compilation, daemon or shell startup changes.

```sh
bash Toolkit/scripts/cli/install-cli.sh
cd SDP/KanBan
kanban status
kanban state
cd active
kanban status
```

Both commands group cards by CardState, sorted alphabetically, with paths relative
to the selected directory. At board root they inspect the seven lifecycle
subdirectories; within one status directory they inspect only its cards. An
optional directory overrides the working directory: `kanban state /path/to/KanBan`.
They do not recurse into other projects, follow symlink cards or modify files.
Paths are shell-quoted for safe terminal display and copying, including spaces/#.
On a terminal (except TERM=dumb), each path is an OSC 8 hyperlink to its absolute,
percent-encoded file URL. Pipes and redirected output stay plain text. The terminal
opens links using its configured file handler; for example, Kitty can route .md
files to XFMD. The command does not launch a viewer or change that configuration.

Only `#*.md` cards with a CardState row in the first visible Field/Value table
immediately after the title are listed. Body examples and templates are ignored.
Cards without the field are skipped. If none match, stderr is exactly
`no KanBan cards found` and exit status is 1. Invalid/duplicate states or arguments
return 2 without partial output; success/help return 0. State/folder consistency
and ledger validation remain workflow checks, not mutation features of the lister.

## Installation and verification

Requires Bash, awk, find, sort and GNU coreutils on Linux. Installed commands need
`~/bin` in PATH. `--bin-dir DIRECTORY` selects a different destination. Existing
regular files with different content are preserved as uniquely named backups;
identical files are left unchanged except making them executable. Symlink/directory
destinations are rejected. Each script replacement is atomic; installation of all
scripts is not a transactional package upgrade.

The installed installer can refresh from any checkout:

```sh
install-cli --source /path/to/SDP/Toolkit/scripts/cli
python3 Toolkit/scripts/cli/test_cli.py
bash -n Toolkit/scripts/cli/*.sh
```

Only shell scripts are installed; tests and this README stay in the repository.
This is a local utility, not `sdptool`, a KanBan graph or a released compatibility
contract. Shared versioning/template distribution has separate backlog scope.

## Optional Sprint and Scrum grouping — PM1

```sh
kanban status --group-by sprint
kanban status --group-by scrum
kanban status --sprint SPR-SDP-0001
kanban status --scrum SCRUM-SDP-0001 --group-by sprint
kanban status /path/to/KanBan/active --sprint SPR-SDP-0001
```

SprintId and ScrumId are optional metadata rows. Grouping displays current
CardState beside each file; unassigned cards appear under `(none)`. Filters match
full IDs exactly and combine with AND. No matching cards retains the existing
`no KanBan cards found` error. Invalid/duplicate grouping metadata returns 2 with
no partial output. Only the first visible metadata table is read; body examples
cannot assign a card to a sprint. Default status/state behavior and directory
scope are unchanged, including terminal OSC 8 links and plain redirected output.

These options read metadata; they do not validate the ledger, edit membership,
start Sprints or move cards. The [management validator](../../../SDP/ProjectManagement/validate.py)
checks current membership against the common ledger and Sprint records.
The example Sprint above is illustrative; Scrum-0001 selected Maintenance and did
not create an actual Sprint. No script is required to create a management event.
