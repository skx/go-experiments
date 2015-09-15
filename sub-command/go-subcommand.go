//
// This is a simple program which is designed to be a base for
// a go-binary which handles several different sub-commands.
//
// For example we might have a script which supports:
//
//    foo login username $password
//
//    foo logout
//
//    foo list_load_balancers
//
//  This stub program allows the simple addition of sub-commands, via
// the implementation of an interface for each.
//
//  It is work-in-progress and very simple because go is still a little
// weird to me!
//
// Steve
// --
//

package main

import "fmt"
import "path"
import "os"

//
// Here's a basic interface for sub-commands.
//
type subcommand interface {
	// Show one-line help
	help(bool) string

	// Get the public-facing name of the command.
	name() string

	// Execute the command, with an array of arguments.
	execute(...string) int
}

//
// Gloal array of registered sub-commands
//
var CMDS = []subcommand{}

//
// Entry-point.
//
func main() {

	//
	// If we have at least one argument
	//
	if len(os.Args) >= 2 {

		//
		// Get the sub-command
		//
		sc := os.Args[1]

		//
		// And execute it, if we found the matching class.
		//
		for _, ent := range CMDS {
			if sc == ent.name() {
				os.Exit(ent.execute(os.Args[2:]...))
			}
		}

	}

	//
	// Are we looking for extended help?
	//
	if (len(os.Args) == 3) && (os.Args[1] == "help") {
		sc := os.Args[2]

		//
		// And execute it, if we found the matching class.
		//
		for _, ent := range CMDS {
			if sc == ent.name() {
				ent.help(true)
				os.Exit(0)
			}
		}
	}

	//
	// Otherwise show the commands and their help
	//
	fmt.Printf("Usage: %s subcommand [options]\n\nSubcommands include:\n\n", path.Base(os.Args[0]))
	fmt.Printf("\t% 8s - %s\n", "help", "Show extended help for the named subcommand")

	for _, ent := range CMDS {
		fmt.Printf("\t% 8s - %s\n", ent.name(), ent.help(false))
	}
}
