package main

import (
	"fmt"

	"github.com/lonelycode/wasmy/interfaces"

	// This MUST be imported to provide boilerplate exports
	// Note this is not suitable for parallel execution, as you may end up overwriting i/o buffers
	// for that, implement a Proto for each export
	module_params "github.com/lonelycode/wasmy/module-params"
)

// SAMPLE USAGE
// ------------

// sample imported func (see exports/exports.go and example/main.go)
func PrintHello(int32) int32

// this is the function signature for all exported functions managed by the prototype
func myFunction(args ...interface{}) (interface{}, map[string]string, error) {
	name := args[0].(string)

	dt := fmt.Sprintf("hello %s", name)

	fmt.Printf("inside module: %s\n", dt)

	// For demo purposes, let's call an imported function here,
	// this is defined in the exports/exports.go file
	doStuff()

	return dt, nil, nil
}

// MyExport is a function stub to export the wrapped and managed version
// of our actual method (don't forget the `//export <foo>` tag otherwise
// the method will remain unexported)
//export myExport
func MyExport(inputLen int) int {
	return interfaces.WrapExport(module_params.Proto, inputLen, myFunction)()
}

// doStuff is an unexported module function that calls an imported method from the host
func doStuff() {
	ret, err := interfaces.CallImport(module_params.Proto, PrintHello, "anderson")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("host function output from inside module: %s\n", ret.(string))
}

func main() {}
