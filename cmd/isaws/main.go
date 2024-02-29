package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/cybercdh/isaws/awschecker"
)

func main() {
	var concurrency int

	flag.IntVar(&concurrency, "c", 50, "set the concurrency level")
	flag.Parse()

	prefixes, err := awschecker.GetAWSPrefixes()
	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup
	ips := make(chan string, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ipStr := range ips {
				ip := net.ParseIP(ipStr)
				if ip == nil {
					continue
				}

				matchingPrefixes, err := awschecker.IsAWSIPAddress(ip, prefixes)
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
