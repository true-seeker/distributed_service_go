package services

import "fmt"

func GenerateTask(taskSize int) {
	newTask := GenerateRandomTask(taskSize)
	task := newTask.GetBackpackTaskParts()
	task = SaveNewTaskParts(task)
	fmt.Printf("Generated Task with ID: %d\n", task.ID)
}
