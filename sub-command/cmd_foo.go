package main

import "fmt"

type cmd_foo struct{}

//
// Implementation for "foo"
//
func (r cmd_foo) name() string {
	return "foo"
}
func (r cmd_foo) help() string {
	return "run sub-command foo"
}
func (r cmd_foo) execute(args ...string) int {
	fmt.Println("I am foo")
	fmt.Println("I received ", len(args), " arguments")

	for _, ent := range args {
		fmt.Println("Argument: ", ent)
	}

	return 0
}
