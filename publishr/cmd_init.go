/**
 * Implementation for 'publishr init'.
 *
 * Generate a random secret for TOTP-authentication.
 */

package main

import (
	"crypto/rand"
	"encoding/base32"
	"flag"
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

func (r cmd_init) help(extended bool) string {
	short := "Initialize our secure secret and state."
	if extended {
		fmt.Printf("%s\n", short)
		fmt.Printf("\nExtra Options:\n")
		fmt.Printf("  --force  For overwriting any existing secret and state.\n")
	}

	return short
}

/**
 * Write ~/.publishr.json with a secret for TOTP.
 */
func (r cmd_init) execute(args ...string) int {

	f1 := flag.NewFlagSet("f1", flag.ContinueOnError)
	force := f1.Bool("force", false, "Force overwriting existing details.")
	f1.Parse(args)

	path := os.Getenv("HOME") + "/.publishr.json"
	if !Exists(path) || *force {

		sec := make([]byte, 6)
		_, err := rand.Read(sec)
		if err != nil {
			fmt.Printf("Error creating random secret key: %s", err)
		}

		secret := base32.StdEncoding.EncodeToString(sec)

		state := PublishrState{Secret: secret, Count: 0}

		SaveState(state)

	} else {
		fmt.Printf("Already initialized, and running without --force\n")
	}

	path = os.Getenv("HOME") + "/public"
	if !Exists(path) {
		os.Mkdir(path, 0755)
	}

	return 0
}

func init() {
	CMDS = append(CMDS, cmd_init{})
}
