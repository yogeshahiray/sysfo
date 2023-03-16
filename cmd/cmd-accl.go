package cmd

import (
	common2 "github.com/stmcginnis/gofish/common"
     "sysfo/common"
	"sysfo/log"
	"sysfo/table"
)

//This function collects the information of Acclerators installed in the system
func ProcessAcclCommand() {
	GetRedfishVersion()
	pullAcclInventory()
	processAcclsCommandOutput()

}
func checkIfAccl(dev common.PCIDevice, adp *common.AcclAdapter) bool {
	for _, e := range common.KnownAcclAdapters {
		if dev.DeviceId == e.Id {
			adp.Desc = e.Desc
			adp.Name = e.Name
			return true
		}
	}
	return false
}
func getAcclAdapters(entity common2.Entity, info *common.Sysinfo, ID string) {
	//In this function, all we need to do is iterate through the PCI devices and see if get any accl
	for _, e := range info.PCIDevices.Devices {
		dev := common.AcclAdapter{}
		if checkIfAccl(e, &dev) {
			dev.Name = e.Name + "-" + dev.Name
			dev.DeviceID = e.DeviceId
			dev.Id = e.Id
			info.AllAcclAdapters = append(info.AllAcclAdapters, dev)
		}
	}
}

//Get the list of Accls
func getAllAccls(entity common2.Entity, info *common.Sysinfo, ID string) {
	getAcclAdapters(entity, info, ID)

}

func pullAcclInventory() {
	for i, e := range common.Systems {
		if e.Info.RedfishStatus {
			service := e.Info.Connection.Service
			chassis, err := service.Chassis()
			if err != nil {
				panic(err)
			}
			for _, chass := range chassis {
				sys, _ := chass.ComputerSystems()
				if len(sys) > 1 {
					log.Info("More than one system exist")
				}
				if len(sys) <= 0 {
					continue
				}
				common.Systems[i].Info.SerialNo = chass.SerialNumber
				common.Systems[i].Info.Mnf.Model = sys[0].Model
				getAcclAdapters(sys[0].Entity, &common.Systems[i].Info, sys[0].ID)
			}
		}
	}
}

// Print/Write the output
func processAcclsCommandOutput() {
	var rows = [][]interface{}{}
	for _, e := range common.Systems {
		if e.Info.Reachable && e.Info.RedfishStatus {
			if len(e.Info.AllAcclAdapters) <= 0 {
				rows = append(rows, []interface{}{e.Ip, e.Info.SerialNo, e.Info.Mnf.Model, len(e.Info.AllAcclAdapters), "None"})
			} else {
				for _, a := range e.Info.AllAcclAdapters {
					data := ""
					data += a.Name + "(" + a.Desc + ")" + "." + " "
					rows = append(rows, []interface{}{e.Ip, e.Info.SerialNo, e.Info.Mnf.Model, len(e.Info.AllAcclAdapters), data})
				}
			}

		}
	}
	if len(rows) > 0 {
		table.WriteCommandOutput(rows, []string{"BMC", "SN", "Model", "# Accls", "Accelerators"})
	} else {
		log.Info("No data to display")
	}

}
