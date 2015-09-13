package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

/**
 * This is the state our server itself keeps:
 *
 *  Secret - Used for TOTP authentication.
 *
 *  Count  - The number of files uploaded.
 */
type PublishrState struct {
	Secret string `json:"secret"`
	Count  int    `json:"count"`
}

/**
 * Report whether the named file or directory exists.
 */
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

/**
 * Save state
 */
func SaveState(state PublishrState) {
	state_pth := os.Getenv("HOME") + "/.publishr.json"
	state_json, _ := json.Marshal(state)
	f, _ := os.Create(state_pth)
	defer f.Close()
	f.WriteString(string(state_json))
}

/**
 * Load state
 */
func LoadState() (PublishrState, error) {

	state_pth := os.Getenv("HOME") + "/.publishr.json"
	state_cnt, _ := ioutil.ReadFile(state_pth)

	var state PublishrState

	if err := json.Unmarshal(state_cnt, &state); err != nil {
		return state, err
	}
	return state, nil

}
