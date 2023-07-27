package main

import (
	"activemq/cmd/service"
	"activemq/internal/app"
	"github.com/go-stomp/stomp"
	"log"
	"os"
	"os/signal"
	"sync"
)

func main() {
	// 创建连接
	conn, err := stomp.Dial("tcp", app.BrokerAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Disconnect()
	log.Printf("已连接到 ActiveMQ: %s", app.BrokerAddr)

	// 订阅队列
	subQueue, err := service.ConsumerInstance.SubscribeQueue(conn, app.QueueName)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("已订阅队列: %s", app.QueueName)

	// 使用 WaitGroup 来等待消费者 goroutine 完成
	var wg sync.WaitGroup

	// 启动消费者
	wg.Add(1)
	log.Printf("正在启动消费者")
	go service.ConsumerInstance.ConsumeQueueMessages(subQueue, &wg)

	// 等待中断信号，优雅地关闭连接
	log.Printf("按 CTRL+C 退出")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	// 等待消费者 goroutine 完成
	wg.Wait()
	log.Printf("已关闭通道")

}
