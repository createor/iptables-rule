#!/bin/bash
# @Desc: 启动脚本
# @Time: 2024/05/17

BASE_DIR=$(cd $(dirname $0);pwd)
PID_FILE=${BASE_DIR}/rule.pid
nohup rule > /dev/null 2>&1 &
echo "$!" > ${PID_FILE}
exit 0