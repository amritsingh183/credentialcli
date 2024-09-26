package util

import (
	"bytes"
	cryptoRand "crypto/rand"
	"encoding/binary"
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
	letterBytes := LetterBytesAlnum
	if includeSpecials {
		letterBytes = LetterBytesAlnum + LetterSpecials
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
