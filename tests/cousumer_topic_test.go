package tests

import (
	"activemq/internal/app"
	"fmt"
	"github.com/go-stomp/stomp"
	"log"
	"testing"
	"time"
)

func TestActiveMQ(t *testing.T) {
	conn, err := stomp.Dial("tcp", "localhost:61613")
	if err != nil {
		t.Fatalf("Failed to connect to ActiveMQ: %v", err)
	}
	defer conn.Disconnect()

	sub, err := conn.Subscribe(app.TopicTest, stomp.AckClient, stomp.SubscribeOpt.Header("id", app.ConsumerName2))
	if err != nil {
		t.Fatalf("Failed to subscribe to topic: %v", err)
	}
	defer sub.Unsubscribe()

	for {
		msg := <-sub.C
		if msg.Err != nil {
			fmt.Printf("接收消息失败, err: %v\n", msg.Err)
			err = conn.Nack(msg)
			if err != nil {
				fmt.Printf("NACK失败, err: %v\n", err)
				break
			}
		}

		log.Printf("收到消息: %s", string(msg.Body))
		err = conn.Ack(msg)
		if err != nil {
			fmt.Printf("ACK失败, err: %v\n", err)
			break
		}
		fmt.Sprintf("Received incorrect message.  Actual: %s", msg.Body)
		time.Sleep(2 * time.Second)
	}
}
