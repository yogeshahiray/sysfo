// Implements system command
package cmd

import (
	"fmt"
	"github.com/yogeshahiray/sysfo/common"
)

func ProcessSystemCommand() {
	for _, e := range common.Systems {
		if e.Info.RedfishStatus {
			service := e.Info.Connection.Service
			chassis, err := service.Chassis()
			fmt.Println(chassis)
			if err != nil {
				panic(err)
			}
			for _, chass := range chassis {
				fmt.Printf("Chassis: %#v\n\n", chass)
			}

		}
	}

}
