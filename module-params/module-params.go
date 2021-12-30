package module_params

import "github.com/lonelycode/wasmy/interfaces"

//==========  START BOILERPLATE ==========//

// This variable will provide all the i/o we need for this module
var Proto *interfaces.WasmModulePrototype = &interfaces.WasmModulePrototype{}

// This is required so we can do I/O

//export inputBuffer
func InputBuffer() *[interfaces.FUNCBUFFER_SIZE]uint8 {
	return Proto.GetInputPtr()
}

//export outputBuffer
func OutputBuffer() *[interfaces.FUNCBUFFER_SIZE]uint8 {
	return Proto.GetOutputPtr()
}

//export hostInputBuffer
func HostInputBuffer() *[interfaces.FUNCBUFFER_SIZE]uint8 {
	return Proto.GetHostInputPtr()
}

//export hostOutputBuffer
func HostOutputBuffer() *[interfaces.FUNCBUFFER_SIZE]uint8 {
	return Proto.GetHostOutputPtr()
}

//==========  END BOILERPLATE ==========//
