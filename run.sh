#!/bin/bash





# 定义清理函数
cleanup() {
    echo -e "\n正在停止服务..."
    if [ -n "$GO_PID" ]; then
        kill $GO_PID 2>/dev/null
        echo "Go后端服务已停止"
    fi
    if [ -n "$NOTIFY_PID" ]; then
        kill $NOTIFY_PID 2>/dev/null
        echo "Notify模拟回调服务已停止"
    fi
    exit
}

# 捕获 SIGINT (Ctrl+C) 和 SIGTERM 信号
trap cleanup SIGINT SIGTERM

read -p "请输入后端服务端口号(默认8080): " SERVER_PORT
SERVER_PORT=${SERVER_PORT:-8080}

# 启动 Go 后端服务
echo "正在启动 Go 后端服务..."
go run cmd/server/main.go --port $SERVER_PORT &
GO_PID=$!

# 检查 Go 服务是否成功启动
sleep 2
if ps -p $GO_PID > /dev/null; then
   echo "Go 后端服务启动成功 (PID: $GO_PID)"
else
   echo "Go 后端服务启动失败"
   exit 1
fi

read -p "是否需要启动模拟回调端口(y/n)?" START_NOTIFY
echo "" # 换行

# 启动 Notify 服务 (如果用户选择 y)
if [[ "$START_NOTIFY" == "y" || "$START_NOTIFY" == "Y" ]]; then
    read -p "请输入模拟回调端口号(默认8081): " NOTIFY_PORT
    NOTIFY_PORT=${NOTIFY_PORT:-8081}
    echo "正在启动 Notify 模拟回调服务..."
    go run examples/notify/notify.go --port $NOTIFY_PORT &
    NOTIFY_PID=$!
    
    # 简单的启动检查
    sleep 1
    if ps -p $NOTIFY_PID > /dev/null; then
        echo "Notify 服务启动成功 (PID: $NOTIFY_PID)"
    else
        echo "Notify 服务启动失败"
        # 不强制退出，继续启动主服务
    fi
fi

# 启动前端服务
echo "正在启动前端服务..."
cd web
npm run dev

# 脚本结束时执行清理
cleanup
