package main

import (
	mq "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	mq, err := mq.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer mq.Close()

	ch, err := mq.Channel()
	failOnError(err, "Failed to open a channel")

	//In order to avoid any error happen,we need to make sure the queue exist.
	//So we declare here again

	q, err := ch.QueueDeclare("Hello MQ", false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	msg, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")
	//infinite loop for receiving message from queue
	var signal chan struct{}
	go func() {
		for m := range msg {
			log.Printf("Received a message : %s", m.Body)
		}
	}()

	log.Printf("Waiting for message. To exist press CTRL+C")
	<-signal
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s : %s", msg, err)
	}
}
