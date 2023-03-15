package main

import (
	"github.com/yogeshahiray/sysfo/cmd"
	"github.com/yogeshahiray/sysfo/common"
	"github.com/yogeshahiray/sysfo/log"
)

func main() {
	log.LogInit()
	common.ProcessInputParams()
	cmd.ProcessCommand(common.InputParams.Cmd)
	//	common.Init("10.20.30.30", "10.20.30.40", 2)

	//ip := "192.168.36.150"
	//ver, err := cmd.GetRedfishVersion(ip, "admin", "redhat123")
	//if err != nil {
	//	fmt.Println("Unable to get version.... ", err)
	//} else {
	//	cmd.PrintRedfishVersion(ip, ver)
	//}
}
