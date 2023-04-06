package main

import (
	"context"
	mq "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"strings"
	"time"
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

	msg := bodyFrom(os.Args)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	err = ch.PublishWithContext(ctx, "logs", "", false, false, mq.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})

	failOnError(err, "Failed to publish a message")

	log.Printf("[x] sent a message")
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
