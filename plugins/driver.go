package main

import (
	"fmt"
	"path/filepath"
	"plugin"
)

//
// Given the path to a plugin, and a string argument, invoke the plugin
// method.
//
func run_plugin(path string, arg string) {
	p, err := plugin.Open(path)
	if err != nil {
		panic(err)
	}
	v, err := p.Lookup("INPUT")
	if err != nil {
		panic(err)
	}
	f, err := p.Lookup("RUN_TEST")
	if err != nil {
		panic(err)
	}
	*v.(*string) = arg
	f.(func())() // call the plugin

}

//
// Call the `HANDLES` method of the specified plugin.
//
func plugin_handles(path string) string {
	p, err := plugin.Open(path)
	if err != nil {
		panic(err)
	}
	f, err := p.Lookup("HANDLES")
	if err != nil {
		panic(err)
	}
	x := f.(func() string)()
	return (x)
}

func main() {

	//
	// The plugin association array.
	//
	plugins := make(map[string]string)

	//
	// Find each plugin.
	//
	files, _ := filepath.Glob("*/*.so")

	//
	// Each plugin will handle a specific test-type
	//
	// Invoke each one so we can store a list of known-types.
	//
	for _, f := range files {
		plugins[f] = plugin_handles(f)
	}

	//
	// Show the plugins we've loaded, as well as what type of
	// test is accepted
	//
	fmt.Println("Plugin Dump")
	for plugin_path, plugin_type := range plugins {
		fmt.Printf("\tPlugin: %s - Handles: %s\n",
			plugin_path, plugin_type)
	}

	//
	// Invoke each plugin.
	//
	fmt.Println("Calling Plugins")
	for plugin_path, _ := range plugins {
		run_plugin(plugin_path, "Hello World!")
	}

}
