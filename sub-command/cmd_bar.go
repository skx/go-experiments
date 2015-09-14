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
func (r cmd_bar) help(extended bool) string {
	short := "run sub-command bar"
	if extended {
		fmt.Printf("%s\n\n", short)
		fmt.Printf("Extra Options:\n\n\tNone\n\n")
	}

	return short
}
func (r cmd_bar) execute(args ...string) int {
	fmt.Println("I am bar")
	fmt.Println("I received ", len(args), " arguments")

	for _, ent := range args {
		fmt.Println("Argument: ", ent)
	}

	return 0
}
