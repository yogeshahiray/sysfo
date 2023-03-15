// This file includes function(s) to ping the remote host
package common

import (
	"net"
	"time"
)

func PingRemoteHost(ip string) bool {
	port := "80"
	timeout := time.Duration(1 * time.Second)
	_, err := net.DialTimeout("tcp", ip+":"+port, timeout)

	return err == nil
}
