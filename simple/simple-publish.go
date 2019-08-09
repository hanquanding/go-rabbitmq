package main

import (
	"go-rabbitmq/rabbitMQ"
	"log"
)

func main() {
	r := rabbitMQ.NewRabbitMQSimple("test-simple")
	r.PublishSimple("Hello Simple！！！")
	log.Println("发送成功！")
}
