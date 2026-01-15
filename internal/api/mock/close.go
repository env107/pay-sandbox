package mock

import (
	"net/http"
	"wepay-sandbox/internal/core"
	"wepay-sandbox/internal/model"

	"github.com/gin-gonic/gin"
)

// CloseOrder 关闭订单
func CloseOrder(c *gin.Context) {
	outTradeNo := c.Param("out_trade_no")

	// 读取 Body 中的 mchid (虽然文档要求，但为了简化 Mock，我们可以只从 DB 查)
	var req struct {
		MchID string `json:"mchid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "PARAM_ERROR",
			"message": "Invalid request body",
		})
		return
	}

	var tx model.Transaction
	if result := core.DB.Where("out_trade_no = ? AND mch_id = ?", outTradeNo, req.MchID).First(&tx); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "ORDER_NOT_EXIST",
			"message": "Order not found",
		})
		return
	}

	if tx.Status == "SUCCESS" {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    "ORDERPAID",
			"message": "Order paid",
		})
		return
	}

	if tx.Status == "CLOSED" {
		c.Status(http.StatusNoContent)
		return
	}

	// 更新状态为 CLOSED
	tx.Status = "CLOSED"
	core.DB.Save(&tx)

	c.Status(http.StatusNoContent)
}
