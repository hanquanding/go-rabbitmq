package main

import (
	"fmt"
	"go-rabbitmq/rabbitMQ"
	"strconv"
	"time"
)

func main() {
	r := rabbitMQ.NewRabbitMQSimple("test-work")

	for i := 1; i <= 20; i++ {
		r.PublishSimple("Hello Work! " + strconv.Itoa(i))
		time.Sleep(time.Second)
		fmt.Println(i)
	}

}
