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

	for i := 0; i < TaskIterationCount; i++ {
		PutTaskInQueue(task, qc)
	}
}

func PutTaskInQueue(task Task, qc queueConnection) {
	if qc.conn == nil {
		qc = getQueueConnection()
		defer qc.conn.Close()
		defer qc.channel.Close()
		defer qc.ctxCancel()
	}

	marshaledTask, _ := json.Marshal(task)
	err := qc.channel.PublishWithContext(qc.ctx,
		"",
		qc.queue.Name,
		false,
		false,
		amqp.Publishing{
			Body: marshaledTask,
		})
	if err != nil {
		fmt.Println("Cant requeue message", err)
	}
}

func GetTaskPartFromQueue() *Task {
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
		task := Task{}

		err = json.Unmarshal(msg.Body, &task)
		FailOnError(err, "Failed to unmarshal task")
		return &task
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
