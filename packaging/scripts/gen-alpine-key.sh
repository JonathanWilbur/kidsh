#!/usr/bin/env bash
# Generate a stable Alpine abuild signing key.
#
# The public key is copied into packaging/alpine/ and should be committed.
# The private key stays on your machine; paste it into the GitHub Actions
# repository secret ALPINE_ABUILD_PRIVKEY.
set -euo pipefail

KEY_NAME="${1:-kidsh@wilbur.space.rsa}"
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
OUTDIR="${ABUILD_KEY_DIR:-$HOME/.abuild}"

mkdir -p "$OUTDIR" "$ROOT/packaging/alpine"
chmod 700 "$OUTDIR"

PRIV="$OUTDIR/$KEY_NAME"
PUB="$PRIV.pub"

if [[ -e "$PRIV" ]]; then
  echo "refusing to overwrite $PRIV" >&2
  exit 1
fi

openssl genrsa -out "$PRIV" 4096
chmod 600 "$PRIV"
openssl rsa -in "$PRIV" -pubout -out "$PUB"
cp "$PUB" "$ROOT/packaging/alpine/$KEY_NAME.pub"

cat <<EOF
Created:
  private: $PRIV
  public:  $ROOT/packaging/alpine/$KEY_NAME.pub

Add a GitHub Actions repository secret named ALPINE_ABUILD_PRIVKEY whose
value is the full PEM contents of the private key file.

If you used a non-default key name, also add repository variable
ALPINE_ABUILD_KEY_NAME=$KEY_NAME

Commit the public key so installs can trust it:
  git add packaging/alpine/$KEY_NAME.pub
EOF
