package services

import "fmt"

var TaskIterationCount = 3

func GenerateTask(taskSize int) {
	newTask := GenerateRandomTask(taskSize)
	//task := newTask.GetBackpackTaskParts()
	task := SaveNewTaskParts(newTask)
	PutNewTasksInQueue(task)

	fmt.Printf("Generated Task with ID: %d\n", task.ID)
}
