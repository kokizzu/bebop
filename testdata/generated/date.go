// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"bytes"
	"io"
	"time"

	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
)

var _ bebop.Record = &MyObj{}

type MyObj struct {
	Start *time.Time
	End *time.Time
}

func (bbp MyObj) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func (bbp MyObj) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.Start != nil {
		buf[at] = 1
		at++
		if *bbp.Start != (time.Time{}) {
			iohelp.WriteInt64Bytes(buf[at:], ((*bbp.Start).UnixNano()/100))
		} else {
			iohelp.WriteInt64Bytes(buf[at:], 0)
		}
		at += 8
	}
	if bbp.End != nil {
		buf[at] = 2
		at++
		if *bbp.End != (time.Time{}) {
			iohelp.WriteInt64Bytes(buf[at:], ((*bbp.End).UnixNano()/100))
		} else {
			iohelp.WriteInt64Bytes(buf[at:], 0)
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
				 return iohelp.ErrTooShort
			}
			(*bbp.Start) = iohelp.ReadDateBytes(buf[at:])
			at += 8
		case 2:
			at += 1
			bbp.End = new(time.Time)
			if len(buf[at:]) < 8 {
				 return iohelp.ErrTooShort
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

func (bbp MyObj) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.Start != nil {
		w.Write([]byte{1})
		if *bbp.Start != (time.Time{}) {
			iohelp.WriteInt64(w, ((*bbp.Start).UnixNano()/100))
		} else {
			iohelp.WriteInt64(w, 0)
		}
	}
	if bbp.End != nil {
		w.Write([]byte{2})
		if *bbp.End != (time.Time{}) {
			iohelp.WriteInt64(w, ((*bbp.End).UnixNano()/100))
		} else {
			iohelp.WriteInt64(w, 0)
		}
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *MyObj) DecodeBebop(ior io.Reader) (err error) {
	er := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(er)
	body := make([]byte, bodyLen)
	er.Read(body)
	r := iohelp.NewErrorReader(bytes.NewReader(body))
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.Start = new(time.Time)
			*bbp.Start = iohelp.ReadDate(r)
		case 2:
			bbp.End = new(time.Time)
			*bbp.End = iohelp.ReadDate(r)
		default:
			return er.Err
		}
	}
}

func (bbp MyObj) Size() int {
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

func makeMyObj(r iohelp.ErrorReader) (MyObj, error) {
	v := MyObj{}
	err := v.DecodeBebop(r)
	return v, err
}

func makeMyObjFromBytes(buf []byte) (MyObj, error) {
	v := MyObj{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakeMyObjFromBytes(buf []byte) MyObj {
	v := MyObj{}
	v.MustUnmarshalBebop(buf)
	return v
}

