Name:           kidsh
Version:        1.0.0
Release:        1
Summary:        A lightweight secure shell for young children

License:        MIT
URL:            https://github.com/JonathanWilbur/kidsh
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang >= 1.22
BuildRequires:  gcc

%global debug_package %{nil}
%global __strip /bin/true

%description
Kid Shell is a lightweight, low-capability secure shell designed for young
children to safely explore and interact with a computer. It provides simple
commands for basic tasks that children can understand, such as displaying
colors, showing days of the week, and taking notes. The shell supports
only built-in commands to ensure security.

%prep
%autosetup

%build
export CGO_ENABLED=%{?go_cgo:%{go_cgo}}%{!?go_cgo:1}
export GOOS=linux
export CC=%{?go_cc:%{go_cc}}%{!?go_cc:gcc}
%if 0%{?go_arch:1}
export GOARCH=%{go_arch}
%endif
%if 0%{?go_sysroot:1}
export CGO_CFLAGS="${CGO_CFLAGS:+${CGO_CFLAGS} }--sysroot=%{go_sysroot}"
# Fedora's cross GCC 16 *lib spec always passes -latomic_asneeded, but the
# kernel-oriented *-linux-gnu toolchain does not ship that virtual library.
export CGO_LDFLAGS="${CGO_LDFLAGS:+${CGO_LDFLAGS} }--sysroot=%{go_sysroot} -fno-link-libatomic"
%endif
go build -trimpath -buildvcs=false -ldflags "-s -w -X main.version=%{version}" -o kidsh ./src

%check
native="$(go env GOHOSTARCH)"
%if 0%{?go_arch:1}
target="%{go_arch}"
%else
target="$(go env GOARCH)"
%endif
# rpmbuild --target sets CC to the cross compiler; tests run on the host.
CGO_ENABLED=1 CC=gcc GOARCH="$native" go test -buildvcs=false ./...
if [ "$target" = "$native" ]; then
  ./kidsh -version
fi

%install
install -D -m 755 kidsh %{buildroot}%{_bindir}/kidsh
install -D -m 644 man/kidsh.1 %{buildroot}%{_mandir}/man1/kidsh.1

%files
%license LICENSE.txt
%doc README.md packaging/kidsh.json.example
%{_bindir}/kidsh
%{_mandir}/man1/kidsh.1*

%changelog
* Sat Aug 22 2026 Jonathan M. Wilbur <jonathan@wilbur.space> - 1.0.4-1
- Make the weather API URL configurable
- Add unit tests
- Fix Alpine packages failing because man pages were uncompressed

* Sat Aug 22 2026 Jonathan M. Wilbur <jonathan@wilbur.space> - 1.0.3-1
- Include the man page in packages

* Sat Aug 22 2026 Jonathan M. Wilbur <jonathan@wilbur.space> - 1.0.2-1
- Build dynamically linked executables

* Sat Aug 22 2026 Jonathan M. Wilbur <jonathan@wilbur.space> - 1.0.1-1
- Fix Alpine package builds
- Fix RPM package builds

* Sat Aug 22 2026 Jonathan M. Wilbur <jonathan@wilbur.space> - 1.0.0-1
- Initial release
- Lightweight secure shell for young children
- Built-in commands for safe exploration, including calendar, colors,
  notes, weather, news, math, stacks, queues, todos, home, birthday,
  age, countdown, counting game, and a Nock interpreter
- Pass configuration into commands
- Package as Debian, RPM, Alpine, Arch, Homebrew, Flatpak, and Snap
- Add Alpine package signing public key
