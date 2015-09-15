package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"reflect"
	"syscall"
)

type Information struct {
	FQDN       string
	Interfaces []string
	IPv4       []string
	IPv6       []string
}

type Utsname syscall.Utsname

func uname() (*syscall.Utsname, error) {
	uts := &syscall.Utsname{}

	if err := syscall.Uname(uts); err != nil {
		return nil, err
	}
	return uts, nil
}

func CharsToString(ca [65]int8) string {
	s := make([]byte, len(ca))
	var lens int
	for ; lens < len(ca); lens++ {
		if ca[lens] == 0 {
			break
		}
		s[lens] = uint8(ca[lens])
	}
	return string(s[0:lens])
}

func GetInformation() Information {
	interfaces := make([]string, 0)
	ipv4 := make([]string, 0)
	ipv6 := make([]string, 0)

	/**
	 * Get network interfaces
	 */
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		interfaces = append(interfaces, i.Name)
	}

	/**
	 * Get IPv4 & IPv6 addresses.
	 */
	addrs, _ := net.InterfaceAddrs()
	for _, address := range addrs {

		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {

			if ipnet.IP.To4() != nil {
				ipv4 = append(ipv4, ipnet.IP.String())
			} else {
				ipv6 = append(ipv6, ipnet.IP.String())
			}
		}
	}

	hostname, _ := os.Hostname()

	fmt.Printf("HOSTNAME: %s\n", hostname)
	/**
	 * Get hostname
	 */
	unme, _ := uname()

	host_name := CharsToString(unme.Nodename)
	domain_name := CharsToString(unme.Domainname)

	i := Information{FQDN: host_name + "." + domain_name,
		Interfaces: interfaces,
		IPv4:       ipv4,
		IPv6:       ipv6}
	return i
}

func hello(res http.ResponseWriter, req *http.Request) {

	info := GetInformation()

	res.Header().Set("Content-Type", "text/plain")

	/* Get the path */
	key := req.URL.Path[1:]

	if key == "" {
		jsn, err := json.Marshal(info)
		if err == nil {
			io.WriteString(res, string(jsn))
		} else {
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(res, "Failed encode to JSON")
		}
		return
	}

	/* Perform the reflection */
	r := reflect.ValueOf(&info)

	/* Get the field. */
	v := reflect.Indirect(r).FieldByName(key)

	/* If it isn't invalid */
	if v.Kind().String() != "invalid" {

		/* Get the value in JSON */
		j, err := json.Marshal(v.Interface())

		if err == nil {

			/* Success */
			fmt.Fprintf(res, "%s", string(j))
			return
		}
	}

	res.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(res, "Failed to lookup value of %s", key)
}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8000", nil)
}
