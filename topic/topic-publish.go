package main

import (
	"fmt"
	"go-rabbitmq/rabbitMQ"
	"strconv"
	"time"
)

func main() {
	r1 := rabbitMQ.NewRabbitMQTopic("exchange-topic", "test.topic.one")
	r2 := rabbitMQ.NewRabbitMQTopic("exchange-topic", "test.topic.two")

	for i := 1; i <= 20; i++ {
		r1.PublishTopic("test topic one " + strconv.Itoa(i))
		r2.PublishTopic("test topic two " + strconv.Itoa(i))

		time.Sleep(time.Second)
		fmt.Println(i)
	}
}
