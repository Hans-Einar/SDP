#!/usr/bin/env bash
# SDP CLI: install CLI scripts without shell startup changes or compilation.
set -euo pipefail
source_dir=$(dirname -- "$(readlink -f -- "${BASH_SOURCE[0]}")")
bin_dir=$HOME/bin
usage() { printf '%s\n' 'Usage: install-cli.sh [--bin-dir DIRECTORY] [--source DIRECTORY]' 'Copy every *.sh CLI script without its suffix (including install-cli).'; }
while (( $# )); do
  case $1 in
    --bin-dir|--source)
      (( $# >= 2 )) || { usage >&2; exit 2; }
      if [[ $1 == --bin-dir ]]; then bin_dir=$2; else source_dir=$2; fi
      shift 2 ;;
    -h|--help) usage; exit 0 ;;
    *) usage >&2; exit 2 ;;
  esac
done
[[ -d $source_dir ]] || { printf 'invalid CLI source: %s\n' "$source_dir" >&2; exit 2; }
source_dir=$(cd -- "$source_dir" && pwd -P)
shopt -s nullglob
scripts=("$source_dir"/*.sh)
(( ${#scripts[@]} )) || { printf '%s\n' 'no CLI scripts found; use --source with Toolkit/scripts/cli' >&2; exit 1; }
mkdir -p -- "$bin_dir"
bin_dir=$(cd -- "$bin_dir" && pwd -P)
for script in "${scripts[@]}"; do
  name=${script##*/}; name=${name%.sh}
  [[ $name =~ ^[a-zA-Z0-9_-]+$ ]] || { printf 'invalid CLI name: %s\n' "$name" >&2; exit 2; }
  target=$bin_dir/$name
  if [[ -L $target || ( -e $target && ! -f $target ) ]]; then
    printf 'refusing non-regular destination: %s\n' "$target" >&2; exit 2
  fi
done
tmp=
trap '[[ -z $tmp ]] || rm -f -- "$tmp"' EXIT
for script in "${scripts[@]}"; do
  name=${script##*/}; name=${name%.sh}; target=$bin_dir/$name
  if [[ -f $target ]] && cmp -s -- "$script" "$target"; then
    chmod u+x -- "$target"; printf 'Unchanged: %s\n' "$target"; continue
  fi
  if [[ -f $target ]]; then
    backup=$(mktemp "$target.backup.XXXXXXXX")
    cp -p -- "$target" "$backup"
    printf 'Preserved previous file: %s\n' "$backup"
  fi
  tmp=$(mktemp "$bin_dir/.sdp-cli.XXXXXXXX")
  install -m 755 -- "$script" "$tmp"
  mv -f -- "$tmp" "$target"; tmp=
  printf 'Installed: %s\n' "$target"
done
