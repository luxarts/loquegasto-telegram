package base64

import "encoding/base64"

func EncodeString(data string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(data))
}

func DecodeString(encoded string) string {
	data, _ := base64.RawURLEncoding.DecodeString(encoded)
	return string(data)
}

func Encode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func Decode(encoded string) []byte {
	data, _ := base64.RawURLEncoding.DecodeString(encoded)
	return data
}
