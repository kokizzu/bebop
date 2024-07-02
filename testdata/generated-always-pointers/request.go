// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

type FurnitureFamily uint32

const (
	FurnitureFamily_Bed FurnitureFamily = 0
	FurnitureFamily_Table FurnitureFamily = 1
	FurnitureFamily_Shoe FurnitureFamily = 2
)

var _ bebop.Record = &Furniture{}

type Furniture struct {
	name string
	price uint32
	family FurnitureFamily
}

func (bbp *Furniture) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.name)))
	copy(buf[at+4:at+4+len(bbp.name)], []byte(bbp.name))
	at += 4 + len(bbp.name)
	iohelp.WriteUint32Bytes(buf[at:], bbp.price)
	at += 4
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.family))
	at += 4
	return at
}

func (bbp *Furniture) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	bbp.name, err = iohelp.ReadStringBytes(buf[at:])
	if err != nil {
		return err
	}
	at += 4 + len(bbp.name)
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.price = iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	bbp.family = FurnitureFamily(iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	return nil
}

func (bbp *Furniture) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.name = iohelp.MustReadStringBytes(buf[at:])
	at += 4 + len(bbp.name)
	bbp.price = iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	bbp.family = FurnitureFamily(iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
}

func (bbp *Furniture) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(len(bbp.name)))
	w.Write([]byte(bbp.name))
	iohelp.WriteUint32(w, bbp.price)
	iohelp.WriteUint32(w, uint32(bbp.family))
	return w.Err
}

func (bbp *Furniture) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.name = iohelp.ReadString(r)
	bbp.price = iohelp.ReadUint32(r)
	bbp.family = FurnitureFamily(iohelp.ReadUint32(r))
	return r.Err
}

func (bbp *Furniture) Size() int {
	bodyLen := 0
	bodyLen += 4 + len(bbp.name)
	bodyLen += 4
	bodyLen += 4
	return bodyLen
}

func (bbp *Furniture) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeFurniture(r *iohelp.ErrorReader) (Furniture, error) {
	v := Furniture{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeFurnitureFromBytes(buf []byte) (Furniture, error) {
	v := Furniture{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeFurnitureFromBytes(buf []byte) Furniture {
	v := Furniture{}
	v.MustUnmarshalBebop(buf)
	return v
}

func (bbp *Furniture) GetName() string {
	return bbp.name
}

func (bbp *Furniture) GetPrice() uint32 {
	return bbp.price
}

func (bbp *Furniture) GetFamily() FurnitureFamily {
	return bbp.family
}

func NewFurniture(
		name string,
		price uint32,
		family FurnitureFamily,
) Furniture {
	return Furniture{
		name: name,
		price: price,
		family: family,
	}
}

const RequestResponseOpCode = 0x31323334

var _ bebop.Record = &RequestResponse{}

type RequestResponse struct {
	availableFurniture []Furniture
}

func (bbp *RequestResponse) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.availableFurniture)))
	at += 4
	for _, v1 := range bbp.availableFurniture {
		(v1).MarshalBebopTo(buf[at:])
		{
			tmp := (v1)
			at += tmp.Size()
		}
		
	}
	return at
}

func (bbp *RequestResponse) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	if len(buf[at:]) < 4 {
		return io.ErrUnexpectedEOF
	}
	bbp.availableFurniture = make([]Furniture, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.availableFurniture {
		(bbp.availableFurniture)[i1], err = MakeFurnitureFromBytes(buf[at:])
		if err != nil {
			return err
		}
		{
			tmp := ((bbp.availableFurniture)[i1])
			at += tmp.Size()
		}
		
	}
	return nil
}

func (bbp *RequestResponse) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.availableFurniture = make([]Furniture, iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	for i1 := range bbp.availableFurniture {
		(bbp.availableFurniture)[i1] = MustMakeFurnitureFromBytes(buf[at:])
		{
			tmp := ((bbp.availableFurniture)[i1])
			at += tmp.Size()
		}
		
	}
}

func (bbp *RequestResponse) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(len(bbp.availableFurniture)))
	for _, elem := range bbp.availableFurniture {
		err = (elem).EncodeBebop(w)
		if err != nil {
			return err
		}
	}
	return w.Err
}

func (bbp *RequestResponse) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.availableFurniture = make([]Furniture, iohelp.ReadUint32(r))
	for i1 := range bbp.availableFurniture {
		((bbp.availableFurniture[i1])), err = MakeFurniture(r)
		if err != nil {
			return err
		}
	}
	return r.Err
}

func (bbp *RequestResponse) Size() int {
	bodyLen := 0
	bodyLen += 4
	for _, elem := range bbp.availableFurniture {
		{
			tmp := (elem)
			bodyLen += tmp.Size()
		}
		
	}
	return bodyLen
}

func (bbp *RequestResponse) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeRequestResponse(r *iohelp.ErrorReader) (RequestResponse, error) {
	v := RequestResponse{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeRequestResponseFromBytes(buf []byte) (RequestResponse, error) {
	v := RequestResponse{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeRequestResponseFromBytes(buf []byte) RequestResponse {
	v := RequestResponse{}
	v.MustUnmarshalBebop(buf)
	return v
}

func (bbp *RequestResponse) GetAvailableFurniture() []Furniture {
	return bbp.availableFurniture
}

func NewRequestResponse(
		availableFurniture []Furniture,
) RequestResponse {
	return RequestResponse{
		availableFurniture: availableFurniture,
	}
}

const RequestCatalogOpCode = 0x41454b49

var _ bebop.Record = &RequestCatalog{}

type RequestCatalog struct {
	Family *FurnitureFamily
	// Deprecated: Nobody react to what I'm about to say...
	SecretTunnel *string
}

func (bbp *RequestCatalog) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.Family != nil {
		buf[at] = 1
		at++
		iohelp.WriteUint32Bytes(buf[at:], uint32(*bbp.Family))
		at += 4
	}
	return at
}

func (bbp *RequestCatalog) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Family = new(FurnitureFamily)
			(*bbp.Family) = FurnitureFamily(iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
		case 2:
			at += 1
			bbp.SecretTunnel = new(string)
			(*bbp.SecretTunnel), err = iohelp.ReadStringBytes(buf[at:])
			if err != nil {
				return err
			}
			at += 4 + len((*bbp.SecretTunnel))
		default:
			return nil
		}
	}
}

func (bbp *RequestCatalog) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Family = new(FurnitureFamily)
			(*bbp.Family) = FurnitureFamily(iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
		case 2:
			at += 1
			bbp.SecretTunnel = new(string)
			(*bbp.SecretTunnel) = iohelp.MustReadStringBytes(buf[at:])
			at += 4 + len((*bbp.SecretTunnel))
		default:
			return
		}
	}
}

func (bbp *RequestCatalog) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.Family != nil {
		w.Write([]byte{1})
		iohelp.WriteUint32(w, uint32(*bbp.Family))
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *RequestCatalog) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	limitReader := &io.LimitedReader{R: r.Reader, N: int64(bodyLen)}
	for {
		r.Reader = limitReader
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.Family = new(FurnitureFamily)
			*bbp.Family = FurnitureFamily(iohelp.ReadUint32(r))
		case 2:
			bbp.SecretTunnel = new(string)
			*bbp.SecretTunnel = iohelp.ReadString(r)
		default:
			r.Drain()
			return r.Err
		}
	}
}

func (bbp *RequestCatalog) Size() int {
	bodyLen := 5
	if bbp.Family != nil {
		bodyLen += 1
		bodyLen += 4
	}
	return bodyLen
}

func (bbp *RequestCatalog) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeRequestCatalog(r *iohelp.ErrorReader) (RequestCatalog, error) {
	v := RequestCatalog{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeRequestCatalogFromBytes(buf []byte) (RequestCatalog, error) {
	v := RequestCatalog{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeRequestCatalogFromBytes(buf []byte) RequestCatalog {
	v := RequestCatalog{}
	v.MustUnmarshalBebop(buf)
	return v
}

