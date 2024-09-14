package main

import (
	"os"

	"amritsingh183/credentialcli/cmd/credential/password"

	"github.com/spf13/cobra"
)

var version = "0.0.1"

// FIXME: the project structure might be good.
/*
	I would have used the cobra-cli command to initialize it.
	Start by creating a go.mod with the command 'go mod init'.
	Then, issue the 'cobra-cli init'.
	I cannot see the root.go file.
	Then, keep adding commands & subcommands with the 'cobra-cli add <name of command>
*/

func main() {
	// Build the root command that runs the credential generator
	rootCmd := &cobra.Command{
		Use:     "credentials",
		Short:   "credentials is a utility to generate credentials",
		Version: version,
	}
	rootCmd.SetOutput(os.Stdout)

	rootCmd.AddCommand(password.Cmd)

	// Execute the Cobra command tree, parsing args and identifying the command
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
