package tests

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"testing"
	"time"

	"activemq/internal/app"

	"github.com/go-stomp/stomp"
)

func TestPublishMessage(t *testing.T) {
	Publisher()
}

func Publisher() {
	// 创建连接
	fmt.Println("启动ActiveMQ消费者...")
	time.Sleep(1 * time.Second)

	// Connect to broker
	client, err := stomp.Dial("tcp", app.Url,
		stomp.ConnOpt.Login(app.Username, app.Password),
		stomp.ConnOpt.AcceptVersion(stomp.V12),
		stomp.ConnOpt.HeartBeat(5*time.Second, 0*time.Second))
	if err != nil {
		fmt.Printf("连接到ActiveMQ失败, err: %v\n", err)
	}

	// 使用 WaitGroup 来等待消费者 goroutine 完成
	var wg sync.WaitGroup
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		// 循环发送消息
		log.Println("正在发送消息")
		i := 0
		for {
			// 发送消息
			msg := fmt.Sprintf("Number[ %d ] Hello, ActiveMQ Hong!", i)
			err = client.Send(
				app.QueueTest, // 发送到的队列名称
				"text/plain",
				[]byte(msg), // 消息内容
				nil,
			)
			if err != nil {
				log.Printf("消息发送失败 %s %s %s", app.QueueTest, msg, err.Error())
			} else {
				log.Printf("消息发送成功 %s %s", app.QueueTest, msg)
			}

			// 延迟一段时间再发送下一条消息
			time.Sleep(1 * time.Second)
			i++
		}
	}(&wg)

	// 等待中断信号，优雅地关闭连接
	log.Printf("按 CTRL+C 退出")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	// 等待消费者 goroutine 完成
	wg.Wait()
	log.Printf("已关闭通道")
}
