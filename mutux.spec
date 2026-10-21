# Go binaries built outside go-rpm-macros leave debugsource empty, which rpm rejects
%global debug_package %{nil}

Name:           mutux
Version:        1.0.0
Release:        1%{?dist}
Summary:        Terminal music player with album art

License:        MIT
URL:            https://github.com/Jeeshan-Shrestha/mutux
Source0:        %{url}/archive/v%{version}/%{name}-%{version}.tar.gz
# COPR builds have no network access, so the Go modules are vendored:
#   git clone https://github.com/Jeeshan-Shrestha/mutux && cd mutux
#   git checkout v%%{version} && go mod vendor
#   tar czf ../mutux-%%{version}-vendor.tar.gz vendor
Source1:        %{name}-%{version}-vendor.tar.gz

BuildRequires:  golang >= 1.26.1
BuildRequires:  gcc
BuildRequires:  pkgconfig(alsa)

%description
mutux is a minimal terminal music player with a low memory footprint. It
browses ~/Music, plays mp3, flac, wav and ogg files, and shows track details
and album art.

%prep
%autosetup -a1

%build
export CGO_ENABLED=1
export GOFLAGS=-mod=vendor
go build \
    -buildmode=pie \
    -ldflags "-linkmode=external -compressdwarf=false -B 0x$(head -c20 /dev/urandom | od -An -tx1 | tr -d ' \n') -extldflags '%{build_ldflags}'" \
    -o %{name} .

%install
install -Dpm 0755 %{name} %{buildroot}%{_bindir}/%{name}

%files
%license LICENSE
%doc README.md
%{_bindir}/%{name}

%changelog
* Wed Oct 21 2026 Jeeshan Shrestha <shresthajeesan@gmail.com> - 1.0.0-1
- Initial package
