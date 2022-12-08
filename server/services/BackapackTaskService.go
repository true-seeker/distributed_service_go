package services

import (
	"math/rand"
)

type BackpackTask struct {
	items []BackpackItem
}

type BackpackItem struct {
	Id     int
	weight float64
	price  float64
}

var maxPrice = 100.0
var maxWeight = 100.0

func GenerateTask(size int) BackpackTask {
	newTask := BackpackTask{items: nil}

	for i := 0; i < size; i++ {
		newItem := BackpackItem{
			Id:     i,
			weight: rand.Float64() * maxWeight,
			price:  rand.Float64() * maxPrice,
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
	return newTask
}
