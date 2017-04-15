package lua

import "encoding/hex"

type Hex struct {
	EncodedLen     interface{}
	ErrLength      interface{}
	Dumper         interface{}
	Encode         interface{}
	Dump           interface{}
	EncodeToString interface{}
	DecodeString   interface{}
	DecodedLen     interface{}
	Decode         interface{}
}

func NewHex() *Hex {
	return &Hex{
		EncodeToString: hex.EncodeToString,
		EncodedLen:     hex.EncodedLen,
		Encode:         hex.Encode,
		Decode:         hex.Decode,
		ErrLength:      hex.ErrLength,
		Dump:           hex.Dump,
		DecodeString:   hex.DecodeString,
		Dumper:         hex.Dumper,
		DecodedLen:     hex.DecodedLen,
	}
}
