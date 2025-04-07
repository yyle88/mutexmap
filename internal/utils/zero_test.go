package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/mutexmap/internal/utils"
)

func TestZero(t *testing.T) {
	res := utils.Zero[int]()
	t.Log(res)
	require.Zero(t, res)

	require.Zero(t, utils.Zero[string]())
	require.Zero(t, utils.Zero[float64]())
	require.Zero(t, utils.Zero[uint32]())
	require.Zero(t, utils.Zero[bool]())
	require.Zero(t, utils.Zero[int64]())
}
