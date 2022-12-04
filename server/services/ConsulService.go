package services

import (
	"fmt"
	"time"

	consul "github.com/hashicorp/consul/api"
)

type Service struct {
	Name               string
	TTL                time.Duration
	ConsulAgent        *consul.Agent
	RegisteredServices []consul.AgentServiceRegistration
}

func RegisterService() (*Service, error) {
	s := Service{
		Name: "hehe",
	}

	c, err := consul.NewClient(&consul.Config{
		Address: fmt.Sprintf("%s:%s",
			GetProperty("Consul", "address"),
			GetProperty("Consul", "port")),
	})
	FailOnError(err, "RegisterService consul client creation failed")

	s.ConsulAgent = c.Agent()
	ips := GetNetworkAddresses()

	isServiceRegistered := false
	var serviceRegistrationErrors []error
	for i, ip := range ips {
		serviceDef := &consul.AgentServiceRegistration{
			ID:   fmt.Sprintf("%s_%d", s.Name, i),
			Name: s.Name,
			Tags: []string{
				"BackpackServer",
			},
			Address: ip.String(),
			Port:    80,
			Check: &consul.AgentServiceCheck{
				HTTP:     "https://google.com",
				Interval: "10s",
			},
		}
		err = s.ConsulAgent.ServiceRegister(serviceDef)
		if err == nil {
			isServiceRegistered = true
			s.RegisteredServices = append(s.RegisteredServices, *serviceDef)
		} else {
			serviceRegistrationErrors = append(serviceRegistrationErrors, err)
		}
	}

	if !isServiceRegistered {
		FailOnError(serviceRegistrationErrors[0], "Error with service register")
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
