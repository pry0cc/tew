package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/n0ncetonic/nmapxml"
)

type parse struct {
	nmapxml.Run
}

func main() {
	var inputArg = flag.String("x", "", "Nmap XML Input File (Required)")
	var ipportArg = flag.String("i", "", "IP:Port Input File (Optional)")
	var dnsxArg = flag.String("dnsx", "", "dnsx -resp output data (Optional)")
	var vhostRep = flag.Bool("vhost", false, "Use dnsx data to insert vhosts (Optional)")
	var urlArg = flag.Bool("urls", false, "Guess HTTP URLs from input (Optional)")
	var outputArg = flag.String("o", "", "Output filename (Optional)")

	flag.Parse()
	var results []string

	input := *inputArg
	ipport := *ipportArg
	output := *outputArg
	dnsx := *dnsxArg
	vhost := *vhostRep
	urls := *urlArg

	if (input == "") && (ipport == "") {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if input != "" {
		results = parse.parseNmap(parse{}, input, dnsx, vhost, urls)
	} else if ipport != "" {
		results = parseIpport(ipport, dnsx, vhost, urls)
	}

	for _, line := range results {
		fmt.Println(line)
	}

	if output != "" {
		file, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}

		datawriter := bufio.NewWriter(file)

		for _, data := range results {
			datawriter.WriteString(data + "\n")
		}

		datawriter.Flush()
		file.Close()
	}

}

func unique(slice []string) []string {
	// create a map with all the values as key
	uniqMap := make(map[string]struct{})
	for _, v := range slice {
		uniqMap[v] = struct{}{}
	}

	// turn the map keys into a slice
	uniqSlice := make([]string, 0, len(uniqMap))
	for v := range uniqMap {
		uniqSlice = append(uniqSlice, v)
	}
	return uniqSlice
}

func parseIpport(input string, dnsx string, vhost bool, urls bool) []string {
	var index map[string][]string
	var output []string

	if input != "-" {
		if _, err := os.Stat(input); err != nil {
			fmt.Printf("File does not exist\n")
			os.Exit(1)
		}
	} else {
		// todo
	}

	if dnsx != "" {
		if _, err := os.Stat(dnsx); err != nil {
			fmt.Printf("dnsx file does not exist\n")
		} else {
			index = parseDnsx(dnsx)
		}
	}

	file, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ":")
		ip, port := s[0], s[1]
		service := ""

		if strings.Contains(port, "80") {
			service = "http"
		} else if strings.Contains(port, "443") {
			service = "https"
		}

		resp := processData(ip, port, service, vhost, urls, index)
		for _, line := range resp {
			output = append(output, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	uniq := unique(output)
	return uniq
}

func (p parse) parseNmap(input string, dnsx string, vhost bool, urls bool) []string {
	/* ParseNmap parses a Nmap XML file */
	var index map[string][]string
	var output []string

	if input != "-" {
		if _, err := os.Stat(input); err != nil {
			fmt.Printf("File does not exist\n")
			os.Exit(1)
		}

		p.Run, _ = nmapxml.Readfile(input)
	} else {
		bytes, _ := ioutil.ReadAll(os.Stdin)
		xml.Unmarshal(bytes, &p.Run)
	}

	if dnsx != "" {
		if _, err := os.Stat(dnsx); err != nil {
			fmt.Printf("dnsx file does not exist\n")
		} else {
			index = parseDnsx(dnsx)
		}
	}

	hosts := p.Host
	for _, host := range hosts {
		ipAddr := host.Address.Addr
		if host.Ports.Port != nil {
			for _, portData := range *host.Ports.Port {
				if portData.State.State == "open" {
					portID := portData.PortID
					service := portData.Service.Name

					resp := processData(ipAddr, portID, service, vhost, urls, index)

					for _, line := range resp {
						output = append(output, line)
					}
				}
			}
		}
	}

	uniq := unique(output)
	return uniq
}

func processData(ipAddr string, port string, service string, vhost bool, urls bool, index map[string][]string) []string {
	var output []string
	if vhost {
		indexed := index[ipAddr]
		for _, dom := range indexed {
			line := ""
			if urls {
				line = generateURL(dom, port, service)
			} else {
				line = dom + ":" + port
			}

			if line != "" {
				output = append(output, line)
			}
		}

	} else {
		line := ""
		if urls {
			line = generateURL(ipAddr, port, service)
		} else {
			line = ipAddr + ":" + port
		}

		if line != "" {
			output = append(output, line)
		}
		//fmt.Println(ipAddr + ":" + portID)
	}

	return output
}

func generateURL(host string, port string, service string) string {
	/* GenURl generates a URL for a given sequence */
	url := ""
	if (port == "80" && service == "http") || (port == "443" && service == "https") {
		url = service + "://" + host
	} else if strings.Contains(service, "http") {
		if strings.Contains(port, "80") {
			service = "http"
		} else if strings.Contains(port, "443") {
			service = "https"
		} else {
			service = "http"
		}
		url = service + "://" + host + ":" + port
	}

	return url
}

func parseDnsx(filename string) map[string][]string {
	/* ParseDnsx parses a DNSX JSON file */
	var data = map[string][]string{}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var result map[string]interface{}
		err = json.Unmarshal([]byte(scanner.Text()), &result)
		if err != nil {
			continue
		}
		host := result["host"].(string)

		if val, ok := result["a"]; ok {
			aRecords := val.([]interface{})

			ip := ""

			for _, record := range aRecords {
				ip = record.(string)
			}

			data[ip] = append(data[ip], host)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return data
}
