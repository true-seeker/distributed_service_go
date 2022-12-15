package services

import (
	"net"
)

func GetNetworkAddresses() []net.IP {
	ifaces, err := net.Interfaces()
	FailOnError(err, "Error on getting net.Interfaces")

	var ips []net.IP

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		FailOnError(err, "Error on getting i.Addrs")
		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPAddr:
				if !v.IP.IsLoopback() {
					ips = append(ips, v.IP)
				}

			case *net.IPNet:
				if !v.IP.IsLoopback() && v.IP.To4() != nil {
					ips = append(ips, v.IP)
				}
			}

		}
	}
	return ips
}
