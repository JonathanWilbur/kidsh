#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 3 ]]; then
  echo "usage: $0 <version> <tag> <sha256>" >&2
  exit 1
fi

VERSION="${1#v}"
TAG="$2"
SHA256="$3"
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

cat > "${ROOT}/kidsh.rb" <<EOF
class Kidsh < Formula
  desc "Lightweight secure shell for young children"
  homepage "https://github.com/JonathanWilbur/kidsh"
  url "https://github.com/JonathanWilbur/kidsh/releases/download/${TAG}/kidsh-${VERSION}.tar.gz"
  sha256 "${SHA256}"
  license "MIT"
  head "https://github.com/JonathanWilbur/kidsh.git", branch: "master"

  depends_on "go" => :build

  def install
    ldflags = "-s -w -X main.version=#{version}"
    system "go", "build", *std_go_args(ldflags: ldflags), "./src"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/kidsh -version")
  end
end
EOF
