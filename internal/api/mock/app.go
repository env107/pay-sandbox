package mock

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"wepay-sandbox/internal/core"
	"wepay-sandbox/internal/model"

	"github.com/gin-gonic/gin"
)

// AppPrepayRequest APP下单请求参数
type AppPrepayRequest struct {
	AppID       string `json:"appid"`
	Mchid       string `json:"mchid"`
	Description string `json:"description"`
	OutTradeNo  string `json:"out_trade_no"`
	NotifyUrl   string `json:"notify_url"`
	Amount      struct {
		Total    int64  `json:"total"`
		Currency string `json:"currency"`
	} `json:"amount"`
}

// AppPrepay APP 下单接口
func AppPrepay(c *gin.Context) {
	var req AppPrepayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "PARAM_ERROR", "message": err.Error()})
		return
	}

	// 校验商户是否存在
	var mch model.Merchant
	if result := core.DB.Where("mch_id = ?", req.Mchid).First(&mch); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "MCH_NOT_FOUND", "message": "Merchant not configured in sandbox"})
		return
	}

	// 生成 Mock PrepayID
	prepayID := fmt.Sprintf("wx%s%06d", time.Now().Format("20060102150405"), rand.Intn(100000))
	transactionID := fmt.Sprintf("420000%s%06d", time.Now().Format("20060102150405"), rand.Intn(100000))

	// 保存交易记录
	tx := model.Transaction{
		AppID:         req.AppID,
		MchID:         req.Mchid,
		Description:   req.Description,
		OutTradeNo:    req.OutTradeNo,
		TransactionID: transactionID,
		PrepayID:      prepayID,
		Amount:        req.Amount.Total,
		Currency:      req.Amount.Currency,
		Status:        "CREATED",
		NotifyUrl:     req.NotifyUrl,
		TradeType:     "WX:APP",
	}

	if err := core.DB.Create(&tx).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "SYSTEM_ERROR", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"prepay_id": prepayID})
}
