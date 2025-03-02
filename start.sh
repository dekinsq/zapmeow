#!/bin/bash

# 自动检测操作系统
if [[ "$(uname -s)" == "Darwin" ]]; then
    APP_NAME="app-mac"    # macOS 专用版本
    PLATFORM="macOS"
else
    APP_NAME="app-linux"  # 其他系统默认使用 Linux 版本
    PLATFORM="Linux"
fi

# 程序路径（假设脚本与程序在同一目录）
APP_PATH="./$APP_NAME"
# 日志文件路径
LOG_FILE="./$APP_NAME.log"

# 启动前环境检测
check_environment() {
    # 检查可执行文件是否存在
    if [ ! -f "$APP_PATH" ]; then
        echo "错误：未找到 [$PLATFORM] 可执行文件 $APP_NAME"
        echo "请先编译对应平台的程序（macOS需编译为app-mac）"
        exit 1
    fi

    # 检查执行权限
    if [ ! -x "$APP_PATH" ]; then
        echo "正在添加可执行权限..."
        chmod +x "$APP_PATH"
    fi
}

# 启动应用
start() {
    check_environment  # 先执行环境检测

    if pgrep -x "$APP_NAME" > /dev/null
    then
        echo "应用已在运行中"
    else
        echo "正在启动 $APP_NAME ($PLATFORM)..."
        nohup $APP_PATH > $LOG_FILE 2>&1 &
        echo "应用已启动，日志输出至: $LOG_FILE"
        echo "可以使用以下命令查看日志：tail -f $LOG_FILE"
    fi
}

# 停止应用
stop() {
    echo "正在停止 $APP_NAME..."
    pkill -x "$APP_NAME"
    sleep 1  # 等待进程结束
    if pgrep -x "$APP_NAME" > /dev/null
    then
        echo "停止失败，尝试强制终止..."
        pkill -9 -x "$APP_NAME"
        rm -f $LOG_FILE  # 强制终止后清理日志
    else
        echo "应用已停止"
    fi
}

# 重启应用
restart() {
    stop
    sleep 2
    start
}

# 使用说明
usage() {
    echo "用法: $0 {start|stop|restart}"
    echo "检测到当前系统为: $PLATFORM"
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    *)
        usage
        exit 1
esac

exit 0
