Name:           kidsh
Version:        1.0.0
Release:        1
Summary:        A lightweight secure shell for young children

License:        MIT
URL:            https://github.com/JonathanWilbur/kidsh
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang >= 1.22

%global debug_package %{nil}
%global __strip /bin/true

%description
Kid Shell is a lightweight, low-capability secure shell designed for young
children to safely explore and interact with a computer. It provides simple
commands for basic tasks that children can understand, such as displaying
colors, showing days of the week, and taking notes. The shell compiles as a
static binary and supports only built-in commands to ensure security.

%prep
%autosetup

%build
export CGO_ENABLED=0
export GOOS=linux
%if 0%{?go_arch:1}
export GOARCH=%{go_arch}
%endif
go build -trimpath -buildvcs=false -ldflags "-s -w -X main.version=%{version}" -o kidsh ./src

%check
export CGO_ENABLED=0
%if 0%{?go_arch:1}
export GOARCH=%{go_arch}
%endif
native="$(go env GOHOSTARCH)"
target="$(go env GOARCH)"
GOARCH="$native" go test -buildvcs=false ./...
if [ "$target" = "$native" ]; then
  ./kidsh -version
fi

%install
install -D -m 755 kidsh %{buildroot}%{_bindir}/kidsh

%files
%license LICENSE.txt
%doc README.md packaging/kidsh.json.example
%{_bindir}/kidsh

%changelog
* Sat Aug 22 2026 Jonathan M. Wilbur <jonathan@wilbur.space> - 1.0.0-1
- Initial packaging
