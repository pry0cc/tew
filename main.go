package main

import (
	"flag"
	"fmt"
	"github.com/n0ncetonic/nmapxml"
	"os"
)

func main() {
	var inputArg = flag.String("x", "", "Nmap XML Input File (Required)")
	var dnsxArg = flag.String("dnsx", "", "dnsx -resp output data (TODO)")
	flag.Parse()

	input := *inputArg
	dnsx := *dnsxArg

	if input == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if dnsx != "" {
		if _, err := os.Stat(dnsx); err != nil {
			fmt.Printf("dnsx file does not exist\n")
		}
	}

	if _, err := os.Stat(input); err != nil {
		fmt.Printf("File does not exist\n")
	}

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
