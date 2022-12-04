package services

import (
	"encoding/json"
	"fmt"
	consul "github.com/hashicorp/consul/api"
)

func GetAvailableServices() {
	c, err := consul.NewClient(&consul.Config{
		Address: fmt.Sprintf("%s:%s",
			GetProperty("Consul", "address"),
			GetProperty("Consul", "port")),
	})
	FailOnError(err, "GetAvailableServices consul client creation failed")

	services, err := c.Agent().Services()
	FailOnError(err, "Get services error")

	var agentServices []consul.AgentService
	for _, value := range services {
		a, _ := json.Marshal(value)
		fmt.Println(string(a))
		agentServices = append(agentServices, *value)
	}

	fmt.Println(agentServices)
}
