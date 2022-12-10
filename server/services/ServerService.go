package services

import "fmt"

var TASK_ITERATION_COUNT = 3

func GenerateTask(taskSize int) {
	newTask := GenerateRandomTask(taskSize)
	task := newTask.GetBackpackTaskParts()
	task = SaveNewTaskParts(task)
	PutNewTasksInQueue(task)

	fmt.Printf("Generated Task with ID: %d\n", task.ID)
}
