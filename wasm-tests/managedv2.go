package main

import (
	"fmt"

	"github.com/lonelycode/wasmy/interfaces"
)

var exp *interfaces.WasmModulePrototype = &interfaces.WasmModulePrototype{}

// This is required so we can do I/O

//export inputBuffer
func InputBuffer() *[interfaces.FUNCBUFFER_SIZE]uint8 {
	return exp.GetInputPtr()
}

//export outputBuffer
func OutputBuffer() *[interfaces.FUNCBUFFER_SIZE]uint8 {
	return exp.GetOutputPtr()
}

// The wrapper will wrap execution of the function by grabbing input
// variables from the input buffer and providing them to the function
// as arguments, and then trapping the output of the function and
// writing it to a serialised output buffer

//export myExport
func MyExport(inputLen int) int {
	return interfaces.WrapExport(exp, inputLen, func(args ...interface{}) (interface{}, error) {
		name := args[0].(string)
		dt := fmt.Sprintf("hello %s", name)
		fmt.Printf("inside module: %s\n", dt)
		return dt, nil
	})()
}

func main() {}
