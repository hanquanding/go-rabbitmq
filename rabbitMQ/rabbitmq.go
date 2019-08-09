package rabbitMQ

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const AMQPURI = "amqp://guest:guest@127.0.0.1:5672/test"

type RabbitMQ struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	ExchangeName string // 交换机名称
	QueueName    string // 队列名称
	RoutingKey   string // 路由
	AmqpUri      string // 链接信息
}

// 创建RabbitMQ结构体实例
func NewRabbitMQ(exchangeName, queueName, routingKey string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		ExchangeName: exchangeName,
		QueueName:    queueName,
		RoutingKey:   routingKey,
		AmqpUri:      AMQPURI,
	}

	var err error
	log.Printf("dialing: [%q]", rabbitmq.AmqpUri)
	rabbitmq.conn, err = amqp.Dial(rabbitmq.AmqpUri)
	rabbitmq.checkError(err, "connect to rabbitMQ failed！")
	log.Printf("got connection, getting channel")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.checkError(err, "open channel failed！")
	log.Printf("got channel,exchange: [%q], queue: [%q], routing key: [%q]", exchangeName, queueName, routingKey)
	return rabbitmq
}

// 错误检查
func (r *RabbitMQ) checkError(err error, message string) {
	if err != nil {
		log.Fatalf("[%s]: %v", message, err)
		panic(fmt.Sprintf("[%s]: %v", message, err))
	}

}

// 创建simple模式RabbitMQ结构体实例
// simple模式exchange和key为空
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", queueName, "")
	return rabbitmq
}

// simple模式下的producer
// 1.申请队列
// 2.发送消息到队列
func (r *RabbitMQ) PublishSimple(message string) {
	_, err := r.channel.QueueDeclare(
		r.QueueName, // name 队列名称
		false,       // durable 是否持久化
		false,       // autoDelete 是否自动删除
		false,       // exclusive 是否具有排他性
		false,       // noWait 是否阻塞
		nil,
	)
	if err != nil {
		fmt.Errorf("queue declare: [%s]", err)
	}

	r.channel.Publish(
		r.ExchangeName, // exchange
		r.QueueName,    // routingKey
		false,          // mandatory 根据exchange类型和routkey规则，如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,          // immediate 如果为true，当exchange发送消息到队列发现队列上没有绑定消费者，则会把消息发还给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// simple模式下的consumer
// 1.申请队列
// 2.接收消息
func (r *RabbitMQ) ConsumerSimple() {
	_, err := r.channel.QueueDeclare(
		r.QueueName, // name 队列名称
		false,       // durable 是否持久化
		false,       // autoDelete 是否自动删除
		false,       // exclusive 是否具有排他性
		false,       // noWait 是否阻塞
		nil,
	)
	if err != nil {
		fmt.Errorf("queue declare: [%s]", err)
	}

	msgs, err := r.channel.Consume(
		r.QueueName,
		"",    // 用来区分多个消费者
		true,  // autoACK 是否自动应答
		false, // exclusive 是否具有排他性
		false, // noLocal 如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false, // noWait 是否阻塞
		nil,
	)
	if err != nil {
		fmt.Errorf("queue consume: [%s]", err)
	}

	done := make(chan bool)
	// 启动协程处理消息
	go func() {
		for msg := range msgs {
			log.Printf("received a message: [%s]", msg.Body)
		}
	}()

	log.Printf("waiting for message===>")
	<-done
}

// 订阅模式创建RabbitMQ实例
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(exchangeName, "", "")
	return rabbitmq
}

// 订阅模式下的producer
// 1.创建交换机
// 2.发送消息
func (r *RabbitMQ) PublishPub(message string) {
	err := r.channel.ExchangeDeclare(
		r.ExchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.checkError(err, "failed to declare an exchange")
	err = r.channel.Publish(
		r.ExchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// 订阅模式consumer
// 1.创建交换机
// 2.创建队列
// 3.绑定队列到交换机
// 4.消费消息
func (r *RabbitMQ) RecieveSub() {
	err := r.channel.ExchangeDeclare(
		r.ExchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.checkError(err, "failed to declare an exchange")
	q, err := r.channel.QueueDeclare(
		"", // 随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.checkError(err, "failed to declare a queue")
	err = r.channel.QueueBind(
		q.Name,
		"", // pub/sub模式下，key为空
		r.ExchangeName,
		false,
		nil,
	)
	if err != nil {
		fmt.Errorf("queue Bind: [%s]", err)
	}
	message, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Errorf("queue consume: [%s]", err)
	}

	done := make(chan bool)
	go func() {
		for msg := range message {
			log.Printf("received a message: [%s]", msg.Body)
		}
	}()
	log.Printf("waiting for message===>")
	<-done
}

// 路由模式创建RabbitMQ实例
func NewRabbitMQRouting(exchangeName string, routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(exchangeName, "", routingKey)
	return rabbitmq
}

// 路由模式producer
// 1.创建交换机
// 2.发送消息
func (r *RabbitMQ) PublishRouting(message string) {
	err := r.channel.ExchangeDeclare(
		r.ExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.checkError(err, "failed to declare an exchange")
	err = r.channel.Publish(
		r.ExchangeName,
		r.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// 路由模式consumer
func (r *RabbitMQ) RecieveRouting() {
	err := r.channel.ExchangeDeclare(
		r.ExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.checkError(err, "failed to declare an exchange")
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	r.checkError(err, "failed to declare a queue")

	err = r.channel.QueueBind(
		q.Name,
		r.RoutingKey,
		r.ExchangeName,
		false,
		nil,
	)
	if err != nil {
		fmt.Errorf("queue Bind: [%s]", err)
	}

	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Errorf("queue consume: [%s]", err)
	}

	done := make(chan bool)
	go func() {
		for msg := range messages {
			log.Printf("received a message: [%s]", msg.Body)
		}
	}()
	log.Printf("waiting for message===>")
	<-done
}

// topic模式创建RabbitMQ实例
func NewRabbitMQTopic(exchangeName string, routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(exchangeName, "", routingKey)
	return rabbitmq
}

// topic模式producer
func (r *RabbitMQ) PublishTopic(message string) {
	err := r.channel.ExchangeDeclare(
		r.ExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.checkError(err, "failed to declare a exchange")

	err = r.channel.Publish(
		r.ExchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func (r *RabbitMQ) RecieveTopic() {
	err := r.channel.ExchangeDeclare(
		r.ExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.checkError(err, "failed to declare an exchange")
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	r.checkError(err, "failed to declare a queue")

	err = r.channel.QueueBind(
		q.Name,
		r.RoutingKey,
		r.ExchangeName,
		false,
		nil,
	)
	if err != nil {
		fmt.Errorf("queue Bind: [%s]", err)
	}

	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Errorf("queue consume: [%s]", err)
	}

	done := make(chan bool)
	go func() {
		for msg := range messages {
			log.Printf("received a message: [%s]", msg.Body)
		}
	}()
	log.Printf("waiting for message===>")
	<-done
}
