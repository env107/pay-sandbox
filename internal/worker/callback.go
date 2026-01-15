package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"wepay-sandbox/internal/api"
	"wepay-sandbox/internal/core"
	"wepay-sandbox/internal/model"
)

var (
	// callbackLocks 支付回调并发锁，key 为 TransactionID
	callbackLocks sync.Map
)

// NotifyConfig 回调配置结构
type NotifyConfig struct {
	Interval   string `json:"interval"`    // e.g. "5s", "1m"
	MaxRetries int    `json:"max_retries"` // e.g. 3
}

// TriggerCallback 触发回调
func TriggerCallback(tx model.Transaction) {
	go func() {
		// 1. 检查是否已经回调成功，避免重复发送
		var currentTx model.Transaction
		if err := core.DB.First(&currentTx, tx.ID).Error; err == nil {
			if currentTx.CallbackStatus == "SUCCESS" {
				fmt.Printf("Transaction %s callback already SUCCESS, skip.\n", tx.TransactionID)
				return
			}
			// 使用最新数据
			tx = currentTx
		}

		// 默认策略
		maxRetries := 3
		retryInterval := 5 * time.Second

		// 查询商户配置
		var mch model.Merchant
		if err := core.DB.Where("mch_id = ?", tx.MchID).First(&mch).Error; err == nil {
			var config NotifyConfig
			if json.Unmarshal([]byte(mch.NotifyConfig), &config) == nil {
				if config.MaxRetries > 0 {
					maxRetries = config.MaxRetries
				}
				if d, err := time.ParseDuration(config.Interval); err == nil {
					retryInterval = d
				}
			}
		}

		// 获取已尝试次数 (从 CallbackLog 中查询)
		var existingLogsCount int64
		core.DB.Model(&model.CallbackLog{}).Where("transaction_id = ?", tx.TransactionID).Count(&existingLogsCount)

		remainingRetries := maxRetries - int(existingLogsCount)
		if remainingRetries <= 0 {
			fmt.Printf("Transaction %s already reached max retries (%d), skip.\n", tx.TransactionID, maxRetries)
			return
		}

		payload := map[string]interface{}{
			"id":             tx.TransactionID, // 通知ID
			"create_time":    time.Now().Format(time.RFC3339),
			"resource_type":  "encrypt-resource",
			"event_type":     "TRANSACTION.SUCCESS",
			"summary":        "支付成功",
			"original_type":  "transaction",
			"transaction_id": tx.TransactionID,
			"out_trade_no":   tx.OutTradeNo,
			"appid":          tx.AppID,
			"mchid":          tx.MchID,
			"trade_state":    "SUCCESS",
			"amount": map[string]interface{}{
				"total":    tx.Amount,
				"currency": tx.Currency,
			},
			"success_time": time.Now().Format(time.RFC3339),
			// 注意：这里省略了真实的加密逻辑，直接返回明文以便调试，或者模拟加密结构
			"resource": map[string]interface{}{
				"original_type":   "transaction",
				"algorithm":       "AEAD_AES_256_GCM",
				"ciphertext":      "mock_ciphertext", // 真实场景需加密
				"associated_data": "",
				"nonce":           "",
			},
		}

		jsonBody, _ := json.Marshal(payload)

		for i := 0; i < remainingRetries; i++ {
			// 如果不是第一次尝试，先等待
			if i > 0 {
				time.Sleep(retryInterval)
			}

			// 在实际发起 HTTP 请求前加锁
			if _, loaded := callbackLocks.LoadOrStore(tx.TransactionID, true); loaded {
				fmt.Printf("Transaction %s individual callback attempt is already in progress, skip this loop.\n", tx.TransactionID)
				continue // 如果当前有请求正在发，跳过本次循环进入下一次重试等待
			}

			resp, err := http.Post(tx.NotifyUrl, "application/json", bytes.NewBuffer(jsonBody))

			status := "FAIL"
			statusCode := 0
			respBody := ""

			if err == nil {
				statusCode = resp.StatusCode
				if statusCode >= 200 && statusCode < 300 {
					status = "SUCCESS"
				}
				resp.Body.Close()
			} else {
				respBody = err.Error()
			}

			// 请求结束，释放锁
			callbackLocks.Delete(tx.TransactionID)

			// 记录日志
			log := model.CallbackLog{
				TransactionID: tx.TransactionID,
				NotifyUrl:     tx.NotifyUrl,
				RequestBody:   string(jsonBody),
				ResponseBody:  respBody,
				StatusCode:    statusCode,
				Status:        status,
				RetryCount:    int(existingLogsCount) + i + 1, // 总计第几次尝试
			}
			core.DB.Create(&log)

			// 广播事件
			api.GlobalEventChan <- api.Event{
				Type: "callback",
				Payload: map[string]interface{}{
					"transaction_id": tx.TransactionID,
					"out_trade_no":   tx.OutTradeNo,
					"status":         status,
					"message":        fmt.Sprintf("新的回调产生，商户订单号：%s", tx.OutTradeNo),
				},
			}

			// 更新订单的回调状态
			core.DB.Model(&tx).Updates(map[string]interface{}{
				"callback_status": status,
				"callback_msg":    respBody,
			})

			if status == "SUCCESS" {
				break
			}
		}
	}()
}
