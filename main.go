package main

import (
	"fmt"
	"github.com/n0ncetonic/nmapxml"
	"os"
)

func main() {
	input := os.Args[1]
	scanData, _ := nmapxml.Readfile(input)
	ParseNmap(scanData)
}

func ParseNmap(r nmapxml.Run) {
	hostS := r.Host
	for _, host := range hostS {
		ipAddr := host.Address.Addr
		if host.Ports.Port != nil {
			for _, portData := range *host.Ports.Port {
				if portData.State.State == "open" {
					portID := portData.PortID
					fmt.Println(ipAddr + ":" + portID)
				}
			}
		}
	}
}
