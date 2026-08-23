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
GO_CGO=1
GO_SYSROOT=""
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
  # Fedora's *-linux-gnu GCC is a kernel/bare-metal toolchain. User-space CGO
  # needs the matching glibc sysroot (Debian's equivalent of libc6-dev-*-cross).
  # GCC 16 also injects -latomic_asneeded; kidsh.spec disables that via
  # -fno-link-libatomic because the cross compiler does not ship libatomic.
  fedora_ver="$(rpm -E %fedora)"
  sysroot_pkg="sysroot-${RPM_ARCH}-fc${fedora_ver}-glibc"
  if dnf install -y "${sysroot_pkg}" || dnf install -y "sysroot-${RPM_ARCH}-glibc"; then
    GO_SYSROOT="/usr/${RPM_ARCH}-redhat-linux/sys-root/fc${fedora_ver}"
    if [[ ! -d "${GO_SYSROOT}/usr/include" ]]; then
      echo "sysroot headers missing at ${GO_SYSROOT}/usr/include" >&2
      exit 1
    fi
  else
    echo "no user-space glibc sysroot for ${RPM_ARCH}; building with CGO disabled" >&2
    GO_CGO=0
  fi
fi

TOPDIR="${ROOT}/rpmbuild"
rm -rf "${TOPDIR}"
mkdir -p "${TOPDIR}"/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

cp "${TARBALL}" "${TOPDIR}/SOURCES/"
cp "${ROOT}/kidsh.spec" "${TOPDIR}/SPECS/"

CC="${CC:-gcc}"

rpmbuild_args=(
  -bb
  --target "${RPM_ARCH}"
  --define "_topdir ${TOPDIR}"
  --define "debug_package %{nil}"
  --define "go_arch ${GOARCH}"
  --define "go_cc ${CC}"
  --define "go_cgo ${GO_CGO}"
)
if [[ -n "${GO_SYSROOT}" ]]; then
  rpmbuild_args+=(--define "go_sysroot ${GO_SYSROOT}")
fi

rpmbuild "${rpmbuild_args[@]}" "${TOPDIR}/SPECS/kidsh.spec"

find "${TOPDIR}/RPMS" -name '*.rpm' -exec cp {} "${ROOT}/" \;
find "${ROOT}" -maxdepth 1 -name '*.rpm' -print
