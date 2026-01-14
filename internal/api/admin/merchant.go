package admin

import (
	"net/http"
	"wepay-sandbox/internal/core"
	"wepay-sandbox/internal/model"

	"github.com/gin-gonic/gin"
)

// ListMerchants 获取商户列表
func ListMerchants(c *gin.Context) {
	var merchants []model.Merchant
	result := core.DB.Find(&merchants)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, merchants)
}

// CreateMerchant 创建商户
func CreateMerchant(c *gin.Context) {
	var m model.Merchant
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := core.DB.Create(&m); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, m)
}

// UpdateMerchant 更新商户
func UpdateMerchant(c *gin.Context) {
	id := c.Param("id")
	var m model.Merchant
	if result := core.DB.First(&m, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merchant not found"})
		return
	}

	var input model.Merchant
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	core.DB.Model(&m).Updates(input)
	c.JSON(http.StatusOK, m)
}

// DeleteMerchants 批量删除商户 (硬删除)
func DeleteMerchants(c *gin.Context) {
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(ids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No IDs provided"})
		return
	}

	// Unscoped() 用于硬删除
	if result := core.DB.Unscoped().Delete(&model.Merchant{}, ids); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
