package tests

import (
	"activemq/internal/app"
	"fmt"
	"github.com/go-stomp/stomp"
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
		fmt.Sprintf("Received incorrect message.  Actual: %s", msg.Body)
		time.Sleep(2 * time.Second)
	}
}
