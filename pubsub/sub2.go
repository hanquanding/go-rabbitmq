package main

import "go-rabbitmq/rabbitMQ"

func main() {
	r := rabbitMQ.NewRabbitMQPubSub("product")
	r.RecieveSub()
}
