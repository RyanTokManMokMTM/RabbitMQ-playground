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
	conn, err := mq.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "Failed to connect to rabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	key := severityFrom(os.Args)
	msg := bodyFrom(os.Args)

	err = ch.PublishWithContext(ctx, "log-direct", key, false, false, mq.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
	failOnError(err, "Failed to publish a message")

	log.Printf("[x] Sent a message %s for key %s", msg, key)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s : %s", msg, err)
	}
}

//go run send.go [form] [message...]
func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}
	return s
}
