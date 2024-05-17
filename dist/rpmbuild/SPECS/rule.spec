Name: rule
Version: 0.0.1
Release: 1%{?dist}
Summary: Rule RPM package
Group: Applications/System

License: GPL
URL: https://github.com/createor/rule.git
Source0: rule-0.0.1.tar.gz
BuildRoot: %{_tmppath}/${name}-%{version}-root

# 描述信息
%description
add or delete and manage iptables rule by html

# 关闭debug
%define debug_package %{nil}

# 清空文件并解压源码包
%prep
rm -rf $RPM_BUILD_ROOT/*
%setup -q

%build

%pre

# 安装
%install
mkdir -p $RPM_BUILD_ROOT/usr/local/bin/
mkdir -p $RPM_BUILD_ROOT/usr/local/share/rule
cp -a %{_builddir}/%{name}-%{version}/rule $RPM_BUILD_ROOT/usr/local/bin/
cp -a %{_builddir}/%{name}-%{version}/html $RPM_BUILD_ROOT/usr/local/share/rule/
cp -a %{_builddir}/%{name}-%{version}/startup.sh $RPM_BUILD_ROOT/usr/local/share/rule/
cp -a %{_builddir}/%{name}-%{version}/shutdown.sh $RPM_BUILD_ROOT/usr/local/share/rule/
cp -a %{_builddir}/%{name}-%{version}/rule.service $RPM_BUILD_ROOT/usr/local/share/rule/

# 安装后执行的命令
%post
if [ -d "/usr/lib/systemd/system" ]; then
  cp -f /usr/local/share/rule/rule.service /usr/lib/systemd/system/
else
  cp -f /usr/local/share/rule/rule.service /lib/systemd/system/
fi
systemctl daemon-reload

# 卸载后执行的命令
%postun
rm -rf /usr/local/share/rule
if [ -f "/usr/lib/systemd/system/rule.service" ]; then
  rm -f /usr/lib/systemd/system/rule.service
else
  rm -f /lib/systemd/system/rule.service
fi

# 文件属性
%files
%defattr(-,root,root,0755)
%attr(0644,root,root)
/usr/local/bin/rule
/usr/local/share/rule

%changelog
