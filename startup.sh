#!/bin/bash
# @Desc: 启动脚本
# @Time: 2024/05/17

DEFAULT_PORT=8089  # 服务监听端口
BASE_DIR=$(cd $(dirname $0);pwd)
PID_FILE=${BASE_DIR}/rule.pid
nohup rule -p ${DEFAULT_PORT} > /dev/null 2>&1 &
echo "$!" > ${PID_FILE}
exit 0
