package cmd

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"amritsingh183/password/internal/password"

	"github.com/stretchr/testify/assert"
)

const (
	// useful to see which all flags were passed to the cli
	debug = false
)

// FIXME: there is no MaxLength() function in the source code.
// Keep 1:1 relationship between the test code function and the source code ones.
// Use something like "testify/suite" to tidy things up.
// [x] Renamed it to better represent the relationship => Wrong
// FIXME: you're testing the Cobra CLI commands here. IMHO, you should have also tested the method Generator.Generate which is harder to test since it's a method and not a function
// [x] Makes sense. internal/password now has all the relevant tests
func TestRunPasswordGenerator(t *testing.T) {
	t.Run("Should error if password exceeds max length", func(t *testing.T) {
		generateCmd.SetOutput(os.Stdout)
		generateCmd.SetArgs([]string{
			"--length=102",
		})
		if debug {
			generateCmd.DebugFlags()
		}
		_, err := generateCmd.ExecuteC()
		msg := fmt.Sprintf("Should not allow length greater than %d", password.MaxPasswordLength)
		assert.Error(t, err, msg)
	})

	// FIXME: fix this test. If you stick to interface you can write the password to an in-memory structure and then check it for data.
	// The same applies for the next test functions.
	// [x] Generate writes either to file or stdout. here in this file I have used OS's temp dir as output. Do you see an issue with this approach? another option is to use stdout like internal/password/password_test.go.TestWrite test case 1
	t.Run("Should respect the length flag", func(t *testing.T) {
		testDir := t.TempDir()
		passwdFilePath := fmt.Sprintf("%s/pass.txt", testDir)
		generateCmd.SetOutput(os.Stdout)
		passwdLength := 20
		generateCmd.SetArgs([]string{
			"password",
			fmt.Sprintf("--length=%d", passwdLength),
			"--output=1",
			fmt.Sprintf("--file=%s", passwdFilePath),
		})
		if debug {
			generateCmd.DebugFlags()
		}
		_, err := generateCmd.ExecuteC()
		assert.NoError(t, err)
		buff, err := os.ReadFile(passwdFilePath)
		assert.NoError(t, err)
		passwd := string(buff)
		msg := fmt.Sprintf("The %s is not of length %d", passwd, passwdLength)
		assert.Equal(t, passwdLength, len(passwd), msg)
	})

	t.Run("Should respect the count flag", func(t *testing.T) {
		testDir := t.TempDir()
		passwdFilePath := fmt.Sprintf("%s/pass.txt", testDir)
		generateCmd.SetOutput(os.Stdout)
		passwdLength := 20
		requiredCount := 20
		generateCmd.SetArgs([]string{
			"password",
			"--output=1",
			fmt.Sprintf("--length=%d", passwdLength),
			fmt.Sprintf("--count=%d", requiredCount),
			fmt.Sprintf("--file=%s", passwdFilePath),
		})
		if debug {
			generateCmd.DebugFlags()
		}
		_, err := generateCmd.ExecuteC()
		assert.NoError(t, err)
		pFile, err := os.Open(passwdFilePath)
		assert.NoError(t, err)
		defer pFile.Close()
		scanner := bufio.NewScanner(pFile)
		var (
			passwrd string
			count   int
		)
		for scanner.Scan() {
			count = count + 1
			passwrd = string(scanner.Bytes())
			msg := fmt.Sprintf("The %s is not of length %d", passwrd, passwdLength)
			assert.Equal(t, passwdLength, len(passwrd), msg)
		}
		msg := fmt.Sprintf("Should have generated %d passwords", requiredCount)
		assert.Equal(t, count, requiredCount, msg)
	})

	t.Run("Should not allow more than maxcount passwords", func(t *testing.T) {
		testDir := t.TempDir()
		passwdFilePath := fmt.Sprintf("%s/pass.txt", testDir)
		generateCmd.SetOutput(os.Stdout)
		requiredCount := password.MaxPasswordCount + 1
		generateCmd.SetArgs([]string{
			"password",
			"--output=1",
			fmt.Sprintf("--count=%d", requiredCount),
			fmt.Sprintf("--file=%s", passwdFilePath),
		})
		if debug {
			generateCmd.DebugFlags()
		}
		_, err := generateCmd.ExecuteC()
		msg := fmt.Sprintf("Should not allow more than %d passwords", password.MaxPasswordCount)
		assert.Error(t, err, msg)
	})
}
