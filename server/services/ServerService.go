package services

import "fmt"

var TaskIterationCount = 3
var DefaultTaskSize = 5000

func GenerateTask(taskSize int) {
	newTask := GenerateRandomTask(taskSize)
	task := SaveNewTaskParts(newTask)
	PutNewTasksInQueue(task)

	fmt.Printf("Generated Task with ID: %d\n", task.ID)
}
