#!/usr/bin/env bash
# SDP CLI: read-only listing of visible CardState metadata.
set -euo pipefail
export LC_ALL=C
usage() { printf '%s\n' 'Usage: kanban status|state [DIRECTORY] [--group-by state|sprint|scrum] [--sprint ID] [--scrum ID]' 'Default directory: current working directory. Unassigned cards group under (none).'; }
case ${1:-} in
  -h|--help|help) usage; exit 0 ;;
  status|state) shift ;;
  *) usage >&2; exit 2 ;;
esac
scope=.
has_scope=0
group_by=state
sprint_filter=
scrum_filter=
while (( $# )); do
  case $1 in
    --group-by|--sprint|--scrum)
      (( $# >= 2 )) || { usage >&2; exit 2; }
      case $1 in
        --group-by) group_by=$2 ;;
        --sprint) sprint_filter=$2 ;;
        --scrum) scrum_filter=$2 ;;
      esac
      [[ -n $2 ]] || { usage >&2; exit 2; }
      shift 2 ;;
    --)
      shift
      (( $# == 1 && has_scope == 0 )) || { usage >&2; exit 2; }
      scope=$1; has_scope=1; shift ;;
    -*) usage >&2; exit 2 ;;
    *)
      (( has_scope == 0 )) || { usage >&2; exit 2; }
      scope=$1; has_scope=1; shift ;;
  esac
done
case $group_by in state|sprint|scrum) ;; *) usage >&2; exit 2 ;; esac
[[ -z $sprint_filter || $sprint_filter =~ ^SPR-[A-Z][A-Z0-9]*-[0-9]+$ ]] || { usage >&2; exit 2; }
[[ -z $scrum_filter || $scrum_filter =~ ^SCRUM-[A-Z][A-Z0-9]*-[0-9]+$ ]] || { usage >&2; exit 2; }
if [[ ! -d $scope ]]; then printf 'not a directory: %s\n' "$scope" >&2; exit 2; fi
scope=$(cd -- "$scope" && pwd -P)
# A board covers its own status directories. A status directory covers itself.
# Never search sibling projects or arbitrary descendants of the working directory.
dirs=("$scope")
case ${scope##*/} in
  backlog|active|onHold|completed|canceled|superseded|irrelevant) ;;
  *)
    for name in backlog active onHold completed canceled superseded irrelevant; do
      [[ ! -d $scope/$name ]] || dirs+=("$scope/$name")
    done ;;
esac
rows=()
invalid=0
links=0
if [[ -t 1 && ${TERM:-} != dumb ]]; then links=1; fi
file_uri() {
  local path=$1 encoded= byte hex i
  # LC_ALL=C makes this byte-wise, including percent-encoding UTF-8 filenames.
  for (( i=0; i<${#path}; i++ )); do
    byte=${path:i:1}
    case $byte in
      [a-zA-Z0-9/._~-]) encoded+=$byte ;;
      *) printf -v hex '%%%02X' "'$byte"; encoded+=$hex ;;
    esac
  done
  printf 'file://%s' "$encoded"
}
for directory in "${dirs[@]}"; do
  while IFS= read -r -d '' card; do
    # Only the first visible metadata table is authoritative; ignore examples/body.
    result=$(awk '
      { sub(/\r$/, "") }
      /^\|[[:space:]]*Field[[:space:]]*\|[[:space:]]*Value[[:space:]]*\|[[:space:]]*$/ { table=1; next }
      !table && (/^#[[:space:]]/ || /^[[:space:]]*$/) { next }
      !table { exit }
      table && !/^\|/ { exit }
      table {
        key=$2; value=$3
        gsub(/^[[:space:]]+|[[:space:]]+$/, "", key)
        gsub(/^[[:space:]]+|[[:space:]]+$/, "", value)
        if (key == "CardState") { count++; state=value }
        if (key == "SprintId") { sprint_count++; sprint=value }
        if (key == "ScrumId") { scrum_count++; scrum=value }
      }
      END { if (count) printf "%d|%s|%s|%s|%d|%d", count, state, sprint, scrum, sprint_count, scrum_count }
    ' FS='|' "$card")
    [[ -n $result ]] || continue
    IFS='|' read -r count state sprint scrum sprint_count scrum_count <<< "$result"
    case $state in backlog|queued|ready|in-progress|gate-review|onHold|completed|canceled|superseded|irrelevant) ;; *) count=invalid ;; esac
    if [[ $count != 1 ]]; then
      printf 'invalid or duplicate CardState: %s\n' "$card" >&2; invalid=1; continue
    fi
    if (( sprint_count > 1 || scrum_count > 1 )) ||
       { (( sprint_count )) && [[ ! $sprint =~ ^SPR-[A-Z][A-Z0-9]*-[0-9]+$ ]]; } ||
       { (( scrum_count )) && [[ ! $scrum =~ ^SCRUM-[A-Z][A-Z0-9]*-[0-9]+$ ]]; }; then
      printf 'invalid or duplicate SprintId/ScrumId: %s\n' "$card" >&2; invalid=1; continue
    fi
    [[ -z $sprint_filter || $sprint == "$sprint_filter" ]] || continue
    [[ -z $scrum_filter || $scrum == "$scrum_filter" ]] || continue
    group=$state
    case $group_by in sprint) group=${sprint:-(none)} ;; scrum) group=${scrum:-(none)} ;; esac
    relative=${card#"$scope"/}
    # Shell quoting keeps tabs/newlines/control bytes in filenames printable.
    printf -v printable '%q' "$relative"
    uri=
    if (( links )); then uri=$(file_uri "$card"); fi
    [[ $group_by == state ]] || printable+=" [$state]"
    rows+=("$group"$'\t'"$printable"$'\t'"$uri")
  done < <(find "$directory" -maxdepth 1 -type f -name '#*.md' -print0)
done
(( invalid == 0 )) || exit 2
if (( ${#rows[@]} == 0 )); then printf '%s\n' 'no KanBan cards found' >&2; exit 1; fi
printf '%s\n' "${rows[@]}" | sort -t $'\t' -k1,1 -k2,2 | awk -F '\t' -v links="$links" '
  $1 != previous { if (NR > 1) print ""; print $1 ":"; previous=$1 }
  {
    if (links) printf "  \033]8;;%s\033\\%s\033]8;;\033\\\n", $3, $2
    else print "  " $2
  }
'
