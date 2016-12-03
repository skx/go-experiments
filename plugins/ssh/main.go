package main

import "fmt"

//
// This is the input to the RUN_TEST method.
//
var INPUT string

//
// This function is called by our driver.
//
func RUN_TEST() {
	fmt.Printf("I'm a SSH-plugin called with: %s\n", INPUT)
}

//
// What kind of protocols will this plugin handle?
func HANDLES() string {
	return "ssh"
}
