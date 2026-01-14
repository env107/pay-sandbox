package model

import (
	"time"

	"gorm.io/gorm"
)

// Merchant 商户配置
type Merchant struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	AppID           string         `gorm:"uniqueIndex;not null" json:"appid"`
	MchID           string         `gorm:"uniqueIndex;not null" json:"mchid"`
	APIV3Key        string         `gorm:"not null" json:"api_v3_key"`
	Description     string         `json:"description"`
	NotifyConfig    string         `gorm:"type:text" json:"notify_config"` // JSON string: {"interval": "1m", "max_retries": 3}
	NotifyUrl       string         `json:"notify_url"`                     // 默认回调地址
	RefundNotifyUrl string         `json:"refund_notify_url"`              // 退款回调地址
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// Transaction 交易订单
type Transaction struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	AppID          string     `gorm:"index" json:"appid"`
	MchID          string     `gorm:"index" json:"mchid"`
	Description    string     `json:"description"`
	OutTradeNo     string     `gorm:"uniqueIndex;not null" json:"out_trade_no"`
	TransactionID  string     `gorm:"uniqueIndex;not null" json:"transaction_id"` // 微信侧单号
	PrepayID       string     `gorm:"index" json:"prepay_id"`                     // 预支付ID
	Amount         int64      `json:"amount"`                                     // 分
	Currency       string     `json:"currency"`
	PayerOpenID    string     `json:"payer_openid"`
	Status         string     `gorm:"index" json:"status"` // CREATED, SUCCESS, REFUND, CLOSED
	NotifyUrl      string     `json:"notify_url"`
	CallbackStatus string     `json:"callback_status"` // SUCCESS, FAIL
	CallbackMsg    string     `json:"callback_msg"`    // 失败原因
	TradeType      string     `json:"trade_type"`      // JSAPI
	PaidAt         *time.Time `json:"paid_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// CallbackLog 回调日志
type CallbackLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TransactionID string    `gorm:"index" json:"transaction_id"`
	NotifyUrl     string    `json:"notify_url"`
	RequestBody   string    `gorm:"type:text" json:"request_body"`
	ResponseBody  string    `gorm:"type:text" json:"response_body"`
	StatusCode    int       `json:"status_code"`
	Status        string    `json:"status"` // SUCCESS, FAIL
	RetryCount    int       `json:"retry_count"`
	CreatedAt     time.Time `json:"created_at"`
}

// Refund 退款记录
type Refund struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	RefundID       string    `gorm:"uniqueIndex;not null" json:"refund_id"`     // 微信退款单号
	OutRefundNo    string    `gorm:"uniqueIndex;not null" json:"out_refund_no"` // 商户退款单号
	TransactionID  string    `gorm:"index;not null" json:"transaction_id"`      // 关联支付订单号
	MchID          string    `gorm:"index" json:"mchid"`
	Amount         int64     `json:"amount"` // 退款金额
	Total          int64     `json:"total"`  // 原订单总金额
	Currency       string    `json:"currency"`
	Reason         string    `json:"reason"`
	Status         string    `json:"status"` // SUCCESS, PROCESSING, ABNORMAL
	NotifyUrl      string    `json:"notify_url"`
	CallbackStatus string    `json:"callback_status"` // SUCCESS, FAIL
	CallbackMsg    string    `json:"callback_msg"`    // 失败原因
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
