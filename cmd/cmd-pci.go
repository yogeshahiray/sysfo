package cmd

import (
	"encoding/json"
	common2 "github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"
	"github.com/yogeshahiray/sysfo/common"
	"github.com/yogeshahiray/sysfo/log"
	"io"
)

// This function will check if the PCI function already part of the inventory
func isPCIFnExists(info *common.Sysinfo, name string) bool {
	for _, e := range info.PCIDevices.Devices {
		if name == e.Name {
			return true
		}
	}
	return false
}

// Dell Redfish implementation is little different
func getDellPCIDevices(entity common2.Entity, info *common.Sysinfo, system *redfish.ComputerSystem) {
	// In case of Dell, first get the list of PCI devices URLs
	url := "/redfish/v1/Systems/" + system.ID

	resp, err := entity.Client.Get(url)
	if err != nil {
		log.Fatal("Failed to get PCI device urls")
	} else {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Failed to read response bytes")
		}
		json.Unmarshal(bodyBytes, &info.PCIDevices.DellPCIUrls)
	}

	// Get the details of each PCI device

	for _, e := range info.PCIDevices.DellPCIUrls.PCIeDevices {
		links := common.DellPCILink{}
		resp, err := entity.Client.Get(e.PCIDeviceUrl)
		if err != nil {
			log.Fatal("Failed to get PCI device details")
		} else {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Failed to read response bytes")
			}
			json.Unmarshal(bodyBytes, &links)
		}
		// Now, get the details of each PCI function
		for _, d := range links.DellPCIFnsUrls.PCIeFunctions {
			dev := common.PCIDevice{}
			resp, err := entity.Client.Get(d.PCIDeviceUrl)
			if err != nil {
				log.Fatal("Failed to get PCI device details")
			} else {
				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Fatal("Failed to read response bytes")
				}
				json.Unmarshal(bodyBytes, &dev)
				// Check if the function already added into inventory
				if !isPCIFnExists(info, dev.Name) {
					info.PCIDevices.Devices = append(info.PCIDevices.Devices, dev)
				}
			}
		}
	}

}

func getHPEPCIDevices(entity common2.Entity, info *common.Sysinfo, system *redfish.ComputerSystem) {
	url := "/redfish/v1/Systems/" + system.ID + "/PCIDevices"
	resp, err := entity.Client.Get(url)
	if err != nil {
		log.Fatal("Failed to get PCI devices")
	} else {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Failed to read response bytes")
		}
		json.Unmarshal(bodyBytes, &info.PCIDevices.Urls)
	}
	// Next, get the details of all PCI devices
	for _, e := range info.PCIDevices.Urls.Members {
		dev := common.PCIDevice{}
		resp, err := entity.Client.Get(e.PCIDeviceUrl)
		if err != nil {
			log.Fatal("Failed to get PCI device details")
		}
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Failed to read response bytes")
		}
		err = json.Unmarshal(bodyBytes, &dev)
		if err != nil {
			log.Fatal("%v", err)
		}
		info.PCIDevices.Devices = append(info.PCIDevices.Devices, dev)
	}
}

// Fetch the list of PCI devices
func getAllPCIDevices(entity common2.Entity, info *common.Sysinfo, system *redfish.ComputerSystem) {
	switch info.Mnf.Vendor {
	case "Dell":
		getDellPCIDevices(entity, info, system)
	case "HPE":
		getHPEPCIDevices(entity, info, system)
	}

}
