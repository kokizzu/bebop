package generated_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/200sc/bebop/testdata/generated"
	"github.com/200sc/bebop/testdata/generated/protos"
	"google.golang.org/protobuf/proto"
)

var benchArray = &generated.BasicArrays{
	A_bool: []bool{
		true, false, true, false,
	},
	A_byte: []byte{
		0, 1, 2, 3, 4, 5, 6, 7, 8,
	},
	A_int16: []int16{
		0, 1, 2, 3, 4, 5, 6, 7, 8,
	},
	A_uint16: []uint16{
		0, 1, 2, 3, 4, 5, 6, 7, 8,
	},
	A_int32: []int32{
		0, 1, 234436345, 3, 4, 5, 634, 7, 8,
	},
	A_uint32: []uint32{
		0, 1, 2, 33453566, 4, 5, 634634, 7, 8,
	},
	A_int64: []int64{
		3436453450, 346345346531, 3463453452, 3, 4, 5346345345, 34634566, 7, 8,
	},
	A_uint64: []uint64{
		0, 1, 2, 3, 34634563454, 5, 6334534634, 7, 8,
	},
	A_float32: []float32{
		0, 341, 2, 34563453, 4, 5, 6, 7, 8,
	},
	A_float64: []float64{
		0, 1, 2, 345343, 3453464, 3453453635, 353453456, 7, 8555555555,
	},
	A_string: []string{
		"0123151234123123", "11234125123415124", "223412512341512341254", "31245123151234125123413", "1231251231512315124", "124123151234151234125", "61231512341541234123", "12315123412512341257", "81231451241234151234151",
	},
}
var benchArrayBytes = []byte{4, 0, 0, 0, 1, 0, 1, 0, 9, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 0, 0, 0, 0, 1, 0, 2, 0,
	3, 0, 4, 0, 5, 0, 6, 0, 7, 0, 8, 0, 9, 0, 0, 0, 0, 0, 1, 0, 2, 0, 3, 0, 4, 0, 5, 0, 6, 0, 7, 0, 8, 0, 9, 0, 0, 0, 0,
	0, 0, 0, 1, 0, 0, 0, 249, 54, 249, 13, 3, 0, 0, 0, 4, 0, 0, 0, 5, 0, 0, 0, 122, 2, 0, 0, 7, 0, 0, 0, 8, 0, 0, 0, 9,
	0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 254, 117, 254, 1, 4, 0, 0, 0, 5, 0, 0, 0, 10, 175, 9, 0, 7, 0, 0, 0, 8,
	0, 0, 0, 9, 0, 0, 0, 74, 30, 212, 204, 0, 0, 0, 0, 227, 149, 202, 163, 80, 0, 0, 0, 12, 27, 112, 206, 0, 0, 0, 0, 3,
	0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 129, 193, 170, 62, 1, 0, 0, 0, 70, 123, 16, 2, 0, 0, 0, 0, 7, 0, 0, 0,
	0, 0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,
	0, 3, 0, 0, 0, 0, 0, 0, 0, 126, 127, 97, 16, 8, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 234, 83, 145, 121, 1, 0, 0, 0, 7, 0,
	0, 0, 0, 0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 9, 0, 0, 0, 0, 0, 0, 0, 0, 128, 170, 67, 0, 0, 0, 64, 95, 217, 3, 76, 0, 0,
	128, 64, 0, 0, 160, 64, 0, 0, 192, 64, 0, 0, 224, 64, 0, 0, 0, 65, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 240, 63, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 252, 19, 21, 65, 0, 0, 0, 0, 12, 89, 74, 65, 0, 0, 96, 168, 240, 186,
	233, 65, 0, 0, 0, 144, 69, 17, 181, 65, 0, 0, 0, 0, 0, 0, 28, 64, 0, 0, 48, 174, 54, 223, 255, 65, 9, 0, 0, 0, 16, 0,
	0, 0, 48, 49, 50, 51, 49, 53, 49, 50, 51, 52, 49, 50, 51, 49, 50, 51, 17, 0, 0, 0, 49, 49, 50, 51, 52, 49, 50, 53, 49,
	50, 51, 52, 49, 53, 49, 50, 52, 21, 0, 0, 0, 50, 50, 51, 52, 49, 50, 53, 49, 50, 51, 52, 49, 53, 49, 50, 51, 52, 49, 50,
	53, 52, 23, 0, 0, 0, 51, 49, 50, 52, 53, 49, 50, 51, 49, 53, 49, 50, 51, 52, 49, 50, 53, 49, 50, 51, 52, 49, 51, 19, 0,
	0, 0, 49, 50, 51, 49, 50, 53, 49, 50, 51, 49, 53, 49, 50, 51, 49, 53, 49, 50, 52, 21, 0, 0, 0, 49, 50, 52, 49, 50, 51,
	49, 53, 49, 50, 51, 52, 49, 53, 49, 50, 51, 52, 49, 50, 53, 20, 0, 0, 0, 54, 49, 50, 51, 49, 53, 49, 50, 51, 52, 49,
	53, 52, 49, 50, 51, 52, 49, 50, 51, 20, 0, 0, 0, 49, 50, 51, 49, 53, 49, 50, 51, 52, 49, 50, 53, 49, 50, 51, 52, 49,
	50, 53, 55, 23, 0, 0, 0, 56, 49, 50, 51, 49, 52, 53, 49, 50, 52, 49, 50, 51, 52, 49, 53, 49, 50, 51, 52, 49, 53, 49,
	0, 0, 0, 0}

var benchArrayJSON = []byte(`{"A_bool":[true,false,true,false],"A_byte":"AAECAwQFBgcI","A_int16":[0,1,2,3,4,5,6,7,8],"A_uint16":[0,1,2,3,4,5,6,7,8],"A_int32":[0,1,234436345,3,4,5,634,7,8],"A_uint32":[0,1,2,33453566,4,5,634634,7,8],"A_int64":[3436453450,346345346531,3463453452,3,4,5346345345,34634566,7,8],"A_uint64":[0,1,2,3,34634563454,5,6334534634,7,8],"A_float32":[0,341,2,34563452,4,5,6,7,8],"A_float64":[0,1,2,345343,3453464,3453453635,353453456,7,8555555555],"A_string":["0123151234123123","11234125123415124","223412512341512341254","31245123151234125123413","1231251231512315124","124123151234151234125","61231512341541234123","12315123412512341257","81231451241234151234151"],"A_guid":null}`)

var benchArray2 *generated.BasicArrays

func BenchmarkMarshalBasicArrays(b *testing.B) {
	for i := 0; i < b.N; i++ {
		out = benchArray.MarshalBebop()
	}
}

func BenchmarkMarshalUnmarshalBasicArrays(b *testing.B) {
	var w = &bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		benchArray.EncodeBebop(w)
		benchArray2 = &generated.BasicArrays{}
		benchArray2.DecodeBebop(w)
	}
}

func BenchmarkUnmarshalBasicArrays(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchArray2 = &generated.BasicArrays{}
		benchArray2.DecodeBebop(bytes.NewBuffer(benchArrayBytes))
	}
}

func BenchmarkMarshalBasicArraysJSON(b *testing.B) {
	var w = &bytes.Buffer{}
	encoder := json.NewEncoder(w)
	for i := 0; i < b.N; i++ {
		encoder.Encode(benchArray)
	}
}

func BenchmarkMarshalUnmarshalBasicArraysJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		decoder := json.NewDecoder(bytes.NewBuffer(benchArrayJSON))
		decoder.Decode(&benchArray2)
	}
}

func BenchmarkUnmarshalBasicArraysJSON(b *testing.B) {
	var w = &bytes.Buffer{}
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(w)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		err := encoder.Encode(benchArray)
		if err != nil {
			panic(err)
		}
		b.StartTimer()
		decoder.Decode(&benchArray2)
	}
}

var basicArraysProto = &protos.BasicArrays{
	ABool: []bool{
		true, false, true, false,
	},
	AByte: []byte{
		0, 1, 2, 3, 4, 5, 6, 7, 8,
	},
	AInt16: []int32{
		0, 1, 2, 3, 4, 5, 6, 7, 8,
	},
	AUint16: []uint32{
		0, 1, 2, 3, 4, 5, 6, 7, 8,
	},
	AInt32: []int32{
		0, 1, 234436345, 3, 4, 5, 634, 7, 8,
	},
	AUint32: []uint32{
		0, 1, 2, 33453566, 4, 5, 634634, 7, 8,
	},
	AInt64: []int64{
		3436453450, 346345346531, 3463453452, 3, 4, 5346345345, 34634566, 7, 8,
	},
	AUint64: []uint64{
		0, 1, 2, 3, 34634563454, 5, 6334534634, 7, 8,
	},
	AFloat32: []float32{
		0, 341, 2, 34563453, 4, 5, 6, 7, 8,
	},
	AFloat64: []float64{
		0, 1, 2, 345343, 3453464, 3453453635, 353453456, 7, 8555555555,
	},
	AString: []string{
		"0123151234123123", "11234125123415124", "223412512341512341254", "31245123151234125123413", "1231251231512315124", "124123151234151234125", "61231512341541234123", "12315123412512341257", "81231451241234151234151",
	},
}

var out []byte

func BenchmarkMarshalBasicArraysProto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		out, _ = proto.Marshal(basicArraysProto)
	}
}

var benchMaps = &generated.SomeMaps{
	M1: map[bool]bool{
		true:  true,
		false: false,
	},
	M2: map[string]map[string]string{
		"hello": {
			"world": "!",
		},
		"foo": {
			"bar": "bizz",
		},
		"ursula": {
			"k": "leguin",
		},
		"mario": {
			"mario":    "",
			"luigi":    "",
			"brothers": "",
		},
	},
	M3: []map[int32][]map[bool]generated.S{
		{
			0: []map[bool]generated.S{{
				true:  generated.S{},
				false: generated.S{},
			}},
		}, {
			1: []map[bool]generated.S{{
				true:  generated.S{},
				false: generated.S{},
			}},
			2: []map[bool]generated.S{{
				true:  generated.S{},
				false: generated.S{},
			}},
			3: []map[bool]generated.S{{
				true:  generated.S{},
				false: generated.S{},
			}},
		}, {
			41111: []map[bool]generated.S{{
				true:  generated.S{},
				false: generated.S{},
			}},
		},
	},
	M4: []map[string][]float32{
		{
			"a": []float32{1321, 1423, 1423, 540, 12314, 1231, 4123, 1412, 1230, 4123, 123},
		},
	},
	M5: map[[16]byte]generated.M{
		{5: 3}: {B: float64p(0.0000002)},
	},
}

var benchMapsBytes = []byte{2, 0, 0, 0, 1, 1, 0, 0, 4, 0, 0, 0, 5, 0, 0, 0, 104, 101, 108, 108, 111, 1, 0, 0, 0, 5, 0,
	0, 0, 119, 111, 114, 108, 100, 1, 0, 0, 0, 33, 3, 0, 0, 0, 102, 111, 111, 1, 0, 0, 0, 3, 0, 0, 0, 98, 97, 114, 4, 0,
	0, 0, 98, 105, 122, 122, 6, 0, 0, 0, 117, 114, 115, 117, 108, 97, 1, 0, 0, 0, 1, 0, 0, 0, 107, 6, 0, 0, 0, 108, 101,
	103, 117, 105, 110, 5, 0, 0, 0, 109, 97, 114, 105, 111, 3, 0, 0, 0, 5, 0, 0, 0, 109, 97, 114, 105, 111, 0, 0, 0, 0,
	5, 0, 0, 0, 108, 117, 105, 103, 105, 0, 0, 0, 0, 8, 0, 0, 0, 98, 114, 111, 116, 104, 101, 114, 115, 0, 0, 0, 0, 3,
	0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0,
	0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 1, 0, 0,
	0, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 1, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 151, 160, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 97, 11, 0, 0, 0, 0, 32, 165, 68, 0, 224, 177,
	68, 0, 224, 177, 68, 0, 0, 7, 68, 0, 104, 64, 70, 0, 224, 153, 68, 0, 216, 128, 69, 0, 128, 176, 68, 0, 192, 153,
	68, 0, 216, 128, 69, 0, 0, 246, 66, 1, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 2, 72,
	175, 188, 154, 242, 215, 138, 62, 0}

var benchMaps2 *generated.SomeMaps

func BenchmarkMarshalSomeMap(b *testing.B) {
	var w = &bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		benchMaps.EncodeBebop(w)
	}
}

func BenchmarkMarshalUnmarshalSomeMap(b *testing.B) {
	var w = &bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		benchMaps.EncodeBebop(w)
		benchMaps2 = &generated.SomeMaps{}
		benchMaps2.DecodeBebop(w)
	}
}

func BenchmarkUnmarshalSomeMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchMaps2 = &generated.SomeMaps{}
		benchMaps2.DecodeBebop(bytes.NewBuffer(benchMapsBytes))
	}
}

// SomeMap cannot be parsed by json as its structure is unsupported
