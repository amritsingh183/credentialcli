package password

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

const (
	DefaultPasswordLength            = 7
	DefaultIncludeSpecialChars       = true
	FlagNameLength                   = "length"
	FlagNameIncludeSpecialCharacters = "includeSpecialCharacters"
)

var (
	passwordLength      string
	includeSpecialChars bool

	// Cmd uses a parent cobra.Command to invoke the run function when subcommand `password` is issued
	Cmd = &cobra.Command{
		Use:     fmt.Sprintf("password [-h] [%s] [%s]", FlagNameLength, FlagNameIncludeSpecialCharacters),
		Aliases: []string{"pass"},
		Short:   "generate secure passwords",
		RunE:    runPasswordGenerator,
	}
)

func init() {
	// Local flags that are only available to this command.
	Cmd.Flags().StringVar(
		&passwordLength,
		FlagNameLength,
		fmt.Sprintf("%d", DefaultPasswordLength),
		"How long the passwords should be?",
	)
	Cmd.Flags().BoolVar(
		&includeSpecialChars,
		FlagNameIncludeSpecialCharacters,
		DefaultIncludeSpecialChars,
		"Whether to include special characters [for example: $ # @ ^]",
	)
}

func runPasswordGenerator(cmd *cobra.Command, args []string) error {
	log.Println("running the PasswordGenerator")
	logger := log.New(cmd.OutOrStdout(), "creds: ", log.Ldate|log.Ltime|log.LUTC)
	lengthVal, shouldIncludeSpecialCharsVal := args[1], args[3]
	logger.Println("givenPasswordLength", lengthVal)
	logger.Println("shouldIncludeSpecialChars", shouldIncludeSpecialCharsVal)
	return nil
}
