package main

import (
	"fmt"

	"github.com/lonelycode/wasmy/shared_types"
)

// PrintHello is a sample function exported from the host and imported by the WASM module
// the function signature is fixed.
func PrintHello(args *shared_types.Args) (interface{}, error) {
	val := fmt.Sprintf("From Host: Hello Mr. %s", args.Args[0].(string))
	fmt.Println(val)

	return val, nil
}
