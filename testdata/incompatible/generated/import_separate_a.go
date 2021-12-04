// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop/testdata/incompatible/generatedtwo"
	"io"
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
)

const (
	Go_package = "github.com/200sc/bebop/testdata/incompatible/generated"
)

var _ bebop.Record = &UsesImport{}

type UsesImport struct {
	Imported generatedtwo.ImportedType
}

func (bbp UsesImport) MarshalBebopTo(buf []byte) int {
	at := 0
	(bbp.Imported).MarshalBebopTo(buf[at:])
	at += (bbp.Imported).Size()
	return at
}

func (bbp *UsesImport) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	bbp.Imported, err = generatedtwo.MakeImportedTypeFromBytes(buf[at:])
	if err != nil{
		return err
	}
	at += (bbp.Imported).Size()
	return nil
}

func (bbp *UsesImport) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.Imported = generatedtwo.MustMakeImportedTypeFromBytes(buf[at:])
	at += (bbp.Imported).Size()
}

func (bbp UsesImport) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	err = (bbp.Imported).EncodeBebop(w)
	if err != nil{
		return err
	}
	return w.Err
}

func (bbp *UsesImport) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	(bbp.Imported), err = generatedtwo.MakeImportedType(r)
	if err != nil{
		return err
	}
	return r.Err
}

func (bbp UsesImport) Size() int {
	bodyLen := 0
	bodyLen += (bbp.Imported).Size()
	return bodyLen
}

func (bbp UsesImport) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeUsesImport(r iohelp.ErrorReader) (UsesImport, error) {
	v := UsesImport{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeUsesImportFromBytes(buf []byte) (UsesImport, error) {
	v := UsesImport{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeUsesImportFromBytes(buf []byte) UsesImport {
	v := UsesImport{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &UsesImportMsg{}

type UsesImportMsg struct {
	Imported *generatedtwo.ImportedType
}

func (bbp UsesImportMsg) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.Imported != nil {
		buf[at] = 1
		at++
		(*bbp.Imported).MarshalBebopTo(buf[at:])
		at += (*bbp.Imported).Size()
	}
	return at
}

func (bbp *UsesImportMsg) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Imported = new(generatedtwo.ImportedType)
			(*bbp.Imported), err = generatedtwo.MakeImportedTypeFromBytes(buf[at:])
			if err != nil{
				return err
			}
			at += ((*bbp.Imported)).Size()
		default:
			return nil
		}
	}
}

func (bbp *UsesImportMsg) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Imported = new(generatedtwo.ImportedType)
			(*bbp.Imported) = generatedtwo.MustMakeImportedTypeFromBytes(buf[at:])
			at += ((*bbp.Imported)).Size()
		default:
			return
		}
	}
}

func (bbp UsesImportMsg) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.Imported != nil {
		w.Write([]byte{1})
		err = (*bbp.Imported).EncodeBebop(w)
		if err != nil{
			return err
		}
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *UsesImportMsg) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.Imported = new(generatedtwo.ImportedType)
			(*bbp.Imported), err = generatedtwo.MakeImportedType(r)
			if err != nil{
				return err
			}
		default:
			io.ReadAll(r)
			return r.Err
		}
	}
}

func (bbp UsesImportMsg) Size() int {
	bodyLen := 5
	if bbp.Imported != nil {
		bodyLen += 1
		bodyLen += (*bbp.Imported).Size()
	}
	return bodyLen
}

func (bbp UsesImportMsg) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeUsesImportMsg(r iohelp.ErrorReader) (UsesImportMsg, error) {
	v := UsesImportMsg{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeUsesImportMsgFromBytes(buf []byte) (UsesImportMsg, error) {
	v := UsesImportMsg{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeUsesImportMsgFromBytes(buf []byte) UsesImportMsg {
	v := UsesImportMsg{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &UnionStruct{}

type UnionStruct struct {
	Hello generatedtwo.ImportedEnum
}

func (bbp UnionStruct) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Hello))
	at += 4
	return at
}

func (bbp *UnionStruct) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	bbp.Hello = generatedtwo.ImportedEnum(iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	return nil
}

func (bbp *UnionStruct) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.Hello = generatedtwo.ImportedEnum(iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
}

func (bbp UnionStruct) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Hello))
	return w.Err
}

func (bbp *UnionStruct) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.Hello = generatedtwo.ImportedEnum(iohelp.ReadUint32(r))
	return r.Err
}

func (bbp UnionStruct) Size() int {
	bodyLen := 0
	bodyLen += 4
	return bodyLen
}

func (bbp UnionStruct) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeUnionStruct(r iohelp.ErrorReader) (UnionStruct, error) {
	v := UnionStruct{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeUnionStructFromBytes(buf []byte) (UnionStruct, error) {
	v := UnionStruct{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeUnionStructFromBytes(buf []byte) UnionStruct {
	v := UnionStruct{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &UnionMessage{}

type UnionMessage struct {
	Goodbye *generatedtwo.ImportedEnum
}

func (bbp UnionMessage) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.Goodbye != nil {
		buf[at] = 1
		at++
		iohelp.WriteUint32Bytes(buf[at:], uint32(*bbp.Goodbye))
		at += 4
	}
	return at
}

func (bbp *UnionMessage) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Goodbye = new(generatedtwo.ImportedEnum)
			(*bbp.Goodbye) = generatedtwo.ImportedEnum(iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
		default:
			return nil
		}
	}
}

func (bbp *UnionMessage) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Goodbye = new(generatedtwo.ImportedEnum)
			(*bbp.Goodbye) = generatedtwo.ImportedEnum(iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
		default:
			return
		}
	}
}

func (bbp UnionMessage) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.Goodbye != nil {
		w.Write([]byte{1})
		iohelp.WriteUint32(w, uint32(*bbp.Goodbye))
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *UnionMessage) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.Goodbye = new(generatedtwo.ImportedEnum)
			*bbp.Goodbye = generatedtwo.ImportedEnum(iohelp.ReadUint32(r))
		default:
			io.ReadAll(r)
			return r.Err
		}
	}
}

func (bbp UnionMessage) Size() int {
	bodyLen := 5
	if bbp.Goodbye != nil {
		bodyLen += 1
		bodyLen += 4
	}
	return bodyLen
}

func (bbp UnionMessage) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeUnionMessage(r iohelp.ErrorReader) (UnionMessage, error) {
	v := UnionMessage{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeUnionMessageFromBytes(buf []byte) (UnionMessage, error) {
	v := UnionMessage{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeUnionMessageFromBytes(buf []byte) UnionMessage {
	v := UnionMessage{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &UsesImportUnion{}

type UsesImportUnion struct {
	UnionStruct *UnionStruct
	UnionMessage *UnionMessage
}

func (bbp UsesImportUnion) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-5))
	at += 4
	if bbp.UnionStruct != nil {
		buf[at] = 1
		at++
		(*bbp.UnionStruct).MarshalBebopTo(buf[at:])
		at += (*bbp.UnionStruct).Size()
		return at
	}
	if bbp.UnionMessage != nil {
		buf[at] = 2
		at++
		(*bbp.UnionMessage).MarshalBebopTo(buf[at:])
		at += (*bbp.UnionMessage).Size()
		return at
	}
	return at
}

func (bbp *UsesImportUnion) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	if len(buf) == 0 {
		return iohelp.ErrUnpopulatedUnion
	}
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.UnionStruct = new(UnionStruct)
			(*bbp.UnionStruct), err = MakeUnionStructFromBytes(buf[at:])
			if err != nil{
				return err
			}
			at += ((*bbp.UnionStruct)).Size()
			return nil
		case 2:
			at += 1
			bbp.UnionMessage = new(UnionMessage)
			(*bbp.UnionMessage), err = MakeUnionMessageFromBytes(buf[at:])
			if err != nil{
				return err
			}
			at += ((*bbp.UnionMessage)).Size()
			return nil
		default:
			return nil
		}
	}
}

func (bbp *UsesImportUnion) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.UnionStruct = new(UnionStruct)
			(*bbp.UnionStruct) = MustMakeUnionStructFromBytes(buf[at:])
			at += ((*bbp.UnionStruct)).Size()
			return
		case 2:
			at += 1
			bbp.UnionMessage = new(UnionMessage)
			(*bbp.UnionMessage) = MustMakeUnionMessageFromBytes(buf[at:])
			at += ((*bbp.UnionMessage)).Size()
			return
		default:
			return
		}
	}
}

func (bbp UsesImportUnion) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-5))
	if bbp.UnionStruct != nil {
		w.Write([]byte{1})
		err = (*bbp.UnionStruct).EncodeBebop(w)
		if err != nil{
			return err
		}
		return w.Err
	}
	if bbp.UnionMessage != nil {
		w.Write([]byte{2})
		err = (*bbp.UnionMessage).EncodeBebop(w)
		if err != nil{
			return err
		}
		return w.Err
	}
	return w.Err
}

func (bbp *UsesImportUnion) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	r.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)+1}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.UnionStruct = new(UnionStruct)
			(*bbp.UnionStruct), err = MakeUnionStruct(r)
			if err != nil{
				return err
			}
			io.ReadAll(r)
			return r.Err
		case 2:
			bbp.UnionMessage = new(UnionMessage)
			(*bbp.UnionMessage), err = MakeUnionMessage(r)
			if err != nil{
				return err
			}
			io.ReadAll(r)
			return r.Err
		default:
			io.ReadAll(r)
			return r.Err
		}
	}
}

func (bbp UsesImportUnion) Size() int {
	bodyLen := 4
	if bbp.UnionStruct != nil {
		bodyLen += 1
		bodyLen += (*bbp.UnionStruct).Size()
		return bodyLen
	}
	if bbp.UnionMessage != nil {
		bodyLen += 1
		bodyLen += (*bbp.UnionMessage).Size()
		return bodyLen
	}
	return bodyLen
}

func (bbp UsesImportUnion) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeUsesImportUnion(r iohelp.ErrorReader) (UsesImportUnion, error) {
	v := UsesImportUnion{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeUsesImportUnionFromBytes(buf []byte) (UsesImportUnion, error) {
	v := UsesImportUnion{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeUsesImportUnionFromBytes(buf []byte) UsesImportUnion {
	v := UsesImportUnion{}
	v.MustUnmarshalBebop(buf)
	return v
}

