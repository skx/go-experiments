package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"
	"reflect"
	"strings"
)

/**
 * The information structure we return to callers.
 *
 * Note this is deliberately "flat" and "fat".
 */
type Information struct {
	ARCH         string
	FQDN         string
	LSB_Codename string
	LSB_Release  string
	LSB_Version  string
	Interfaces   []string
	IPv4         []string
	IPv6         []string
}

/**
 * Run a command, and return the output without any newlines.
 */
func runCommand(cmd string, args ...string) string {

	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		log.Panic(err)
		return ""
	}
	return strings.Trim(string(out), "\r\n")
}

/**
 * Populate the structure.
 *
 * Currently this is called once per-request.  Ideally it'd be done
 * only at startup-time.
 */
func GetInformation() Information {

	/**
	 * Some fields in our structure are arrays.
	 *
	 * Create them.
	 */
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

	i := Information{
		ARCH:         runCommand("arch"),
		FQDN:         runCommand("/bin/hostname", "--fqdn"),
		LSB_Codename: runCommand("/usr/bin/lsb_release", "--short", "--codename"),
		LSB_Release:  runCommand("/usr/bin/lsb_release", "--short", "--id"),
		LSB_Version:  runCommand("/usr/bin/lsb_release", "--short", "--release"),
		Interfaces:   interfaces,
		IPv4:         ipv4,
		IPv6:         ipv6}
	return i
}

func hello(res http.ResponseWriter, req *http.Request) {

	info := GetInformation()

	res.Header().Set("Content-Type", "text/plain")

	/* Get the path */
	key := req.URL.Path[1:]

	/**
	 *  This allows:
	 *   http://example.com:800/LSB/Release -> LSB_Release
	 */
	key = strings.Replace(key, "/", "_", -1)

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
