// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

var _ bebop.Record = &messageInline{}

type messageInline struct {
	onOff *bool
}

func (bbp messageInline) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.onOff != nil {
		buf[at] = 1
		at++
		iohelp.WriteBoolBytes(buf[at:], *bbp.onOff)
		at += 1
	}
	return at
}

func (bbp *messageInline) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.onOff = new(bool)
			if len(buf[at:]) < 1 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.onOff) = iohelp.ReadBoolBytes(buf[at:])
			at += 1
		default:
			return nil
		}
	}
}

func (bbp *messageInline) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.onOff = new(bool)
			(*bbp.onOff) = iohelp.ReadBoolBytes(buf[at:])
			at += 1
		default:
			return
		}
	}
}

func (bbp messageInline) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.onOff != nil {
		w.Write([]byte{1})
		iohelp.WriteBool(w, *bbp.onOff)
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *messageInline) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	baseReader := r.Reader
	r.Reader = &io.LimitedReader{R: baseReader, N: int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.onOff = new(bool)
			*bbp.onOff = iohelp.ReadBool(r)
		default:
			r.Drain()
			r.Reader = baseReader
			return r.Err
		}
	}
}

func (bbp messageInline) Size() int {
	bodyLen := 5
	if bbp.onOff != nil {
		bodyLen += 1
		bodyLen += 1
	}
	return bodyLen
}

func (bbp messageInline) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makemessageInline(r *iohelp.ErrorReader) (messageInline, error) {
	v := messageInline{}
	err := v.DecodeBebop(r)
	return v, err
}

func makemessageInlineFromBytes(buf []byte) (messageInline, error) {
	v := messageInline{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakemessageInlineFromBytes(buf []byte) messageInline {
	v := messageInline{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &messageInlineWrapper{}

type messageInlineWrapper struct {
	bla *messageInline
}

func (bbp messageInlineWrapper) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.bla != nil {
		buf[at] = 1
		at++
		(*bbp.bla).MarshalBebopTo(buf[at:])
		{
			tmp := (*bbp.bla)
			at += tmp.Size()
		}
		
	}
	return at
}

func (bbp *messageInlineWrapper) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.bla = new(messageInline)
			(*bbp.bla), err = makemessageInlineFromBytes(buf[at:])
			if err != nil {
				return err
			}
			{
				tmp := ((*bbp.bla))
				at += tmp.Size()
			}
			
		default:
			return nil
		}
	}
}

func (bbp *messageInlineWrapper) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.bla = new(messageInline)
			(*bbp.bla) = mustMakemessageInlineFromBytes(buf[at:])
			{
				tmp := ((*bbp.bla))
				at += tmp.Size()
			}
			
		default:
			return
		}
	}
}

func (bbp messageInlineWrapper) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.bla != nil {
		w.Write([]byte{1})
		err = (*bbp.bla).EncodeBebop(w)
		if err != nil {
			return err
		}
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *messageInlineWrapper) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	baseReader := r.Reader
	r.Reader = &io.LimitedReader{R: baseReader, N: int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.bla = new(messageInline)
			(*bbp.bla), err = makemessageInline(r)
			if err != nil {
				return err
			}
		default:
			r.Drain()
			r.Reader = baseReader
			return r.Err
		}
	}
}

func (bbp messageInlineWrapper) Size() int {
	bodyLen := 5
	if bbp.bla != nil {
		bodyLen += 1
		{
			tmp := (*bbp.bla)
			bodyLen += tmp.Size()
		}
		
	}
	return bodyLen
}

func (bbp messageInlineWrapper) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makemessageInlineWrapper(r *iohelp.ErrorReader) (messageInlineWrapper, error) {
	v := messageInlineWrapper{}
	err := v.DecodeBebop(r)
	return v, err
}

func makemessageInlineWrapperFromBytes(buf []byte) (messageInlineWrapper, error) {
	v := messageInlineWrapper{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakemessageInlineWrapperFromBytes(buf []byte) messageInlineWrapper {
	v := messageInlineWrapper{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &messageInlineWrapper2{}

type messageInlineWrapper2 struct {
	bla *[]messageInline
}

func (bbp messageInlineWrapper2) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.bla != nil {
		buf[at] = 1
		at++
		iohelp.WriteUint32Bytes(buf[at:], uint32(len(*bbp.bla)))
		at += 4
		for _, v2 := range *bbp.bla {
			(v2).MarshalBebopTo(buf[at:])
			{
				tmp := (v2)
				at += tmp.Size()
			}
			
		}
	}
	return at
}

func (bbp *messageInlineWrapper2) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.bla = new([]messageInline)
			if len(buf[at:]) < 4 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.bla) = make([]messageInline, iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
			for i3 := range (*bbp.bla) {
				((*bbp.bla))[i3], err = makemessageInlineFromBytes(buf[at:])
				if err != nil {
					return err
				}
				{
					tmp := (((*bbp.bla))[i3])
					at += tmp.Size()
				}
				
			}
		default:
			return nil
		}
	}
}

func (bbp *messageInlineWrapper2) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.bla = new([]messageInline)
			(*bbp.bla) = make([]messageInline, iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
			for i3 := range (*bbp.bla) {
				((*bbp.bla))[i3] = mustMakemessageInlineFromBytes(buf[at:])
				{
					tmp := (((*bbp.bla))[i3])
					at += tmp.Size()
				}
				
			}
		default:
			return
		}
	}
}

func (bbp messageInlineWrapper2) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.bla != nil {
		w.Write([]byte{1})
		iohelp.WriteUint32(w, uint32(len(*bbp.bla)))
		for _, elem := range *bbp.bla {
			err = (elem).EncodeBebop(w)
			if err != nil {
				return err
			}
		}
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *messageInlineWrapper2) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	baseReader := r.Reader
	r.Reader = &io.LimitedReader{R: baseReader, N: int64(bodyLen)}
	for {
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.bla = new([]messageInline)
			*bbp.bla = make([]messageInline, iohelp.ReadUint32(r))
			for i := range *bbp.bla {
				((*bbp.bla)[i]), err = makemessageInline(r)
				if err != nil {
					return err
				}
			}
		default:
			r.Drain()
			r.Reader = baseReader
			return r.Err
		}
	}
}

func (bbp messageInlineWrapper2) Size() int {
	bodyLen := 5
	if bbp.bla != nil {
		bodyLen += 1
		bodyLen += 4
		for _, elem := range *bbp.bla {
			{
				tmp := (elem)
				bodyLen += tmp.Size()
			}
			
		}
	}
	return bodyLen
}

func (bbp messageInlineWrapper2) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func makemessageInlineWrapper2(r *iohelp.ErrorReader) (messageInlineWrapper2, error) {
	v := messageInlineWrapper2{}
	err := v.DecodeBebop(r)
	return v, err
}

func makemessageInlineWrapper2FromBytes(buf []byte) (messageInlineWrapper2, error) {
	v := messageInlineWrapper2{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func mustMakemessageInlineWrapper2FromBytes(buf []byte) messageInlineWrapper2 {
	v := messageInlineWrapper2{}
	v.MustUnmarshalBebop(buf)
	return v
}

