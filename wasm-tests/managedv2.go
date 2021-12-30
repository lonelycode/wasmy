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

//export hostInputBuffer
func HostInputBuffer() *[interfaces.FUNCBUFFER_SIZE]uint8 {
	return exp.GetHostInputPtr()
}

//export hostOutputBuffer
func HostOutputBuffer() *[interfaces.FUNCBUFFER_SIZE]uint8 {
	return exp.GetHostOutputPtr()
}

// The wrapper will wrap execution of the function by grabbing input
// variables from the input buffer and providing them to the function
// as arguments, and then trapping the output of the function and
// writing it to a serialised output buffer

// sample imported func
func PrintHello(int32) int32

//export myExport
func MyExport(inputLen int) int {
	return interfaces.WrapExport(exp, inputLen, func(args ...interface{}) (interface{}, error) {
		name := args[0].(string)
		dt := fmt.Sprintf("hello %s", name)
		fmt.Printf("inside module: %s\n", dt)

		doStuff()

		return dt, nil
	})()
}

func doStuff() {
	ret, err := interfaces.CallImport(exp, PrintHello, "anderson")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("host function output from inside module: %v\n", ret.Data)
}

func main() {}
