# Syntax Message Queue Workshop 2021

In this workshop we will learn the basic concepts of a Message Queue, and we will learn how to write a small Golang
application to interact with RabbitMQ as our message broker.

## Prerequisites to participate

- You have installed Docker (https://docs.docker.com/get-docker/)
- You have copied `.env.example` to `.env` and have entered the variables we will provide during the workshop

## Note for Windows users

Please use powershell (and not cmd) and change the volume mount as shown below

```shell
--volume="${pwd}:/app"
```

instead of:

```shell
--volume="$PWD:/app"
```

## Step 1 - get a single message from a shared queue

Lets get a message from a queue!

This application will try to get a single message from a queue called "results".

```shell
docker run --rm -it --volume="$PWD:/app" -w /app golang:1-alpine go run cmd/step1/main.go
```

## Step 2 - consume messages from a shared queue

We can do better than polling for every single message. Let's consume from a queue instead.

```shell
docker run --rm -it --volume="$PWD:/app" -w /app golang:1-alpine go run cmd/step2/main.go
```

## Step 3 - create your own exclusive queue

No more shared queue, we want to receive all messages.

This program will create an exclusive queue and bind it to the exchange called "results" with routing key "#".

The routing key determines what messages your queue will receive, a # means all messages.

```shell
docker run --rm -it --volume="$PWD:/app" -w /app golang:1-alpine go run cmd/step3/main.go
```

## Step 4 - create your own exclusive queue and only receive a specific set of messages

Up until now we have been receiving all messages. In this final example we will use the routing key to tell RabbitMQ we
only want to receive LAeq messages.

```shell
docker run --rm -it --volume="$PWD:/app" -w /app golang:1-alpine go run cmd/step4/main.go
```

### Links

- https://github.com/rabbitmq/rabbitmq-tutorials/blob/master/go/receive.go
- https://www.rabbitmq.com/tutorials/amqp-concepts.html
- https://www.rabbitmq.com/tutorials/tutorial-one-go.html
- https://www.rabbitmq.com/getstarted.html
- https://www.rabbitmq.com/amqp-0-9-1-reference.html
