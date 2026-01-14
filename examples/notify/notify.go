package main

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 设置 Gin 为发布模式，减少干扰日志
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 定义回调接收接口
	r.POST("/notify", func(c *gin.Context) {
		// 读取请求体
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("读取请求失败: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		// 打印接收到的回调内容
		log.Printf("\n========== 收到微信支付回调 ==========\n%s\n======================================\n", string(body))

		// 按照微信支付V3文档要求，返回 200 或 204 即可代表处理成功
		// 这里返回 200 OK 和一个简单的 JSON 结构
		c.JSON(http.StatusOK, gin.H{
			"code":    "SUCCESS",
			"message": "OK",
		})
	})

	addr := ":8081"
	log.Printf("回调服务已启动，监听地址 %s/notify", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
