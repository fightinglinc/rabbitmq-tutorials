package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// url format = amqp://username:password@rabbitmq service address/port number/vhost
const MQURL = "amqp://imoocuser:imoocuser@127.0.0.1:5672/imooc"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// name of queue
	QueueName string
	// exchanges
	Exchange string
	// key
	Key string
	// conn info
	Mqurl string
}

// create RabbitMQ struct instance
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     MQURL,
	}
	// create RabbitMQ connection
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect RabbitMQ!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open the channel")
	return rabbitmq
}

// disconnect
func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

// error handle function
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// Step 1: Create RabbitMQ instance for simple mode
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "")
}

// Step 2: Sending-Simple Mode
func (r *RabbitMQ) PublishSimple(message string) {
	// 1. Declare a queue for us to send.
	// Declaring a queue is idempotent-only create a queue if it doesn't exist.
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	// 2.Send messages, publish a message to the queue
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		// Publishings can be undeliverable
		// when the mandatory flag is true and no queue is bound that matches the routing key,
		// or when the immediate flag is true
		// and no consumer on the matched queue is ready to accept the delivery.
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func (r *RabbitMQ) ConsumeSimple() {
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	// Receive messages
	msgs, err := r.channel.Consume(
		r.QueueName,
		"",
		true, // auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	// Read the message from the channel in a goroutine
	go func() {
		for d:= range msgs {
			log.Printf("Received a message: %s", d.Body)
			fmt.Println(d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
