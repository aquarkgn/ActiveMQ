package service

import (
	"activemq/internal/app"
	"fmt"
	"github.com/go-stomp/stomp"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func Publisher() {
	// 创建连接
	conn, err := stomp.Dial("tcp", app.BrokerAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Disconnect()

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
			err = conn.Send(
				app.QueueName, // 发送到的队列名称
				"text/plain",
				[]byte(msg), // 消息内容
				nil,
			)
			if err != nil {
				log.Printf("消息发送失败 %s %s %s", app.QueueName, msg, err.Error())
			} else {
				log.Printf("消息发送成功 %s %s", app.QueueName, msg)
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
