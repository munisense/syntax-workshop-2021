package main

import (
	"fmt"
	"github.com/munisense/syntax-workshop-2021/internal/pkg/config"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/streadway/amqp"
)

const (
	exchange = "results"
	// Our routing key is no longer defined here, it is random generated at runtime
)

func main() {
	c, err := config.LoadConfig()
	failOnError(err, "Failed to load config")

	rand.Seed(time.Now().UnixNano())

	conn, err := amqp.Dial(fmt.Sprintf("%s://%s:%s@%s:%d%s", c.Protocol, c.Username, c.Password, c.Host, c.Port, c.VHost))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// We declare an exchange, (this exchange should already exist, but this code checks that it does)
	err = ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	failOnError(err, "Failed to declare an Exchange")
	hostname, _ := os.Hostname()
	pid := os.Getpid()

	// For the purpose of this demo we will generate a random routingKey feel free to changes this to your name for example!
	routingKey := fmt.Sprintf("%d", rand.Intn(math.MaxInt32))
	log.Printf("Using routingkey %s to publish our messages", routingKey)

	go func() {
		for _ = range time.Tick(time.Second * 5) {
			log.Print("Publishing a new message...")
			msg := amqp.Publishing{
				Timestamp: time.Now(),
				Body:      []byte(fmt.Sprintf("Message from %s.%d: %d", hostname, pid, rand.Intn(1337))),
			}
			err = ch.Publish(exchange, routingKey, false, false, msg)
			failOnError(err, "Failed to publish to an Exchange")
		}
	}()

	// For purposes of this demo we will also use the code from step 4 to get our messages in the same app, the code below has no changes
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
