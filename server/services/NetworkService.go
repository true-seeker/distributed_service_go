package services

import (
	"net"
)

func GetNetworkAddresses() []net.IP {
	var ips []net.IP

	isDocker := GetProperty("Docker", "is_docker")
	if isDocker == "true" {
		ip, _, err := net.ParseCIDR(GetProperty("Docker", "host_ip"))
		FailOnError(err, "Error on parsing host address")

		ips = append(ips, ip)
		return ips
	}

	ifaces, err := net.Interfaces()
	FailOnError(err, "Error on getting net.Interfaces")

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
				if !v.IP.IsLoopback() && !v.IP.IsPrivate() && v.IP.To4() != nil {
					ips = append(ips, v.IP)
				}
			}

		}
	}
	return ips
}
