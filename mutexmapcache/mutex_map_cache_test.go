package mutexmapcache

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
		res1, erx1 := cache.Getset(constK, calc)
		t.Log(res1, erx1)
		res2, erx2 := cache.Getset(constK, calc)
		t.Log(res2, erx2)

		if erx1 != nil || erx2 != nil {
			require.ErrorIs(t, erx1, erx2)
		}
		require.Equal(t, res1, res2)

		res3, erx3, done := cache.Get(constK)
		require.True(t, done)
		if erx2 != nil || erx3 != nil {
			require.ErrorIs(t, erx2, erx3)
		}
		require.Equal(t, res2, res3)

		cache.Delete(constK)

		t.Log(cache.Len())
	}
}

func TestMap_SetVE(t *testing.T) {
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
		var erx0 error
		if res0 < 50 {
			res0 = -1
			erx0 = errors.New("wrong")
		}
		cache.SetVE(constK, res0, erx0)

		res1, erx1 := cache.Getset(constK, calc)
		t.Log(res1, erx1)

		if erx0 != nil || erx1 != nil {
			t.Log(erx0)
			t.Log(erx1)
			require.ErrorIs(t, erx0, erx1)
		}
		require.Equal(t, res0, res1)

		res2, erx2 := cache.Getset(constK, calc)
		t.Log(res2, erx2)

		if erx1 != nil || erx2 != nil {
			require.ErrorIs(t, erx1, erx2)
		}
		require.Equal(t, res1, res2)

		res3, erx3, done := cache.Get(constK)
		require.True(t, done)
		if erx2 != nil || erx3 != nil {
			require.ErrorIs(t, erx2, erx3)
		}
		require.Equal(t, res2, res3)

		cache.Delete(constK)

		t.Log(cache.Len())
	}
}
