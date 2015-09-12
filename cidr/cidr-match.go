package main

import (
	"fmt"
	"net"
	"os"
	"path"
	"regexp"
)

func test_range(rng string, ip string) bool {
	match, _ := regexp.MatchString("/", rng)

	if match {
		_, cidrnet, err := net.ParseCIDR(rng)
		if err != nil {
			panic(err)
		}
		myaddr := net.ParseIP(ip)
		if cidrnet.Contains(myaddr) {
			fmt.Printf("Range %s contains IP %s.\n", rng, ip)
			return true
		}
	} else {
		if ip == rng {
			fmt.Printf("Literal match: %s == %s\n", ip, rng)
			return true
		}
	}

	fmt.Printf("Failed to match IP %s\n", ip)
	return false
}

func main() {

	if len(os.Args) >= 3 {
		test_range(os.Args[1], os.Args[2])
	} else {
		fmt.Printf("Usage %s ip-range|ip1 ip2\n", path.Base(os.Args[0]))
	}
}
