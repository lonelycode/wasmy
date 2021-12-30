package main

import (
	"fmt"

	"github.com/lonelycode/wasmy/runner"
)

func main() {
	r := &runner.Runner{}
	// TODO: make this use application arguments
	err := r.WarmUp("/home/vmuser/wasmy/wasm-tests/managedv2.wasm", "myExport")
	if err != nil {
		panic(err)
	}

	// TODO: This breaks when there are no args, handle 0 len args
	out, err := r.Run("myExport", "martin")

	if err != nil {
		panic(err)
	}

	fmt.Printf("function output (from runner): %v \n", out)

}
