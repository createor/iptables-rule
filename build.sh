#!/bin/bash
# @Desc：编译并制作deb\rpm安装包
# @Time: 2024/05/13

VERSION=0.0.1  # 版本
DEFAULT_ARCH=x86_64  # 默认架构

if [ $# -ne 2 ]; then
  echo "bash build.sh 系统架构[i386/amd64/arm64] 系统[ubuntu/debian/centos]"
  exit 0
fi

# go build打包文件
export CGO_ENABLED=0
export GOOS=linux

# uname -i 获取cpu架构
case "${1}" in
    "i386")  # 32位
      export AECH=i386
    ;;
    "amd64")
      export AECH=amd64
    ;;
    "arm64")
      export AECH=arm64
    ;;
    *)
      export AECH=${DEFAULT_ARCH}
    ;;
esac

case "${2}" in
    "ubuntu"|"debian")
      # 生成deb安装包
      # 请在Ubuntu/Debain系统上执行
      go build .  # 编译
      sed -i "s/Version:[[:space:]]\S*/Version: ${VERSION}/" dist/rule/DEBIAN/control   # 替换版本字段
      sed -i "s/Architecture:[[:space:]]\S*/Architecture: ${1:-${DEFAULT_ARCH}}/" dist/rule/DEBIAN/control  # 替换架构
      # 拷贝文件
      cp -f rule dist/rule/usr/local/bin/rule
      cp -f startup.sh dist/rule/usr/local/share/rule
      cp -f shutdown.sh dist/rule/usr/local/share/rule
      cp -f rule.service dist/rule/usr/local/share/rule
      cp -rf html dist/rule/usr/local/share/rule/
      # 修改权限
      chmod -R 0755 dist/rule/
      sudo dpkg -b dist/rule rule_${VERSION}-${1:-${DEFAULT_ARCH}}.deb
      # dpkg-deb --build dist/rule rule_${VERSION}_arm64.deb
    ;;
    "centos")
      # 生成rpm安装包
      # 请在Centos系统上执行
      # go build -ldflags="-buildid" -o rule
      go build -ldflags="-s -w" -o rule
      sed -i "s/Version[[:space:]]\S*/Version: ${VERSION}/" dist/rpmbuild/SPECS/rule.spec  # 替换版本字段
      mkdir rule-${VERSION}
      cp -rf rule html startup.sh shutdown.sh rule.service rule-${VERSION}/  # 拷贝文件
      chmod -R 0755 rule-${VERSION}  # 修改权限
      tar zcvf rule-${VERSION}.tar.gz rule-${VERSION}  # 压缩文件
      mv rule-${VERSION}.tar.gz dist/rpmbuild/SOURCES/  # 移动文件
      # rpmbuild只能从/root/目录下操作
      cp -rf dist/rpmbuild ~/
      sudo rpmbuild -ba ~/rpmbuild/SPECS/rule.spec  # rpm生成在/root/rpmbuild/RPMS目录下
    ;;
    *)
      echo "错误的'系统'参数"
      exit 0
    ;;
esac
