package main

import (
	"fmt"
	"github.com/munisense/syntax-workshop-2021/internal/pkg/config"
	"log"

	"github.com/streadway/amqp"
)

// This is the name of a queue that we have already created. We could have called it anything.
const queue = "results"

func main() {
	c, err := config.LoadConfig()
	failOnError(err, "Failed to load config")

	conn, err := amqp.Dial(fmt.Sprintf("%s://%s:%s@%s:%d%s", c.Protocol, c.Username, c.Password, c.Host, c.Port, c.VHost))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	msgs, err := ch.Consume(
		queue, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	// Using the power of Go, we now consume from the queue while the rest of the application continues.
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// failOnError is a tiny helper function that outputs the error and terminates the program when a non-nil error is supplied
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
