package util

import "encoding/base64"

// Base64URLEncode Base64URLEncode
// more efficient than hex encoding
func Base64URLEncode(b []byte) string {
	return base64.URLEncoding.EncodeToString(b)
}

// Base64URLDecode Base64URLDecode
func Base64URLDecode(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}
