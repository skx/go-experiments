package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//
// types for each sub-command
//
type cmd_secret struct{}

//
// Implementation for "secret"
//
func (r cmd_secret) name() string {
	return "secret"
}

func (r cmd_secret) help() string {
	return "Show our authentication-secret."
}

/**
 * Read ~/.publishr.json and show the secret
 */
func (r cmd_secret) execute(args ...string) int {

	path := os.Getenv("HOME") + "/.publishr.json"
	if Exists(path) {

		state_cnt, _ := ioutil.ReadFile(path)

		var state PublishrState

		if err := json.Unmarshal(state_cnt, &state); err != nil {
			fmt.Printf("Failed to parse %s:\n", path)
			return 0
		}

		fmt.Printf("Configure your authenticator-client with the secret %s\n", state.Secret)
		fmt.Printf("For example:\n")
		fmt.Printf("\toathtool --totp -b %s\n", state.Secret)
	} else {
		fmt.Printf("Not initialized - Please run 'publishr init'\n")
	}
	return 0
}
