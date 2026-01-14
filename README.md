# 微信支付 JSAPI 模拟沙箱平台

## 简介
这是一个本地化的微信支付沙箱环境，用于模拟 JSAPI 支付流程。包含 Go 后端服务和 Vue3 前端管理控制台。

## 功能特性
- **管理控制台**: 配置商户信息，查看交易流水和回调日志。
- **Mock API**: 兼容微信支付 V3 接口规范 (`/v3/pay/transactions/jsapi`)。
- **支付预览**: 仿微信客户端的 H5 支付确认页，支持模拟输入密码支付。
- **自动回调**: 支付成功后自动根据配置策略（如每1分钟重试）向 `notify_url` 发送回调通知。

## 快速开始

### 1. 启动后端
```bash
go mod tidy
go run cmd/server/main.go
```
后端服务运行在 `http://localhost:8080`

### 2. 启动前端
```bash
cd web
npm install
npm run dev
```
前端控制台运行在 `http://localhost:3000`

## 使用流程
1. 打开控制台 `http://localhost:3000/admin/merchants`，添加一个测试商户（MchID, AppID 等）。
2. 配置您的业务后端，将微信支付 API 域名指向 `http://localhost:8080`。
3. 发起下单请求，获取 `prepay_id`。
4. 拼接预览链接 `http://localhost:3000/pay/preview/{prepay_id}` 并在浏览器打开。
5. 输入任意6位密码完成支付。
6. 查看您的业务后端是否收到回调通知。

## 目录结构
- `cmd/server`: 后端启动入口
- `internal/api`: API 接口实现 (Admin & Mock)
- `internal/core`: 数据库与核心逻辑
- `internal/worker`: 回调任务调度
- `web`: 前端 Vue3 项目
