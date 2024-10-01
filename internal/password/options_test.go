package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Run("Should honor MaxPasswordLength", func(t *testing.T) {
		options := &Options{
			Length: MaxPasswordLength + 1,
		}
		message := "Error should not be nil"
		err := options.Validate()
		assert.NotNil(t, err, message)
	})

	t.Run("Should honor MaxPasswordCount", func(t *testing.T) {
		options := &Options{
			Length: MaxPasswordLength - 1,
			Count:  MaxPasswordCount + 1,
		}
		message := "Error should not be nil"
		err := options.Validate()
		assert.NotNil(t, err, message)
	})

	t.Run("Should honor MinPasswordCount", func(t *testing.T) {
		options := &Options{
			Length: MaxPasswordLength - 1,
			Count:  MinPasswordCount - 1,
		}
		message := "Error should not be nil"
		err := options.Validate()
		assert.NotNil(t, err, message)
	})
	t.Run("Destination should be valid", func(t *testing.T) {
		options := &Options{
			DestinationType: 5,
			Length:          MaxPasswordLength - 1,
			Count:           MaxPasswordCount - 1,
		}
		message := "Error should not be nil"
		err := options.Validate()
		assert.NotNil(t, err, message)
	})

	t.Run("when destination is a file, file path must not be empty", func(t *testing.T) {
		options := &Options{
			DestinationType: ToFile,
			Length:          MaxPasswordLength - 1,
			Count:           MaxPasswordCount - 1,
		}
		message := "Error should not be nil"
		err := options.Validate()
		assert.NotNil(t, err, message)
	})

	t.Run("when all options are valid, error mus be nil", func(t *testing.T) {
		options := &Options{
			DestinationType: ToFile,
			Filepath:        "./out/passwords.txt",
			Length:          MaxPasswordLength - 1,
			Count:           MaxPasswordCount - 1,
		}
		message := "Error should be nil"
		err := options.Validate()
		assert.Nil(t, err, message)
	})

}
