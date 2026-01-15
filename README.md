# Pay Sandbox - 在线支付模拟沙箱平台

[![Go Version](https://img.shields.io/badge/go-1.20+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/vue-3.x-4FC08D?style=flat&logo=vue.js)](https://vuejs.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## 1. 项目概述

### 1.1 项目简介
**Pay Sandbox** 是一个专为开发者设计的本地化在线支付模拟环境，旨在解决在开发过程中真实环境测试成本高、官方沙箱环境不稳定或配置复杂等痛点。通过本项目，开发者可以在本地轻松调试支付下单、查询、关闭、退款以及回调通知等完整业务流。

### 1.2 目前支持的平台
- **微信支付 V3**

### 1.3 项目图
<img width="3819" height="1611" alt="1" src="https://github.com/user-attachments/assets/595dd56e-34a0-49ba-9b10-a0230580dd6d" />
<img width="3822" height="1383" alt="2" src="https://github.com/user-attachments/assets/306caa74-204b-4a59-8ccb-cfdc31e3013e" />
<img width="3827" height="1077" alt="3" src="https://github.com/user-attachments/assets/5f4277ff-78ca-4208-9919-54107fd7ddcb" />
<img width="563" height="1095" alt="4" src="https://github.com/user-attachments/assets/4d79d008-58da-422a-972d-f6b52e378362" />

---


## 2. 功能说明

### 微信支付v3

#### 1.1 商户管理
- **配置多商户**: 支持配置多个商户号 (MchID) 和应用 ID (AppID)。
- **自定义回调**: 可独立配置每个商户的支付回调地址 (`notify_url`) 和退款回调地址 (`refund_notify_url`)。
- **配置直观化**: 前端提供 JSON 配置编辑界面，方便设置重试间隔和最大重试次数。

#### 1.2 支付与退款模拟
- **JSAPI 预下单**: 模拟 `/v3/pay/transactions/jsapi` 接口，生成 `prepay_id`。
- **移动端模拟页**: 提供高仿微信支付确认页，支持手动输入 6 位密码触发支付。
- **订单管理**: 支持通过微信支付单号或商户订单号查询订单状态、手动关闭订单。
- **模拟退款**: 支持对已支付订单发起退款，可指定退款金额和原因。

#### 1.3 回调通知系统
- **异步自动重试**: 支付或退款成功后，系统会根据商户配置自动发起 HTTP 回调。
- **幂等性保证**: 内部逻辑确保如果某笔流水已回调成功，则不再重复发送。
- **手动重试**: 对于因业务服务异常导致的回调失败，支持在管理后台点击“重试”按钮手动触发。
- **详细日志**: 完整记录每次回调的 Request Body、Response Body 及 HTTP 状态码，方便排查业务端接口问题。

#### 1.4 监控与通知
- **实时控制台**: 页面右上角通过 SSE 实时弹出新的回调提醒。
- **批量操作**: 支持对商户记录、交易流水和退款记录进行批量删除（硬删除）。

#### 1.5 Mock 接口使用
本项目提供了一套高度兼容微信支付 V3 规范的 Mock 接口，开发者只需将业务代码中的微信支付域名替换为本地沙箱地址（如 `http://localhost:8080`）即可开始调试：
- **JSAPI 下单**: `POST /v3/pay/transactions/jsapi`
- **订单查询**: 
  - 通过微信支付单号: `GET /v3/pay/transactions/id/{transaction_id}`
  - 通过商户订单号: `GET /v3/pay/transactions/out-trade-no/{out_trade_no}`
- **关闭订单**: `POST /v3/pay/transactions/out-trade-no/{out_trade_no}/close`
- **无需签名**: 为了方便本地调试，所有的 Mock 接口均跳过了复杂的微信支付 V3 签名验证。

---

## 3. 使用指南

### 3.1 环境要求
- **Go**: 1.20 或更高版本
- **Node.js**: 16.x 或更高版本 (推荐使用 18+)
- **npm/yarn**: 前端依赖管理

### 3.2 安装与启动

#### 第一步：克隆仓库
```bash
git clone https://github.com/env107/pay-sandbox.git
cd pay-sandbox
```

#### 第二步：启动后端服务
```bash
# 整理依赖
go mod tidy
# 运行服务 (可通过 -port 参数指定端口，默认 8080)
go run cmd/server/main.go -port 8080
```

#### 第三步：启动前端管理后台
```bash
cd web
npm install
npm run dev
```
启动成功后，访问 `http://localhost:3000` 进入管理后台。

---

## 4. 开发说明

### 4.1 目录结构
```text
f:\develop\go\src\wepay-sandbox
├── cmd/server          # 后端入口 (main.go)
├── examples/           # 使用示例 (Go 客户端 & 业务回调接收端)
├── internal/           # 内部核心逻辑
│   ├── api/            # API 处理层 (Admin 管理接口 & Mock 模拟接口)
│   ├── core/           # 核心组件 (数据库初始化等)
│   ├── model/          # 数据模型 (GORM 模型定义)
│   └── worker/         # 异步任务处理 (回调发送逻辑)
├── web/                # 前端 Vue 3 项目
│   ├── src/views/admin # 管理后台视图
│   └── src/views/mobile# 移动端模拟视图
└── README.md           # 项目文档
```

### 4.2 核心代码逻辑参考
- 支付模拟接口实现: [jsapi.go](pay-sandbox/internal/api/mock/jsapi.go)
- 回调任务调度: [callback.go](pay-sandbox/internal/worker/callback.go)
- 前端交易列表: [TransactionList.vue](pay-sandbox/web/src/views/admin/TransactionList.vue)

---

## 5. 其他信息

### 5.1 许可证说明
本项目采用 [MIT License](LICENSE) 开源。

### 5.2 版本更新记录
- **v1.1.0 (2026-01-14)**: 
  - 新增退款回调幂等性校验。
  - 新增退款流水记录批量硬删除功能。
  - 优化管理后台 Material Design UI 风格。
- **v1.0.0 (2026-01-10)**: 
  - 初始版本发布，支持 JSAPI 基本支付流程。

### 5.3 联系方式
如有疑问或建议，请提交 [Issue](https://github.com/env107/pay-sandbox/issues) 或通过邮件联系：`env107@126.com`
