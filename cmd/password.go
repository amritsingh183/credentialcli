package cmd

import (
	"fmt"
	"log"
	"os"

	"amritsingh183/password/internal/password"

	"github.com/spf13/cobra"
)

var (
	generateCmd = &cobra.Command{
		Use:     fmt.Sprintf("generate [-h] [-v] [%s] [%s] [%s] [%s] [%s]", FlagNameLength, FlagNameIncludeSpecialCharacters, FlagNameOutput, FlagNamePasswordCount, FlagNameFilePath),
		Aliases: []string{"gen"},
		Short:   "generate secure passwords",
		RunE:    runPasswordGenerator,
	}
	passwordLength      uint
	passwordCount       uint
	includeSpecialChars bool
	destinationType     uint
	destinationFilePath string
	logger              *log.Logger
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
	generateCmd.Flags().UintVar(
		&passwordLength,
		FlagNameLength,
		password.DefaultPasswordLength,
		fmt.Sprintf("How long the passwords should be? (max limit %d)", password.MaxPasswordLength),
	)
	generateCmd.Flags().UintVar(
		&passwordCount,
		FlagNamePasswordCount,
		password.DefaultPasswordCount,
		fmt.Sprintf("How many passwords to generate? (max limit %d)", password.MaxPasswordCount),
	)
	generateCmd.Flags().BoolVar(
		&includeSpecialChars,
		FlagNameIncludeSpecialCharacters,
		password.DefaultIncludeSpecialChars,
		"Whether to include special characters [for example: $ # @ ^]",
	)
	generateCmd.Flags().UintVar(
		&destinationType,
		FlagNameOutput,
		password.DefaultOutput,
		fmt.Sprintf("Device for dumping the password. %d for console, %d for file (filepath must be specified with %s)", password.ToStdOut, password.ToFile, FlagNameFilePath),
	)
	generateCmd.Flags().StringVar(
		&destinationFilePath,
		FlagNameFilePath,
		password.DefaultFilePath,
		fmt.Sprintf("filepath (when %d is provided for %s)", password.ToFile, FlagNameOutput),
	)
	logOpts := log.LstdFlags | log.Lshortfile | log.Ldate | log.Ltime | log.LUTC
	logger = log.New(os.Stderr, "password generator: ", logOpts)
}

func runPasswordGenerator(cmd *cobra.Command, args []string) error {
	log.Println("running the PasswordGenerator")
	passOptions := password.Options{
		Length:              passwordLength,
		Count:               passwordCount,
		IncludeSpecialChars: includeSpecialChars,
		DestinationType:     destinationType,
		Filepath:            destinationFilePath,
	}
	if passwordCount > 1 {
		// if there are more than 1 passwords
		// each will be separated by the Delimiter
		passOptions.Delimiter = []byte{'\n'}
	}
	err := passOptions.Validate()
	if err != nil {
		return err
	}
	logMesg := `generating password(s) with the following options %#v`
	logger.Printf(logMesg, passOptions)
	passwords := password.Generate(&passOptions)
	return password.Write(passwords, &passOptions)
}
