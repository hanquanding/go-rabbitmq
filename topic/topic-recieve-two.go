package main

import "go-rabbitmq/rabbitMQ"

func main() {
	r := rabbitMQ.NewRabbitMQTopic("exchange-topic", "test.*.two")
	r.RecieveTopic()
}
