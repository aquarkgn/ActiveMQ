package service

import (
	"fmt"
	"github.com/go-stomp/stomp"
	"sync"

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
func (consumer *Consumer) ConsumeQueueMessages(sub *stomp.Subscription, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		msg := <-sub.C
		fmt.Printf("接收到队列消息：%s\n", string(msg.Body))
		// 处理接收到的消息，可以根据业务需求进行逻辑处理
	}
}
