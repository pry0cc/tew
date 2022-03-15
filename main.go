package main

import (
	"fmt"
	"github.com/n0ncetonic/nmapxml"
	"os"
)

func main() {
	input := os.Args[1]
	ParseXml(input)
}

func ParseXml(xmlFileName string) {
	scanData, _ := nmapxml.Readfile(xmlFileName)
	dealWithRun(scanData)
}

func dealWithRun(r nmapxml.Run) {
	hostSlice := r.Host
	for _, host := range hostSlice {
		ipAddr := host.Address.Addr
		if host.Ports.Port != nil {
			for _, portInfo := range *host.Ports.Port {
				if portInfo.State.State == "open" {
					portID := portInfo.PortID
					fmt.Println(ipAddr + ":" + portID)
				}
			}
		}
	}
}
