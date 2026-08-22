#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <version>" >&2
  exit 1
fi

VERSION="${1#v}"
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT"

sed -i "s/^Version:.*/Version:        ${VERSION}/" kidsh.spec
sed -i "s/^pkgver=.*/pkgver=${VERSION}/" packaging/alpine/APKBUILD packaging/arch/PKGBUILD
sed -i "s/^version:.*/version: '${VERSION}'/" snap/snapcraft.yaml

printf '%s\n' "${VERSION}" > VERSION
