// This file will implement nics command
package cmd

import (
	"encoding/json"
	"fmt"
	common2 "github.com/stmcginnis/gofish/common"
	"github.com/yogeshahiray/sysfo/common"
	"github.com/yogeshahiray/sysfo/log"
	"github.com/yogeshahiray/sysfo/table"
	"io"
	"strconv"
)

func getNetworkAdaptersInfo(info common.Sysinfo) string {
	s := ""
	n := 1

	for _, e := range info.NetworkAdapters.Adapters {
		s += strconv.Itoa(n) + "." + "\"" + e.Name + "\"" + " "
		n++
	}
	/*
		for _, e := range info.PCIDevices.Devices {
			if e.DeviceType == "NIC" {
				s += strconv.Itoa(n) + "." + "\"" + e.Name + "\"" + " "
				n++
			}
		} */
	return s
}

func getDellNetworkAdapters(entity common2.Entity, info *common.Sysinfo, ID string) {
	// First get the adapter count
	url := "/redfish/v1/Chassis/" + ID + "/NetworkAdapters"
	//fmt.Println("URL prepared : ", url)
	resp, err := entity.Client.Get(url)
	if err != nil {
		log.Fatal("Failed to get Network adapters : error [%s]", err)
	} else {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Failed to read response bytes")
		}
		json.Unmarshal(bodyBytes, &info.NetworkAdapters.Urls)
	}
	//Get the details of each adapter
	for _, e := range info.NetworkAdapters.Urls.Members {
		dev := common.NetworkAdapter{}
		log.Info("Adapter URL := %s", e.AdapterUrl)
		resp, err := entity.Client.Get(e.AdapterUrl)
		if err != nil {
			log.Fatal("Failed to get network adapter details")
		}
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Failed to read response bytes")
		}
		err = json.Unmarshal(bodyBytes, &dev)
		if err != nil {
			log.Fatal("%v", err)
		}
		info.NetworkAdapters.Adapters = append(info.NetworkAdapters.Adapters, dev)
	}

}

//func getDellNetworkAdapters(entity common2.Entity, info *common.Sysinfo, ID string) {
//	// In case of Dell, get the network interfaces first
//	n, err := info.System.NetworkInterfaces()
//	if err != nil {
//		log.Fatal("Failed to get network devices")
//	}
//	for _, e := range n {
//		dev := common.NetworkAdapter{}
//		ports, err := e.NetworkPorts()
//		dev.Name = e.Name
//		if err != nil {
//			log.Fatal("Failed to get network devices")
//		}
//		for _, p := range ports {
//			port := common.NetworkPort{}
//			port.LinkStatus = string(p.LinkStatus)
//			port.MacAddress = p.AssociatedNetworkAddresses[0]
//			port.Status.Health = string(p.Status.Health)
//
//			dev.PhysicalPorts = append(dev.PhysicalPorts, port)
//
//		}
//		info.NetworkAdapters.Adapters = append(info.NetworkAdapters.Adapters, dev)
//	}
//	// For Dell, look into PCI devices as well
//	for _, e := range info.PCIDevices.Devices {
//		if e.DeviceClass == "NetworkController" {
//			dev := common.NetworkAdapter{}
//
//			dev.Name = e.Name
//			info.NetworkAdapters.Adapters = append(info.NetworkAdapters.Adapters, dev)
//		}
//	}
//
//}

func getHPENetworkAdapters(entity common2.Entity, info *common.Sysinfo, ID string) {
	// First get the adapter count
	url := "/redfish/v1/systems/" + ID + "/BaseNetworkAdapters"
	//	fmt.Println("URL prepared : ", url)
	resp, err := entity.Client.Get(url)
	if err != nil {
		log.Fatal("Failed to get Network adapters")
	} else {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Failed to read response bytes")
		}
		json.Unmarshal(bodyBytes, &info.NetworkAdapters.Urls)
	}
	//Get the details of each adapter
	for _, e := range info.NetworkAdapters.Urls.Members {
		dev := common.NetworkAdapter{}
		resp, err := entity.Client.Get(e.AdapterUrl)
		if err != nil {
			log.Fatal("Failed to get network adapter details")
		}
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Failed to read response bytes")
		}
		err = json.Unmarshal(bodyBytes, &dev)
		if err != nil {
			log.Fatal("%v", err)
		}
		info.NetworkAdapters.Adapters = append(info.NetworkAdapters.Adapters, dev)
	}

}

//Get the list of Network Adapters
func getAllNetworkAdapters(entity common2.Entity, info *common.Sysinfo, ID string) {
	switch info.Mnf.Vendor {
	case "Dell":
		getDellNetworkAdapters(entity, info, ID)
	case "HPE":
		getHPENetworkAdapters(entity, info, ID)
	}
}
func pullNicInventory() {
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
					fmt.Println("More than one system exist")
				}
				common.Systems[i].Info.SerialNo = chass.SerialNumber
				common.Systems[i].Info.Mnf.Model = sys[0].Model
				getAllNetworkAdapters(sys[0].Entity, &common.Systems[i].Info, sys[0].ID)
				log.Info("Adapters after getAllNetworkAdapters:= %d", len(common.Systems[i].Info.NetworkAdapters.Adapters))
			}
		}
	}
}

// Print/Write the output
func processNICsCommandOutput() {
	var rows = [][]interface{}{}
	var header = false
	for _, e := range common.Systems {
		header = false
		if e.Info.Reachable && e.Info.RedfishStatus {
			//	log.Info("Adapter now here := %d", len(e.Info.NetworkAdapters.Adapters))
			for _, a := range e.Info.NetworkAdapters.Adapters {
				ports := strconv.Itoa(len(a.PhysicalPorts))
				macs := ""
				for i, m := range a.PhysicalPorts {
					if i == len(a.PhysicalPorts) {
						macs += m.MacAddress
					} else {
						macs += m.MacAddress + ","
					}
					//log.Info("Mac so far := [%s]", macs)
				}
				if header == false {
					rows = append(rows, []interface{}{e.Ip, e.Info.SerialNo, e.Info.Mnf.Model,
						a.Name, a.Location, a.Firmware, ports, macs})
					header = true
				} else {
					rows = append(rows, []interface{}{" ", " ", " ",
						a.Name, a.Location, a.Firmware, ports, macs})
				}

			}

		} else {
			//	rows = append(rows, []interface{}{e.Ip, "NA", "NA", "NA", "NA", "NA", "NA", "Down"})
			log.Info("No systems detected")
		}
	}
	table.WriteCommandOutput(rows, []string{"BMC", "SN", "Model", "NIC", "Location", "Fw Ver", "# Ports", "MAC(s)"})

}
func ProcessNICsCommand() {
	// Get details of all NICs
	GetRedfishVersion()
	//pullNicInventory()
	processNICsCommandOutput()

}
