package password

import (
	"amritsingh183/password/internal/util"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	t.Run("Should respect the length option", func(t *testing.T) {
		passwordLength := 40
		options := Options{
			Length: uint(passwordLength),
			Count:  1,
		}
		passwords, err := Generate(&options)
		assert.NoError(t, err)
		usPtr := unsafe.Pointer(&passwords[0])
		strPtr := (*string)(usPtr)
		stringPassword := *strPtr
		assert.Equal(t, len(stringPassword), passwordLength)
	})

	t.Run("Should respect the default count option (case of 1 password)", func(t *testing.T) {
		passwordLength := 40
		count := 1
		options := Options{
			Length: uint(passwordLength),
			Count:  1,
		}
		passwords, err := Generate(&options)
		assert.NoError(t, err)
		assert.Equal(t, len(passwords), count)
	})
	t.Run("Should respect the default count option (case of >1 passwords)", func(t *testing.T) {
		passwordLength := 40
		count := 10
		options := Options{
			Length: uint(passwordLength),
			Count:  uint(count),
		}
		passwords, err := Generate(&options)
		assert.NoError(t, err)
		assert.Equal(t, len(passwords), count)
	})

	t.Run("Should respect IncludeSpecialChars=true option ", func(t *testing.T) {
		passwordLength := 40
		count := 1
		options := Options{
			Length:              uint(passwordLength),
			Count:               uint(count),
			IncludeSpecialChars: true,
		}
		passwords, err := Generate(&options)
		assert.NoError(t, err)
		usPtr := unsafe.Pointer(&passwords[0])
		strPtr := (*string)(usPtr)
		stringPassword := *strPtr
		var specialFound bool
		for _, c := range stringPassword {
			if strings.Contains(util.LetterSpecials, string(c)) {
				specialFound = true
			}
		}
		assert.Equal(t, true, specialFound)
	})

	t.Run("Should respect IncludeSpecialChars=false option ", func(t *testing.T) {
		passwordLength := 40
		count := 1
		options := Options{
			Length: uint(passwordLength),
			Count:  uint(count),
		}
		passwords, err := Generate(&options)
		assert.NoError(t, err)
		usPtr := unsafe.Pointer(&passwords[0])
		strPtr := (*string)(usPtr)
		stringPassword := *strPtr
		var specialFound bool
		for _, c := range stringPassword {
			if strings.Contains(util.LetterSpecials, string(c)) {
				specialFound = true
			}
		}
		assert.Equal(t, false, specialFound)
	})
}

func TestWrite(t *testing.T) {
	t.Run("Must write to stdout when asked to do so", func(t *testing.T) {
		tapper := tapStdOut{}
		tapper.Start()
		passwordLength := 40
		count := 1
		options := Options{
			Length:          uint(passwordLength),
			Count:           uint(count),
			DestinationType: ToStdOut,
		}
		passwords, err := Generate(&options)
		assert.NoError(t, err)
		Write(passwords, &options)
		output, err := tapper.Flush()
		assert.NoError(t, err)
		msg := fmt.Sprintf("Expected password to be of length %d", passwordLength)
		assert.Equal(t, passwordLength, len(output), msg)
	})

	t.Run("Must write to file when asked to do so", func(t *testing.T) {
		testDir := t.TempDir()
		passwdFilePath := fmt.Sprintf("%s/pass.txt", testDir)
		passwordLength := 40
		count := 1
		options := Options{
			Length:          uint(passwordLength),
			Count:           uint(count),
			DestinationType: ToFile,
			Filepath:        passwdFilePath,
		}
		passwords, err := Generate(&options)
		assert.NoError(t, err)
		Write(passwords, &options)
		buff, err := os.ReadFile(passwdFilePath)
		if err != nil {
			t.Error(err)
		}
		passwd := string(buff)
		assert.NoError(t, err)
		msg := fmt.Sprintf("Expected password to be of length %d", passwordLength)
		assert.Equal(t, passwordLength, len(passwd), msg)
	})

	t.Run("Must write to DefaultFilePath when no filepath is given but output to file is chosen", func(t *testing.T) {
		if _, err := os.Stat(DefaultFilePath); err == nil {
			t.Fatalf("%s file already exists. Please remove it first to avoid accidental data loss", DefaultFilePath)
		}
		defer os.Remove(DefaultFilePath)
		passwordLength := 40
		count := 1
		options := Options{
			Length:          uint(passwordLength),
			Count:           uint(count),
			DestinationType: ToFile,
			Filepath:        DefaultFilePath,
		}
		passwords, err := Generate(&options)
		assert.NoError(t, err)
		Write(passwords, &options)
		buff, err := os.ReadFile(DefaultFilePath)
		if err != nil {
			t.Error(err)
		}
		passwd := string(buff)
		assert.NoError(t, err)
		msg := fmt.Sprintf("Expected password to be of length %d", passwordLength)
		assert.Equal(t, passwordLength, len(passwd), msg)
	})

}

// tapStdOut provides mechanism to tap into the stdout
// useful for testing only
type tapStdOut struct {
	outChan chan string
	errChan chan error

	writeTo      *os.File
	stdOutbackup *os.File
}

// Start starts the tapping process and backsup stdout
func (tapper *tapStdOut) Start() error {
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
func (tapper *tapStdOut) read(readFrom *os.File) {
	var output bytes.Buffer
	_, err := io.Copy(&output, readFrom)
	tapper.errChan <- err
	tapper.outChan <- output.String()
}

// Flush Stops the tapping process, sends the stored output and restores stdout
func (tapper *tapStdOut) Flush() (string, error) {
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
