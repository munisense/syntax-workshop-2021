package main

import (
	"fmt"
	"github.com/munisense/syntax-workshop-2021/internal/pkg/config"
	"log"

	"github.com/streadway/amqp"
)

const (
	exchange = "results"
	// A routing key is a dot-separated string. A * can substitute exactly one word, a # can substitute zero or more words.
	routingKey = "#.Sound2.LAeq"
)

func main() {
	c, err := config.LoadConfig()
	failOnError(err, "Failed to load config")

	conn, err := amqp.Dial(fmt.Sprintf("%s://%s:%s@%s:%d%s", c.Protocol, c.Username, c.Password, c.Host, c.Port, c.VHost))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// We supply an empty name to make RabbitMQ generate a random name for us.
	q, err := ch.QueueDeclare("", false, true, true, false, nil)
	failOnError(err, "Failed to declare a queue")

	// And bind it to the exchange called "results".
	// Be sure to retrieve the randomly generated queue name from the "q" struct!
	err = ch.QueueBind(q.Name, routingKey, exchange, false, nil)
	failOnError(err, "Failed to bind the queue")

	log.Printf("Our exclusive queue name is '%s' and is now bound to exchange '%s' using routing key '%s'", q.Name, exchange, routingKey)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		true,   // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	// Using the power of Go, we now consume from the queue while the rest of the application continues.
	go func() {
		for d := range msgs {
			log.Printf("Received a message, routing key: %s body: %s", d.RoutingKey, d.Body)
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
