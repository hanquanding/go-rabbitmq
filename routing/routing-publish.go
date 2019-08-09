package main

import (
	"fmt"
	"go-rabbitmq/rabbitMQ"
	"strconv"
	"time"
)

func main() {
	r1 := rabbitMQ.NewRabbitMQRouting("test-exchange", "test-routingKey-1")
	r2 := rabbitMQ.NewRabbitMQRouting("test-exchange", "test-routingKey-2")

	for i := 1; i <= 20; i++ {
		r1.PublishRouting("test-exchange : test-routingKey-1: " + strconv.Itoa(i))
		r2.PublishRouting("test-exchange : test-routingKey-2: " + strconv.Itoa(i))
		time.Sleep(time.Second)
		fmt.Println("生产消息：" + strconv.Itoa(i))
	}
}
