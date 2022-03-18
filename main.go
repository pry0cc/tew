package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/n0ncetonic/nmapxml"
	"log"
	"os"
	"strings"
)

func main() {
	var inputArg = flag.String("x", "", "Nmap XML Input File (Required)")
	var dnsxArg = flag.String("dnsx", "", "dnsx -resp output data (Optional)")
	var vhostRep = flag.Bool("vhost", false, "Use dnsx data to insert vhosts (Optional)")
	var urlArg = flag.Bool("urls", false, "Guess HTTP URLs from input (Optional)")
	var outputArg = flag.String("o", "", "Output filename (Optional)")
	flag.Parse()

	input := *inputArg
	output := *outputArg
	dnsx := *dnsxArg
	vhost := *vhostRep
	urls := *urlArg

	if input == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	results := ParseNmap(input, dnsx, vhost, urls)

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
			_, _ = datawriter.WriteString(data + "\n")
		}

		datawriter.Flush()
		file.Close()
	}

}

func Unique(slice []string) []string {
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

func ParseNmap(input string, dnsx string, vhost bool, urls bool) []string {
	/* ParseNmap parses a Nmap XML file */
	var index map[string][]string
	var output []string
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
					service := portData.Service.Name

					if vhost {
						for _, ipp := range index {
							domains := ipp

							for _, dom := range domains {
								line := ""
								if urls {
									line = GenUrl(dom, portID, service)
								} else {
									line = dom + ":" + portID
								}

								if line != "" {
									output = append(output, line)
								}
							}
						}

					} else {
						line := ""
						if urls {
							line = GenUrl(ipAddr, portID, service)
						} else {
							line = ipAddr + ":" + portID
						}

						if line != "" {
							output = append(output, line)
						}
						//fmt.Println(ipAddr + ":" + portID)
					}
				}
			}
		}
	}

	uniq := Unique(output)
	return uniq
}

func GenUrl(host string, port string, service string) string {
	/* GenURl generates a URL for a given sequence */
	url := ""
	if service == "http" || service == "https" {
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

func ParseDnsx(filename string) map[string][]string {
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
		json.Unmarshal([]byte(scanner.Text()), &result)
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
