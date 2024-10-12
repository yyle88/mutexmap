package mutexmap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap_Getset(t *testing.T) {
	a := NewMap[int, string](0)
	{
		v, created := a.Getset(0, func() string {
			return "abc"
		})
		require.True(t, created)
		require.Equal(t, v, "abc")
	}
	{
		v, created := a.Getset(0, func() string {
			return "xyz"
		})
		require.False(t, created)
		require.Equal(t, v, "abc") // not change value when already exist.
	}
}
