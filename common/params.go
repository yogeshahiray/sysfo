// This file will process the input parameters
package common

import (
	"flag"
	"strings"
)

type Params struct {
	Username string
	Password string
	StartIp  string
	EndIp    string
	Step     int
	Cmd      string
	OutFile  string
	CmdArgs  string
}

var InputParams = Params{
	Step: 1,
}

// Process input parameters
func ProcessInputParams() bool {
	step := flag.Int("s", 1, "IP addr increment step")
	ipaddr := flag.String("i", "", "IP address range (x.x.x.x-x.x.x.x or x.x.x.x)")
	command := flag.String("c", "ping", "Command from [ping, status, system, nics, accls]")
	username := flag.String("u", "", "Username")
	password := flag.String("p", "", "Password")
	outfile := flag.String("o", "", "Output file")
	cmdargs := flag.String("a", "", "Command argument(s)")

	flag.Parse()
	ips := strings.Split(*ipaddr, "-")
	ipa := ""
	ipb := ""
	if len(ips) <= 1 {
		ipa = *ipaddr
		ipb = ipa
	} else {
		ipa = ips[0]
		ipb = ips[1]
	}

	if !IsValidIP(ipa) {
		PrintError("Please enter valid IP address")
	}
	// See if we got range of ip addresses
	if len(ipb) > 0 && !IsValidIP(ipb) {
		PrintError("Please enter valid IP address")
	}
	if !IsValidIPRange(ipa, ipb) {
		PrintError("Please enter valid range of IP addresses")
	}

	InputParams.Step = *step
	InputParams.CmdArgs = *cmdargs
	if InputParams.Step > 254 {
		PrintError("Please enter valid increment step")
	}
	// At this point, IP addresses are valid
	InputParams.StartIp = ipa
	InputParams.EndIp = ipb

	InputParams.Cmd = *command

	//Check if username and password is provided
	if *username == "" {
		PrintError("Please enter username")
	}
	if *password == "" {
		PrintError("Please enter password")
	}
	InputParams.Username = *username
	InputParams.Password = *password
	InputParams.OutFile = *outfile
	// Everything looks good here, go ahead and init all hosts
	Init(InputParams.StartIp, InputParams.EndIp, InputParams.Step, InputParams.Username, InputParams.Password)
	return true
}
