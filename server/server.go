package main

import (
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

	services.Migrate()

	service, err := services.RegisterService()
	services.FailOnError(err, "Failed to register consul service")
	defer service.DeregisterServices()

	fmt.Sprintf("Successfully registered Consul service with name %s", service.Name)

	go services.StartGRPCListener()
	go services.StartWebServerListener()

	select {
	case <-c:
		service.DeregisterServices()
	}
}
