package main

import (
	"fmt"
	"server/services"
)

func main() {
	service, err := services.RegisterService()
	services.FailOnError(err, "Error on RegisterService")
	fmt.Scanln() // wait for Enter Key

	defer service.DeregisterServices()
}
