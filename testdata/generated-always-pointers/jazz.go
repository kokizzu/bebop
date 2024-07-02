// Code generated by bebopc-go; DO NOT EDIT.

package generated

import (
	"github.com/200sc/bebop"
	"github.com/200sc/bebop/iohelp"
	"io"
)

type Instrument uint32

const (
	Instrument_Sax Instrument = 0
	Instrument_Trumpet Instrument = 1
	Instrument_Clarinet Instrument = 2
)

var _ bebop.Record = &Musician{}

type Musician struct {
	name string
	plays Instrument
}

func (bbp *Musician) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.name)))
	copy(buf[at+4:at+4+len(bbp.name)], []byte(bbp.name))
	at += 4 + len(bbp.name)
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.plays))
	at += 4
	return at
}

func (bbp *Musician) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	bbp.name, err = iohelp.ReadStringBytes(buf[at:])
	if err != nil {
		return err
	}
	at += 4 + len(bbp.name)
	bbp.plays = Instrument(iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
	return nil
}

func (bbp *Musician) MustUnmarshalBebop(buf []byte) {
	at := 0
	bbp.name = iohelp.MustReadStringBytes(buf[at:])
	at += 4 + len(bbp.name)
	bbp.plays = Instrument(iohelp.ReadUint32Bytes(buf[at:]))
	at += 4
}

func (bbp *Musician) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(len(bbp.name)))
	w.Write([]byte(bbp.name))
	iohelp.WriteUint32(w, uint32(bbp.plays))
	return w.Err
}

func (bbp *Musician) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bbp.name = iohelp.ReadString(r)
	bbp.plays = Instrument(iohelp.ReadUint32(r))
	return r.Err
}

func (bbp *Musician) Size() int {
	bodyLen := 0
	bodyLen += 4 + len(bbp.name)
	bodyLen += 4
	return bodyLen
}

func (bbp *Musician) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeMusician(r *iohelp.ErrorReader) (Musician, error) {
	v := Musician{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeMusicianFromBytes(buf []byte) (Musician, error) {
	v := Musician{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeMusicianFromBytes(buf []byte) Musician {
	v := Musician{}
	v.MustUnmarshalBebop(buf)
	return v
}

func (bbp *Musician) GetName() string {
	return bbp.name
}

func (bbp *Musician) GetPlays() Instrument {
	return bbp.plays
}

func NewMusician(
		name string,
		plays Instrument,
) Musician {
	return Musician{
		name: name,
		plays: plays,
	}
}

var _ bebop.Record = &Library{}

type Library struct {
	Songs map[[16]byte]Song
}

func (bbp *Library) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(len(bbp.Songs)))
	at += 4
	for k1, v1 := range bbp.Songs {
		iohelp.WriteGUIDBytes(buf[at:], k1)
		at += 16
		(v1).MarshalBebopTo(buf[at:])
		{
			tmp := (v1)
			at += tmp.Size()
		}
		
	}
	return at
}

func (bbp *Library) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	ln1 := iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	bbp.Songs = make(map[[16]byte]Song,ln1)
	for i := uint32(0); i < ln1; i++ {
		if len(buf[at:]) < 16 {
			return io.ErrUnexpectedEOF
		}
		k1 := iohelp.ReadGUIDBytes(buf[at:])
		at += 16
		(bbp.Songs)[k1], err = MakeSongFromBytes(buf[at:])
		if err != nil {
			return err
		}
		{
			tmp := ((bbp.Songs)[k1])
			at += tmp.Size()
		}
		
	}
	return nil
}

func (bbp *Library) MustUnmarshalBebop(buf []byte) {
	at := 0
	ln2 := iohelp.ReadUint32Bytes(buf[at:])
	at += 4
	bbp.Songs = make(map[[16]byte]Song,ln2)
	for i := uint32(0); i < ln2; i++ {
		k1 := iohelp.ReadGUIDBytes(buf[at:])
		at += 16
		(bbp.Songs)[k1] = MustMakeSongFromBytes(buf[at:])
		{
			tmp := ((bbp.Songs)[k1])
			at += tmp.Size()
		}
		
	}
}

func (bbp *Library) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(len(bbp.Songs)))
	for k1, v1 := range bbp.Songs {
		iohelp.WriteGUID(w, k1)
		err = (v1).EncodeBebop(w)
		if err != nil {
			return err
		}
	}
	return w.Err
}

func (bbp *Library) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	ln1 := iohelp.ReadUint32(r)
	bbp.Songs = make(map[[16]byte]Song, ln1)
	for i1 := uint32(0); i1 < ln1; i1++ {
		k1 := iohelp.ReadGUID(r)
		((bbp.Songs[k1])), err = MakeSong(r)
		if err != nil {
			return err
		}
	}
	return r.Err
}

func (bbp *Library) Size() int {
	bodyLen := 0
	bodyLen += 4
	for _, v1 := range bbp.Songs {
		bodyLen += 16
		{
			tmp := (v1)
			bodyLen += tmp.Size()
		}
		
	}
	return bodyLen
}

func (bbp *Library) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeLibrary(r *iohelp.ErrorReader) (Library, error) {
	v := Library{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeLibraryFromBytes(buf []byte) (Library, error) {
	v := Library{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeLibraryFromBytes(buf []byte) Library {
	v := Library{}
	v.MustUnmarshalBebop(buf)
	return v
}

var _ bebop.Record = &Song{}

type Song struct {
	Title *string
	Year *uint16
	Performers *[]Musician
}

func (bbp *Song) MarshalBebopTo(buf []byte) int {
	at := 0
	iohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))
	at += 4
	if bbp.Title != nil {
		buf[at] = 1
		at++
		iohelp.WriteUint32Bytes(buf[at:], uint32(len(*bbp.Title)))
		copy(buf[at+4:at+4+len(*bbp.Title)], []byte(*bbp.Title))
		at += 4 + len(*bbp.Title)
	}
	if bbp.Year != nil {
		buf[at] = 2
		at++
		iohelp.WriteUint16Bytes(buf[at:], *bbp.Year)
		at += 2
	}
	if bbp.Performers != nil {
		buf[at] = 3
		at++
		iohelp.WriteUint32Bytes(buf[at:], uint32(len(*bbp.Performers)))
		at += 4
		for _, v2 := range *bbp.Performers {
			(v2).MarshalBebopTo(buf[at:])
			{
				tmp := (v2)
				at += tmp.Size()
			}
			
		}
	}
	return at
}

func (bbp *Song) UnmarshalBebop(buf []byte) (err error) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Title = new(string)
			(*bbp.Title), err = iohelp.ReadStringBytes(buf[at:])
			if err != nil {
				return err
			}
			at += 4 + len((*bbp.Title))
		case 2:
			at += 1
			bbp.Year = new(uint16)
			if len(buf[at:]) < 2 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.Year) = iohelp.ReadUint16Bytes(buf[at:])
			at += 2
		case 3:
			at += 1
			bbp.Performers = new([]Musician)
			if len(buf[at:]) < 4 {
				return io.ErrUnexpectedEOF
			}
			(*bbp.Performers) = make([]Musician, iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
			for i3 := range (*bbp.Performers) {
				((*bbp.Performers))[i3], err = MakeMusicianFromBytes(buf[at:])
				if err != nil {
					return err
				}
				{
					tmp := (((*bbp.Performers))[i3])
					at += tmp.Size()
				}
				
			}
		default:
			return nil
		}
	}
}

func (bbp *Song) MustUnmarshalBebop(buf []byte) {
	at := 0
	_ = iohelp.ReadUint32Bytes(buf[at:])
	buf = buf[4:]
	for {
		switch buf[at] {
		case 1:
			at += 1
			bbp.Title = new(string)
			(*bbp.Title) = iohelp.MustReadStringBytes(buf[at:])
			at += 4 + len((*bbp.Title))
		case 2:
			at += 1
			bbp.Year = new(uint16)
			(*bbp.Year) = iohelp.ReadUint16Bytes(buf[at:])
			at += 2
		case 3:
			at += 1
			bbp.Performers = new([]Musician)
			(*bbp.Performers) = make([]Musician, iohelp.ReadUint32Bytes(buf[at:]))
			at += 4
			for i3 := range (*bbp.Performers) {
				((*bbp.Performers))[i3] = MustMakeMusicianFromBytes(buf[at:])
				{
					tmp := (((*bbp.Performers))[i3])
					at += tmp.Size()
				}
				
			}
		default:
			return
		}
	}
}

func (bbp *Song) EncodeBebop(iow io.Writer) (err error) {
	w := iohelp.NewErrorWriter(iow)
	iohelp.WriteUint32(w, uint32(bbp.Size()-4))
	if bbp.Title != nil {
		w.Write([]byte{1})
		iohelp.WriteUint32(w, uint32(len(*bbp.Title)))
		w.Write([]byte(*bbp.Title))
	}
	if bbp.Year != nil {
		w.Write([]byte{2})
		iohelp.WriteUint16(w, *bbp.Year)
	}
	if bbp.Performers != nil {
		w.Write([]byte{3})
		iohelp.WriteUint32(w, uint32(len(*bbp.Performers)))
		for _, elem := range *bbp.Performers {
			err = (elem).EncodeBebop(w)
			if err != nil {
				return err
			}
		}
	}
	w.Write([]byte{0})
	return w.Err
}

func (bbp *Song) DecodeBebop(ior io.Reader) (err error) {
	r := iohelp.NewErrorReader(ior)
	bodyLen := iohelp.ReadUint32(r)
	limitReader := &io.LimitedReader{R: r.Reader, N: int64(bodyLen)}
	for {
		r.Reader = limitReader
		switch iohelp.ReadByte(r) {
		case 1:
			bbp.Title = new(string)
			*bbp.Title = iohelp.ReadString(r)
		case 2:
			bbp.Year = new(uint16)
			*bbp.Year = iohelp.ReadUint16(r)
		case 3:
			bbp.Performers = new([]Musician)
			*bbp.Performers = make([]Musician, iohelp.ReadUint32(r))
			for i := range *bbp.Performers {
				((*bbp.Performers)[i]), err = MakeMusician(r)
				if err != nil {
					return err
				}
			}
		default:
			r.Drain()
			return r.Err
		}
	}
}

func (bbp *Song) Size() int {
	bodyLen := 5
	if bbp.Title != nil {
		bodyLen += 1
		bodyLen += 4 + len(*bbp.Title)
	}
	if bbp.Year != nil {
		bodyLen += 1
		bodyLen += 2
	}
	if bbp.Performers != nil {
		bodyLen += 1
		bodyLen += 4
		for _, elem := range *bbp.Performers {
			{
				tmp := (elem)
				bodyLen += tmp.Size()
			}
			
		}
	}
	return bodyLen
}

func (bbp *Song) MarshalBebop() []byte {
	buf := make([]byte, bbp.Size())
	bbp.MarshalBebopTo(buf)
	return buf
}

func MakeSong(r *iohelp.ErrorReader) (Song, error) {
	v := Song{}
	err := v.DecodeBebop(r)
	return v, err
}

func MakeSongFromBytes(buf []byte) (Song, error) {
	v := Song{}
	err := v.UnmarshalBebop(buf)
	return v, err
}

func MustMakeSongFromBytes(buf []byte) Song {
	v := Song{}
	v.MustUnmarshalBebop(buf)
	return v
}

