package mutexmap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap_GetOrzSet(t *testing.T) {
	a := NewMap[int, string](0)
	{
		v, created := a.GetOrzSet(0, func() string {
			return "abc"
		})
		require.True(t, created)
		require.Equal(t, v, "abc")
	}
	{
		v, created := a.GetOrzSet(0, func() string {
			return "xyz"
		})
		require.False(t, created)
		require.Equal(t, v, "abc") // not change value when already exist.
	}
}
