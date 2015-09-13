/**
 * Implementation for 'publishr init'.
 *
 * Generate a random secret for TOTP-authentication.
 */

package main

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"os"
)

//
// types for each sub-command
//
type cmd_init struct{}

//
// Implementation for "init"
//
func (r cmd_init) name() string {
	return "init"
}

func (r cmd_init) help() string {
	return "Initialize our secure secret and state."
}

/**
 * Write ~/.publishr.json with a secret for TOTP.
 */
func (r cmd_init) execute(args ...string) int {

	path := os.Getenv("HOME") + "/.publishr.json"
	if !Exists(path) {

		sec := make([]byte, 6)
		_, err := rand.Read(sec)
		if err != nil {
			fmt.Printf("Error creating random secret key: %s", err)
		}

		secret := base32.StdEncoding.EncodeToString(sec)

		state := &PublishrState{Secret: secret,
			Count: 0}

		state_json, _ := json.Marshal(state)

		f, err := os.Create(path)
		defer f.Close()
		f.WriteString(string(state_json))

	} else {
		fmt.Printf("Already initialized - to remove the config please run:\n")
		fmt.Printf("\trm -f ~/.publishr.json\n")
	}

	path = os.Getenv("HOME") + "/public"
	if !Exists(path) {
		os.Mkdir(path, 0755)
	}

	return 0
}
