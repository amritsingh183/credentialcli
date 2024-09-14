package util

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	mathRand "math/rand"
)

var srcForMathRand mathRand.Source

// FIXME: too many comments here.
// Put a section in the README.md file if you need this.

const (
	letterBytesAlnum    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	letterBytesSpecials = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()_+=-/?.,><';:[]{}|`~`"
	// 6 bits can encode 64 distinct characters
	// Hence we choose 6 bits to represent a letter index
	letterIdxBits = 6
	// 1 changed to 1000000 which is a 7 bit number
	// subtracting 1 from this 7 bit number gives
	// the largest 6 bit number ie 111111
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits // # of letter indices fitting in 63 bits
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

// GenerateShortID generates a password or a cryptographic key, etc
// side note: never use math/rand; use crypto/rand
func GenerateKey(n int, includeSpecials bool) []byte {
	letterBytes := letterBytesAlnum
	if includeSpecials {
		letterBytes = letterBytesSpecials
	}
	randBytes := make([]byte, 10240)
	io.ReadFull(cryptoRand.Reader, randBytes) // the error part is already handled in assertAvailablePRNG
	// extracts a 64-bit integer from the first 8 bytes of the random data and uses this as the seed for a math/rand.Source
	randSeed := int64(binary.LittleEndian.Uint64(randBytes[:]))
	srcForMathRand = mathRand.NewSource(randSeed)
	// creates a slice of n bytes to store the generated string
	b := make([]byte, n)
	cache := srcForMathRand.Int63()
	remain := letterIdxMax
	// The loop iterates over the byte slice b and fills it
	// with random alphanumeric characters.
	// It uses the math/rand package to generate random numbers,
	// which are then used to index into the letterBytes string to select a character
	for i := n - 1; i >= 0; {
		if remain == 0 {
			// int64 is a 64-bit signed integer type.
			// That means it has 1 sign bit and 63 significant bits.
			// which means that anything returning a non-negative int64
			// is producing 63 bits of data
			// the 64th bit, the sign bit, will always have the same value
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
