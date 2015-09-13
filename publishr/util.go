package main

import (
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
