#!/bin/bash
# @Desc: 停止脚本
# @Time: 2024/05/17

BASE_DIR=$(cd $(dirname $0);pwd)
PID_FILE=${BASE_DIR}/rule.pid
kill -9 `cat ${PID_FILE}`
exit 0