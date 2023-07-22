package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var prefixes []Prefix
var ips = make(chan string, 100)
var concurrency int

func init() {
	var err error
	prefixes, err = GetAWSPrefixes()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {

	flag.IntVar(&concurrency, "c", 50, "set the concurrency level")
	flag.Parse()

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// iterate the user input
			for ip := range ips {
				ip := net.ParseIP(ip)
				if ip == nil {
					continue
				}

				// check if IP appears in AWS CIDR range
				matchingPrefixes, err := IsAWSIPAddress(ip)
				if err != nil {
					continue
				}

				// output
				for _, prefix := range matchingPrefixes {
					fmt.Printf("%s,%s,%s,%s,%s\n", ip, prefix.IPPrefix, prefix.Region, prefix.Service, prefix.NetworkBorderGroup)
				}

			}
		}()
	}

	// check for input piped to stdin
	info, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if info.Mode()&os.ModeCharDevice != 0 || (info.Mode()&os.ModeNamedPipe == 0 && info.Size() <= 0) {
		print_usage()
	}

	// get user input
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ips <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	// wait for workers
	close(ips)
	wg.Wait()
}
