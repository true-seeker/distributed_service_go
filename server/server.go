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
		services.GenerateTask(*taskSize)
	}
	//service, err := services.RegisterService()
	//services.FailOnError(err, "Error on RegisterService")
	//fmt.Scanln()
	//
	//defer service.DeregisterServices()
}
