package api

import (
	"io"

	"github.com/gin-gonic/gin"
)

// Event 广播事件结构
type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// GlobalEventChan 全局事件通道
var GlobalEventChan = make(chan Event, 100)

// SSEHeaders 设置 SSE 响应头
func SSEHeaders(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}

// StreamEvents SSE 处理函数
func StreamEvents(c *gin.Context) {
	SSEHeaders(c)

	// 创建一个专属的客户端通道
	clientChan := make(chan Event, 10)

	// 注册到全局广播 (这里简化处理，实际生产环境应该用更完善的 Pub/Sub 模型)
	// 为演示方便，我们启动一个 goroutine 监听 GlobalEventChan 并分发
	// 注意：这种简单的全局 channel 分发在多客户端时会有竞争问题，
	// 正确的做法是维护一个 client list。

	// 简单实现：轮询全局 channel 并不适合多客户端，
	// 我们改为：GlobalEventChan 仅作为生产者入口，
	// 我们需要一个 Broker 来管理所有客户端。

	broker.NewClients <- clientChan

	defer func() {
		broker.ClosingClients <- clientChan
	}()

	// 监听客户端通道
	c.Stream(func(w io.Writer) bool {
		if event, ok := <-clientChan; ok {
			c.SSEvent("message", event)
			return true
		}
		return false
	})
}

// 简单的 Broker 实现
type Broker struct {
	Notifier       chan Event
	NewClients     chan chan Event
	ClosingClients chan chan Event
	Clients        map[chan Event]bool
}

var broker = &Broker{
	Notifier:       GlobalEventChan,
	NewClients:     make(chan chan Event),
	ClosingClients: make(chan chan Event),
	Clients:        make(map[chan Event]bool),
}

func init() {
	go broker.listen()
}

func (b *Broker) listen() {
	for {
		select {
		case s := <-b.NewClients:
			b.Clients[s] = true
		case s := <-b.ClosingClients:
			delete(b.Clients, s)
		case event := <-b.Notifier:
			for clientMessageChan := range b.Clients {
				select {
				case clientMessageChan <- event:
				default:
					// 如果客户端阻塞，跳过
				}
			}
		}
	}
}
