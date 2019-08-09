package main

import "go-rabbitmq/rabbitMQ"

func main() {
	r := rabbitMQ.NewRabbitMQRouting("test-exchange", "test-routingKey-1")
	r.RecieveRouting()
}
