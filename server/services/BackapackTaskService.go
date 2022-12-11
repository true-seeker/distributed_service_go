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
var maxBackpackCapacity = 100

func GenerateRandomTask(size int) BackpackTask {
	newTask := BackpackTask{items: nil,
		BackpackCapacity: uint32(rand.Intn(maxBackpackCapacity))}

	for i := 0; i < size; i++ {
		newItem := BackpackItem{
			Id:     i,
			weight: uint32(rand.Intn(maxWeight)) + 1,
			price:  uint32(rand.Intn(maxPrice)),
		}
		newTask.items = append(newTask.items, newItem)
	}

	return newTask
}

func (bt BackpackTask) GetBackpackTaskParts() Task {
	newTask := Task{}
	taskParts := *new([]TaskPart)
	for _, elem := range bt.items {
		newTaskPart := TaskPart{}

		fixedItem := BackpackTaskItem{
			Weight:  elem.weight,
			Price:   elem.price,
			IsFixed: true,
		}
		newTaskPart.Items = append(newTaskPart.Items, fixedItem)

		for _, elem2 := range bt.items {
			if elem2.Id == elem.Id {
				continue
			}
			backpackTaskItem := BackpackTaskItem{
				Weight:  elem2.weight,
				Price:   elem2.price,
				IsFixed: false,
			}
			newTaskPart.Items = append(newTaskPart.Items, backpackTaskItem)
		}
		taskParts = append(taskParts, newTaskPart)
	}
	newTask.TaskParts = taskParts
	newTask.BackpackCapacity = bt.BackpackCapacity
	return newTask
}
