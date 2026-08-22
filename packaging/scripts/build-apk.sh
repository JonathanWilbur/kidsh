#!/bin/sh
# Build a signed Alpine apk from packaging/alpine/APKBUILD. Intended to run inside Alpine.
set -eu

VERSION="${VERSION:?VERSION is required}"
ROOT="$(pwd)"
TARBALL="${ROOT}/kidsh-${VERSION}.tar.gz"
KEY_NAME="${ALPINE_ABUILD_KEY_NAME:-kidsh@wilbur.space.rsa}"
PRIVKEY_SRC="${ALPINE_ABUILD_PRIVKEY_FILE:-}"
CARCH="${CARCH:-}"
GOARCH="${GOARCH:-}"

if [ -z "${CARCH}" ]; then
  CARCH="$(uname -m)"
fi
if [ -z "${GOARCH}" ]; then
  case "${CARCH}" in
    x86_64) GOARCH=amd64 ;;
    aarch64) GOARCH=arm64 ;;
    riscv64) GOARCH=riscv64 ;;
    *) GOARCH="${CARCH}" ;;
  esac
fi
export CARCH GOARCH

if [ ! -f "${TARBALL}" ]; then
  echo "missing source tarball: ${TARBALL}" >&2
  exit 1
fi

if [ -z "${PRIVKEY_SRC}" ] || [ ! -f "${PRIVKEY_SRC}" ]; then
  echo "ALPINE_ABUILD_PRIVKEY_FILE must point to the abuild RSA private key" >&2
  echo "Generate one with packaging/scripts/gen-alpine-key.sh and add it as" >&2
  echo "the ALPINE_ABUILD_PRIVKEY GitHub Actions secret." >&2
  exit 1
fi

apk add --no-cache alpine-sdk go git sudo ca-certificates openssl

if ! id packager >/dev/null 2>&1; then
  adduser -D packager
  addgroup packager abuild
fi

echo "packager ALL=(ALL) NOPASSWD: ALL" >/etc/sudoers.d/packager

ABUILD_DIR="/home/packager/.abuild"
mkdir -p "${ABUILD_DIR}"
cp "${PRIVKEY_SRC}" "${ABUILD_DIR}/${KEY_NAME}"
chmod 600 "${ABUILD_DIR}/${KEY_NAME}"
openssl rsa -in "${ABUILD_DIR}/${KEY_NAME}" -pubout -out "${ABUILD_DIR}/${KEY_NAME}.pub"

REPO_PUB="${ROOT}/packaging/alpine/${KEY_NAME}.pub"
if [ -f "${REPO_PUB}" ] && ! cmp -s "${REPO_PUB}" "${ABUILD_DIR}/${KEY_NAME}.pub"; then
  echo "signing key does not match committed public key ${REPO_PUB}" >&2
  exit 1
fi

printf 'PACKAGER_PRIVKEY=%s\nPACKAGER="Jonathan M. Wilbur <jonathan@wilbur.space>"\n' \
  "${ABUILD_DIR}/${KEY_NAME}" >"${ABUILD_DIR}/abuild.conf"
cp "${ABUILD_DIR}/${KEY_NAME}.pub" "/etc/apk/keys/${KEY_NAME}.pub"
cp "${ABUILD_DIR}/${KEY_NAME}.pub" "${ROOT}/${KEY_NAME}.pub"
chown -R packager:packager /home/packager

WORKDIR="/home/packager/aports/kidsh"
rm -rf "${WORKDIR}"
mkdir -p "${WORKDIR}"
cp "${ROOT}/packaging/alpine/APKBUILD" "${WORKDIR}/APKBUILD"
cp "${TARBALL}" "${WORKDIR}/"
sed -i "s/^pkgver=.*/pkgver=${VERSION}/" "${WORKDIR}/APKBUILD"
chown -R packager:packager /home/packager

su packager -c "cd '${WORKDIR}' && CARCH='${CARCH}' GOARCH='${GOARCH}' abuild checksum && CARCH='${CARCH}' GOARCH='${GOARCH}' abuild -r"

find /home/packager/packages -name 'kidsh-*.apk' ! -name '*-dev-*' ! -name '*-doc-*' -exec cp {} "${ROOT}/" \;
find "${ROOT}" -maxdepth 1 -name '*.apk' -print
