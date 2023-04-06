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
	failOnError(err, "Failed to connect to Rabbit MQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("WorkerQueue", false, false, false, false, nil)
	failOnError(err, "Failed to register a queue")

	body := bodyFrom(os.Args)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	//msg := "Hello Worker..."
	err = ch.PublishWithContext(ctx, "", q.Name, false, false, mq.Publishing{
		DeliveryMode: mq.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(body),
	})

	failOnError(err, "Failed to publish a message")
	log.Printf("[x] sent a message : %s", body)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s : %s", msg, err.Error())
	}
}

func bodyFrom(args []string) string {
	var str string
	if len(args) < 2 || os.Args[1] == "" {
		str = "hello"
	} else {
		str = strings.Join(args[1:], " ")
	}
	return str
}
