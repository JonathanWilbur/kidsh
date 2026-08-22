# Kid Shell

This is a lightweight / low-capability secure shell for young children to fool
around on a computer. It has simple commands for doing easy tasks in a simple
way that children can understand, such as pressing `c` to display colors, or
`d` to display days of the week, or `n` to take notes. It compiles as a static
binary that only supports built-ins, and (not implemented yet) it will be able
to ran with dropped capabilities, so your little one cannot accidentally do
anything to mess up your computer.

## Build

```bash
go build -o kidsh ./src
```

## Install

Publishing a GitHub Release (from a tag such as `v1.0.0`) builds packages and
attaches them to that release as assets. Linux packages are produced for
amd64, ARM64, and RISC-V. macOS archives cover Apple Silicon, Intel, and a
universal binary. Snap and Flatpak are built for amd64 and ARM64. The Homebrew
formula compiles for whatever Mac or Linux machine you install on.

Before tagging a release, update `debian/changelog` and the `%changelog`
section in `kidsh.spec` by hand. The Debian package version is taken from
`debian/changelog`, so that entry must match the git tag.

```bash
# Debian / Ubuntu (also _arm64.deb, _riscv64.deb)
sudo apt install ./kidsh_*_amd64.deb

# Fedora / RHEL (also .aarch64.rpm, .riscv64.rpm)
sudo rpm -i kidsh-*.x86_64.rpm

# Alpine (install the matching *.rsa.pub from the release first)
sudo cp kidsh@wilbur.space.rsa.pub /etc/apk/keys/
sudo apk add ./kidsh-*.apk

# Arch Linux (also aarch64, riscv64)
sudo pacman -U ./kidsh-*-x86_64.pkg.tar.zst

# Snap
sudo snap install --dangerous ./kidsh_*.snap

# Flatpak (also kidsh-aarch64.flatpak)
flatpak install --user ./kidsh-x86_64.flatpak

# macOS tarball (Apple Silicon, Intel, or universal)
tar -xf kidsh-*-darwin-universal.tar.gz

# Homebrew (Linux or macOS)
brew install --formula ./kidsh.rb
```

Alpine packages are signed. Generate a key once with
`packaging/scripts/gen-alpine-key.sh`, store the private key as the
`ALPINE_ABUILD_PRIVKEY` repository secret, and commit the public
`*.rsa.pub` file it prints.

Optional configuration is read from `$KIDSH_CONFIG`, then
`~/.config/kidsh/config.json`, then `/etc/kidsh.json`. An example file ships as
`kidsh.json.example` in the package docs.

## Usage

I expect users to full-screen the window where this shell is running so their
kids cannot easily escape it. Even better would be to define a boot menu
configuration where Linux defines `kidsh` as PID 0, so they cannot escape it.
Alternatively, you could just open up this shell in a real terminal that is not
running a GUI.
