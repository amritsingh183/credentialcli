package cmd

import (
	"amritsingh183/credentialcli/internal/password"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	// passwordCmd uses a parent cobra.Command to invoke the run function when subcommand `password` is issued
	// FIXME: rename this command to "rootCmd" or "passwordCmd"
	// "Cmd" is a little bit unhappy.
	//
	// [x] passwordCmd represents the password command
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

// [Q]: don't get why did you use this interface.
// You didn't implement this interface in your code.
// type PasswordGenerator interface {
// 	Write(io.Writer) error
// 	Generate(int) string
// }

// [x] I have now moved major code to the internal package
// and implemented the interface there as well the struct is also there now

// [Q]: don't get why did you use this struct.
// type PasswordOptions struct {
// 	length              uint
// 	count               uint
// 	includeSpecialChars bool
// 	destination         io.Writer
// }

// FIXME: this doesn't implement the interface "PasswordGenerator"
// func (pg *PasswordOptions) Generate() {
// 	// FIXME: "i = i + 1" => i++
// [x] Goes do not have postfix or prefix increment/descrement operators
// 	for i := 0; i < int(pg.count); i = i + 1 {
// 		bytePassword := util.GenerateKey(int(pg.length), pg.includeSpecialChars)
// 		// FIXME: hard to read here.
// [x] this is now improved in internal/password/password.go
// 		stringPassword := *(*string)(unsafe.Pointer(&bytePassword))
// 		pg.destination.Write([]byte(stringPassword + "\n"))
// 	}
// }

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

// FIXME: this function does too many things.
// break it down into smaller functions.
// Try to use the "internal" pkg and leave exposed only the minimum things.
// [x] The function is now simplified
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
