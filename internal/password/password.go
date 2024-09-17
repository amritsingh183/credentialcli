package password

import (
	"amritsingh183/credentialcli/internal/util"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"unsafe"
)

var logger *log.Logger

// Output Devices
const (
	ToStdOut = iota
	ToFile
)

// Password Generation Rules
const (
	DefaultPasswordLength      = 7
	DefaultPasswordCount       = 1
	DefaultIncludeSpecialChars = true
	DefaultMustBeUrlSafe       = false
	DefaultOutput              = ToStdOut
	DefaultFilePath            = "./passwords.txt"

	MaxPasswordLength = 100
	MaxPasswordCount  = 100
)

func init() {
	logOpts := log.LstdFlags | log.Lshortfile | log.Ldate | log.Ltime | log.LUTC
	logger = log.New(os.Stderr, "password generator: ", logOpts)
}

type Generator struct {
	Length              uint
	Count               uint
	IncludeSpecialChars bool
	DestinationType     uint
	destination         io.Writer

	Filepath     string
	passwordFile *os.File
}

func (g *Generator) Generate() [][]byte {
	logMesg := `generating password(s) with the following options
count=%v
length=%v
destination=%v
filePath=%v
includeSpecialChars=%v`
	logger.Printf(logMesg, g.Count, g.Length, g.DestinationType, g.Filepath, g.IncludeSpecialChars)
	passwds := make([][]byte, g.Count)
	for i := 0; i < int(g.Count); i = i + 1 {
		passwds[i] = util.GenerateKey(int(g.Length), g.IncludeSpecialChars)
	}
	return passwds
}

func (g *Generator) Write(data [][]byte) error {
	defer g.passwordFile.Close()
	var err error
	var stringPassword string
	addNewLine := false
	if g.Count > 1 {
		addNewLine = true
	}
	for _, bytePassword := range data {
		usPtr := unsafe.Pointer(&bytePassword)
		strPtr := (*string)(usPtr)
		stringPassword = *strPtr
		if addNewLine {
			_, err = g.destination.Write([]byte(stringPassword + "\n"))
		} else {
			_, err = g.destination.Write([]byte(stringPassword))
		}
		if err != nil {
			return err
		}
	}
	return nil
}
func (g *Generator) Validate() error {
	if g.Length > MaxPasswordLength {
		return fmt.Errorf("the max length should not exceed %d", MaxPasswordLength)
	}
	if g.Count > MaxPasswordCount {
		return fmt.Errorf("the max count should not exceed %d", MaxPasswordCount)
	}
	switch g.DestinationType {
	case ToStdOut:
		g.destination = os.Stdout
	case ToFile:
		passwordFile, err := util.CreateFile(g.Filepath)
		if err != nil {
			return errors.New("Error opening file " + g.Filepath)
		}
		g.passwordFile = passwordFile
		g.destination = passwordFile
	}
	return nil
}
