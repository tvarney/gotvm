package vmerr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tvarney/gotvm/vmerr"
)

func TestError(t *testing.T) {
	t.Parallel()
	t.Run("Error", func(t *testing.T) {
		assert.Equal(t, vmerr.ConstError("hello").Error(), "hello")
		assert.Equal(t, vmerr.ConstError("world").Error(), "world")
	})
}
