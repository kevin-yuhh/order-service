package utils

import (
	"net"
)

const defaultAddress = "unknown host"

var IpAddress string

func init() {
	IpAddress = localIPv4s()
}

// LocalIPs return all non-loopback IPv4 addresses
func localIPv4s() string {
	var ips []string
	address, err := net.InterfaceAddrs()
	if err != nil {
		return defaultAddress
	}

	for _, a := range address {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			ips = append(ips, ipNet.IP.String())
		}
	}

	if ips == nil || len(ips) == 0 {
		return defaultAddress
	}

	return ips[0]
}

func GetLocalIpAddress() string {
	return IpAddress
}
