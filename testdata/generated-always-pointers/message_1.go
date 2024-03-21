// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

var _ bebop.Record = &ExampleMessage{}

type ExampleMessage struct {
	X *byte
	Y *int16
	Z *int32
}

func (bbp *ExampleMessage) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.X != nil {
		buf[at] = 1
		at++
		iohelp.WriteByteBytes(buf[at:], *bbp.X)
		at += 1
	}
	if bbp.Y != nil {
		buf[at] = 2
		at++
		iohelp.WriteInt16Bytes(buf[at:], *bbp.Y)
		at += 2
	}
	if bbp.Z != nil {
		buf[at] = 3
		at++
		iohelp.WriteInt32Bytes(buf[at:], *bbp.Z)
		at += 4
	}
	return at
}

func (bbp *ExampleMessage) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.X = new(byte)
			if len(buf[at:]) < 1 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.X) = iohelp.ReadByteBytes(buf[at:])
			at += 1
		case 2:
			at += 1
			bbp.Y = new(int16)
			if len(buf[at:]) < 2 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.Y) = iohelp.ReadInt16Bytes(buf[at:])
			at += 2
		case 3:
			at += 1
			bbp.Z = new(int32)
			if len(buf[at:]) < 4 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.Z) = iohelp.ReadInt32Bytes(buf[at:])
			at += 4
		default:
			return nil
		}
	}
}

func (bbp *ExampleMessage) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.X = new(byte)
			(*bbp.X) = iohelp.ReadByteBytes(buf[at:])
			at += 1
		case 2:
			at += 1
			bbp.Y = new(int16)
			(*bbp.Y) = iohelp.ReadInt16Bytes(buf[at:])
			at += 2
		case 3:
			at += 1
			bbp.Z = new(int32)
			(*bbp.Z) = iohelp.ReadInt32Bytes(buf[at:])
			at += 4
		default:
			return
		}
	}
}

func (bbp *ExampleMessage) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.X != nil {
		w.Write([]byte{1})
		iohelp.WriteByte(w, *bbp.X)
	}
	if bbp.Y != nil {
		w.Write([]byte{2})
		iohelp.WriteInt16(w, *bbp.Y)
	}
	if bbp.Z != nil {
		w.Write([]byte{3})
		iohelp.WriteInt32(w, *bbp.Z)
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *ExampleMessage) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.X = new(byte)
			*bbp.X = iohelp.ReadByte(r)
		case 2:
			bbp.Y = new(int16)
			*bbp.Y = iohelp.ReadInt16(r)
		case 3:
			bbp.Z = new(int32)
			*bbp.Z = iohelp.ReadInt32(r)
		default:
			r.Drain()
			return r.Err
		}
	}
}

func (bbp *ExampleMessage) Size() int {
	bodyLen := 5
	if bbp.X != nil {
		bodyLen += 1
		bodyLen += 1
	}
	if bbp.Y != nil {
		bodyLen += 1
		bodyLen += 2
	}
	if bbp.Z != nil {
		bodyLen += 1
		bodyLen += 4
	}
	return bodyLen
}

func (bbp *ExampleMessage) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeExampleMessage(r *iohelp.ErrorReader) (ExampleMessage, error) {
	v := ExampleMessage{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeExampleMessageFromBytes(buf []byte) (ExampleMessage, error) {
	v := ExampleMessage{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeExampleMessageFromBytes(buf []byte) ExampleMessage {
	v := ExampleMessage{}
	v.MustUnmarshalBebop(buf)
	return v
}
