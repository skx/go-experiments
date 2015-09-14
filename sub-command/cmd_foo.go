package main

import "fmt"

type cmd_foo struct{}

//
// Implementation for "foo"
//
func (r cmd_foo) name() string {
	return "foo"
}
func (r cmd_foo) help(extended bool) string {
        short := "run sub-command foo"
        if extended {
                fmt.Printf("%s\n\n", short)
                fmt.Printf("Extra Options:\n\n\tNone\n\n")
        }

        return short
}
func (r cmd_foo) execute(args ...string) int {
	fmt.Println("I am foo")
	fmt.Println("I received ", len(args), " arguments")

	for _, ent := range args {
		fmt.Println("Argument: ", ent)
	}

	return 0
}
