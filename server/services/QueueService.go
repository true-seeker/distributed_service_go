package services

import (
	"context"
	"encoding/json"
	"fmt"
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
	if qc.conn == nil {
		fmt.Println("qc.conn == nil", part)
		qc = getQueueConnection()
		defer qc.conn.Close()
		defer qc.channel.Close()
		defer qc.ctxCancel()
	}

	marshaledPart, _ := json.Marshal(part)
	err := qc.channel.PublishWithContext(qc.ctx,
		"",            // exchange
		qc.queue.Name, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			Body: marshaledPart,
		})
	if err != nil {
		fmt.Println("Cant requeue message", err)
	}
}

func GetTaskPartFromQueue() *TaskPart {
	qc := getQueueConnection()
	defer qc.conn.Close()
	defer qc.channel.Close()
	defer qc.ctxCancel()

	msg, ok, err := qc.channel.Get(
		qc.queue.Name, // queue
		true,          // auto-ack
	)
	FailOnError(err, "Failed to register a consumer")
	if ok {
		taskPart := TaskPart{}

		err = json.Unmarshal(msg.Body, &taskPart)
		FailOnError(err, "Failed to unmarshal taskPart")
		return &taskPart
	} else {
		return nil
	}

}

func GetMessageCountFromChannel() int {
	qc := getQueueConnection()
	defer qc.conn.Close()
	defer qc.channel.Close()
	defer qc.ctxCancel()

	return qc.queue.Messages
}
