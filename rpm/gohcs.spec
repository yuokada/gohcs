# Original Version and Article, Special-Thanx to :
# http://momijiame.tumblr.com/post/32458077768/spec-rpm
%define         prefix  /
%define debug_package %{nil}

Name:      gohcs
Version:   1.2.0
Release:   %{release}
#Release:   1%{?dist}
Group:     abc
License:   MIT
URL:       https://github.com/yuokada/gohcs
Summary:   health check server
BuildArch: x86_64
Source0:   gohcs.tar.gz
# (only create temporary directory name, for RHEL5 compat environment)
# see : http://fedoraproject.org/wiki/Packaging:Guidelines#BuildRoot_tag
Prefix:    /
BuildRoot: %(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)
Requires: glibc
BuildRequires: golang

#%define INSTALLDIR %{buildroot}/gohcs
%define INSTALLDIR %{buildroot}

%description
gohcs is simple health check server by Go.

%prep
#%setup -q -n %{name}
%setup -q -n gohcs

mkdir -p $RPM_BUILD_ROOT/usr/local/{bin,man/man1}
echo $RPM_BUILD_ROOT
echo %{INSTALLDIR}

%build
make build

%install
rm   -rf      %{INSTALLDIR}
mkdir -p %{buildroot}/var/run/gohcs
mkdir -p %{buildroot}/etc/gohcs
%{__install} -Dp -m0755 src/gohcs                              %{buildroot}/usr/local/bin/%{name}
%{__install} -Dp -m0644 etc/tmpfiles.d/gohcs.conf              %{buildroot}/etc/tmpfiles.d/gohcs.conf
%{__install} -Dp -m0644 etc/systemd/system/gohcs.service       %{buildroot}/usr/lib/systemd/system/gohcs.service
%{__install} -Dp -m0644 etc/systemd/system/gohcs-file@.service %{buildroot}/usr/lib/systemd/system/gohcs-file@.service
%{__install} -Dp -m0644 etc/gohcs/checklist.json               %{buildroot}/etc/gohcs/checklist.json

# Instructions to clean out the build root.
%clean
#rm -rf %{buildroot}
# Avoid Disastarous Damage : http://dev.tapweb.co.jp/2010/12/273
[ "$RPM_BUILD_ROOT" != "/" ] && rm -rf $RPM_BUILD_ROOT

%files
%defattr(0644,root,root)
%config %{prefix}/usr/lib/systemd/system/gohcs.service
%config %{prefix}/usr/lib/systemd/system/gohcs-file@.service
%config %{prefix}/etc/tmpfiles.d/gohcs.conf
%config %{prefix}/etc/gohcs/checklist.json

# ## auto include child files under the directory
# #%{prefix}/etc/

%defattr(0755,root,root)
%{prefix}/usr/local/bin/gohcs

# directory only
%dir %attr(0755,-,-) /var/run/gohcs
%dir %attr(0775,-,-) %{prefix}/etc/gohcs
#%exclude   /Makefile
#%exclude   /src/gohcs
#%exclude   /src/server.go

%pre
if [ "$1" = "1" ]; then
  echo "pre-install script : Initial installation."
elif [ "$1" = "2" ]; then
    echo "pre-install script : Upgrade installation."
fi

%post
if [ "$1" = "1" ]; then
  echo "post-install script : Initial installation."
elif [ "$1" = "2" ]; then
    echo "post-install script : Upgrade installation."
fi

%changelog
* Mon Jan 16 2017 yuokada
- initial release
