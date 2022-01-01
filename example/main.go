package main

import (
	"fmt"
	"time"

	"github.com/lonelycode/wasmy/runner"
)

func main() {
	r := &runner.Runner{}

	// Let's export a sample function into the module
	r.HostFunctions = map[string]runner.ExportFunc{
		"PrintHello": r.WrapExport(PrintHello),
	}

	// TODO: make this use application arguments
	start := time.Now()
	err := r.WarmUp("/home/vmuser/wasmy/wasm-tests/managedv2.wasm", "myExport")
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	fmt.Printf("wasm load took %s\n", elapsed)

	// TODO: This breaks when there are no args, handle 0 len args
	start = time.Now()
	out, err := r.Run("myExport", "martin")
	if err != nil {
		panic(err)
	}
	elapsed = time.Since(start)
	fmt.Printf("wasm run took %s\n", elapsed)
	fmt.Printf("function output (from runner): %v \n", out)

}
