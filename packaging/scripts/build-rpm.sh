#!/usr/bin/env bash
# Build an RPM from kidsh.spec. Intended to run inside Fedora/RHEL.
set -euo pipefail

VERSION="${VERSION:?VERSION is required}"
GOARCH="${GOARCH:-$(go env GOHOSTARCH 2>/dev/null || echo amd64)}"
RPM_ARCH="${RPM_ARCH:-}"
ROOT="$(pwd)"
TARBALL="${ROOT}/kidsh-${VERSION}.tar.gz"

if [[ -z "${RPM_ARCH}" ]]; then
  case "${GOARCH}" in
    amd64) RPM_ARCH=x86_64 ;;
    arm64) RPM_ARCH=aarch64 ;;
    riscv64) RPM_ARCH=riscv64 ;;
    *) RPM_ARCH="${GOARCH}" ;;
  esac
fi

if [[ ! -f "${TARBALL}" ]]; then
  echo "missing source tarball: ${TARBALL}" >&2
  exit 1
fi

dnf install -y rpm-build golang tar gzip ca-certificates

TOPDIR="${ROOT}/rpmbuild"
rm -rf "${TOPDIR}"
mkdir -p "${TOPDIR}"/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

cp "${TARBALL}" "${TOPDIR}/SOURCES/"
cp "${ROOT}/kidsh.spec" "${TOPDIR}/SPECS/"

rpmbuild -bb \
  --target "${RPM_ARCH}" \
  --define "_topdir ${TOPDIR}" \
  --define "debug_package %{nil}" \
  --define "go_arch ${GOARCH}" \
  "${TOPDIR}/SPECS/kidsh.spec"

find "${TOPDIR}/RPMS" -name '*.rpm' -exec cp {} "${ROOT}/" \;
find "${ROOT}" -maxdepth 1 -name '*.rpm' -print
