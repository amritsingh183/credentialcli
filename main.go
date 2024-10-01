package main

import (
	"amritsingh183/password/cmd"
	"os"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
