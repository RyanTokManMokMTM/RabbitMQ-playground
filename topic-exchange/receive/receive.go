package main

import (
	mq "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

func main() {
	conn, err := mq.Dial("amqp://guest:guest@localhost:5672/")
	failedOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failedOnError(err, "Failed to open channel")
	defer conn.Close()

	err = ch.ExchangeDeclare("log-topic", "topic", true, false, false, false, nil)
	failedOnError(err, "Failed to declare exchange")

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	failedOnError(err, "Failed to declare a queue")

	//bind the key to the queue
	if len(os.Args) < 2 {
		log.Printf("Usage: %s [binding_key]...", os.Args[0])
		os.Exit(1)
	}

	for _, key := range os.Args[1:] {
		log.Printf("Bining the queue %s to exchange %s with routing key %s\n.", q.Name, "log-topic", key)
		err = ch.QueueBind(q.Name, key, "log-topic", false, nil)
		failedOnError(err, "Failed to bind the queue to exchange")
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failedOnError(err, "Failed to register a consumer")

	var signal chan struct{}
	for msg := range msgs {
		log.Printf("receive a message %s.\n", msg.Body)
	}

	log.Println("Waiting to logs message.To exist press CTRL+C.")
	<-signal
}

func failedOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s : %s", msg, err.Error())
	}
}
