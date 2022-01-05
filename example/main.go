package main

import (
	"fmt"
	"time"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/lonelycode/wasmy/runner"
)

func call(engine *wasmtime.Engine, module *wasmtime.Module, arg string) {
	r := &runner.Runner{}

	// Let's export a sample function into the module
	r.HostFunctions = map[string]runner.ExportFunc{
		"PrintHello": r.WrapExport(PrintHello),
	}

	t3 := time.Now()
	err := r.WarmUp(engine, module, "myExport")
	if err != nil {
		panic(err)
	}
	t4 := time.Since(t3)
	fmt.Printf("warmup took %s\n", t4)

	// TODO: This breaks when there are no args, handle 0 len args
	t5 := time.Now()
	out, err := r.Run("myExport", arg)
	if err != nil {
		panic(err)
	}
	t6 := time.Since(t5)
	fmt.Printf("wasm run took %s\n", t6)
	fmt.Printf("function output (from runner): %v \n", out)
}

func main() {
	// TODO: make this use application arguments
	t1 := time.Now()
	engine := wasmtime.NewEngine()
	module, err := runner.GetModule("/home/vmuser/wasmy/wasm-tests/managedv2.wasm", engine)
	if err != nil {
		panic(err)
	}
	t2 := time.Since(t1)
	fmt.Printf("wasm load took %s\n", t2)

	go call(engine, module, "martin")
	go call(engine, module, "bianca")
	go call(engine, module, "ilrud")

	time.Sleep(5 * time.Second)

}
