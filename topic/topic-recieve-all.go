package main

import "go-rabbitmq/rabbitMQ"

func main() {
	r := rabbitMQ.NewRabbitMQTopic("exchange-topic", "#")
	r.RecieveTopic()
}
