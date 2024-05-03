package cerr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tvarney/gotvm/cerr"
)

func TestError(t *testing.T) {
	t.Parallel()
	t.Run("Error", func(t *testing.T) {
		assert.Equal(t, cerr.Error("hello").Error(), "hello")
		assert.Equal(t, cerr.Error("world").Error(), "world")
	})
}
