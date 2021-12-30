// package runner provides types and methods that make it easier for a host to run wasm
// modules with managed I/O for better data passing between guest and host.
package runner

import (
	"bytes"
	"fmt"
	"os"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
	"github.com/lonelycode/wasmy/exports"
	"github.com/lonelycode/wasmy/shared_types"
	"github.com/tinylib/msgp/msgp"
)

// Runner provides methods that manages all the components of running multiple WASM
// modules and calling arbitrary functions from them using managed I/O
type Runner struct {
	mem                *wasmtime.Memory
	store              *wasmtime.Store
	instance           *wasmtime.Instance
	inputBufferFn      *wasmtime.Func
	outputBufferFn     *wasmtime.Func
	hostInputBufferFn  *wasmtime.Func
	hostOutputBufferFn *wasmtime.Func
	funcMap            map[string]*wasmtime.Func
}

// exportFun represents the signature needed for any function exported by
// the host and imported by the WASM file
type exportFunc func(int32, int32, int32) int32

// WrapExport will wrap any function to be exported by the HOST and used by the WASM
// module in order to capture the input args for the function and the output from the
// function and pass the data cleanly to the WASM module (see the exports package
// for an example function that can be wrapped).
func (r *Runner) WrapExport(fn func(*shared_types.Args) (interface{}, error)) exportFunc {
	return func(dataLen int32, t2 int32, t3 int32) int32 {

		ptr, err := r.hostInputBufferFn.Call(r.store)
		if err != nil {
			// fmt.Printf("failed to call inputBufferFn: %v\n", err)
			os.Stderr.WriteString(err.Error())
			return -1
		}

		outPtr, err := r.hostOutputBufferFn.Call(r.store)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			// fmt.Printf("failed to call hostOutputBufferFn: %v\n", err)
			return -1
		}

		hostArgs := &shared_types.Args{}
		outDat := make([]byte, dataLen)
		copy(outDat[:], r.mem.UnsafeData(r.store)[int(ptr.(int32)):int(ptr.(int32))+int(dataLen)])
		buf := bytes.NewBuffer(outDat)
		err = msgp.Decode(buf, hostArgs)
		if err != nil {
			// fmt.Printf("failed to decode: %v\n", err)
			os.Stderr.WriteString(err.Error())
			return -1
		}

		// call the actual functions
		ret, err := fn(hostArgs)
		if err != nil {
			// fmt.Printf("failed to call function: %v\n", err)
			os.Stderr.WriteString(err.Error())
			return -1
		}

		// Encode the output back into the guest VM
		out := &shared_types.Payload{Data: ret}

		enc, err := out.MarshalMsg(nil)
		if err != nil {
			// fmt.Printf("failed to call marshal for output: %v\n", err)
			os.Stderr.WriteString(err.Error())
			return -1
		}

		outputLen := copy(r.mem.UnsafeData(r.store)[int(outPtr.(int32)):int(outPtr.(int32))+len(enc)], enc)

		// return how much we wrote
		return int32(outputLen)
	}

}

// AddHostFunctions adds functions that can be imported into the WASM module,
// multiple funcs can be added, they all live in the `env` namespace
func (r *Runner) AddHostFunctions(linker *wasmtime.Linker, funcs map[string]exportFunc) {
	for name, fn := range funcs {
		linker.DefineFunc(r.store, "env", fmt.Sprintf("main.%s", name), fn)
	}
}

// GetInstance provides a WASM VM instance from the file name. It enables WASI,
// but only shares stdout and stderr for easier logging.
func (r *Runner) GetInstance(filename string) (*wasmtime.Instance, *wasmtime.Store, error) {
	engine := wasmtime.NewEngine()
	r.store = wasmtime.NewStore(engine)

	module, err := wasmtime.NewModuleFromFile(r.store.Engine, filename)
	if err != nil {
		return nil, nil, err
	}

	wConf := wasmtime.NewWasiConfig()
	wConf.InheritStdout()
	wConf.InheritStderr()
	r.store.SetWasi(wConf)
	if err != nil {
		return nil, nil, err
	}

	// Create a linker with WASI functions defined within it
	linker := wasmtime.NewLinker(engine)
	err = linker.DefineWasi()
	if err != nil {
		return nil, nil, err
	}

	// Set up the host functions we want to import
	r.AddHostFunctions(linker, map[string]exportFunc{
		"PrintHello": r.WrapExport(exports.PrintHello),
	})

	// Next up we instantiate a module which is where we link in all our
	// imports.
	r.instance, err = linker.Instantiate(r.store, module)
	if err != nil {
		return nil, nil, err
	}

	return r.instance, r.store, nil
}

// GetRequiredExports gets the expoerted WASM functions needed to make managed I/O work,
// these functions MUST be declared in the WASM module as exported functions as boilerplate,
// they are provided by a WasmModulePrototype instance.
func (r *Runner) GetRequiredExports(instance *wasmtime.Instance, store *wasmtime.Store) {
	r.mem = instance.GetExport(store, "memory").Memory()
	r.inputBufferFn = instance.GetExport(store, "inputBuffer").Func()
	r.outputBufferFn = instance.GetExport(store, "outputBuffer").Func()

	r.hostInputBufferFn = instance.GetExport(store, "hostInputBuffer").Func()
	r.hostOutputBufferFn = instance.GetExport(store, "hostOutputBuffer").Func()
}

// WarmUp will load and prepare a WASM module instance and create a call map for the
// runner to call, this means the wasm module can be warmed up in advance to minimise
// execution time of WASM funcs.
func (r *Runner) WarmUp(wasmFileName string, funcNames ...string) error {
	_, _, err := r.GetInstance(wasmFileName)
	if err != nil {
		return err
	}

	r.GetRequiredExports(r.instance, r.store)

	r.funcMap = make(map[string]*wasmtime.Func)
	for i, _ := range funcNames {
		r.funcMap[funcNames[i]] = r.instance.GetExport(r.store, funcNames[i]).Func()
	}

	return nil
}

// Run will call a function in the WASM module
func (r *Runner) Run(name string, args ...interface{}) (*shared_types.Payload, error) {
	fn, ok := r.funcMap[name]
	if !ok {
		return nil, fmt.Errorf("function name not found")
	}

	out := &shared_types.Payload{}
	err := ManagedCall(r.store, r.mem, r.inputBufferFn, r.outputBufferFn, fn, out, args...)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// ManagedCall handles all the I/O for calling an exported WASM mmodule function by reading
// and writing from the required WASM memory buffers and unmarshalling the output.
func ManagedCall(store wasmtime.Storelike, mem *wasmtime.Memory, inputBufferFn *wasmtime.Func, outputBufferFn *wasmtime.Func, guestFn *wasmtime.Func, output *shared_types.Payload, args ...interface{}) error {
	ptr, err := inputBufferFn.Call(store)
	if err != nil {
		return err
	}

	outPtr, err := outputBufferFn.Call(store)
	if err != nil {
		return err
	}

	stArgs := &shared_types.Args{
		Args: args,
	}

	enc, err := stArgs.MarshalMsg(nil)
	if err != nil {
		return err
	}

	inputLen := copy(mem.UnsafeData(store)[int(ptr.(int32)):int(ptr.(int32))+len(enc)], enc)

	dataLen, err := guestFn.Call(store, inputLen)
	if err != nil {
		return err
	}

	outDat := make([]byte, dataLen.(int32))
	copy(outDat[:], mem.UnsafeData(store)[int(outPtr.(int32)):int(outPtr.(int32))+int(dataLen.(int32))])

	buf := bytes.NewBuffer(outDat)
	err = msgp.Decode(buf, output)
	if err != nil {
		return err
	}

	return nil
}
