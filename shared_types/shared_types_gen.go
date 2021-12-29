package shared_types

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Args) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "args":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "Args")
				return
			}
			if cap(z.Args) >= int(zb0002) {
				z.Args = (z.Args)[:zb0002]
			} else {
				z.Args = make([]interface{}, zb0002)
			}
			for za0001 := range z.Args {
				z.Args[za0001], err = dc.ReadIntf()
				if err != nil {
					err = msgp.WrapError(err, "Args", za0001)
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Args) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "args"
	err = en.Append(0x81, 0xa4, 0x61, 0x72, 0x67, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Args)))
	if err != nil {
		err = msgp.WrapError(err, "Args")
		return
	}
	for za0001 := range z.Args {
		err = en.WriteIntf(z.Args[za0001])
		if err != nil {
			err = msgp.WrapError(err, "Args", za0001)
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Args) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "args"
	o = append(o, 0x81, 0xa4, 0x61, 0x72, 0x67, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Args)))
	for za0001 := range z.Args {
		o, err = msgp.AppendIntf(o, z.Args[za0001])
		if err != nil {
			err = msgp.WrapError(err, "Args", za0001)
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Args) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "args":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Args")
				return
			}
			if cap(z.Args) >= int(zb0002) {
				z.Args = (z.Args)[:zb0002]
			} else {
				z.Args = make([]interface{}, zb0002)
			}
			for za0001 := range z.Args {
				z.Args[za0001], bts, err = msgp.ReadIntfBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Args", za0001)
					return
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Args) Msgsize() (s int) {
	s = 1 + 5 + msgp.ArrayHeaderSize
	for za0001 := range z.Args {
		s += msgp.GuessSize(z.Args[za0001])
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Payload) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "data":
			z.Data, err = dc.ReadIntf()
			if err != nil {
				err = msgp.WrapError(err, "Data")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Payload) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "data"
	err = en.Append(0x81, 0xa4, 0x64, 0x61, 0x74, 0x61)
	if err != nil {
		return
	}
	err = en.WriteIntf(z.Data)
	if err != nil {
		err = msgp.WrapError(err, "Data")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Payload) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "data"
	o = append(o, 0x81, 0xa4, 0x64, 0x61, 0x74, 0x61)
	o, err = msgp.AppendIntf(o, z.Data)
	if err != nil {
		err = msgp.WrapError(err, "Data")
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Payload) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "data":
			z.Data, bts, err = msgp.ReadIntfBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Data")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Payload) Msgsize() (s int) {
	s = 1 + 5 + msgp.GuessSize(z.Data)
	return
}
