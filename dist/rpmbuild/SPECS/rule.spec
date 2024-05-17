Name: rule
Version: 0.0.1
Release: 1%{?dist}
Summary: Rule RPM package
Group: Applications/System

License: GPL
URL: https://github.com/createor/rule.git
Source0: rule-0.0.1.tar.gz
BuildRoot: %{_tmppath}/${name}-%{version}-root

%description
add or delete and manage iptables rule by html

%define debug_package %{nil}

%prep
rm -rf $RPM_BUILD_ROOT/*
%setup -q

%build

%pre

%install
mkdir -p $RPM_BUILD_ROOT/usr/local/bin/
mkdir -p $RPM_BUILD_ROOT/usr/local/share/rule
cp -a %{_builddir}/%{name}-%{version}/rule $RPM_BUILD_ROOT/usr/local/bin/
cp -a %{_builddir}/%{name}-%{version}/html $RPM_BUILD_ROOT/usr/local/share/rule/
cp -a %{_builddir}/%{name}-%{version}/startup.sh $RPM_BUILD_ROOT/usr/local/share/rule/
cp -a %{_builddir}/%{name}-%{version}/shutdown.sh $RPM_BUILD_ROOT/usr/local/share/rule/
cp -a %{_builddir}/%{name}-%{version}/rule.service $RPM_BUILD_ROOT/usr/local/share/rule/

%post
if [ -d "/usr/lib/systemd/system" ]; then
  cp -f /usr/local/share/rule/rule.service /usr/lib/systemd/system/
else
  cp -f /usr/local/share/rule/rule.service /lib/systemd/system/
fi
systemctl daemon-reload

%postun
rm -rf /usr/local/share/rule

%files
%defattr(-,root,root,0755)
%attr(0644,root,root)
/usr/local/bin/rule
/usr/local/share/rule

%changelog