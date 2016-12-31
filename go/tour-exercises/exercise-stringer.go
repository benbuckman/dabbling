package main

import "fmt"
import "strings"
import "strconv"

type IPAddr [4]byte

// TODO: Add a "String() string" method to IPAddr.

func (ip IPAddr) String() string {
	numstrings := []string{}
	for _, b := range(ip) {
		numstrings = append(numstrings, strconv.Itoa(int(b)))
	}
	return strings.Join(numstrings, ".")
}

func main() {
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}
}
