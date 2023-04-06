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
	failedOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failedOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare("log-topic", "topic", true, false, false, false, nil)
	failedOnError(err, "Failed to declare a exchange")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	body := bodyFrom(os.Args)
	err = ch.PublishWithContext(ctx, "log-topic", severityFrom(os.Args), false, false, mq.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})

	failedOnError(err, "Failed to publish the message")
	log.Printf("[x] sent message %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if len(args) < 3 || os.Args[2] == " " {
		s = "hello"
	} else {
		s = strings.Join(os.Args[2:], " ")
	}

	return s
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "anonymous.info"
	} else {
		s = os.Args[1]
	}
	return s
}

func failedOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s : %s", msg, err.Error())
	}
}
