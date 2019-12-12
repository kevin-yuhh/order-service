package utils

import (
	"testing"
)

func TestLocalIPv4s(t *testing.T) {
	ipAddress := GetLocalIpAddress()
	t.Log(ipAddress)
}
