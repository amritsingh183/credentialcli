package password

import (
	"amritsingh183/password/internal/util"
	"errors"
	"io"
	"os"
	"unsafe"
)

// // FIXME: let's get rid of this struct. Let's group the parameters of the password generator in another struct with a different name. For the target, you can leverage the interface io.Writer that can be accepted as a parameter of the function that has to write the password somewhere.
// // FIXME: I saw to many methods around and this does not make too much sense in this kind of application (where we're not relying on external dependencies).
// [x] removing the Generator struct and defining
// a new Options struct in internal/password/options.go

// type Generator struct {
// 	destination  io.Writer
// 	passwordFile *os.File
// }

// Generate generate password(s) according to the
// given options.
// For now, since all the code is for internal use
// we do not need to make sure that
// Options.Validate() was called
func Generate(o *Options) ([][]byte, error) {
	passwds := make([][]byte, o.Count)
	var err error
	for i := 0; i < int(o.Count); i = i + 1 {
		passwds[i], err = util.GenerateKey(int(o.Length), o.IncludeSpecialChars)
		if err != nil {
			return nil, err
		}
		if len(o.Delimiter) > 0 {
			passwds[i] = append(passwds[i], o.Delimiter...)
		}
	}
	return passwds, nil
}

// Write writes the password to destination.
// It closes the file when destination is a file
func Write(data [][]byte, o *Options) error {
	var w io.Writer
	var err error
	switch o.DestinationType {
	case ToFile:
		passwdFile, err := os.OpenFile(o.Filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return errors.New("Error opening file " + o.Filepath)
		}
		defer passwdFile.Close()
		w = passwdFile
	case ToStdOut:
		w = os.Stdout
	}
	var stringPassword string
	for _, bytePassword := range data {
		usPtr := unsafe.Pointer(&bytePassword)
		strPtr := (*string)(usPtr)
		stringPassword = *strPtr
		_, err = w.Write([]byte(stringPassword))
		if err != nil {
			return err
		}
	}
	return nil
}
