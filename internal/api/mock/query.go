package mock

import (
	"net/http"
	"time"
	"wepay-sandbox/internal/core"
	"wepay-sandbox/internal/model"

	"github.com/gin-gonic/gin"
)

// QueryByTransactionID 微信支付订单号查询
func QueryByTransactionID(c *gin.Context) {
	transactionID := c.Param("transaction_id")
	mchid := c.Query("mchid")

	if mchid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "PARAM_ERROR",
			"message": "mchid is required",
		})
		return
	}

	var tx model.Transaction
	if result := core.DB.Where("transaction_id = ? AND mch_id = ?", transactionID, mchid).First(&tx); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "RESOURCE_NOT_EXISTS",
			"message": "Transaction not found",
		})
		return
	}

	c.JSON(http.StatusOK, buildTransactionResponse(tx))
}

// QueryByOutTradeNo 商户订单号查询
func QueryByOutTradeNo(c *gin.Context) {
	outTradeNo := c.Param("out_trade_no")
	mchid := c.Query("mchid")

	if mchid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "PARAM_ERROR",
			"message": "mchid is required",
		})
		return
	}

	var tx model.Transaction
	if result := core.DB.Where("out_trade_no = ? AND mch_id = ?", outTradeNo, mchid).First(&tx); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "RESOURCE_NOT_EXISTS",
			"message": "Transaction not found",
		})
		return
	}

	c.JSON(http.StatusOK, buildTransactionResponse(tx))
}

// buildTransactionResponse 构建标准响应结构
func buildTransactionResponse(tx model.Transaction) map[string]interface{} {
	resp := map[string]interface{}{
		"appid":            tx.AppID,
		"mchid":            tx.MchID,
		"out_trade_no":     tx.OutTradeNo,
		"transaction_id":   tx.TransactionID,
		"trade_type":       "JSAPI",
		"trade_state":      tx.Status,
		"trade_state_desc": "支付成功", // 简化描述
		"bank_type":        "OTHERS",
		"attach":           "",
		"payer": map[string]interface{}{
			"openid": "mock_openid_123", // 模拟 OpenID
		},
		"amount": map[string]interface{}{
			"total":          tx.Amount,
			"payer_total":    tx.Amount,
			"currency":       tx.Currency,
			"payer_currency": tx.Currency,
		},
	}

	if tx.Status == "SUCCESS" {
		resp["success_time"] = tx.UpdatedAt.Format(time.RFC3339)
	}

	return resp
}
