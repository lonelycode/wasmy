package runner

import (
	"bytes"
	"fmt"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
	"github.com/lonelycode/wasmy/shared_types"
	"github.com/tinylib/msgp/msgp"
)

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
