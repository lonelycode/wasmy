package main

import (
	"bytes"
	"fmt"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
	"github.com/lonelycode/wasmy/shared_types"
	"github.com/tinylib/msgp/msgp"
)

func sample_main() {

	engine := wasmtime.NewEngine()
	// Almost all operations in wasmtime require a contextual `store`
	// argument to share, so create that first
	store := wasmtime.NewStore(engine)

	// Once we have our binary `wasm` we can compile that into a `*Module`
	// which represents compiled JIT code.
	module, err := wasmtime.NewModuleFromFile(store.Engine, "wasm-tests/module.wasm")
	check(err)

	wConf := wasmtime.NewWasiConfig()
	wConf.SetEnv([]string{"doo"}, []string{"dod"})
	wConf.InheritStdout()
	store.SetWasi(wConf)

	check(err)

	// Create a linker with WASI functions defined within it
	linker := wasmtime.NewLinker(engine)

	var mem *wasmtime.Memory
	var ptrFn *wasmtime.Func

	// this is how we define external reference functions!
	linker.DefineFunc(store, "env", "main.hello", func(t1 int32, t2 int32) {
		fmt.Println("hello from the outside called from the inside")
		ptr, err := ptrFn.Call(store)
		check(err)
		fmt.Println("read from host (from inside host)")
		fmt.Println(mem.UnsafeData(store)[int(ptr.(int32))])
	})

	err = linker.DefineWasi()
	check(err)

	// Next up we instantiate a module which is where we link in all our
	// imports.
	instance, err := linker.Instantiate(store, module)
	check(err)

	fmt.Println("getting externs")
	fmt.Println("--memory")
	mem = instance.GetExport(store, "memory").Memory()
	fmt.Println("--write")
	wrt := instance.GetExport(store, "storeValueInWasmMemoryBufferIndexZero").Func()
	ptrFn = instance.GetExport(store, "getWasmMemoryBufferPointer").Func()
	red := instance.GetExport(store, "readWasmMemoryBufferAndReturnIndexOne").Func()

	fmt.Println("write in the wasm function")
	_, err = wrt.Call(store, 42)
	check(err)

	ptr, err := ptrFn.Call(store)
	check(err)

	// fmt.Println("read from host")
	// fmt.Println(mem.UnsafeData(store)[int(ptr.(int32))])

	fmt.Println("write from host")
	mem.UnsafeData(store)[int(ptr.(int32))+1] = 12
	//fmt.Printf("%v\n", val)

	fmt.Println("read from guest")
	val, err := red.Call(store)
	check(err)

	fmt.Println(val)

}

func main_v2() {
	engine := wasmtime.NewEngine()
	// Almost all operations in wasmtime require a contextual `store`
	// argument to share, so create that first
	store := wasmtime.NewStore(engine)

	// Once we have our binary `wasm` we can compile that into a `*Module`
	// which represents compiled JIT code.
	module, err := wasmtime.NewModuleFromFile(store.Engine, "wasm-tests/managed.wasm")
	check(err)

	wConf := wasmtime.NewWasiConfig()
	wConf.InheritStdout()
	store.SetWasi(wConf)
	check(err)

	// Create a linker with WASI functions defined within it
	linker := wasmtime.NewLinker(engine)

	var mem *wasmtime.Memory
	var ptrFn *wasmtime.Func

	err = linker.DefineWasi()
	check(err)

	// Next up we instantiate a module which is where we link in all our
	// imports.
	instance, err := linker.Instantiate(store, module)
	check(err)

	fmt.Println("getting externs")
	fmt.Println("--memory")
	mem = instance.GetExport(store, "memory").Memory()

	fmt.Println("-- write test")
	wrt := instance.GetExport(store, "storeValueInWasmMemoryBuffer").Func()
	fmt.Println("--memory pointer")
	ptrFn = instance.GetExport(store, "getMemoryPtr").Func()
	fmt.Println("-- read test")
	red := instance.GetExport(store, "readWasmMemoryBuffer").Func()

	fmt.Println("-- cursor map pointer")
	getCursorForDataFn := instance.GetExport(store, "getCursorForData").Func()
	fmt.Println("-- cursor start map pointer")
	getStartCursorForDataFn := instance.GetExport(store, "getStartCursorForData").Func()
	fmt.Println("-- cursor end map pointer")
	getEndCursorForDataFn := instance.GetExport(store, "getEndCursorForData").Func()

	fmt.Println("TEST: write a value in the wasm module")
	_, err = wrt.Call(store)
	check(err)

	iPtr, err := ptrFn.Call(store)
	check(err)

	fmt.Println("TEST: read from host")

	// the data is at id 1
	locPtr, err := getCursorForDataFn.Call(store, 1)
	check(err)
	fmt.Printf("map data location: %v\n", locPtr)

	st, err := getStartCursorForDataFn.Call(store, 1)
	check(err)
	ln, err := getEndCursorForDataFn.Call(store, 1)
	check(err)

	relStart := st.(int32) + iPtr.(int32)
	relEnd := relStart + ln.(int32)
	js := mem.UnsafeData(store)[relStart:relEnd]

	fmt.Printf("data from memory buffer (called from host): %v\n", js)
	fmt.Println(string(js))

	fmt.Println("TEST: read mem location from guest")
	_, err = red.Call(store)
	check(err)
}

func main() {
	engine := wasmtime.NewEngine()
	// Almost all operations in wasmtime require a contextual `store`
	// argument to share, so create that first
	store := wasmtime.NewStore(engine)

	// Once we have our binary `wasm` we can compile that into a `*Module`
	// which represents compiled JIT code.
	module, err := wasmtime.NewModuleFromFile(store.Engine, "wasm-tests/managedv2.wasm")
	check(err)

	wConf := wasmtime.NewWasiConfig()
	wConf.InheritStdout()
	wConf.InheritStderr()
	store.SetWasi(wConf)
	check(err)

	// Create a linker with WASI functions defined within it
	linker := wasmtime.NewLinker(engine)

	var mem *wasmtime.Memory
	err = linker.DefineWasi()
	check(err)

	// Next up we instantiate a module which is where we link in all our
	// imports.
	instance, err := linker.Instantiate(store, module)
	check(err)

	fmt.Println("getting externs")
	fmt.Println("--memory")
	mem = instance.GetExport(store, "memory").Memory()
	funcInputPtrFn := instance.GetExport(store, "inputBuffer").Func()
	funcOutputPtrFn := instance.GetExport(store, "outputBuffer").Func()
	testFn := instance.GetExport(store, "myExport").Func()

	out := &shared_types.Payload{}
	err = ManagedCall(store, mem, funcInputPtrFn, funcOutputPtrFn, testFn, out, "martin")
	check(err)

	fmt.Println(out.Data)

}

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
		fmt.Println("output decode failed")
		return err
	}

	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
