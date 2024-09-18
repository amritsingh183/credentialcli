package cmd

import (
	"fmt"
	"log"

	"amritsingh183/credentialcli/internal/password"

	"github.com/spf13/cobra"
)

var (
	passwordCmd = &cobra.Command{
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
	FlagNameLength                   = "length"
	FlagNameIncludeSpecialCharacters = "includeSpecialCharacters"
	FlagNameOutput                   = "output"
	FlagNameFilePath                 = "file"
	FlagNamePasswordCount            = "count"
)

func init() {
	// Local flags that are only available to this command.
	passwordCmd.Flags().UintVar(
		&passwordLength,
		FlagNameLength,
		password.DefaultPasswordLength,
		fmt.Sprintf("How long the passwords should be? (max limit %d)", password.MaxPasswordLength),
	)
	passwordCmd.Flags().UintVar(
		&passwordCount,
		FlagNamePasswordCount,
		password.DefaultPasswordCount,
		fmt.Sprintf("How many passwords to generate? (max limit %d)", password.MaxPasswordCount),
	)
	passwordCmd.Flags().BoolVar(
		&includeSpecialChars,
		FlagNameIncludeSpecialCharacters,
		password.DefaultIncludeSpecialChars,
		"Whether to include special characters [for example: $ # @ ^]",
	)
	passwordCmd.Flags().UintVar(
		&destination,
		FlagNameOutput,
		password.DefaultOutput,
		fmt.Sprintf("Device for dumping the password. %d for console, %d for file (filepath must be specified with %s)", password.ToStdOut, password.ToFile, FlagNameFilePath),
	)
	passwordCmd.Flags().StringVar(
		&destinationFilePath,
		FlagNameFilePath,
		password.DefaultFilePath,
		fmt.Sprintf("filepath (when %d is provided for %s)", password.ToFile, FlagNameOutput),
	)
}

type Generator interface {
	Generate() [][]byte
	Write([][]byte) error
	Validate() error
}

func runPasswordGenerator(cmd *cobra.Command, args []string) error {
	log.Println("running the PasswordGenerator")
	passGen := password.Generator{
		Length:              passwordLength,
		Count:               passwordCount,
		IncludeSpecialChars: includeSpecialChars,
		DestinationType:     destination,
		Filepath:            destinationFilePath,
	}
	return GeneratePassword(&passGen)
}

func GeneratePassword(passGen Generator) error {
	err := passGen.Validate()
	if err != nil {
		return err
	}
	passwords := passGen.Generate()
	return passGen.Write(passwords)
}
