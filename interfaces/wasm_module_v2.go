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
)

type WasmModulePrototype struct {
	guestFnInputBfr  [FUNCBUFFER_SIZE]uint8
	guestFnOutputBfr [FUNCBUFFER_SIZE]uint8

	hostFnInputBfr  [FUNCBUFFER_SIZE]uint8
	hostFnOutputBfr [FUNCBUFFER_SIZE]uint8
}

func (d *WasmModulePrototype) GetInputPtr() *[FUNCBUFFER_SIZE]uint8 {
	return &d.guestFnInputBfr
}

func (d *WasmModulePrototype) GetOutputPtr() *[FUNCBUFFER_SIZE]uint8 {
	return &d.guestFnOutputBfr
}

func (d *WasmModulePrototype) GetHostOutputPtr() *[FUNCBUFFER_SIZE]uint8 {
	return &d.hostFnOutputBfr
}

func (d *WasmModulePrototype) GetHostInputPtr() *[FUNCBUFFER_SIZE]uint8 {
	return &d.hostFnInputBfr
}

func (d *WasmModulePrototype) ReadGuestFnInput(length int) ([]interface{}, error) {
	dat := make([]byte, length)
	copy(dat, d.guestFnInputBfr[:length])

	args := &shared_types.Args{
		Args: make([]interface{}, 0),
	}

	buf := bytes.NewBuffer(dat)
	err := msgp.Decode(buf, args)
	if err != nil {
		return nil, err
	}

	return args.Args, nil
}

func (d *WasmModulePrototype) ReadHostFnOutput(length int, output *shared_types.Payload) error {
	dat := make([]byte, length)
	copy(dat, d.hostFnOutputBfr[:length])

	buf := bytes.NewBuffer(dat)
	err := msgp.Decode(buf, output)
	if err != nil {
		return err
	}

	return nil
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

func (d *WasmModulePrototype) WriteHostFnInput(args []interface{}) (int, error) {
	out := &shared_types.Args{Args: args}

	enc, err := out.MarshalMsg(nil)
	if err != nil {
		return 0, err
	}
	copy(d.hostFnInputBfr[:len(enc)], enc)
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

func CallImport(proto *WasmModulePrototype, fn func(int32) int32, args ...interface{}) (*shared_types.Payload, error) {
	// Write our args to the host input buffer
	lenInp, err := proto.WriteHostFnInput(args)
	if err != nil {
		return nil, err
	}

	// call the imported function with the length of the input data
	lenOut := fn(int32(lenInp))

	// Read the host output buffer for the return value
	output := &shared_types.Payload{}
	err = proto.ReadHostFnOutput(int(lenOut), output)
	if err != nil {
		return nil, err
	}

	return output, nil
}
