package services

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"strings"
	"time"
)

var AvailableServices []consul.AgentService

// GetAvailableServices Получение доступных сервисов из консула
func GetAvailableServices() {
	c, err := consul.NewClient(&consul.Config{
		Address: fmt.Sprintf("%s:%s",
			GetProperty("Consul", "address"),
			GetProperty("Consul", "port")),
	})
	FailOnError(err, "GetAvailableServices consul client creation failed")

	isFirstMessage := true

	services, err := c.Agent().Services()
	for err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			if isFirstMessage {
				fmt.Println("Нет подключения к серверу Consul. Переподключаемся")
				isFirstMessage = false
			}
			services, err = c.Agent().Services()
			time.Sleep(5 * time.Second)
		}
	}

	AvailableServices = nil

	for _, value := range services {
		AvailableServices = append(AvailableServices, *value)
	}
}
