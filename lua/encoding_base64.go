package lua

import "encoding/base64"

type Base64 struct {
	NewEncoder     interface{}
	RawURLEncoding interface{}
	NewEncoding    interface{}
	NewDecoder     interface{}
	StdEncoding    interface{}
	RawStdEncoding interface{}
	URLEncoding    interface{}
}

func NewBase64() *Base64 {
	return &Base64{
		RawStdEncoding: base64.RawStdEncoding,
		StdEncoding:    base64.StdEncoding,
		NewEncoding:    base64.NewEncoding,
		URLEncoding:    base64.URLEncoding,
		NewEncoder:     base64.NewEncoder,
		NewDecoder:     base64.NewDecoder,
		RawURLEncoding: base64.RawURLEncoding,
	}
}
