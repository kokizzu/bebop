// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"io"
	"time"

	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
)

var _ bebop.Record = &BasicTypes{}

type BasicTypes struct {
	A_bool bool
	A_byte byte
	A_int16 int16
	A_uint16 uint16
	A_int32 int32
	A_uint32 uint32
	A_int64 int64
	A_uint64 uint64
	A_float32 float32
	A_float64 float64
	A_string string
	A_guid [16]byte
	A_date time.Time
}

func (bbp BasicTypes) MarshalBebop() []byte {
	buf := make([]byte, bbp.bodyLen())
	bbp.MarshalBebopTo(buf)
	return buf
}

func (bbp BasicTypes) MarshalBebopTo(buf []byte) {
	at := 0
	iohelp.WriteBoolBytes(buf[at:], bbp.A_bool)
	at += 1
	iohelp.WriteByteBytes(buf[at:], bbp.A_byte)
	at += 1
	iohelp.WriteInt16Bytes(buf[at:], bbp.A_int16)
	at += 2
	iohelp.WriteUint16Bytes(buf[at:], bbp.A_uint16)
	at += 2
	iohelp.WriteInt32Bytes(buf[at:], bbp.A_int32)
	at += 4
	iohelp.WriteUint32Bytes(buf[at:], bbp.A_uint32)
	at += 4
	iohelp.WriteInt64Bytes(buf[at:], bbp.A_int64)
	at += 8
	iohelp.WriteUint64Bytes(buf[at:], bbp.A_uint64)
	at += 8
	iohelp.WriteFloat32Bytes(buf[at:], bbp.A_float32)
	at += 4
	iohelp.WriteFloat64Bytes(buf[at:], bbp.A_float64)
	at += 8
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.A_string)))
	at += 4
	copy(buf[at:at+len(bbp.A_string)], []byte(bbp.A_string))
	at += len(bbp.A_string)
	iohelp.WriteGUIDBytes(buf[at:], bbp.A_guid)
	at += 16
	if bbp.A_date != (time.Time{}) {
		iohelp.WriteInt64Bytes(buf[at:], (bbp.A_date.UnixNano()/100))
	} else {
		iohelp.WriteInt64Bytes(buf[at:], 0)
	}
	at += 8
}

func (bbp BasicTypes) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteBool(w, bbp.A_bool)
	iohelp.WriteByte(w, bbp.A_byte)
	iohelp.WriteInt16(w, bbp.A_int16)
	iohelp.WriteUint16(w, bbp.A_uint16)
	iohelp.WriteInt32(w, bbp.A_int32)
	iohelp.WriteUint32(w, bbp.A_uint32)
	iohelp.WriteInt64(w, bbp.A_int64)
	iohelp.WriteUint64(w, bbp.A_uint64)
	iohelp.WriteFloat32(w, bbp.A_float32)
	iohelp.WriteFloat64(w, bbp.A_float64)
	iohelp.WriteUint32(w, uint32(len(bbp.A_string)))
	w.Write([]byte(bbp.A_string))
	iohelp.WriteGUID(w, bbp.A_guid)
	if bbp.A_date != (time.Time{}) {
		iohelp.WriteInt64(w, (bbp.A_date.UnixNano()/100))
	} else {
		iohelp.WriteInt64(w, 0)
	}
	return w.Err
}

func (bbp *BasicTypes) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.A_bool = iohelp.ReadBool(r)
	bbp.A_byte = iohelp.ReadByte(r)
	bbp.A_int16 = iohelp.ReadInt16(r)
	bbp.A_uint16 = iohelp.ReadUint16(r)
	bbp.A_int32 = iohelp.ReadInt32(r)
	bbp.A_uint32 = iohelp.ReadUint32(r)
	bbp.A_int64 = iohelp.ReadInt64(r)
	bbp.A_uint64 = iohelp.ReadUint64(r)
	bbp.A_float32 = iohelp.ReadFloat32(r)
	bbp.A_float64 = iohelp.ReadFloat64(r)
	bbp.A_string = iohelp.ReadString(r)
	bbp.A_guid = iohelp.ReadGUID(r)
	bbp.A_date = iohelp.ReadTime(r)
	return r.Err
}

func (bbp *BasicTypes) bodyLen() int {
	bodyLen := 0
	bodyLen += 1
	bodyLen += 1
	bodyLen += 2
	bodyLen += 2
	bodyLen += 4
	bodyLen += 4
	bodyLen += 8
	bodyLen += 8
	bodyLen += 4
	bodyLen += 8
	bodyLen += 4
	bodyLen += len(bbp.A_string)
	bodyLen += 16
	bodyLen += 8
	return bodyLen
}

func makeBasicTypes(r iohelp.ErrorReader) (BasicTypes, error) {
	v := BasicTypes{}
	err := v.DecodeBebop(r)
	return v, err
}

