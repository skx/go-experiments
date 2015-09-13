/**
 * Utility functions which we require, but didn't otherwise
 * have an obvious home.
 */

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

var mutex = &sync.Mutex{}

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

	mutex.Lock()
	state_pth := os.Getenv("HOME") + "/.publishr.json"
	state_json, _ := json.Marshal(state)
	f, _ := os.Create(state_pth)
	defer f.Close()
	f.WriteString(string(state_json))

	mutex.Unlock()
}

/**
 * Load state
 */
func LoadState() (PublishrState, error) {
	mutex.Lock()

	state_pth := os.Getenv("HOME") + "/.publishr.json"
	state_cnt, _ := ioutil.ReadFile(state_pth)

	var state PublishrState

	if err := json.Unmarshal(state_cnt, &state); err != nil {
		mutex.Unlock()
		return state, err
	}
	mutex.Unlock()
	return state, nil

}
