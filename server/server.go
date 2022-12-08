package main

import (
	"flag"
	"fmt"
	"server/services"
)

func main() {
	generateTask := flag.Bool("g", false, "generate new task")
	taskSize := flag.Int("s", 5, "new task size")
	flag.Parse()
	fmt.Println(*generateTask, *taskSize)

	services.Migrate()

	if *generateTask {
		newTask := services.GenerateTask(*taskSize)
		task := newTask.GetBackpackTaskParts()
		task = services.SaveNewTaskParts(task)
		fmt.Printf("Generated Task with ID: %d\n", task.ID)
	}
	//service, err := services.RegisterService()
	//services.FailOnError(err, "Error on RegisterService")
	//fmt.Scanln()
	//
	//defer service.DeregisterServices()
}
