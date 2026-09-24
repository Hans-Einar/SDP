#!/usr/bin/env bash
# SDP CLI: read-only listing of visible CardState metadata.
set -euo pipefail
export LC_ALL=C
usage() { printf '%s\n' 'Usage: kanban status|state [DIRECTORY]' 'List cards by CardState. Default directory: current working directory.'; }
case ${1:-} in
  -h|--help|help) usage; exit 0 ;;
  status|state) shift ;;
  *) usage >&2; exit 2 ;;
esac
if (( $# > 1 )); then usage >&2; exit 2; fi
scope=${1:-.}
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
      }
      END { if (count) printf "%d\t%s", count, state }
    ' FS='|' "$card")
    [[ -n $result ]] || continue
    count=${result%%$'\t'*}; state=${result#*$'\t'}
    case $state in backlog|queued|ready|in-progress|gate-review|onHold|completed|canceled|superseded|irrelevant) ;; *) count=invalid ;; esac
    if [[ $count != 1 ]]; then
      printf 'invalid or duplicate CardState: %s\n' "$card" >&2; invalid=1; continue
    fi
    relative=${card#"$scope"/}
    # Shell quoting keeps tabs/newlines/control bytes in filenames printable.
    printf -v printable '%q' "$relative"
    uri=
    if (( links )); then uri=$(file_uri "$card"); fi
    rows+=("$state"$'\t'"$printable"$'\t'"$uri")
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
