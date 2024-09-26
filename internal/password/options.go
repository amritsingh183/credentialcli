package password

import (
	"errors"
	"fmt"
)

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
	DefaultOutput              = ToStdOut
	DefaultFilePath            = "./passwords.txt"

	MaxPasswordLength = 100
	MaxPasswordCount  = 100
)

type Options struct {
	Length              uint
	Count               uint
	IncludeSpecialChars bool
	DestinationType     uint
	Delimiter           []byte

	Filepath string
}

// Validate vlaidates the options available for
// password generator
func (o *Options) Validate() error {
	if o.Length > MaxPasswordLength {
		return fmt.Errorf("the max length should not exceed %d", MaxPasswordLength)
	}
	if o.Count > MaxPasswordCount {
		return fmt.Errorf("the max count should not exceed %d", MaxPasswordCount)
	}
	switch o.DestinationType {
	case ToFile:
		if len(o.Filepath) == 0 {
			return errors.New("filepath should not be empty")
		}
	case ToStdOut:
	default:
		return fmt.Errorf("invalid destination %d", o.DestinationType)
	}
	return nil
}
