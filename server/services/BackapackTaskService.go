package services

import (
	"math/rand"
)

type BackpackTask struct {
	items            []BackpackItem
	BackpackCapacity uint32
}

type BackpackItem struct {
	Id     int
	weight uint32
	price  uint32
}

var maxPrice = 100
var maxWeight = 100
var maxBackpackCapacity = 750

func GenerateRandomTask(size int) Task {
	newTask := Task{Items: nil,
		BackpackCapacity: uint32(rand.Intn(maxBackpackCapacity))}

	for i := 0; i < size; i++ {
		newItem := BackpackTaskItem{
			Weight: uint32(rand.Intn(maxWeight)) + 40,
			Price:  uint32(rand.Intn(maxPrice)),
		}
		newTask.Items = append(newTask.Items, newItem)
	}

	return newTask
}

func (bt BackpackTask) GetBackpackTaskParts() Task {
	newTask := Task{}

	for _, elem := range bt.items {
		if elem.Id == elem.Id {
			continue
		}
		backpackTaskItem := BackpackTaskItem{
			Weight: elem.weight,
			Price:  elem.price,
		}
		newTask.Items = append(newTask.Items, backpackTaskItem)
	}
	newTask.BackpackCapacity = bt.BackpackCapacity
	return newTask
}
