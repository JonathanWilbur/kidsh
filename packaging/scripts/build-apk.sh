#!/bin/sh
# Build a signed Alpine apk from packaging/alpine/APKBUILD.
# Intended to run inside Alpine on the target architecture (use docker --platform).
set -eu

VERSION="${VERSION:?VERSION is required}"
ROOT="$(pwd)"
TARBALL="${ROOT}/kidsh-${VERSION}.tar.gz"
KEY_NAME="${ALPINE_ABUILD_KEY_NAME:-kidsh@wilbur.space.rsa}"
PRIVKEY_SRC="${ALPINE_ABUILD_PRIVKEY_FILE:-}"
CARCH="${CARCH:-$(uname -m)}"
HOST_ARCH="$(uname -m)"
if [ "${CARCH}" != "${HOST_ARCH}" ]; then
  echo "Alpine packages are built natively (cgo PIE). Container is ${HOST_ARCH}, requested CARCH=${CARCH}." >&2
  echo "Run this script in an Alpine image for ${CARCH} (docker --platform)." >&2
  exit 1
fi
export CARCH

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

# Keep APKINDEX in cache so abuild does not warn about missing CDN indexes.
apk add alpine-sdk go git sudo ca-certificates openssl

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

printf '%s\n' \
  "PACKAGER_PRIVKEY=${ABUILD_DIR}/${KEY_NAME}" \
  'PACKAGER="Jonathan M. Wilbur <jonathan@wilbur.space>"' \
  "REPODEST=/home/packager/packages" \
  >"${ABUILD_DIR}/abuild.conf"
cp "${ABUILD_DIR}/${KEY_NAME}.pub" "/etc/apk/keys/${KEY_NAME}.pub"
cp "${ABUILD_DIR}/${KEY_NAME}.pub" "${ROOT}/${KEY_NAME}.pub"
chown -R packager:packager /home/packager

# Seed a signed empty local repo so first abuild -r does not warn on APKINDEX.
REPO_ARCH="/home/packager/packages/aports/${CARCH}"
mkdir -p "${REPO_ARCH}"

WORKDIR="/home/packager/aports/kidsh"
rm -rf "${WORKDIR}"
mkdir -p "${WORKDIR}"
cp "${ROOT}/packaging/alpine/APKBUILD" "${WORKDIR}/APKBUILD"
cp "${TARBALL}" "${WORKDIR}/"
sed -i "s/^pkgver=.*/pkgver=${VERSION}/" "${WORKDIR}/APKBUILD"
chown -R packager:packager /home/packager

su packager -c "cd '${REPO_ARCH}' && apk index -o APKINDEX.tar.gz && abuild-sign APKINDEX.tar.gz"
su packager -c "cd '${WORKDIR}' && CARCH='${CARCH}' abuild checksum && CARCH='${CARCH}' abuild -r"

find /home/packager/packages -name 'kidsh-*.apk' ! -name '*-dev-*' -exec cp {} "${ROOT}/" \;
find "${ROOT}" -maxdepth 1 -name '*.apk' -print
