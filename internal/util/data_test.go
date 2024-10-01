package util

import (
	"fmt"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// TestGenerateKey Although this function is extensively tested by the password package
// we test it here as an extra safety net
func TestGenerateKey(t *testing.T) {
	t.Run("Must generate n length key", func(t *testing.T) {
		passwordLength := 10
		passBytes, err := GenerateKey(passwordLength, false)
		assert.NoError(t, err)
		assert.NotNil(t, passBytes, "Generated password must not be nil")
		assert.NotEqual(t, 0, len(passBytes))
		usPtr := unsafe.Pointer(&passBytes)
		strPtr := (*string)(usPtr)
		stringPassword := *strPtr
		assert.Equal(t, passwordLength, len(stringPassword))
	})

	t.Run("Zero length should not be allowed", func(t *testing.T) {
		passwordLength := 0
		passBytes, err := GenerateKey(passwordLength, false)
		assert.Error(t, err)
		assert.Nil(t, passBytes, "Generated password must be nil")
		assert.Equal(t, 0, len(passBytes))
	})

	t.Run("should not allow passwords longer than the allowed maxlength", func(t *testing.T) {
		passwordLength := MaxKeyLength + 1
		passBytes, err := GenerateKey(passwordLength, false)
		assert.Error(t, err)
		assert.Nil(t, passBytes, "Generated password must be nil")
		assert.Equal(t, 0, len(passBytes))
	})
	t.Run("should not allow passwords less than the min length", func(t *testing.T) {
		passwordLength := MinKeyLength - 1
		passBytes, err := GenerateKey(passwordLength, false)
		assert.Error(t, err)
		assert.Nil(t, passBytes, "Generated password must be nil")
		assert.Equal(t, 0, len(passBytes))
	})

	t.Run("Should respect IncludeSpecialChars=true option ", func(t *testing.T) {
		passwordLength := 10
		passBytes, err := GenerateKey(passwordLength, true)
		assert.NoError(t, err)
		assert.NotNil(t, passBytes, "Generated password must not be nil")
		assert.NotEqual(t, 0, len(passBytes))

		usPtr := unsafe.Pointer(&passBytes)
		strPtr := (*string)(usPtr)
		stringPassword := *strPtr
		var specialFound bool
		for _, c := range stringPassword {
			if strings.Contains(LetterSpecials, string(c)) {
				specialFound = true
			}
		}
		assert.Equal(t, true, specialFound, fmt.Sprintf("generated password was %v", stringPassword))
	})

	t.Run("Should respect IncludeSpecialChars=false option ", func(t *testing.T) {
		passwordLength := 10
		passBytes, err := GenerateKey(passwordLength, false)
		assert.NoError(t, err)
		assert.NotNil(t, passBytes, "Generated password must not be nil")
		assert.NotEqual(t, 0, len(passBytes))

		usPtr := unsafe.Pointer(&passBytes)
		strPtr := (*string)(usPtr)
		stringPassword := *strPtr
		var specialFound bool
		for _, c := range stringPassword {
			if strings.Contains(LetterSpecials, string(c)) {
				specialFound = true
			}
		}
		assert.Equal(t, false, specialFound, fmt.Sprintf("generated password was %v", stringPassword))
	})
}
