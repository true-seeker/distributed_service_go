package services

import (
	"fmt"
	"github.com/google/uuid"
	consul "github.com/hashicorp/consul/api"
	"os"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	Name               string
	TTL                time.Duration
	ConsulAgent        *consul.Agent
	RegisteredServices []consul.AgentServiceRegistration
}

func RegisterService() (*Service, error) {
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
	ips := GetNetworkAddresses()
	if len(ips) == 0 {
		fmt.Println("No available interfaces to host the server")
		os.Exit(0)
	}

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
		err = s.ConsulAgent.ServiceRegister(serviceDef)
		if err == nil {
			s.RegisteredServices = append(s.RegisteredServices, *serviceDef)
		} else {
			if strings.Contains(err.Error(), "connection refused") {
				fmt.Println("Cant connect to Consul. Exiting")
				os.Exit(0)
			}
		}
	}

	if len(s.RegisteredServices) == 0 {
		FailOnError(nil, "Cant register any Consul service. Exiting")
	}

	return &s, nil
}

func (s *Service) DeregisterServices() {
	for _, serviceDef := range s.RegisteredServices {
		err := s.ConsulAgent.ServiceDeregister(serviceDef.ID)
		FailOnError(err, "Error on deregisterService "+serviceDef.ID)
		fmt.Printf("Service with ID %s deregistered\n", serviceDef.ID)
	}
}
