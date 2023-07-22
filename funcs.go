package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
)

func print_usage() {
	log.Fatalln("Expected usage: echo <ip_address> | isaws")
}

func IsAWSIPAddress(ip net.IP) ([]Prefix, error) {
	matchingPrefixes := []Prefix{}

	for _, prefix := range prefixes {
		_, ipNet, err := net.ParseCIDR(prefix.IPPrefix)
		if err != nil {
			return nil, err
		}
		if ipNet.Contains(ip) {
			matchingPrefixes = append(matchingPrefixes, prefix)
		}
	}

	if len(matchingPrefixes) > 0 {
		return matchingPrefixes, nil
	}

	return nil, fmt.Errorf("IP %s not found in any AWS prefix\n", ip.String())
}

func GetAWSPrefixes() ([]Prefix, error) {

	url := "https://ip-ranges.amazonaws.com/ip-ranges.json"

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var data Response
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	var out []Prefix

	for _, prefix := range data.Prefixes {
		out = append(out, prefix)
	}
	return out, nil

}
