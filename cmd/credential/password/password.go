package password

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"unsafe"

	"amritsingh183/credentialcli/util"

	"github.com/spf13/cobra"
)

var (
	// Cmd uses a parent cobra.Command to invoke the run function when subcommand `password` is issued
	// FIXME: rename this command to "rootCmd" or "passwordCmd"
	// "Cmd" is a little bit unhappy.
	Cmd = &cobra.Command{
		Use:     fmt.Sprintf("password [-h] [-v] [%s] [%s] [%s] [%s] [%s]", FlagNameLength, FlagNameIncludeSpecialCharacters, FlagNameOutput, FlagNamePasswordCount, FlagNameFilePath),
		Aliases: []string{"pass"},
		Short:   "generate secure passwords",
		RunE:    runPasswordGenerator,
	}
	passwordLength      uint
	passwordCount       uint
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
	DefaultPasswordCount       = 1
	DefaultIncludeSpecialChars = true
	DefaultOutput              = ToStdOut
	DefaultFilePath            = "./passwords.txt"

	MaxPasswordLength = 100
	MaxPasswordCount  = 100

	FlagNameLength                   = "length"
	FlagNameIncludeSpecialCharacters = "includeSpecialCharacters"
	FlagNameOutput                   = "output"
	FlagNameFilePath                 = "file"
	FlagNamePasswordCount            = "count"
)

// [Q]: don't get why did you use this interface.
// You didn't implement this interface in your code.
type PasswordGenerator interface {
	Write(io.Writer) error
	Generate(int) string
}

// [Q]: don't get why did you use this struct.
type PasswordOptions struct {
	length              uint
	count               uint
	includeSpecialChars bool
	destination         io.Writer
}

// FIXME: this doesn't implement the interface "PasswordGenerator"
func (pg *PasswordOptions) Generate() {
	// FIXME: "i = i + 1" => i++
	for i := 0; i < int(pg.count); i = i + 1 {
		bytePassword := util.GenerateKey(int(pg.length), pg.includeSpecialChars)
		// FIXME: hard to read here.
		stringPassword := *(*string)(unsafe.Pointer(&bytePassword))
		pg.destination.Write([]byte(stringPassword + "\n"))
	}
}

func init() {
	// Local flags that are only available to this command.
	Cmd.Flags().UintVar(
		&passwordLength,
		FlagNameLength,
		DefaultPasswordLength,
		fmt.Sprintf("How long the passwords should be? (max limit %d)", MaxPasswordLength),
	)
	Cmd.Flags().UintVar(
		&passwordCount,
		FlagNamePasswordCount,
		DefaultPasswordCount,
		fmt.Sprintf("How many passwords to generate? (max limit %d)", MaxPasswordCount),
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

// FIXME: this function does too many things.
// break it down into smaller functions.
// Try to use the "internal" pkg and leave exposed only the minimum things.
func runPasswordGenerator(cmd *cobra.Command, args []string) error {
	log.Println("running the PasswordGenerator")
	logger := log.New(cmd.OutOrStdout(), "creds: ", log.Ldate|log.Ltime|log.LUTC)
	logger.Println("givenPasswordLength", passwordLength)
	logger.Println("passwordCount", passwordCount)
	logger.Println("shouldIncludeSpecialChars", includeSpecialChars)
	myPg := PasswordOptions{
		length:              passwordLength,
		includeSpecialChars: includeSpecialChars,
		count:               passwordCount,
	}
	if passwordLength > MaxPasswordLength {
		return fmt.Errorf("the max length should not exceed %d", MaxPasswordLength)
	}
	if passwordCount > MaxPasswordCount {
		return fmt.Errorf("the max count should not exceed %d", MaxPasswordCount)
	}
	var humanReadableDestName string

	switch destination {
	case ToStdOut:
		myPg.destination = os.Stdout
		humanReadableDestName = "console"
	case ToFile:
		passwordFile, err := os.OpenFile(destinationFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
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
