package main

import "go-rabbitmq/rabbitMQ"

func main() {
	r := rabbitMQ.NewRabbitMQSimple("test-work")
	r.ConsumerSimple()
}
