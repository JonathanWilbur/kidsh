Name:           kidsh
Version:        1.0.0
Release:        1%{?dist}
Summary:        A lightweight secure shell for young children

License:        MIT
URL:            https://github.com/yourusername/kidsh
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang >= 1.20
Requires:       glibc
Requires:       systemd

%description
Kid Shell is a lightweight, low-capability secure shell designed for young children
to safely explore and interact with a computer. It provides simple commands for
basic tasks that children can understand, such as displaying colors, showing days
of the week, and taking notes. The shell compiles as a static binary and supports
only built-in commands to ensure security.

%prep
%autosetup

%build
go build -o kidsh src/main.go src/cmd.go src/ansi.go src/config.go

%install
mkdir -p %{buildroot}%{_bindir}
install -m 755 kidsh %{buildroot}%{_bindir}/kidsh

%files
%license LICENSE.txt
%doc README.md
%{_bindir}/kidsh

%changelog
* %(date '+%a %b %d %Y') %{packager} - %{version}-%{release}
- Initial RPM package for kidsh 