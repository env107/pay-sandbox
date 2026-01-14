package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 沙箱配置
const (
	SandboxHost = "http://localhost:8080"
	AppID       = "wx88888888" // 请确保在控制台配置了这个 AppID
	MchID       = "12345678"   // 请确保在控制台配置了这个 MchID
)

func main() {
	// 1. 启动一个本地服务接收回调
	go startCallbackServer()

	// 等待服务启动
	time.Sleep(1 * time.Second)

	// 2. 发起下单请求
	log.Println(">>> 开始 JSAPI 下单...")
	prepayID, err := jsapiPrepay()
	if err != nil {
		log.Fatalf("下单失败: %v", err)
	}
	log.Printf(">>> 下单成功! PrepayID: %s", prepayID)
	log.Printf(">>> 请在浏览器打开以下链接进行模拟支付: http://localhost:3000/pay/preview/%s", prepayID)

	// 阻塞主进程
	select {}
}

func startCallbackServer() {
	r := gin.Default()
	r.POST("/notify", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		log.Printf("\n[Callback] 收到支付通知:\n%s\n", string(body))
		c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "OK"})
	})
	log.Println("业务回调服务启动在 :8081")
	r.Run(":8081")
}

func jsapiPrepay() (string, error) {
	url := fmt.Sprintf("%s/v3/pay/transactions/jsapi", SandboxHost)

	reqBody := map[string]interface{}{
		"appid":        AppID,
		"mchid":        MchID,
		"description":  "测试商品",
		"out_trade_no": fmt.Sprintf("ORDER_%d", time.Now().Unix()),
		"notify_url":   "http://localhost:8081/notify", // 回调地址
		"amount": map[string]interface{}{
			"total":    100,
			"currency": "CNY",
		},
		"payer": map[string]interface{}{
			"openid": "oUpF8uMuAJO_M2pxb1Q9zNjWeS6o",
		},
	}

	jsonBytes, _ := json.Marshal(reqBody)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status: %d, body: %s", resp.StatusCode, string(body))
	}

	var res map[string]string
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	return res["prepay_id"], nil
}
