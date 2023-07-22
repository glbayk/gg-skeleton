package services

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var ch *amqp.Channel
var AmqpClient *RabbitMQClient

type AmqpConfirmationMessage struct {
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

type AmqpForgotPasswordMessage struct {
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type RabbitMQClient struct{}

func (rmq *RabbitMQClient) failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("[%s]: %s", msg, err)
	}
}

func (rmq *RabbitMQClient) SendMessage(qName string, message string) {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	rmq.failOnError(err, "Failed to connect to RabbitMQ")

	defer conn.Close()
	log.Println("Connected to RabbitMQ successfully")

	ch, err = conn.Channel()
	rmq.failOnError(err, "Failed to open a channel")
	defer ch.Close()
	log.Println("Channel opened successfully")

	q, err := ch.QueueDeclare(
		qName, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	rmq.failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	rmq.failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", message)
}
