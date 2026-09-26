#!/usr/bin/env bash
# Maintainer build only. Viewing documentation never invokes this script.
set -euo pipefail
if [[ $# != 1 ]]; then
    echo "Usage: package.sh NEW_OUTPUT_DIRECTORY" >&2
    exit 2
fi
sdp_package_module=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
sdp_package_output=$(realpath -m -- "$1")
if [[ -e "$sdp_package_output" ]]; then
    echo "Package destination already exists: $sdp_package_output" >&2
    exit 2
fi
sdp_package_go=${SDP_GO:-go}
command -v -- "$sdp_package_go" >/dev/null || { echo "Go build tool unavailable; set SDP_GO" >&2; exit 2; }
sdp_package_parent=$(dirname -- "$sdp_package_output")
mkdir -p -- "$sdp_package_parent"
sdp_package_temp=$(mktemp -d "$sdp_package_parent/.sdptool-package.XXXXXX")
trap 'rm -rf -- "$sdp_package_temp"' EXIT
sdp_package_revision=$(git -C "$sdp_package_module" rev-parse HEAD 2>/dev/null || echo unknown)
if [[ -n "$(git -C "$sdp_package_module" status --porcelain 2>/dev/null)" ]]; then
    sdp_package_revision="$sdp_package_revision-dirty"
fi
(
    cd -- "$sdp_package_module"
    GOTOOLCHAIN=local "$sdp_package_go" build -trimpath -mod=readonly         -ldflags "-X github.com/Hans-Einar/SDP/Toolkit/SDPTool.BuildRevision=$sdp_package_revision"         -o "$sdp_package_temp/sdptool" ./cmd/sdptool
)
"$sdp_package_temp/sdptool" --version > "$sdp_package_temp/sdptool.manifest.json"
(
    cd -- "$sdp_package_temp"
    sha256sum sdptool sdptool.manifest.json > SHA256SUMS
)
mv -- "$sdp_package_temp" "$sdp_package_output"
trap - EXIT
echo "$sdp_package_output"
