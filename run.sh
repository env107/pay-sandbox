#!/bin/bash

# 清理函数
cleanup() {
    echo -e "\nStopping services..."
    if [ -n "$GO_PID" ]; then
        kill $GO_PID 2>/dev/null
        echo "Go backend service stopped"
    fi
    if [ -n "$NOTIFY_PID" ]; then
        kill $NOTIFY_PID 2>/dev/null
        echo "Notify mock callback service stopped"
    fi
    exit
}

# 捕获 SIGINT (Ctrl+C) 和 SIGTERM 信号
trap cleanup SIGINT SIGTERM

read -p "Please enter backend server port (default 8080): " SERVER_PORT
SERVER_PORT=${SERVER_PORT:-8080}

# 启动 Go 后端服务
echo "Starting Go backend service..."
go run cmd/server/main.go --port $SERVER_PORT &
GO_PID=$!

# 检查 Go 服务是否成功启动
sleep 2
if ps -p $GO_PID > /dev/null; then
   echo "Go backend service started successfully (PID: $GO_PID)"
else
   echo "Go backend service failed to start"
   exit 1
fi

read -p "Start mock callback service? (y/n) " START_NOTIFY
echo "" 

# 启动 Notify 服务 (如果用户选择 y)
if [[ "$START_NOTIFY" == "y" || "$START_NOTIFY" == "Y" ]]; then
    read -p "Please enter mock callback service port (default 8081): " NOTIFY_PORT
    NOTIFY_PORT=${NOTIFY_PORT:-8081}
    echo "Starting Notify mock callback service..."
    go run examples/notify/notify.go --port $NOTIFY_PORT &
    NOTIFY_PID=$!
    
    # 简单的启动检查
    sleep 1
    if ps -p $NOTIFY_PID > /dev/null; then
        echo "Notify service started successfully (PID: $NOTIFY_PID)"
    else
        echo "Notify service failed to start"
        # 不强制退出，继续启动主服务
    fi
fi

# 启动前端服务（将后端端口传给前端，前端代理会使用该端口）
echo "Starting frontend service (backend port: $SERVER_PORT)..."
cd web
export BACKEND_PORT=$SERVER_PORT
npm run dev

# 脚本结束时执行清理
cleanup
