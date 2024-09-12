package main

import (
	"os"
	"passgen/cmd/credential/password"

	"github.com/spf13/cobra"
)

var (
	version = "0.0.1"
)

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
