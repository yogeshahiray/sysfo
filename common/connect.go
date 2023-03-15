// This package will include common code
package common

import (
	"errors"
	"github.com/stmcginnis/gofish"
)

func GetConnection(ip string, username string, password string) error {
	// Check if we can reach remote host
	var s System

	s.Ip = ip
	s.Info.RedfishStatus = false
	if !PingRemoteHost(ip) {
		Systems = append(Systems, s)
		return errors.New("Unable to reach remote host")
	} else {
		s.Info.Reachable = true
	}

	config := gofish.ClientConfig{
		Endpoint: "https://" + ip,
		Username: username,
		Password: password,
		Insecure: true,
	}
	c, err := gofish.Connect(config)
	if err != nil {
		return errors.New("Unable to connect to remote system")
	}

	// Host is reachable, so add it in the list
	s.Info.RedfishStatus = true
	s.Info.Connection = c
	Systems = append(Systems, s)

	return nil
}

func TearDown() {
	for _, e := range Systems {
		if e.Info.RedfishStatus {
			e.Info.Connection.Logout()
		}
	}
}
