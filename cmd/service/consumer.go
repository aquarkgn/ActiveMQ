package service

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-stomp/stomp"

	"activemq/internal/app"
)

var ConsumerInstance = new(Consumer)

type Consumer struct{}

func (consumer *Consumer) Message() {
	// 创建连接
	conn, err := stomp.Dial("tcp", app.BrokerAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Disconnect()
	log.Printf("已连接到 ActiveMQ: %s", app.BrokerAddr)

	// 订阅队列
	subQueue, err := subscribeQueue(conn, app.QueueName)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("已订阅队列: %s", app.QueueName)

	// 创建一个通道用于接收中断信号
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// 启动一个goroutine，从通道接收中断信号并处理
	go func() {
		<-interrupt
		fmt.Println("接收到中断信号，停止消费")
		os.Exit(0)
	}()

	log.Printf("正在启动消费者")
	// 持续消费队列消息
	consumeQueueMessages(subQueue)
}

// 订阅队列
func subscribeQueue(conn *stomp.Conn, queueName string) (*stomp.Subscription, error) {
	sub, err := conn.Subscribe(app.QueueName, stomp.AckAuto,
		stomp.SubscribeOpt.Header("id", app.ConsumerName))
	if err != nil {
		return nil, err
	}
	fmt.Printf("订阅队列 %s 成功\n", queueName)
	return sub, nil
}

// 消费队列消息
func consumeQueueMessages(sub *stomp.Subscription) {
	for {
		msg := <-sub.C
		fmt.Printf("接收到队列消息：%s\n", string(msg.Body))
		// 处理接收到的消息，可以根据业务需求进行逻辑处理
	}
}
