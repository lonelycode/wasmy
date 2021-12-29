package main

import (
	"fmt"

	"github.com/lonelycode/wasmy/runner"
)

func main() {
	r := &runner.Runner{}
	err := r.WarmUp("/home/vmuser/wasmy/wasm-tests/managedv2.wasm", "myExport")
	if err != nil {
		panic(err)
	}

	out, err := r.Run("myExport", "martin")

	if err != nil {
		panic(err)
	}

	fmt.Printf("function output (from runner): %v \n", out.Data)

}
