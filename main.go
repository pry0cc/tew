package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/n0ncetonic/nmapxml"
	"log"
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

	ParseNmap(input, dnsx)
}

func ParseNmap(input string, dnsx string) {
	var index map[string][]string
	r, _ := nmapxml.Readfile(input)

	if _, err := os.Stat(input); err != nil {
		fmt.Printf("File does not exist\n")
	}

	if dnsx != "" {
		if _, err := os.Stat(dnsx); err != nil {
			fmt.Printf("dnsx file does not exist\n")
		} else {
			index = ParseDnsx(dnsx)
		}
	}

	hostS := r.Host
	for _, host := range hostS {
		ipAddr := host.Address.Addr
		if host.Ports.Port != nil {
			for _, portData := range *host.Ports.Port {
				if portData.State.State == "open" {
					portID := portData.PortID

					for _, ipp := range index {
						fmt.Println(ipp)
					}

					fmt.Println(ipAddr + ":" + portID)
				}
			}
		}
	}
}

func ParseDnsx(filename string) map[string][]string {
	/* ParseDnsx parses a DNSX JSON file*/
	var data = map[string][]string{}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var result map[string]interface{}
		json.Unmarshal([]byte(scanner.Text()), &result)
		host := result["host"].(string)
		aRecords := result["a"].([]interface{})
		ip := ""

		for _, record := range aRecords {
			ip = record.(string)
			data[ip] = append(data[ip], host)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return data
}
