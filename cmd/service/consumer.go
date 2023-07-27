package service

import (
	"fmt"
	"github.com/go-stomp/stomp"
	"log"
	"sync"
	"time"

	"activemq/internal/app"
)

var ConsumerInstance = new(Consumer)

type Consumer struct{}

// 订阅队列
func (consumer *Consumer) SubscribeQueue(conn *stomp.Conn, queueName string) (*stomp.Subscription, error) {
	sub, err := conn.Subscribe(
		app.QueueName,
		stomp.AckAuto,
		stomp.SubscribeOpt.Header("id", app.ConsumerName),
	)
	if err != nil {
		return nil, err
	}
	fmt.Printf("订阅队列 %s 成功\n", queueName)
	return sub, nil
}

// 消费队列消息
func (consumer *Consumer) ConsumeQueueMessages(conn *stomp.Conn, sub *stomp.Subscription, wg *sync.WaitGroup) {
	defer wg.Done()

	// 设置定时器，达到超时时间后关闭订阅
	timer := time.NewTimer(10 * time.Second)
	// 持续读取消息
	for {
		select {
		case msg := <-sub.C:
			log.Printf("接收到消息：%s\n", string(msg.Body))
		case <-timer.C:
			log.Printf("订阅超时，关闭订阅")
			err := sub.Unsubscribe()
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}

	// 创建通道来处理消息
	messages := make(chan *stomp.Message)

	// 启动 goroutine 来消费消息
	go func() {
		for {
			msg := <-sub.C
			messages <- msg
		}
	}()

	for {
		select {
		case msg := <-messages:
			// 在此处处理收到的消息
			log.Printf("接收到消息：%s\n", string(msg.Body))

			// 模拟长时间处理任务
			time.Sleep(1 * time.Second)

			log.Printf("消息处理完成")
		case <-time.After(5 * time.Second):
			log.Printf("订阅超时，继续消费")

			// 取消订阅并重新订阅
			err := sub.Unsubscribe()
			if err != nil {
				log.Fatal(err)
			}

			sub, err = conn.Subscribe(
				app.QueueName, // 订阅的队列名称
				stomp.AckAuto, // 设置自动确认模式
				stomp.SubscribeOpt.Header("id", app.ConsumerName), // 设置消费者ID
			)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
