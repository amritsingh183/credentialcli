package util

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	mathRand "math/rand"
	"os"
)

var srcForMathRand mathRand.Source

// FIXME: too many comments here.
// Put a section in the README.md file if you need this.
// [x] removed all comments
const (
	letterBytesAlnum    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	letterBytesSpecials = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()_+=-/?.,><';:[]{}|`~`"
	letterIdxBits       = 6
	letterIdxMask       = 1<<letterIdxBits - 1
	letterIdxMax        = 63 / letterIdxBits
)

func init() {
	// 100 is chosen without any specific reason
	assertAvailablePRNG(100)
}

// assertAvailablePRNG Assert that a CSPRNG is available.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should NOT continue.
func assertAvailablePRNG(n uint) {
	buf := make([]byte, n)
	_, err := io.ReadFull(cryptoRand.Reader, buf)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: Read() failed with %#v", err))
	}
}

// GenerateShortID generates a password or a cryptographic key
func GenerateKey(n int, includeSpecials bool) []byte {
	letterBytes := letterBytesAlnum
	if includeSpecials {
		letterBytes = letterBytesSpecials
	}
	randBytes := make([]byte, 10240)
	io.ReadFull(cryptoRand.Reader, randBytes)
	randSeed := int64(binary.LittleEndian.Uint64(randBytes[:]))
	srcForMathRand = mathRand.NewSource(randSeed)
	b := make([]byte, n)
	cache := srcForMathRand.Int63()
	remain := letterIdxMax
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache = srcForMathRand.Int63()
			remain = letterIdxMax
		}
		idx := int(cache & letterIdxMask)
		if idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return b
}

func CreateFile(destinationFilePath string) (*os.File, error) {
	return os.OpenFile(destinationFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
}
