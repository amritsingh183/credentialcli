package password

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

var (
	// Cmd uses a parent cobra.Command to invoke the run function when subcommand `password` is issued
	Cmd = &cobra.Command{
		Use:     fmt.Sprintf("password [-h] [-v] [%s] [%s] [%s]", FlagNameLength, FlagNameIncludeSpecialCharacters, FlagNameOutput),
		Aliases: []string{"pass"},
		Short:   "generate secure passwords",
		RunE:    runPasswordGenerator,
	}
	passwordLength      uint
	includeSpecialChars bool
	destination         uint
	destinationFilePath string
)

const (
	ToStdOut = iota
	ToFile
)

const (
	DefaultPasswordLength      = 7
	DefaultIncludeSpecialChars = true
	DefaultOutput              = ToStdOut
	DefaultFilePath            = "./passwords.txt"

	MaxPasswordLength = 100

	FlagNameLength                   = "length"
	FlagNameIncludeSpecialCharacters = "includeSpecialCharacters"
	FlagNameOutput                   = "output"
	FlagNameFilePath                 = "file"
)

type PasswordGenerator interface {
	Write(io.Writer) error
	Generate(int) string
}

type PasswordOptions struct {
	length              uint
	includeSpecialChars bool
	destination         io.Writer
}

func (pg *PasswordOptions) Generate() {
	bytePassword := util.GenerateKey(int(pg.length), pg.includeSpecialChars)
	stringPassword := *(*string)(unsafe.Pointer(&bytePassword))
	pg.destination.Write([]byte(stringPassword))
	pg.destination.Write([]byte("\n"))
}

func init() {
	// Local flags that are only available to this command.
	Cmd.Flags().UintVar(
		&passwordLength,
		FlagNameLength,
		DefaultPasswordLength,
		fmt.Sprintf("How long the passwords should be? (max limit %d)", MaxPasswordLength),
	)
	Cmd.Flags().BoolVar(
		&includeSpecialChars,
		FlagNameIncludeSpecialCharacters,
		DefaultIncludeSpecialChars,
		"Whether to include special characters [for example: $ # @ ^]",
	)
	Cmd.Flags().UintVar(
		&destination,
		FlagNameOutput,
		DefaultOutput,
		fmt.Sprintf("Device for dumping the password. %d for console, %d for file (filepath must be specified with %s)", ToStdOut, ToFile, FlagNameFilePath),
	)
	Cmd.Flags().StringVar(
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
	logger.Println("shouldIncludeSpecialChars", includeSpecialChars)

	myPg := PasswordOptions{
		length:              passwordLength,
		includeSpecialChars: includeSpecialChars,
	}
	if passwordLength > MaxPasswordLength {
		return fmt.Errorf("the max length should not exceed %d", MaxPasswordLength)
	}
	var humanReadableDestName string

	switch destination {
	case ToStdOut:
		myPg.destination = os.Stdout
		humanReadableDestName = "console"
	case ToFile:
		passwordFile, err := os.OpenFile(destinationFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return errors.New("Error opening file " + destinationFilePath)
		}
		humanReadableDestName = fmt.Sprintf("File %s", destinationFilePath)
		myPg.destination = passwordFile
		defer passwordFile.Close()
	}
	logger.Println("destination", humanReadableDestName)
	myPg.Generate()
	return nil
}
