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
			Length: MaxPasswordCount + 1,
		}
		message := "Error should not be nil"
		err := options.Validate()
		assert.NotNil(t, err, message)
	})

	t.Run("Destination should be valid", func(t *testing.T) {
		options := &Options{
			DestinationType: 5,
		}
		message := "Error should not be nil"
		err := options.Validate()
		assert.NotNil(t, err, message)
	})

	t.Run("when destination is a file, file path must not be empty", func(t *testing.T) {
		options := &Options{
			DestinationType: ToFile,
		}
		message := "Error should not be nil"
		err := options.Validate()
		assert.NotNil(t, err, message)
	})

}
