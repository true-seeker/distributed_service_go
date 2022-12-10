package services

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type queueConnection struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queue     amqp.Queue
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func getQueueConnection() queueConnection {
	conn, err := amqp.Dial("amqp://lab2:lab2@176.124.200.41:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"taskParts", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	return queueConnection{
		conn:      conn,
		channel:   ch,
		queue:     q,
		ctx:       ctx,
		ctxCancel: cancel,
	}
}

func PutNewTasksInQueue(task Task) {
	qc := getQueueConnection()
	defer qc.conn.Close()
	defer qc.channel.Close()
	defer qc.ctxCancel()

	for _, taskPart := range task.TaskParts {
		for i := 0; i < TASK_ITERATION_COUNT; i++ {
			PutTaskPartInQueue(taskPart, qc)
		}
	}
}

func PutTaskPartInQueue(part TaskPart, qc queueConnection) {
	marshaledPart, _ := json.Marshal(part)
	qc.channel.PublishWithContext(qc.ctx,
		"",            // exchange
		qc.queue.Name, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			Body: marshaledPart,
		})
}
