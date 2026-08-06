#!/usr/bin/env bash
# Downloads the prebuilt zktf-sdk native archive (matching zktf-sdk-version)
# for the host (or requested) target triple, and prints the CGO/linker env
# vars needed to build against it.
#
# Usage:
#   scripts/fetch-native.sh                # fetch for GOHOSTOS/GOHOSTARCH
#   scripts/fetch-native.sh linux amd64    # fetch for an explicit GOOS/GOARCH
#
# Output: writes a sourceable .env file at the repo root (CGO_CFLAGS,
# CGO_LDFLAGS, LD_LIBRARY_PATH) and echoes the same `export` lines to stdout.
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"

version_file="$repo_root/zktf-sdk-version"
if [[ ! -f "$version_file" ]]; then
  echo "fetch-native: missing $version_file" >&2
  exit 1
fi
version="$(tr -d '[:space:]' < "$version_file")"

goos="${1:-$(go env GOHOSTOS 2>/dev/null || go env GOOS)}"
goarch="${2:-$(go env GOHOSTARCH 2>/dev/null || go env GOARCH)}"

case "$goos/$goarch" in
  linux/amd64)  triple=x86_64-unknown-linux-gnu ;;
  linux/arm64)  triple=aarch64-unknown-linux-gnu ;;
  darwin/arm64) triple=aarch64-apple-darwin ;;
  *)
    echo "fetch-native: unsupported GOOS/GOARCH combination: $goos/$goarch" >&2
    echo "fetch-native: supported: linux/amd64, linux/arm64, darwin/arm64" >&2
    exit 1
    ;;
esac

native_root="$repo_root/.zktf-native"
dest_dir="$native_root/$triple-$version"
archive_name="zktf-sdk-$triple-$version.tar.gz"
gcs_https_url="https://storage.googleapis.com/download.joinself.com/zktf-sdk/$archive_name"
gcs_uri="gs://download.joinself.com/zktf-sdk/$archive_name"

if [[ -f "$dest_dir/zktf-sdk.h" ]]; then
  echo "fetch-native: $triple @ $version already present at $dest_dir" >&2
else
  mkdir -p "$dest_dir"
  tarball="$native_root/$archive_name"

  download_ok=0
  if command -v curl >/dev/null 2>&1; then
    echo "fetch-native: downloading via curl: $gcs_https_url" >&2
    if curl -fsSL -o "$tarball" "$gcs_https_url"; then
      download_ok=1
    fi
  fi
  if [[ "$download_ok" -eq 0 ]] && command -v wget >/dev/null 2>&1; then
    echo "fetch-native: downloading via wget: $gcs_https_url" >&2
    if wget -q -O "$tarball" "$gcs_https_url"; then
      download_ok=1
    fi
  fi
  if [[ "$download_ok" -eq 0 ]] && command -v gcloud >/dev/null 2>&1; then
    echo "fetch-native: downloading via gcloud storage cp: $gcs_uri" >&2
    if gcloud storage cp "$gcs_uri" "$tarball"; then
      download_ok=1
    fi
  fi
  if [[ "$download_ok" -eq 0 ]]; then
    echo "fetch-native: failed to download $archive_name via curl, wget, or gcloud storage cp" >&2
    exit 1
  fi

  tar -xzf "$tarball" -C "$dest_dir" --strip-components=1
  rm -f "$tarball"
fi

env_file="$repo_root/.env"
{
  echo "CGO_CFLAGS=-I$dest_dir"
  echo "CGO_LDFLAGS=-L$dest_dir -lzktf_sdk"
  echo "LD_LIBRARY_PATH=$dest_dir"
} > "$env_file"

echo "fetch-native: wrote $env_file" >&2
echo "export CGO_CFLAGS=\"-I$dest_dir\""
echo "export CGO_LDFLAGS=\"-L$dest_dir -lzktf_sdk\""
echo "export LD_LIBRARY_PATH=\"$dest_dir\""
