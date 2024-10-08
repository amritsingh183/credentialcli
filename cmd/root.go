/*
Copyright © 2024 Amrit Singh <amritsingh183@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "0.0.1"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "password",
	Short:   "password is a utility to generate passwords",
	Long:    "",
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	rootCmd.SetOutput(os.Stdout)

	rootCmd.AddCommand(generateCmd)
	// Execute the Cobra command tree, parsing args and identifying the command
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
