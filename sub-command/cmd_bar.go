package main

import "fmt"

//
// types for each sub-command
//
type cmd_bar struct{}

//
// Implementation for "bar"
//
func (r cmd_bar) name() string {
	return "bar"
}
func (r cmd_bar) help() string {
	return "run sub-command bar"
}
func (r cmd_bar) execute(args ...string) int {
	fmt.Println("I am bar")
	fmt.Println("I received ", len(args), " arguments")

	for _, ent := range args {
		fmt.Println("Argument: ", ent)
	}

	return 0
}
