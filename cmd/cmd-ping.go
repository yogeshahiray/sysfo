// This command will reach out to the system and get the version and other details of redfish
package cmd

import (
	"fmt"
	"sysfo/common"
	"sysfo/table"
)

// This function will populate the information of redfish version of all hosts
func GetRedfishVersionOfAll(ip string, u string, p string) {
	//for s := range common.Systems {
	//
	//}
}

func ProcessPingCommandOutput() {
	var rows = [][]interface{}{}
	fmt.Println("Systems : ", len(common.Systems))
	for _, e := range common.Systems {
		if e.Info.Reachable {
			rows = append(rows, []interface{}{e.Ip, "Up"})
		} else {
			rows = append(rows, []interface{}{e.Ip, "Down"})
		}
	}

	if len(common.InputParams.OutFile) > 0 {
		table.WriteTable(rows, []string{"IP", "Status"}, common.InputParams.OutFile)

	} else {
		table.PrintTable(rows, []string{"IP", "Status"})
	}

}
