package runner

import (
	"fmt"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
	"github.com/lonelycode/wasmy/shared_types"
)

type Runner struct {
	mem            *wasmtime.Memory
	store          *wasmtime.Store
	instance       *wasmtime.Instance
	inputBufferFn  *wasmtime.Func
	outputBufferFn *wasmtime.Func
	funcMap        map[string]*wasmtime.Func
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
