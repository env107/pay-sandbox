package admin

import (
	"net/http"
	"wepay-sandbox/internal/core"
	"wepay-sandbox/internal/model"
	"wepay-sandbox/internal/worker"

	"github.com/gin-gonic/gin"
)

// ListTransactions 获取交易列表
func ListTransactions(c *gin.Context) {
	var transactions []model.Transaction
	query := core.DB.Order("created_at desc")

	// 1. 商户号筛选
	if mchid := c.Query("mchid"); mchid != "" {
		query = query.Where("mch_id = ?", mchid)
	}

	// 2. Prepay ID 精确查询
	if prepayID := c.Query("prepay_id"); prepayID != "" {
		query = query.Where("prepay_id = ?", prepayID)
	}

	// 3. 商户订单号模糊查询
	if outTradeNo := c.Query("out_trade_no"); outTradeNo != "" {
		query = query.Where("out_trade_no LIKE ?", "%"+outTradeNo+"%")
	}

	// 4. 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 5. 时间范围筛选
	if startTime := c.Query("start_time"); startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime := c.Query("end_time"); endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	result := query.Find(&transactions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

// DeleteTransactions 批量删除交易记录 (硬删除)
func DeleteTransactions(c *gin.Context) {
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(ids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No IDs provided"})
		return
	}

	// 物理删除交易记录
	if result := core.DB.Unscoped().Delete(&model.Transaction{}, ids); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// 同时删除关联的回调日志 (可选，保持数据整洁)
	// 这里假设 transaction_id 是外键关联的逻辑键，需要先查出 transaction_id 字符串
	// 为简单起见，这里仅删除 Transaction 表记录

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

// SimulatePay 模拟支付成功（手动触发）
func SimulatePay(c *gin.Context) {
	var input struct {
		PrepayID string `json:"prepay_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tx model.Transaction
	if result := core.DB.Where("prepay_id = ?", input.PrepayID).First(&tx); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	// 更新状态
	if tx.Status != "SUCCESS" {
		tx.Status = "SUCCESS"
		core.DB.Save(&tx)
		// 触发回调任务
		worker.TriggerCallback(tx)
	}

	c.JSON(http.StatusOK, tx)
}

// GetTransactionLogs 获取交易回调日志
func GetTransactionLogs(c *gin.Context) {
	transactionID := c.Param("transaction_id")
	var logs []model.CallbackLog
	if result := core.DB.Where("transaction_id = ?", transactionID).Order("created_at desc").Find(&logs); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// RetryTransactionCallback 重试交易回调
func RetryTransactionCallback(c *gin.Context) {
	id := c.Param("transaction_id")
	var tx model.Transaction
	if result := core.DB.First(&tx, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	// 触发回调
	worker.TriggerCallback(tx)

	c.JSON(http.StatusOK, gin.H{"message": "Retry task submitted"})
}
