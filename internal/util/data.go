package util

import (
	"bytes"
	cryptoRand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	mathRand "math/rand"
	"os"
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
	if n > MaxKeyLength {
		return nil, fmt.Errorf("key can not exceed length=%d", MaxKeyLength)
	}
	if n < MinKeyLength {
		return nil, fmt.Errorf("key can not be smaller than length=%d", MinKeyLength)
	}
	if n == 0 {
		return nil, errors.New("key can not of length=0")
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

// TapStdOut provides mechanism to tap into the stdout
type TapStdOut struct {
	outChan chan string
	errChan chan error

	writeTo      *os.File
	stdOutbackup *os.File
}

// Start starts the tapping process and backsup stdout
func (tapper *TapStdOut) Start() error {
	tapper.stdOutbackup = os.Stdout
	rf, wf, err := os.Pipe()
	if err != nil {
		return err
	}
	os.Stdout = wf
	tapper.writeTo = wf
	tapper.outChan = make(chan string)
	tapper.errChan = make(chan error)
	go tapper.read(rf)
	return nil
}

// read reads from readpipe into channel of tapper
// must be called in a go routine to prevent
// blocking writes to writepipe (such as stdout)
func (tapper *TapStdOut) read(readFrom *os.File) {
	var output bytes.Buffer
	_, err := io.Copy(&output, readFrom)
	tapper.errChan <- err
	tapper.outChan <- output.String()
}

// Flush Stops the tapping process, sends the stored output and restores stdout
func (tapper *TapStdOut) Flush() (string, error) {
	err := tapper.writeTo.Close()
	if err != nil {
		return "", err
	}
	err = <-tapper.errChan
	if err != nil {
		return "", err
	}
	os.Stdout = tapper.stdOutbackup
	return <-tapper.outChan, nil
}

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
