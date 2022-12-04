package services

import (
	"fmt"
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
					fmt.Printf("%v : %s\n", i.Name, v)
					ips = append(ips, v.IP)
				}

			case *net.IPNet:
				if !v.IP.IsLoopback() && v.IP.To4() != nil {
					fmt.Printf("%v : %s\n", i.Name, v.IP.To4())
					ips = append(ips, v.IP)
				}
			}

		}
	}
	return ips
}
