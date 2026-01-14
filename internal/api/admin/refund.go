package admin

import (
	"fmt"
	"net/http"
	"time"
	"wepay-sandbox/internal/core"
	"wepay-sandbox/internal/model"
	"wepay-sandbox/internal/worker"

	"github.com/gin-gonic/gin"
)

// SimulateRefund 模拟退款
func SimulateRefund(c *gin.Context) {
	var input struct {
		TransactionID string `json:"transaction_id" binding:"required"`
		Amount        int64  `json:"amount" binding:"required"`
		Reason        string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. 查找原订单
	var tx model.Transaction
	if result := core.DB.Where("transaction_id = ?", input.TransactionID).First(&tx); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if tx.Status != "SUCCESS" && tx.Status != "REFUND" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order status is not SUCCESS or REFUND"})
		return
	}

	// 2. 校验金额 (简化逻辑：不严格校验累计退款是否超额)
	if input.Amount > tx.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refund amount exceeds total amount"})
		return
	}

	// 3. 获取商户退款配置
	var mch model.Merchant
	var notifyUrl string
	if core.DB.Where("mch_id = ?", tx.MchID).First(&mch).Error == nil {
		notifyUrl = mch.RefundNotifyUrl
		if notifyUrl == "" {
			notifyUrl = mch.NotifyUrl // Fallback
		}
	}

	// 4. 创建退款记录
	refundID := fmt.Sprintf("REF_%d", time.Now().UnixNano())
	outRefundNo := fmt.Sprintf("REF_OUT_%d", time.Now().UnixNano())

	refund := model.Refund{
		RefundID:      refundID,
		OutRefundNo:   outRefundNo,
		TransactionID: tx.TransactionID,
		MchID:         tx.MchID,
		Amount:        input.Amount,
		Total:         tx.Amount,
		Currency:      tx.Currency,
		Reason:        input.Reason,
		Status:        "SUCCESS", // 模拟直接成功
		NotifyUrl:     notifyUrl,
	}

	if err := core.DB.Create(&refund).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 5. 更新原订单状态 (标记为 REFUND)
	if tx.Status != "REFUND" {
		tx.Status = "REFUND"
		core.DB.Save(&tx)
	}

	// 6. 触发退款回调
	worker.TriggerRefundCallback(refund)

	c.JSON(http.StatusOK, refund)
}

// ListRefunds 获取退款流水
func ListRefunds(c *gin.Context) {
	var refunds []model.Refund
	query := core.DB.Order("created_at desc")

	if mchid := c.Query("mchid"); mchid != "" {
		query = query.Where("mch_id = ?", mchid)
	}
	if transactionID := c.Query("transaction_id"); transactionID != "" {
		query = query.Where("transaction_id = ?", transactionID)
	}
	if outRefundNo := c.Query("out_refund_no"); outRefundNo != "" {
		query = query.Where("out_refund_no LIKE ?", "%"+outRefundNo+"%")
	}

	if result := query.Find(&refunds); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, refunds)
}

// GetRefundLogs 获取退款回调日志
func GetRefundLogs(c *gin.Context) {
	refundID := c.Param("refund_id")
	var logs []model.CallbackLog
	// 注意：我们在 TriggerRefundCallback 中将 RefundID 存入了 TransactionID 字段
	if result := core.DB.Where("transaction_id = ?", refundID).Order("created_at desc").Find(&logs); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// RetryRefundCallback 重试退款回调
func RetryRefundCallback(c *gin.Context) {
	id := c.Param("refund_id")
	var refund model.Refund
	if result := core.DB.Where("refund_id = ?", id).First(&refund); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Refund not found"})
		return
	}

	// 触发回调
	worker.TriggerRefundCallback(refund)

	c.JSON(http.StatusOK, gin.H{"message": "Retry task submitted"})
}

// DeleteRefunds 批量删除退款记录 (硬删除)
func DeleteRefunds(c *gin.Context) {
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(ids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No IDs provided"})
		return
	}

	// 物理删除退款记录
	if result := core.DB.Unscoped().Delete(&model.Refund{}, ids); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
