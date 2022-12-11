package main

import (
	"flag"
	"fmt"
	"math/rand"
	"server/services"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	generateTask := flag.Bool("g", false, "generate new task")
	taskSize := flag.Int("s", services.DefaultTaskSize, "new task size")
	flag.Parse()
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("%s: %s\n", f.Name, f.Value)
	})

	services.Migrate()

	if *generateTask {
		services.GenerateTask(*taskSize)
	}

	go services.StartGRPCListener()
	select {}

	//service, err := services.RegisterService()
	//services.FailOnError(err, "Error on RegisterService")
	//fmt.Scanln()
	//
	//defer service.DeregisterServices()
}
