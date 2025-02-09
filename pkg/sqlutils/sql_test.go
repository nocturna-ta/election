package sqlutils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewNullString(t *testing.T) {
	t.Run("ShouldReturnValid", func(t *testing.T) {
		s := "test"
		res := NewNullString(&s)
		require.True(t, res.Valid)
		require.Equal(t, s, res.String)
	})

	t.Run("ShouldReturnNull", func(t *testing.T) {
		res := NewNullString(nil)
		require.False(t, res.Valid)
	})
}
