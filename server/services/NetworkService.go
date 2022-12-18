package services

import (
	"net"
)

// GetNetworkAddresses Получение всех возможных IP адресов
func GetNetworkAddresses() []net.IP {
	var ips []net.IP

	ifaces, err := net.Interfaces()
	FailOnError(err, "Error on getting net.Interfaces")

	// Перебираем все интерфейсы
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		FailOnError(err, "Error on getting i.Addrs")
		for _, a := range addrs {
			switch v := a.(type) {
			// Отбрасываем localhost и приватные адреса
			case *net.IPAddr:
				if !v.IP.IsLoopback() && !v.IP.IsPrivate() && v.IP.To4() != nil {
					ips = append(ips, v.IP)
				}

			case *net.IPNet:
				if !v.IP.IsLoopback() && !v.IP.IsPrivate() && v.IP.To4() != nil {
					ips = append(ips, v.IP)
				}
			}
		}
	}
	return ips
}
