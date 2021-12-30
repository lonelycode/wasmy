// package interfaces provides types and functions that help manage complex
// data type passing between WASM imported and WASM-exported functions in
// order to make working with WASM easier for a developer.
package interfaces

import (
	"bytes"
	"fmt"
	"os"

	"github.com/lonelycode/wasmy/shared_types"
	"github.com/tinylib/msgp/msgp"
)

const (
	// FUNCBUFFER_SIZE sets the maximum size of I/O
	// buffers used to pass data between host and
	// guest and vice versa
	FUNCBUFFER_SIZE = 1024
)

// WasmModulePrototype provides a wrapper for managing I/O for WASM modules
type WasmModulePrototype struct {
	guestFnInputBfr  [FUNCBUFFER_SIZE]uint8 // exported Fn input buffer
	guestFnOutputBfr [FUNCBUFFER_SIZE]uint8 // exported Fn output buffer

	hostFnInputBfr  [FUNCBUFFER_SIZE]uint8 // imported Fn input buffer
	hostFnOutputBfr [FUNCBUFFER_SIZE]uint8 // imported Fn input buffer
}

// GetInputPtr will return a pointer to the `guestFnInputBfr` in
// WASM linear memory, this buffer provides arg input to functions
// exported by the WASM module in the form of shared_types.Args
func (d *WasmModulePrototype) GetInputPtr() *[FUNCBUFFER_SIZE]uint8 {
	return &d.guestFnInputBfr
}

// GetOutputPtr will return a pointer to the `guestFnOutputBfr` in
// WASM linear memory, this buffer provides output for functions
// exported by the WASM module in the form of shared_types.Payload
func (d *WasmModulePrototype) GetOutputPtr() *[FUNCBUFFER_SIZE]uint8 {
	return &d.guestFnOutputBfr
}

// GetHostOutputPtr will return a pointer to the `hostFnOutputBfr` in
// WASM linear memory, this buffer provides output for functions
// imported by the WASM module in the form of shared_types.Payload
func (d *WasmModulePrototype) GetHostOutputPtr() *[FUNCBUFFER_SIZE]uint8 {
	return &d.hostFnOutputBfr
}

// GetHostInputPtr will return a pointer to the `hostFnInputBfr` in
// WASM linear memory, this buffer provides input for functions
// imported by the WASM module in the form of shared_types.Args
func (d *WasmModulePrototype) GetHostInputPtr() *[FUNCBUFFER_SIZE]uint8 {
	return &d.hostFnInputBfr
}

// ReadGuestFnInput will read the input buffer for any WASM-exported functions
// as `shared_types.Args` and return the argument array to the caller
func (d *WasmModulePrototype) ReadGuestFnInput(length int) ([]interface{}, error) {
	dat := make([]byte, length)
	copy(dat, d.guestFnInputBfr[:length])

	args := &shared_types.Args{
		Args: make([]interface{}, 0),
	}

	buf := bytes.NewBuffer(dat)
	err := msgp.Decode(buf, args)
	if err != nil {
		return nil, err
	}

	return args.Args, nil
}

// ReadHostFnOutput will read the output buffer for host functions (imported by WASM module)and
// returns an error. The function takes an output interface pointer in order to easily modify
// the payload type by the caller.
func (d *WasmModulePrototype) ReadHostFnOutput(length int, output *shared_types.Payload) error {
	dat := make([]byte, length)
	copy(dat, d.hostFnOutputBfr[:length])

	buf := bytes.NewBuffer(dat)
	err := msgp.Decode(buf, output)
	if err != nil {
		return err
	}

	return nil
}

// WriteGuestFnOutput writes WASM-exported function output into the guest buffer as a Payload
func (d *WasmModulePrototype) WriteGuestFnOutput(data interface{}) (int, error) {
	out := &shared_types.Payload{Data: data}

	enc, err := out.MarshalMsg(nil)
	if err != nil {
		return 0, err
	}
	copy(d.guestFnOutputBfr[:len(enc)], enc)
	return len(enc), nil
}

// WriteHostFnInput will write the args for a WASM-imported function into the host
// input buffer as an Args object
func (d *WasmModulePrototype) WriteHostFnInput(args []interface{}) (int, error) {
	out := &shared_types.Args{Args: args}

	enc, err := out.MarshalMsg(nil)
	if err != nil {
		return 0, err
	}
	copy(d.hostFnInputBfr[:len(enc)], enc)
	return len(enc), nil
}

// externGuestErr is sup;poed to provide an error buffer, TODO: still unsure if this idea
// is worth pursuing for managed error output from called funcs
func (d *WasmModulePrototype) externGuestErr(err error) int {
	errTp := fmt.Sprintf("ERR %s", err.Error())
	os.Stderr.WriteString(errTp)
	copy(d.guestFnOutputBfr[:], []byte(errTp))

	return len([]byte(errTp))
}

// WrapExport takes a prototype (managed buffers) object and a caller function, it will
// write the args to the guest input buffer, run the function, and capture the return data
// from the guest function to write into the output buffer. It returns the length of the data
// written in order for the caller to pull the correct data from the buffer.
func WrapExport(proto *WasmModulePrototype, inputLen int, exportFn func(args ...interface{}) (interface{}, error)) func() int {
	return func() int {
		args, err := proto.ReadGuestFnInput(inputLen)
		if err != nil {
			return proto.externGuestErr(err)
		}

		ret, err := exportFn(args...)
		if err != nil {
			return proto.externGuestErr(err)
		}

		n, err := proto.WriteGuestFnOutput(ret)
		if err != nil {
			return proto.externGuestErr(err)
		}

		return n
	}
}

// CallImport will take a managed buffer prototype, imported function and arguments and
// writes the args to the host input buffer, it will then capture the output of the function
// from the host output buffer, unmarshal it and return it to the caller as a Payload.
// TODO: Maybe it should return the payload.Data interface instead?
func CallImport(proto *WasmModulePrototype, fn func(int32) int32, args ...interface{}) (*shared_types.Payload, error) {
	// Write our args to the host input buffer
	lenInp, err := proto.WriteHostFnInput(args)
	if err != nil {
		return nil, err
	}

	// call the imported function with the length of the input data
	lenOut := fn(int32(lenInp))

	// Read the host output buffer for the return value
	output := &shared_types.Payload{}
	err = proto.ReadHostFnOutput(int(lenOut), output)
	if err != nil {
		return nil, err
	}

	return output, nil
}
