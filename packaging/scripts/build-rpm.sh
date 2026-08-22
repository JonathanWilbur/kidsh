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

dnf install -y rpm-build golang gcc tar gzip ca-certificates

HOST_ARCH="$(uname -m)"
if [[ "${RPM_ARCH}" != "${HOST_ARCH}" ]]; then
  case "${RPM_ARCH}" in
    aarch64)
      dnf install -y gcc-aarch64-linux-gnu
      export CC=aarch64-linux-gnu-gcc
      ;;
    riscv64)
      dnf install -y gcc-riscv64-linux-gnu
      export CC=riscv64-linux-gnu-gcc
      ;;
    *)
      echo "no C cross-compiler configured for RPM_ARCH=${RPM_ARCH} on ${HOST_ARCH}" >&2
      exit 1
      ;;
  esac
fi

TOPDIR="${ROOT}/rpmbuild"
rm -rf "${TOPDIR}"
mkdir -p "${TOPDIR}"/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

cp "${TARBALL}" "${TOPDIR}/SOURCES/"
cp "${ROOT}/kidsh.spec" "${TOPDIR}/SPECS/"

CC="${CC:-gcc}"

rpmbuild -bb \
  --target "${RPM_ARCH}" \
  --define "_topdir ${TOPDIR}" \
  --define "debug_package %{nil}" \
  --define "go_arch ${GOARCH}" \
  --define "go_cc ${CC}" \
  "${TOPDIR}/SPECS/kidsh.spec"

find "${TOPDIR}/RPMS" -name '*.rpm' -exec cp {} "${ROOT}/" \;
find "${ROOT}" -maxdepth 1 -name '*.rpm' -print
