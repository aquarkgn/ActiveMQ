package main

import (
	"activemq/internal/app"
	"fmt"
	"github.com/go-stomp/stomp"
	"log"
	"time"
)

func main() {
Loop:
	for {
		fmt.Println("启动ActiveMQ消费者...")
		time.Sleep(5 * time.Second)

		// Connect to broker
		client, err := stomp.Dial("tcp", app.Url,
			stomp.ConnOpt.Login(app.Username, app.Password),
			stomp.ConnOpt.AcceptVersion(stomp.V12),
			stomp.ConnOpt.HeartBeat(5*time.Second, 0*time.Second))
		if err != nil {
			fmt.Printf("连接到ActiveMQ失败, err: %v\n", err)
			continue Loop
		}
		log.Printf("已连接到 ActiveMQ: %s", app.Url)
		sub, err := client.Subscribe(app.QueueTest, stomp.AckClient, stomp.SubscribeOpt.Header("id", app.ConsumerName))
		if err != nil {
			fmt.Printf("订阅队列失败, err: %v\n", err)
			continue Loop
		}

		fmt.Println("已启动ActiveMQ消费者, 正在等待消息...")
		for {
			msg := <-sub.C
			if msg.Err != nil {
				fmt.Printf("接收消息失败, err: %v\n", msg.Err)
				err = client.Nack(msg)
				if err != nil {
					fmt.Printf("NACK失败, err: %v\n", err)
					break
				}
			}

			fmt.Println(string(msg.Body))
			log.Printf("接收到消息：%s\n", string(msg.Body))
			err = client.Ack(msg)
			if err != nil {
				fmt.Printf("ACK失败, err: %v\n", err)
				break
			}
		}
	}
}
