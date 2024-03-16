// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
	"time"
)

var _ bebop.Record = &MyObj{}

type MyObj struct {
	Start *time.Time
	End *time.Time
}

func (bbp *MyObj) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.Start != nil {
		buf[at] = 1
		at++
		if (*bbp.Start).IsZero() {
			iohelp.WriteInt64Bytes(buf[at:], 0)
		} else {
			iohelp.WriteInt64Bytes(buf[at:], ((*bbp.Start).UnixNano() / 100))
		}
		at += 8
	}
	if bbp.End != nil {
		buf[at] = 2
		at++
		if (*bbp.End).IsZero() {
			iohelp.WriteInt64Bytes(buf[at:], 0)
		} else {
			iohelp.WriteInt64Bytes(buf[at:], ((*bbp.End).UnixNano() / 100))
		}
		at += 8
	}
	return at
}

func (bbp *MyObj) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Start = new(time.Time)
			if len(buf[at:]) < 8 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.Start) = iohelp.ReadDateBytes(buf[at:])
			at += 8
		case 2:
			at += 1
			bbp.End = new(time.Time)
			if len(buf[at:]) < 8 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.End) = iohelp.ReadDateBytes(buf[at:])
			at += 8
		default:
			return nil
		}
	}
}

func (bbp *MyObj) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Start = new(time.Time)
			(*bbp.Start) = iohelp.ReadDateBytes(buf[at:])
			at += 8
		case 2:
			at += 1
			bbp.End = new(time.Time)
			(*bbp.End) = iohelp.ReadDateBytes(buf[at:])
			at += 8
		default:
			return
		}
	}
}

func (bbp *MyObj) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.Start != nil {
		w.Write([]byte{1})
		if (*bbp.Start).IsZero() {
			iohelp.WriteInt64(w, 0)
		} else {
			iohelp.WriteInt64(w, ((*bbp.Start).UnixNano() / 100))
		}
	}
	if bbp.End != nil {
		w.Write([]byte{2})
		if (*bbp.End).IsZero() {
			iohelp.WriteInt64(w, 0)
		} else {
			iohelp.WriteInt64(w, ((*bbp.End).UnixNano() / 100))
		}
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *MyObj) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.Start = new(time.Time)
			*bbp.Start = iohelp.ReadDate(r)
		case 2:
			bbp.End = new(time.Time)
			*bbp.End = iohelp.ReadDate(r)
		default:
			r.Drain()
			return r.Err
		}
	}
}

func (bbp *MyObj) Size() int {
	bodyLen := 5
	if bbp.Start != nil {
		bodyLen += 1
		bodyLen += 8
	}
	if bbp.End != nil {
		bodyLen += 1
		bodyLen += 8
	}
	return bodyLen
}

func (bbp *MyObj) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeMyObj(r *iohelp.ErrorReader) (MyObj, error) {
	v := MyObj{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeMyObjFromBytes(buf []byte) (MyObj, error) {
	v := MyObj{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeMyObjFromBytes(buf []byte) MyObj {
	v := MyObj{}
	v.MustUnmarshalBebop(buf)
	return v
}

