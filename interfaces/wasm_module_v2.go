package interfaces

import (
	"bytes"
	"fmt"
	"os"

	"github.com/lonelycode/wasmy/shared_types"
	"github.com/tinylib/msgp/msgp"
)

const (
	FUNCBUFFER_SIZE = 1024
	DATA_SIZE       = 1024
	OUTPUT_SIZE     = 1024
)

type WasmModulePrototype struct {
	// TODO: use structured IO for guest/host interop
	guestFnInputBfr  [FUNCBUFFER_SIZE]uint8
	guestFnOutputBfr [FUNCBUFFER_SIZE]uint8

	hostFnOutputtBfr [FUNCBUFFER_SIZE]uint8
}

func (d *WasmModulePrototype) GetInputPtr() *[FUNCBUFFER_SIZE]uint8 {
	return &d.guestFnInputBfr
}

func (d *WasmModulePrototype) GetOutputPtr() *[FUNCBUFFER_SIZE]uint8 {
	return &d.guestFnOutputBfr
}

func (d *WasmModulePrototype) ReadGuestFnInput(length int) ([]interface{}, error) {
	fmt.Printf("input length is: %v\n", length)
	dat := make([]byte, length)
	copy(dat, d.guestFnInputBfr[:length])

	fmt.Println("INPUT BUFFER")
	fmt.Println(string(dat))
	// fmt.Println(dat)

	args := &shared_types.Args{
		Args: make([]interface{}, 0),
	}

	fmt.Println("unmarshalling")
	buf := bytes.NewBuffer(dat)
	err := msgp.Decode(buf, args)
	if err != nil {
		return nil, err
	}

	fmt.Println("unmarshalled")

	return args.Args, nil
}

func (d *WasmModulePrototype) WriteGuestFnOutput(data interface{}) (int, error) {
	out := &shared_types.Payload{Data: data}

	enc, err := out.MarshalMsg(nil)
	if err != nil {
		return 0, err
	}
	copy(d.guestFnOutputBfr[:len(enc)], enc)
	return len(enc), nil
}

func (d *WasmModulePrototype) externGuestErr(err error) int {
	errTp := fmt.Sprintf("ERR %s", err.Error())
	os.Stderr.WriteString(errTp)
	copy(d.guestFnOutputBfr[:], []byte(errTp))

	return len([]byte(errTp))
}

func WrapExport(proto *WasmModulePrototype, inputLen int, exportFn func(args ...interface{}) (interface{}, error)) func() int {
	return func() int {
		args, err := proto.ReadGuestFnInput(inputLen)
		if err != nil {
			return proto.externGuestErr(err)
		}

		fmt.Println("calling wrapped function")
		ret, err := exportFn(args...)
		if err != nil {
			return proto.externGuestErr(err)
		}

		n, err := proto.WriteGuestFnOutput(ret)
		if err != nil {
			return proto.externGuestErr(err)
		}

		return n
	}
}
