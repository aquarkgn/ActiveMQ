package service

var ConsumerInstance = new(Consumer)

type Consumer struct{}

func (consumer *Consumer) Message() {

}
