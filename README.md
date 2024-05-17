# rule 一款iptables的管理工具
注意： 只支持对INPUT的管理

## 文件说明

```bash
|--dist (编译安装包目录)
|--html (前端文件夹)
|--build.sh (编译脚本)
|--main.go (主函数文件)
|--rule.service (systemd服务文件)
|--shutdown.sh (停止服务脚本)
|--startup.sh (启动服务脚本)
```

## 介绍

服务默认使用8089端口, 如果需要修改可以在startup.sh中DEFAULT_PORT参数

访问 http://ip:8089

首页

![001.png](https://github.com/createor/iptables-rule/blob/master/img/001.png)


## Debian/Ubuntu下使用

```bash
使用命令查看服务器架构: uname -i
打包安装包: bash build 架构[arm64/amd64] ubuntu
安装: sudo dpkg -i rule_xxx_xx.deb
卸载: sudo dpkg -r rule
```

## Centos下使用

```bash
# yum install -y rpm-build
使用命令查看服务器架构: uname -i
打包安装包: bash build 架构[arm64/amd64] centos
安装: sudo rpm -ivh rule_xxx_xx.rpm
卸载: sudo rpm -e rule
```

## 使用

```bash
启动：systemctl start rule.service
停止：systemctl stop rule.service   【注意:每次修改完规则后,应关闭服务避免造成不必要的安全问题】
状态：systemctl status rule.service
```

## 开机自动加载防火墙规则

首先在页面上导出规则文件rule.json

在/etc/rc.local文件中添加如下命令

```bash
/usr/local/bin/rule -f /path/to/rule.json
```

