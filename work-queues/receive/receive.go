package main

import (
	"bytes"
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

	q, err := ch.QueueDeclare("WorkerQueue", false, false, false, false, nil)
	failOnError(err, "Failed to register a queue")

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	var signal chan struct{}
	go func() {
		for m := range msgs {
			log.Printf("Received a message : %s\n", m.Body)
			waitFor := bytes.Count(m.Body, []byte("."))
			time.Sleep(time.Second * time.Duration(waitFor)) //stand for a busy task
			log.Println("Task Done!")
		}
	}()

	log.Printf("Waiting for message. To exist press CTRL+C.")
	<-signal
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s : %s", msg, err)
	}
}
