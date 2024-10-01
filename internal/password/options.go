package password

import (
	"amritsingh183/password/internal/util"
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
	DefaultPasswordLength      = util.MinKeyLength
	DefaultPasswordCount       = 1
	DefaultIncludeSpecialChars = true
	DefaultOutput              = ToStdOut
	DefaultFilePath            = "./passwords.txt"

	MaxPasswordLength = util.MaxKeyLength
	MaxPasswordCount  = 100
	MinPasswordCount  = 1
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
	if o.Length < util.MinKeyLength {
		return fmt.Errorf("the length should be >= %d", util.MinKeyLength)
	}
	if o.Count > MaxPasswordCount {
		return fmt.Errorf("the max count should not exceed %d", MaxPasswordCount)
	}
	if o.Count < MinPasswordCount {
		return fmt.Errorf("the number of passwords should be >= %d", MinPasswordCount)
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
