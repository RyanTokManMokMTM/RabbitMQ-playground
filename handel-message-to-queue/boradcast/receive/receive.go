package main

import (
	mq "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	conn, err := mq.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	//we use exchange to publish message to the queue
	err = ch.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
	failOnError(err, "Failed to create exchange")

	//queue
	//declare the queue here -> can have multiple queue?
	q, err := ch.QueueDeclare("", false, false, false, true, nil)
	failOnError(err, "Failed to create queue")

	err = ch.QueueBind(q.Name, "", "logs", false, nil)
	failOnError(err, "Failed to bind a queue to exchange")

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer.")

	var single chan struct{}

	go func() {
		for msg := range msgs {
			log.Printf("received a message : %s ", msg.Body)
		}
	}()

	log.Printf("Waiting for message. To Exit press CTRL+Z.")
	<-single
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s : %s", msg, err)
	}
}
