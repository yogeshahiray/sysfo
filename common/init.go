// Initializatiob if hosts
package common

import (
	"github.com/c-robinson/iplib"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/redfish"
	"net"
)

type NetworkAdaptersUrl struct {
	Count   int `json:"Members@odata.count"`
	Members []struct {
		AdapterUrl string `json:"@odata.id"`
	} `json:"Members"`
}
type NetworkPort struct {
	LinkStatus string `json:"LinkStatus"`
	MacAddress string `json:"MacAddress"`
	Status     struct {
		Health string `json:"Health"`
	} `json:"Status"`
}
type NetworkAdpFirmware struct {
	Current struct {
		VersionString string `json:"VersionString"`
	} `json:"Current"`
}

// For Dell iDRAC
type NetworkLinks struct {
	redfish.NetworkDeviceFunction
}
type NetworkControllers struct {
	FirmwarePackageVersion string `json:"FirmwarePackageVersion"`
}
type NetworkAdapter struct {
	Id            string             `json:"Id"`
	Firmware      NetworkAdpFirmware `json:"Firmware"`
	Location      string             `json:"Location"`
	Name          string             `json:"Name"`
	PhysicalPorts []NetworkPort      `json:"PhysicalPorts"`
	Status        struct {
		Health string `json:"Health"`
		State  string `json:"State"`
	} `json:"Status"`
	StructureName string `json:"StructureName"`
}

type AcclAdapter struct {
	Id             string
	Desc           string
	Name           string
	DeviceID       int
	DeviceLocation string
}

type CPUInfo struct {
	CPUModel        string
	CPUCount        int
	LogicalCPUCount int
	TotalCores      int
	TotalThreads    int
}

type MemInfo struct {
	NumDIMs     int
	TotalMemory float32
}

type Manufacturer struct {
	Vendor  string
	Product string
	Model   string
}
type Storage struct {
	Name        string
	NumDisks    int
	StorageSize int64
	DiskSize    []int64
}
type StorageInfo struct {
	DisksInfo  string
	TotalSize  int64
	TotalDisks int
	StorageIns []Storage //Storage instance
}

type DellPCIeFnsUrl struct {
	PCIDeviceUrl string `json:"@odata.id"`
}

type DellPCILinkEntry struct {
	PCIeFunctions []DellPCIeFnsUrl `json:"PCIeFunctions"`
	Count         int              `json:"PCIeFunctions@odata.count"`
}
type DellPCILink struct {
	DellPCIFnsUrls DellPCILinkEntry `json:"Links"`
}
type DellPCIDevicesUrl struct {
	PCIeDevices []struct {
		PCIDeviceUrl string `json:"@odata.id"`
	} `json:"PCIeDevices"`
	Count int `json:"PCIeDevices@odata.count"`
}
type PCIDevicesUrl struct {
	Count   int `json:"Members@odata.count"`
	Members []struct {
		PCIDeviceUrl string `json:"@odata.id"`
	} `json:"Members"`
}

type PCIDevice struct {
	Id                string `json:"Id"`
	BusNumber         int    `json:"BusNumber"`
	ClassCode         int    `json:"ClassCode"`
	DeviceId          int    `json:"DeviceId"`
	DeviceInstance    int    `json:"DeviceInstance"`
	DeviceLocation    string `json:"DeviceLocation"`
	DeviceClass       string `json:"DeviceClass"`
	DeviceNumber      int    `json:"DeviceNumber"`
	DeviceSubInstance int    `json:"DeviceSubInstance"`
	DeviceType        string `json:"DeviceType"`
	FunctionNumber    int    `json:"FunctionNumber"`
	LocationString    string `json:"LocationString"`
	Name              string `json:"Name"`
	StructuredName    string `json:"StructuredName"`
	SubsystemDeviceID int    `json:"SubsystemDeviceID"`
	SubsytemVendorID  int    `json:"SubsytemVendorID"`
	VendorID          int    `json:"VendorID"`
}
type AllNetworkAdapters struct {
	Urls     NetworkAdaptersUrl
	Adapters []NetworkAdapter
}
type AllPCIDevices struct {
	Urls        PCIDevicesUrl
	DellPCIUrls DellPCIDevicesUrl
	Devices     []PCIDevice
}
type Sysinfo struct {
	Connection      *gofish.APIClient
	PCIDevices      AllPCIDevices
	AllAcclAdapters []AcclAdapter
	NetworkAdapters AllNetworkAdapters
	SerialNo        string
	Reachable       bool
	RedfishStatus   bool
	RedfishVer      string
	Mnf             Manufacturer
	Memory          MemInfo
	CPU             CPUInfo
	NetworkInfo     string
	Storage         StorageInfo
	System          *redfish.ComputerSystem
	FwVersion       string
}

// Structure to hold system status information

type SystemStatus struct {
	BMCVersion  string
	PowerState  string
	BIOSVersion string
}

type System struct {
	ID        string
	Ip        string
	Outfile   string
	ErrorCode error
	Info      Sysinfo
	Status    SystemStatus
}

var Systems = []System{}
var KnownAcclAdapters = []struct {
	Id   int
	Name string
	Desc string
}{
	{3420,
		"ACC100",
		"FEC accel"},
}

func Init(startIp string, endIp string, step int, username string, password string) {
	ipa := net.ParseIP(startIp)
	ipb := net.ParseIP(endIp)

	for iplib.CompareIPs(ipa, ipb) <= 0 {
		GetConnection(ipa.String(), username, password)
		ipa = iplib.IncrementIPBy(ipa, uint32(step))
	}

}
