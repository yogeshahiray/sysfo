// Command processing
package cmd

import (
	"github.com/yogeshahiray/sysfo/common"
)

// Command structure
type CommandInfo struct {
	Headers      []string
	DefaultValue []string
}

var Commands = map[string]CommandInfo{
	"ping":    {[]string{"Ping"}, []string{"Down"}},
	"redfish": {[]string{"RedfishVer"}, []string{"NA"}},
}

// Process command to ping the systems
func ProcessPingCommand() {
	ProcessPingCommandOutput()
}

func ProcessCommand(c string) {

	switch c {
	case "ping":
		ProcessPingCommand()
	case "status":
		ProcessRedfishStatusCommand()
	case "system":
		ProcessSystemCommand()
	case "fw-ver":
		ProcessFwCommand()
	case "bmc-ver":
		ProcessFwBMCCommand()
	case "nics":
		ProcessNICsCommand()
	case "accls":
		ProcessAcclCommand()

	}

	common.TearDown()
}
