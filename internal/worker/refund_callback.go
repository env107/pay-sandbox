package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"wepay-sandbox/internal/api"
	"wepay-sandbox/internal/core"
	"wepay-sandbox/internal/model"
)

// TriggerRefundCallback 触发退款回调
func TriggerRefundCallback(refund model.Refund) {
	go func() {
		// 1. 检查是否已经回调成功
		var currentRefund model.Refund
		if err := core.DB.First(&currentRefund, refund.ID).Error; err == nil {
			if currentRefund.CallbackStatus == "SUCCESS" {
				fmt.Printf("Refund %s callback already SUCCESS, skip.\n", refund.RefundID)
				return
			}
			// 使用最新数据
			refund = currentRefund
		}

		// 默认策略
		maxRetries := 3
		retryInterval := 5 * time.Second

		// 查询商户配置
		var mch model.Merchant
		if err := core.DB.Where("mch_id = ?", refund.MchID).First(&mch).Error; err == nil {
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

		payload := map[string]interface{}{
			"id":             refund.RefundID, // 通知ID
			"create_time":    time.Now().Format(time.RFC3339),
			"resource_type":  "encrypt-resource",
			"event_type":     "REFUND.SUCCESS",
			"summary":        "退款成功",
			"original_type":  "refund",
			"refund_id":      refund.RefundID,
			"out_refund_no":  refund.OutRefundNo,
			"transaction_id": refund.TransactionID,
			"mchid":          refund.MchID,
			"refund_status":  "SUCCESS",
			"amount": map[string]interface{}{
				"refund":   refund.Amount,
				"total":    refund.Total,
				"currency": refund.Currency,
			},
			"success_time": time.Now().Format(time.RFC3339),
			"resource": map[string]interface{}{
				"original_type":   "refund",
				"algorithm":       "AEAD_AES_256_GCM",
				"ciphertext":      "mock_refund_ciphertext",
				"associated_data": "",
				"nonce":           "",
			},
		}

		jsonBody, _ := json.Marshal(payload)

		for i := 0; i < maxRetries; i++ {
			if i > 0 {
				time.Sleep(retryInterval)
			}

			// 退款回调地址优先使用 Merchant 配置的 RefundNotifyUrl，如果没有则使用 NotifyUrl (或退款接口传入的)
			// 此处假设使用 Refund.NotifyUrl (在创建 Refund 时已从 Merchant 获取并存入)
			notifyUrl := refund.NotifyUrl
			if notifyUrl == "" {
				notifyUrl = mch.NotifyUrl // Fallback
			}

			resp, err := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonBody))

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

			// 记录日志 (复用 CallbackLog, TransactionID 存 RefundID 方便查询)
			log := model.CallbackLog{
				TransactionID: refund.RefundID, // 注意：这里存的是退款单号，以便在退款流水中查询
				NotifyUrl:     notifyUrl,
				RequestBody:   string(jsonBody),
				ResponseBody:  respBody,
				StatusCode:    statusCode,
				Status:        status,
				RetryCount:    i + 1,
			}
			core.DB.Create(&log)

			// 广播事件
			api.GlobalEventChan <- api.Event{
				Type: "callback",
				Payload: map[string]interface{}{
					"transaction_id": refund.RefundID,
					"out_trade_no":   refund.OutRefundNo,
					"status":         status,
					"message":        fmt.Sprintf("新的退款回调产生，商户退款单号：%s", refund.OutRefundNo),
				},
			}

			// 更新退款单的回调状态
			core.DB.Model(&refund).Updates(map[string]interface{}{
				"callback_status": status,
				"callback_msg":    respBody,
			})

			if status == "SUCCESS" {
				break
			}
		}
	}()
}
