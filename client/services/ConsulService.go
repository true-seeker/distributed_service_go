package services

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
)

var AvailableServices []consul.AgentService

func GetAvailableServices() {
	c, err := consul.NewClient(&consul.Config{
		Address: fmt.Sprintf("%s:%s",
			GetProperty("Consul", "address"),
			GetProperty("Consul", "port")),
	})
	FailOnError(err, "GetAvailableServices consul client creation failed")

	services, err := c.Agent().Services()
	FailOnError(err, "Get services error")

	AvailableServices = nil

	for _, value := range services {
		AvailableServices = append(AvailableServices, *value)
	}
}
