package password

import (
	"amritsingh183/password/internal/util"
	"fmt"
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
		passwords := Generate(&options)
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
		passwords := Generate(&options)
		assert.Equal(t, len(passwords), count)
	})
	t.Run("Should respect the default count option (case of >1 passwords)", func(t *testing.T) {
		passwordLength := 40
		count := 10
		options := Options{
			Length: uint(passwordLength),
			Count:  uint(count),
		}
		passwords := Generate(&options)
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
		passwords := Generate(&options)
		usPtr := unsafe.Pointer(&passwords[0])
		strPtr := (*string)(usPtr)
		stringPassword := *strPtr
		var specialFound bool
		for _, c := range stringPassword {
			// fmt.Println(string(c))
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
		passwords := Generate(&options)
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
		tapper := util.TapStdOut{}
		tapper.Start()
		passwordLength := 40
		count := 1
		options := Options{
			Length:          uint(passwordLength),
			Count:           uint(count),
			DestinationType: ToStdOut,
		}
		passwords := Generate(&options)
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
		passwords := Generate(&options)
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
		passwords := Generate(&options)
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
