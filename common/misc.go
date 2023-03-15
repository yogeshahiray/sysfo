// misc functions
package common

import (
	"flag"
	"fmt"
	"github.com/c-robinson/iplib"
	"net"
	"os"
)

// Return true if valid IP
func IsValidIP(ipaddr string) bool {
	if net.ParseIP(ipaddr) == nil {
		return false
	}
	return true
}
func IsValidIPRange(start string, end string) bool {
	ipa := net.ParseIP(start)
	ipb := net.ParseIP(end)
	if iplib.CompareIPs(ipa, ipb) > 0 {
		return false
	}

	return true
}

func PrintUsage() {
	flag.Usage()
}

func PrintError(m string) {
	fmt.Println(m)
	PrintUsage()
	os.Exit(1)
}
