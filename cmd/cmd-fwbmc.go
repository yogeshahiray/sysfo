package cmd

import (
	"github.com/yogeshahiray/sysfo/common"
	"github.com/yogeshahiray/sysfo/table"
	"strings"
)

// Process FW-BMC command
func ProcessFwBMCCommand() {
	var rows = [][]interface{}{}
	var headers []string
	headers = append(headers, "IP", "Name", "Version")
	for _, e := range common.Systems {
		if e.Info.RedfishStatus {
			service := e.Info.Connection.Service
			s, err := service.UpdateService()
			if err != nil {
				panic(err)
			}
			fws, err := s.FirmwareInventories()
			if err != nil {
				panic(err)
			}

			for _, f := range fws {
				if e.Info.Reachable {
					substr := "iLO"
					str, substr := strings.ToUpper(f.Name), strings.ToUpper(substr)
					if strings.Contains(str, substr) {
						rows = append(rows, []interface{}{e.Ip, f.Name, f.Version})
					}
					substr = "iDRAC"
					str, substr = strings.ToUpper(f.Name), strings.ToUpper(substr)
					if strings.Contains(str, substr) {
						rows = append(rows, []interface{}{e.Ip, f.Name, f.Version})
					}
				}
			}
		}
	}
	table.WriteCommandOutput(rows, headers)
}
