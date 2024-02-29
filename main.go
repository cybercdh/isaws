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
			for ip := range ips {
				ip := net.ParseIP(ip)
				if ip == nil {
					continue
				}

				matchingPrefixes, err := IsAWSIPAddress(ip)
				if err != nil {
					continue
				}

				for _, prefix := range matchingPrefixes {
					fmt.Printf("%s,%s,%s,%s,%s\n", ip, prefix.IPPrefix, prefix.Region, prefix.Service, prefix.NetworkBorderGroup)
				}
			}
		}()
	}

	// Simplified: Directly read from stdin without checking if it's piped or redirected
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ips <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error reading from stdin:", err)
	}

	close(ips)
	wg.Wait()
}
