package services

import (
	"math/rand"
)

var maxPrice = 100
var maxWeight = 100
var maxBackpackCapacity = 7500

// GenerateRandomTask Генерация случайной задачи
func GenerateRandomTask(size int) Task {
	newTask := Task{Items: nil,
		BackpackCapacity: uint32(rand.Intn(maxBackpackCapacity))}

	for i := 0; i < size; i++ {
		newItem := BackpackTaskItem{
			Weight: uint32(rand.Intn(maxWeight)) + 1,
			Price:  uint32(rand.Intn(maxPrice)),
		}
		newTask.Items = append(newTask.Items, newItem)
	}

	return newTask
}
