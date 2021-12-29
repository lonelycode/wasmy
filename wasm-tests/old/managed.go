package main

import (
	"fmt"

	"github.com/lonelycode/wasmy/interfaces"
)

// define your module
type MyModule struct {
	interfaces.DefaultWasmModule
}

// create a global instance for exports
var mod *MyModule

//== Start Boilerplate ==//

//export getMemoryPtr
func GetMemoryPtr() *[interfaces.BUFFER_SIZE]uint8 {
	return mod.GetMemoryPtr()
}

//export getCursor
func GetCursor() uint32 {
	return mod.GetCursor()
}

//export getCursorMapPtr
func GetCursorMapPtr() *[2]uint32 {
	return mod.GetCursorMapPtr()
}

//export getCursorForData
func GetCursorForData(id uint32) *[2]uint32 {
	return mod.GetCursorForData(id)
}

//export getStartCursorForData
func GetStartCursorForData(id uint32) uint32 {
	return mod.GetStartCursorForData(id)
}

//export getEndCursorForData
func GetEndCursorForData(id uint32) uint32 {
	return mod.GetEndCursorForData(id)
}

//export incrementCursor
func IncrementCursor(length uint32) uint32 {
	return mod.IncrementCursor(length)
}

//== End Boilerplate ==//

// Initialise it
func init() {
	mod = &MyModule{}
	mod.Init()
}

//export storeValueInWasmMemoryBuffer
func storeValueInWasmMemoryBuffer() {
	mod.Write(1, "hello")
}

// Function to read from index 1 of our buffer
// And return the value at the index
//export readWasmMemoryBuffer
func readWasmMemoryBuffer() {
	a := mod.GetStartCursorForData(1)
	b := mod.GetEndCursorForData(1)
	fmt.Println(a)
	fmt.Println(b)
}

func main() {}
