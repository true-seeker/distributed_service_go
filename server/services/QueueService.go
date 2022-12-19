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

// GetQueueConnection Подключение к очереди
func GetQueueConnection() *queueConnection {
	fmt.Println(GetProperty("Queue", "user"))
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		GetProperty("Queue", "user"),
		GetProperty("Queue", "password"),
		GetProperty("Queue", "address"),
		GetProperty("Queue", "port")))
	FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"tasks", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	return &queueConnection{
		conn:      conn,
		channel:   ch,
		queue:     q,
		ctx:       ctx,
		ctxCancel: cancel,
	}
}

// PutNewTasksInQueue Сложить новую задачу в очередь
func (qc *queueConnection) PutNewTasksInQueue(task Task) {
	// Кладем одну и ту же задачу TaskIterationCount раз
	for i := 0; i < TaskIterationCount; i++ {
		qc.PutTaskInQueue(task)
	}
}

// PutTaskInQueue Сложить задачу в очередь
func (qc *queueConnection) PutTaskInQueue(task Task) {
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
		fmt.Println("Cant put task in queue", err)
	}
}

// GetTaskFromQueue получить задачу из очереди
func (qc *queueConnection) GetTaskFromQueue() *Task {
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

// GetMessageCountFromChannel Получить количество сообщений в очереди
func (qc *queueConnection) GetMessageCountFromChannel() int {
	return qc.queue.Messages
}
