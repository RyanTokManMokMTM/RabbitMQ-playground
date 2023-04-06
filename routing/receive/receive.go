package main

import (
	mq "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"strings"
)

func main() {
	conn, err := mq.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "Failed to connect to rabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer conn.Close()

	//direct type with a key
	err = ch.ExchangeDeclare("log-direct", "direct", true, false, false, false, nil)
	failOnError(err, "Failed to declare a exchange")

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	failOnError(err, "Failed to declare a queue")

	//bind key
	if len(os.Args) < 2 {
		log.Println("")
		os.Exit(1)
	}

	for _, key := range os.Args[1:] {
		log.Printf("Binding queue %s to exchnage %s with key %s", q.Name, "log-direct", key)
		err = ch.QueueBind(q.Name, key, "log-direct", false, nil)
		failOnError(err, "Failed to bind a key to exchange of the queue")
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")
	var signal chan struct{}
	go func() {
		for msg := range msgs {
			log.Printf("receive a message %s with %s", msg.Body, msg.RoutingKey)
		}
	}()

	log.Println("Waiting for logs message. To exist press CTRL+C.")
	<-signal
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s : %s", msg, err)
	}
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
