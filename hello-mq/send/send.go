package main

import (
	"context"
	mq "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func main() {
	conn, err := mq.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to Rabbit MQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("Hello MQ", false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	msg := "Hello MQ Testing!!"
	err = ch.PublishWithContext(ctx, "", q.Name, false, false, mq.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})

	failOnError(err, "Failed to publish a message")
	log.Printf("[x] sent %s\n", msg)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s : %s", msg, err)
	}
}
