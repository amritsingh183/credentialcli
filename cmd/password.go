/*
Copyright © 2024 Amrit Singh <amritsingh183@gmail.com>
*/
package cmd

import (
	"amritsingh183/credentialcli/util"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"unsafe"

	"github.com/spf13/cobra"
)

const (
	ToStdOut = iota
	ToFile
)

const (
	DefaultPasswordLength      = 7
	DefaultPasswordCount       = 1
	DefaultIncludeSpecialChars = true
	DefaultMustBeUrlSafe       = false
	DefaultOutput              = ToStdOut
	DefaultFilePath            = "./passwords.txt"

	MaxPasswordLength = 100
	MaxPasswordCount  = 100

	FlagNameLength                   = "length"
	FlagNameIncludeSpecialCharacters = "includeSpecialCharacters"
	FlagNameMustBeUrlSafe            = "urlSafe"
	FlagNameOutput                   = "output"
	FlagNameFilePath                 = "file"
	FlagNamePasswordCount            = "count"
)

// passwordCmd represents the password command
var (
	passwordCmd = &cobra.Command{
		Use:     fmt.Sprintf("password [-h] [-v] [%s] [%s] [%s] [%s] [%s]", FlagNameLength, FlagNameIncludeSpecialCharacters, FlagNameOutput, FlagNamePasswordCount, FlagNameFilePath),
		Aliases: []string{"pass"},
		Short:   "generate secure passwords",
		RunE:    runPasswordGenerator,
	}
	passwordLength      uint
	passwordCount       uint
	mustBeUrlSafe       bool
	includeSpecialChars bool
	destination         uint
	outputDevice        io.Writer
	destinationFilePath string
)

func init() {

	// Local flags that are only available to this command.
	passwordCmd.Flags().UintVar(
		&passwordLength,
		FlagNameLength,
		DefaultPasswordLength,
		fmt.Sprintf("How long the passwords should be? (max limit %d)", MaxPasswordLength),
	)
	passwordCmd.Flags().UintVar(
		&passwordCount,
		FlagNamePasswordCount,
		DefaultPasswordCount,
		fmt.Sprintf("How many passwords to generate? (max limit %d)", MaxPasswordCount),
	)
	passwordCmd.Flags().BoolVar(
		&includeSpecialChars,
		FlagNameIncludeSpecialCharacters,
		DefaultIncludeSpecialChars,
		"Whether to include special characters [for example: $ # @ ^]",
	)
	passwordCmd.Flags().BoolVar(
		&mustBeUrlSafe,
		FlagNameMustBeUrlSafe,
		DefaultMustBeUrlSafe,
		"Whether to generate URL safe passwords",
	)
	passwordCmd.Flags().UintVar(
		&destination,
		FlagNameOutput,
		DefaultOutput,
		fmt.Sprintf("Device for dumping the password. %d for console, %d for file (filepath must be specified with %s)", ToStdOut, ToFile, FlagNameFilePath),
	)
	passwordCmd.Flags().StringVar(
		&destinationFilePath,
		FlagNameFilePath,
		DefaultFilePath,
		fmt.Sprintf("filepath (when %d is provided for %s)", ToFile, FlagNameOutput),
	)
}

func runPasswordGenerator(cmd *cobra.Command, args []string) error {
	log.Println("running the PasswordGenerator")
	logger := log.New(cmd.OutOrStdout(), "creds: ", log.Ldate|log.Ltime|log.LUTC)
	logger.Println("givenPasswordLength", passwordLength)
	logger.Println("passwordCount", passwordCount)
	logger.Println("shouldIncludeSpecialChars", includeSpecialChars)
	err := validateOptions()
	if err != nil {
		return err
	}
	var humanReadableDestName string

	switch destination {
	case ToStdOut:
		outputDevice = os.Stdout
		humanReadableDestName = "console"
	case ToFile:
		passwordFile, err := createFile()
		if err != nil {
			return errors.New("Error opening file " + destinationFilePath)
		}
		humanReadableDestName = fmt.Sprintf("File %s", destinationFilePath)
		outputDevice = passwordFile
		defer passwordFile.Close()
	}
	logger.Println("destination", humanReadableDestName)
	return write(generate())
}
func createFile() (*os.File, error) {
	return os.OpenFile(destinationFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
}
func validateOptions() error {
	if passwordLength > MaxPasswordLength {
		return fmt.Errorf("the max length should not exceed %d", MaxPasswordLength)
	}
	if passwordCount > MaxPasswordCount {
		return fmt.Errorf("the max count should not exceed %d", MaxPasswordCount)
	}
	return nil
}

func generate() [][]byte {
	passwds := make([][]byte, passwordCount)
	for i := 0; i < int(passwordCount); i = i + 1 {
		passwds[i] = util.GenerateKey(int(passwordLength), includeSpecialChars)
	}
	return passwds
}

func write(data [][]byte) error {
	var err error
	var stringPassword string
	addNewLine := false
	if passwordCount > 1 {
		addNewLine = true
	}
	for _, bytePassword := range data {
		if mustBeUrlSafe {
			stringPassword = util.Base64URLEncode(bytePassword)
		} else {
			usPtr := unsafe.Pointer(&bytePassword)
			strPtr := (*string)(usPtr)
			stringPassword = *strPtr
		}
		if addNewLine {
			_, err = outputDevice.Write([]byte(stringPassword + "\n"))
		} else {
			_, err = outputDevice.Write([]byte(stringPassword))
		}
		if err != nil {
			return err
		}
	}
	return nil
}
