// Implementation of redfish status command
package cmd

import (
	"fmt"
	"github.com/yogeshahiray/sysfo/common"
	"github.com/yogeshahiray/sysfo/log"
	"github.com/yogeshahiray/sysfo/table"
	"strconv"
)

// Function to get the version of redfish
func GetRedfishVersion() {
	// Get connection

	if len(common.Systems) <= 0 {
		log.Fatal("No systems found")
	}
	for i, e := range common.Systems {
		if e.Info.RedfishStatus {
			service := e.Info.Connection.Service
			common.Systems[i].Info.Mnf.Vendor = service.Vendor
			common.Systems[i].Info.Mnf.Product = service.Product
			chassis, err := service.Chassis()
			if err != nil {
				panic(err)
			}
			fmt.Println("Chassis : ", len(chassis))
			for _, chass := range chassis {
				sys, _ := chass.ComputerSystems()
				if len(sys) > 1 {
					log.Info("More than one system exist")
				}
				if len(sys) == 0 {
					continue
				}
				fmt.Println("Sys ID : ", sys[0].ID)
				common.Systems[i].ID = sys[0].ID
				common.Systems[i].Info.System = sys[0]
				getAllPCIDevices(sys[0].Entity, &common.Systems[i].Info, sys[0])
				getAllNetworkAdapters(sys[0].Entity, &common.Systems[i].Info, sys[0].ID)
				common.Systems[i].Info.SerialNo = chass.SerialNumber
				common.Systems[i].Info.Memory.TotalMemory = sys[0].MemorySummary.TotalSystemMemoryGiB
				common.Systems[i].Info.Mnf.Model = sys[0].Model
				// Fill out CPU details
				common.Systems[i].Info.CPU.CPUModel = sys[0].ProcessorSummary.Model
				common.Systems[i].Info.CPU.CPUCount = sys[0].ProcessorSummary.Count
				common.Systems[i].Info.CPU.LogicalCPUCount = sys[0].ProcessorSummary.LogicalProcessorCount

				common.Systems[i].Info.CPU.TotalCores = 0
				common.Systems[i].Info.CPU.TotalThreads = 0

				cpus, err := sys[0].Processors()
				if err != nil {
					log.Fatal("Failed to get CPU details")
				}

				for _, c := range cpus {
					common.Systems[i].Info.CPU.TotalCores += c.TotalCores
					common.Systems[i].Info.CPU.TotalThreads += c.TotalThreads
				}

				common.Systems[i].Status.BIOSVersion = sys[0].BIOSVersion
				common.Systems[i].Status.PowerState = string(sys[0].PowerState)

				// Capture storage information
				common.Systems[i].Info.Storage.TotalSize = 0
				common.Systems[i].Info.Storage.TotalDisks = 0
				stg, err := sys[0].Storage()
				if err != nil {
					log.Fatal("Failed to get storage information")
				}
				//	fmt.Println("Storage : ", len(stg))
				dinfo := " ("
				for _, e := range stg {
					s := common.Storage{}
					s.Name = e.Name
					s.NumDisks = e.DrivesCount
					common.Systems[i].Info.Storage.TotalDisks += s.NumDisks
					disks, err := e.Drives()
					if err != nil {
						log.Fatal("Failed to get drives list")
					}
					for _, d := range disks {
						s.DiskSize = append(s.DiskSize, d.CapacityBytes)
						common.Systems[i].Info.Storage.TotalSize += d.CapacityBytes
						dinfo += strconv.FormatInt(int64(d.CapacityBytes/1000000000), 10) + "GB "
					}
					dinfo += ")"
					//fmt.Println("Diskinfo :", dinfo)
					common.Systems[i].Info.Storage.StorageIns = append(common.Systems[i].Info.Storage.StorageIns, s)
					common.Systems[i].Info.Storage.DisksInfo = dinfo
				}
				common.Systems[i].Info.Storage.TotalSize = common.Systems[i].Info.Storage.TotalSize / 1000000000
				common.Systems[i].Info.NetworkInfo = getNetworkAdaptersInfo(common.Systems[i].Info)

				/*
					eth, _ := sys[0].EthernetInterfaces()
					inf, _ := sys[0].NetworkInterfaces()

					adp, _ := inf[0].NetworkAdapter()
					ports, _ := adp.NetworkPorts()

					fmt.Printf("Info:  %v\n\n", eth[0].LinkStatus)
					fmt.Printf("Info:  %v\n\n", eth[1].MACAddress)
					fmt.Printf("Info:  %v\n\n", eth[2].SpeedMbps)
					fmt.Printf("Info:  %v\n\n", adp.Manufacturer)
					fmt.Printf("Info:  %v\n\n", ports[1].VendorID)
					fmt.Printf("Info:  %v\n\n", ports[1].Name) // E810 Network port
					fmt.Printf("Info:  %v\n\n", ports[1].ID)
					fmt.Printf("Info:  %v\n\n", ports[1].AssociatedNetworkAddresses)
				*/
			}
			common.Systems[i].Info.RedfishVer = service.RedfishVersion
		}
	}

}

func WriteRedfishStatusOutput() {

}

func ProcessRedfishStatusOutput() {
	var rows = [][]interface{}{}
	fmt.Println("Systems : ", len(common.Systems))
	for _, e := range common.Systems {
		if e.Info.Reachable && e.Info.RedfishStatus {
			mem := fmt.Sprintf("%.0f", e.Info.Memory.TotalMemory)
			rows = append(rows, []interface{}{e.Ip, e.Info.SerialNo, e.Info.Mnf.Model, e.Info.CPU.CPUModel, strconv.Itoa(e.Info.CPU.CPUCount) + "/" + strconv.Itoa(e.Info.CPU.TotalCores) + "/" + strconv.Itoa(e.Info.CPU.TotalThreads),
				mem, strconv.Itoa(e.Info.Storage.TotalDisks) + e.Info.Storage.DisksInfo, strconv.FormatInt(e.Info.Storage.TotalSize, 10), e.Info.NetworkInfo, e.Status.PowerState})
		} else {
			rows = append(rows, []interface{}{e.Ip, "NA", "NA", "NA", "NA", "NA", "NA", "Down"})
		}
	}
	table.WriteCommandOutput(rows, []string{"BMC", "SN", "Model", "CPU Model", "CPUs/Cores/Threads", "Mem(GiB)", "# Disks", "Storage(GB)", "NICs", "Power"})

}

func ProcessRedfishStatusCommand() {
	GetRedfishVersion()
	ProcessRedfishStatusOutput()
}
