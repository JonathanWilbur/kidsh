#!/bin/bash
# Local helper to produce an RPM. On GitHub Actions the same spec is built in Fedora.
set -euo pipefail

if [[ ! -f kidsh.spec ]]; then
  echo "Error: kidsh.spec not found. Run this script from the project root." >&2
  exit 1
fi

VERSION="${VERSION:-$(sed -n 's/^Version:[[:space:]]*//p' kidsh.spec | head -1)}"
VERSION="${VERSION#v}"

echo "Building RPM package for kidsh ${VERSION}"

git archive --format=tar.gz --prefix="kidsh-${VERSION}/" -o "kidsh-${VERSION}.tar.gz" HEAD
VERSION="${VERSION}" packaging/scripts/build-rpm.sh
