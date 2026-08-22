#!/usr/bin/env bash
# Build an Arch package from packaging/arch/PKGBUILD. Intended to run inside Arch.
set -euo pipefail

VERSION="${VERSION:?VERSION is required}"
ROOT="$(pwd)"
TARBALL="${ROOT}/kidsh-${VERSION}.tar.gz"
CARCH="${CARCH:-}"
GOARCH="${GOARCH:-}"

if [[ -z "${CARCH}" ]]; then
  CARCH="$(uname -m)"
fi
if [[ -z "${GOARCH}" ]]; then
  case "${CARCH}" in
    x86_64) GOARCH=amd64 ;;
    aarch64) GOARCH=arm64 ;;
    riscv64) GOARCH=riscv64 ;;
    *) GOARCH="${CARCH}" ;;
  esac
fi
export CARCH GOARCH

if [[ ! -f "${TARBALL}" ]]; then
  echo "missing source tarball: ${TARBALL}" >&2
  exit 1
fi

pacman -Sy --noconfirm archlinux-keyring
pacman -Syu --noconfirm --needed base-devel go git

if [[ "${CARCH}" != "$(uname -m)" ]]; then
  pacman -S --noconfirm --needed "${CARCH}-linux-gnu-gcc"
  export CC="${CARCH}-linux-gnu-gcc"
fi

if ! id packager >/dev/null 2>&1; then
  useradd -m packager
fi

WORKDIR="/home/packager/pkg"
rm -rf "${WORKDIR}"
mkdir -p "${WORKDIR}"
cp "${ROOT}/packaging/arch/PKGBUILD" "${WORKDIR}/PKGBUILD"
cp "${TARBALL}" "${WORKDIR}/"
sed -i "s/^pkgver=.*/pkgver=${VERSION}/" "${WORKDIR}/PKGBUILD"
chown -R packager:packager "${WORKDIR}"

MAKEPKG_ENV="CARCH='${CARCH}' GOARCH='${GOARCH}'"
if [[ -n "${CC:-}" ]]; then
  MAKEPKG_ENV="${MAKEPKG_ENV} CC='${CC}'"
fi
su packager -c "cd '${WORKDIR}' && ${MAKEPKG_ENV} makepkg -f --noconfirm"

find "${WORKDIR}" -maxdepth 1 -name '*.pkg.tar.zst' -exec cp {} "${ROOT}/" \;
find "${ROOT}" -maxdepth 1 -name '*.pkg.tar.zst' -print
