package main

import (
	"bytes"
	mq "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func main() {
	conn, err := mq.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	//Set Qos
	err = ch.Qos(1, 0, false) //each worker only have 1 task before it send an ack
	failOnError(err, "Failed to set Qos")

	//Durable : ture -> no guarantee
	q, err := ch.QueueDeclare("WorkerQueueNoDelete", true, false, false, false, nil)
	failOnError(err, "Failed to register a queue")

	//SET AUTO ACK TO FALSE, COZ WE WILL ACK IT MANUALLY
	//TODO: Be careful un-acked message -> it won't be able to release before it is acked
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	var signal chan struct{}
	go func() {
		for m := range msgs {
			log.Printf("Received a message : %s\n", m.Body)
			waitFor := bytes.Count(m.Body, []byte("."))
			time.Sleep(time.Second * time.Duration(waitFor)) //stand for a busy task
			log.Println("Task Done!")
			m.Ack(false)
		}
	}()

	log.Printf("Waiting for message. To exist press CTRL+C.")
	<-signal

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s : %s", msg, err.Error())
	}
}
