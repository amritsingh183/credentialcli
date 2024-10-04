package util

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	mathRand "math/rand"
)

var srcForMathRand mathRand.Source

const (
	LetterBytesAlnum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	LetterSpecials   = "!@#$%^&*()_+=-/?.,><';:[]{}|`~`"

	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
	MaxKeyLength  = 100
	MinKeyLength  = 7
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
func GenerateKey(n int, includeSpecials bool) ([]byte, error) {
	if err := IsValidKeyLength(n); err != nil {
		return nil, err
	}
	if includeSpecials {
		var b1, b2 []byte
		b1, err := generate(n, LetterBytesAlnum)
		if err != nil {
			return nil, err
		}
		b2, err = generate(n, LetterSpecials)
		if err != nil {
			return nil, err
		}
		return append(b1[:n-5], b2[:5]...), nil
	} else {
		return generate(n, LetterBytesAlnum)
	}
}

// isValid validates the length of key to be generated
// TODO: this could also be replaced by some validation pkg.
// BUG: I'd like to generate unsecure keys. You should allow the users to do so.
// Maybe, you can print a warning that the password is not secure enough and the reason.
func IsValidKeyLength(n int) error {
	if n > MaxKeyLength {
		return fmt.Errorf("maximum key length is %d but %d was provided", MaxKeyLength, n)
	}
	if n < MinKeyLength {
		return fmt.Errorf("minimum key length is %d but %d was provided", MinKeyLength, n)
	}
	return nil
}

// TODO: this is really needed? Or we can solve it with an easier way?
// I mean the key and all of the other stuff. The code will be much more simplified.
func generate(n int, letterBytes string) ([]byte, error) {
	randBytes := make([]byte, 10240)
	_, err := io.ReadFull(cryptoRand.Reader, randBytes)
	if err != nil {
		return nil, err
	}
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
	return b, nil
}
