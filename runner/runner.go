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

type exportFunc func(int32, int32, int32) int32

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

func (r *Runner) AddHostFunctions(linker *wasmtime.Linker, funcs map[string]exportFunc) {
	for name, fn := range funcs {
		linker.DefineFunc(r.store, "env", fmt.Sprintf("main.%s", name), fn)
	}
}

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

func (r *Runner) GetRequiredExports(instance *wasmtime.Instance, store *wasmtime.Store) {
	r.mem = instance.GetExport(store, "memory").Memory()
	r.inputBufferFn = instance.GetExport(store, "inputBuffer").Func()
	r.outputBufferFn = instance.GetExport(store, "outputBuffer").Func()

	r.hostInputBufferFn = instance.GetExport(store, "hostInputBuffer").Func()
	r.hostOutputBufferFn = instance.GetExport(store, "hostOutputBuffer").Func()
}

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
