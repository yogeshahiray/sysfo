package cmd

import (
	"sysfo/common"
	"sysfo/table"
)

// Show firmware version
func ProcessFwCommand() {
	var rows = [][]interface{}{}
	var headers []string

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
			headers = append(headers, e.Ip, "Version")
			for _, f := range fws {
				vals := []interface{}{}
				if e.Info.Reachable {
					vals = append(vals, f.Name, f.Version)
					rows = append(rows, vals)
				}
			}
		}
	}
	table.WriteCommandOutput(rows, headers)
}
