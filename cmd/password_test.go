package cmd

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

// FIXME: there is no MaxLength() function in the source code.
// Keep 1:1 relationship between the test code function and the source code ones.
// Use something like "testify/suite" to tidy things up.
func TestRunPasswordGenerator(t *testing.T) {
	t.Run("Should error if password exceeds max length", func(t *testing.T) {
		passwordCmd.SetOutput(os.Stdout)
		passwordCmd.SetArgs([]string{
			"--length=102",
		})
		passwordCmd.DebugFlags()
		_, err := passwordCmd.ExecuteC()
		fmt.Println(err)
		if err == nil {
			t.Errorf("Should not allow length greater than %d", MaxPasswordLength)
		}
	})

	t.Run("Should respect the length flag", func(t *testing.T) {
		testDir := t.TempDir()
		passwdFilePath := fmt.Sprintf("%s/pass.txt", testDir)
		passwordCmd.SetOutput(os.Stdout)
		passwdLength := 20
		passwordCmd.SetArgs([]string{
			"password",
			fmt.Sprintf("--length=%d", passwdLength),
			"--output=1",
			fmt.Sprintf("--file=%s", passwdFilePath),
		})
		passwordCmd.DebugFlags()
		_, err := passwordCmd.ExecuteC()
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		}
		buff, err := os.ReadFile(passwdFilePath)
		if err != nil {
			t.Error(err)
		}
		passwd := string(buff)
		if len(passwd) != passwdLength {
			t.Errorf("The %s is not of length %d", passwd, passwdLength)
		}
	})

	t.Run("Should respect the url safe flag", func(t *testing.T) {
		testDir := t.TempDir()
		passwdFilePath := fmt.Sprintf("%s/pass.txt", testDir)
		passwordCmd.SetOutput(os.Stdout)
		passwdLength := 20
		passwordCmd.SetArgs([]string{
			"password",
			fmt.Sprintf("--length=%d", passwdLength),
			fmt.Sprintf("--length=%d", passwdLength),
			"--output=1",
			fmt.Sprintf("--file=%s", passwdFilePath),
		})
		passwordCmd.DebugFlags()
		_, err := passwordCmd.ExecuteC()
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		}
		buff, err := os.ReadFile(passwdFilePath)
		if err != nil {
			t.Error(err)
		}
		passwd := string(buff)
		if len(passwd) != passwdLength {
			// FIXME: this test fails
			t.Errorf("The %s is not of length %d", passwd, passwdLength)
		}
	})

	t.Run("Should respect the count flag", func(t *testing.T) {
		testDir := t.TempDir()
		passwdFilePath := fmt.Sprintf("%s/pass.txt", testDir)
		passwordCmd.SetOutput(os.Stdout)
		passwdLength := 20
		requiredCount := 20
		passwordCmd.SetArgs([]string{
			"password",
			"--output=1",
			fmt.Sprintf("--length=%d", passwdLength),
			fmt.Sprintf("--count=%d", requiredCount),
			fmt.Sprintf("--file=%s", passwdFilePath),
		})
		passwordCmd.DebugFlags()
		_, err := passwordCmd.ExecuteC()
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		}
		pFile, err := os.Open(passwdFilePath)
		if err != nil {
			t.Error(err)
		}
		defer pFile.Close()
		scanner := bufio.NewScanner(pFile)
		var (
			passwrd string
			count   int
		)
		for scanner.Scan() {
			count = count + 1
			passwrd = string(scanner.Bytes())
			if len(passwrd) != passwdLength {
				t.Errorf("The %s is not of length %d", passwrd, passwdLength)
			}
		}
		fmt.Println("Password count", count)
		if count != requiredCount {
			t.Errorf("Should have generated %d passwords", requiredCount)
		}
	})

	t.Run("Should not allow more than maxcount passwords", func(t *testing.T) {
		testDir := t.TempDir()
		passwdFilePath := fmt.Sprintf("%s/pass.txt", testDir)
		passwordCmd.SetOutput(os.Stdout)
		requiredCount := MaxPasswordCount + 1
		passwordCmd.SetArgs([]string{
			"password",
			"--output=1",
			fmt.Sprintf("--count=%d", requiredCount),
			fmt.Sprintf("--file=%s", passwdFilePath),
		})
		passwordCmd.DebugFlags()
		_, err := passwordCmd.ExecuteC()
		if err == nil {
			t.Errorf("Should have generated %d passwords", requiredCount)
		}
	})
}
