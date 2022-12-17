package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"server/services"
	"syscall"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	generateTask := flag.Bool("g", false, "generate new task")
	taskSize := flag.Int("s", services.DefaultTaskSize, "new task size")
	flag.Parse()
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("%s: %s\n", f.Name, f.Value)
	})

	services.Migrate()

	service, err := services.RegisterService()
	if len(service.RegisteredServices) == 0 {
		fmt.Println("Cant register any Consul service. Exiting")
		os.Exit(0)
	}

	defer service.DeregisterServices()
	fmt.Sprintf("Successfully registered Consul service with name %s", service.Name)
	services.FailOnError(err, "Failed to register consul service")

	if *generateTask {
		services.GenerateTask(*taskSize)
	}

	go services.StartGRPCListener()
	go services.StartWebServerListener()

	select {
	case <-c:
		service.DeregisterServices()
	}

}
