//
// pretty-printer for JSON.
//

package main

import "encoding/json"
import "io/ioutil"
import "os"
import "io"

//
// Return either a handle to STDIN, or the file on the command-line.
//
func openStdinOrFile() io.Reader {
	var err error
	r := os.Stdin
	if len(os.Args) > 1 {
		r, err = os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
	}
	return r
}

//
// Entry point.
//
func main() {

	//
	// Read the complete contents of the named file, or from STDIN.
	//
	input := openStdinOrFile()
	byt, err := ioutil.ReadAll(input)
	if err != nil {
		panic(err)
	}

	//
	// Unpack the JSON.
	//
	var dat map[string]interface{}

	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}

	//
	// Pretty-Print it.
	//
	b, err := json.MarshalIndent(dat, "", "  ")
	if err != nil {
		panic(err)
	}

	//
	// Write it to STDOUT
	//
	b2 := append(b, '\n')
	os.Stdout.Write(b2)

}
