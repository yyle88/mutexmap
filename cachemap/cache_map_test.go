package cachemap

import (
	"math/rand/v2"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestNewMap(t *testing.T) {
	cache := NewMap[string, int](100)

	calc := func() (int, error) {
		num := rand.IntN(100)
		if num < 50 {
			return -1, errors.New("wrong")
		}
		return num, nil
	}
	for idx := 0; idx < 10; idx++ {
		t.Log("(", idx, ")")
		const constK = "abc"
		res1, err1 := cache.Getset(constK, calc)
		t.Log(res1, err1)
		res2, err2 := cache.Getset(constK, calc)
		t.Log(res2, err2)

		if err1 != nil || err2 != nil {
			require.ErrorIs(t, err1, err2)
		}
		require.Equal(t, res1, res2)

		res3, err3, done := cache.Get(constK)
		require.True(t, done)
		if err2 != nil || err3 != nil {
			require.ErrorIs(t, err2, err3)
		}
		require.Equal(t, res2, res3)

		cache.Delete(constK)

		t.Log(cache.Len())
	}
}

func TestMap_SetVe(t *testing.T) {
	cache := NewMap[string, int](100)

	calc := func() (int, error) {
		num := rand.IntN(100)
		if num < 50 {
			return -1, errors.New("wrong")
		}
		return num, nil
	}
	for idx := 0; idx < 10; idx++ {
		t.Log("(", idx, ")")
		const constK = "abc"
		var res0 = rand.IntN(100)
		var err0 error
		if res0 < 50 {
			res0 = -1
			err0 = errors.New("wrong")
		}
		cache.SetVe(constK, res0, err0)

		res1, err1 := cache.Getset(constK, calc)
		t.Log(res1, err1)

		if err0 != nil || err1 != nil {
			t.Log(err0)
			t.Log(err1)
			require.ErrorIs(t, err0, err1)
		}
		require.Equal(t, res0, res1)

		res2, err2 := cache.Getset(constK, calc)
		t.Log(res2, err2)

		if err1 != nil || err2 != nil {
			require.ErrorIs(t, err1, err2)
		}
		require.Equal(t, res1, res2)

		res3, err3, done := cache.Get(constK)
		require.True(t, done)
		if err2 != nil || err3 != nil {
			require.ErrorIs(t, err2, err3)
		}
		require.Equal(t, res2, res3)

		cache.Delete(constK)

		t.Log(cache.Len())
	}
}
