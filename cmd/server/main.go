package main

import (
	"flag"
	"wepay-sandbox/internal/api"
	"wepay-sandbox/internal/api/admin"
	"wepay-sandbox/internal/api/mock"
	"wepay-sandbox/internal/core"

	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	// 初始化数据库
	core.InitDB("sandbox.db")

	r := gin.Default()

	// 允许跨域
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Mock API (Open)
	v3 := r.Group("/v3")
	{
		v3.POST("/pay/transactions/jsapi", mock.JSAPIPrepay)
		v3.GET("/pay/transactions/id/:transaction_id", mock.QueryByTransactionID)
		v3.GET("/pay/transactions/out-trade-no/:out_trade_no", mock.QueryByOutTradeNo)
		v3.POST("/pay/transactions/out-trade-no/:out_trade_no/close", mock.CloseOrder)
	}

	// Internal API (Admin)
	internal := r.Group("/api/internal")
	{
		internal.GET("/merchants", admin.ListMerchants)
		internal.POST("/merchants", admin.CreateMerchant)
		internal.PUT("/merchants/:id", admin.UpdateMerchant)
		internal.DELETE("/merchants", admin.DeleteMerchants)

		internal.GET("/transactions", admin.ListTransactions)
		internal.DELETE("/transactions", admin.DeleteTransactions)
		internal.GET("/transactions/:transaction_id/logs", admin.GetTransactionLogs)
		internal.POST("/transactions/:transaction_id/retry-callback", admin.RetryTransactionCallback)
		internal.POST("/simulate/pay", admin.SimulatePay)

		internal.POST("/simulate/refund", admin.SimulateRefund)
		internal.GET("/refunds", admin.ListRefunds)
		internal.DELETE("/refunds", admin.DeleteRefunds)
		internal.GET("/refunds/:refund_id/logs", admin.GetRefundLogs)
		internal.POST("/refunds/:refund_id/retry-callback", admin.RetryRefundCallback)

		internal.GET("/events", api.StreamEvents)
	}

	r.Run(":" + *port)
}
