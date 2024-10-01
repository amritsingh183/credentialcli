package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	t.Run("The password cli tool must be runnable", func(t *testing.T) {
		err := Execute()
		assert.Nil(t, err)
	})
}
