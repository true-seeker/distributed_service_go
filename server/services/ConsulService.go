package services

import (
	"fmt"
	"github.com/google/uuid"
	consul "github.com/hashicorp/consul/api"
	"log"
	"strconv"
	"strings"
	"time"
)

// Service Консул сервис
type Service struct {
	Name               string
	TTL                time.Duration
	ConsulAgent        *consul.Agent
	RegisteredServices []consul.AgentServiceRegistration
}

// RegisterService Регистрация сервиса
func RegisterService() *Service {
	s := Service{
		Name: "BackpackServer",
	}

	c, err := consul.NewClient(&consul.Config{
		Address: fmt.Sprintf("%s:%s",
			GetProperty("Consul", "address"),
			GetProperty("Consul", "port")),
	})
	FailOnError(err, "RegisterService consul client creation failed")

	s.ConsulAgent = c.Agent()

	//Получаем все IP адреса текущей машины
	ips := GetNetworkAddresses()
	if len(ips) == 0 {
		log.Fatalln("No available interfaces to host the server")
	}

	// Перебираем все возможные адреса и на каждом из них регистрируем консул сервис
	for _, ip := range ips {
		servicePort, err := strconv.Atoi(GetProperty("gRPC", "server_port"))
		FailOnError(err, "Cant get port from config")

		serviceDef := &consul.AgentServiceRegistration{
			ID:   uuid.Must(uuid.NewRandom()).String(),
			Name: s.Name,
			Tags: []string{
				"BackpackServer",
			},
			Address: ip.String(),
			Port:    servicePort,
			Check: &consul.AgentServiceCheck{
				HTTP:     fmt.Sprintf("http://%s:%s/healthCheck", ip, GetProperty("Consul", "http_check_port")),
				Interval: "10s",
			},
		}

		isFirstMessage := true
		err = s.ConsulAgent.ServiceRegister(serviceDef)
		for err != nil {
			if strings.Contains(err.Error(), "connection refused") {
				if isFirstMessage {
					fmt.Println("Нет подключения к серверу Consul. Переподключаемся")
					isFirstMessage = false
				}
				time.Sleep(5 * time.Second)
				err = s.ConsulAgent.ServiceRegister(serviceDef)
			}
		}
		s.RegisteredServices = append(s.RegisteredServices, *serviceDef)
	}

	if len(s.RegisteredServices) == 0 {
		FailOnError(nil, "Cant register any Consul service. Exiting")
	}

	return &s
}

// DeregisterServices Дерегистрация сервиса Консула
func (s *Service) DeregisterServices() {
	for _, serviceDef := range s.RegisteredServices {
		err := s.ConsulAgent.ServiceDeregister(serviceDef.ID)
		FailOnError(err, "Error on deregisterService "+serviceDef.ID)
		fmt.Printf("Сервис с ID %s дерегистрирован\n", serviceDef.ID)
	}
}
