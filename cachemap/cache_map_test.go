package cachemap_test

import (
	"math/rand/v2"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/mutexmap/cachemap"
)

func TestNew(t *testing.T) {
	cache := cachemap.New[string, int]()
	require.NotNil(t, cache)
	require.Equal(t, 0, cache.Len())

	// Verify default capacity by adding elements
	const constK = "abc"
	cache.Set(constK, 42)
	res, err, done := cache.Get(constK)
	t.Log(res, err)
	require.True(t, done)
	require.NoError(t, err)
	require.Equal(t, 42, res)
}

func TestNewMap(t *testing.T) {
	cache := cachemap.NewMap[string, int](100)

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

func TestGet(t *testing.T) {
	cache := cachemap.NewMap[string, int](100)

	const constK = "abc"
	res, err, done := cache.Get(constK)
	t.Log(res, err)
	require.False(t, done)
	require.Error(t, err)
	require.Equal(t, "not exist", err.Error())
	require.Equal(t, 0, res)

	cache.Set(constK, 42)
	res, err, done = cache.Get(constK)
	t.Log(res, err)
	require.True(t, done)
	require.NoError(t, err)
	require.Equal(t, 42, res)
}

func TestSet(t *testing.T) {
	cache := cachemap.NewMap[string, int](100)

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
		cache.Set(constK, res0)

		res1, err1, done := cache.Get(constK)
		t.Log(res1, err1)
		require.True(t, done)
		require.NoError(t, err1)
		require.Equal(t, res0, res1)

		res2, err2 := cache.Getset(constK, calc)
		t.Log(res2, err2)
		require.NoError(t, err2)
		require.Equal(t, res0, res2)

		cache.Delete(constK)
		t.Log(cache.Len())
	}
}

func TestSetKv(t *testing.T) {
	cache := cachemap.NewMap[string, int](100)

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
		cache.SetKv(constK, res0)

		res1, err1, done := cache.Get(constK)
		t.Log(res1, err1)
		require.True(t, done)
		require.NoError(t, err1)
		require.Equal(t, res0, res1)

		res2, err2 := cache.Getset(constK, calc)
		t.Log(res2, err2)
		require.NoError(t, err2)
		require.Equal(t, res0, res2)

		cache.Delete(constK)
		t.Log(cache.Len())
	}
}

func TestSetKe(t *testing.T) {
	cache := cachemap.NewMap[string, int](100)

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
		err0 := errors.New("wrong")
		cache.SetKe(constK, err0)

		res1, err1, done := cache.Get(constK)
		t.Log(res1, err1)
		require.True(t, done)
		require.ErrorIs(t, err0, err1)
		require.Equal(t, 0, res1) // Zero value for int

		res2, err2 := cache.Getset(constK, calc)
		t.Log(res2, err2)
		require.ErrorIs(t, err0, err2)
		require.Equal(t, 0, res2)

		cache.Delete(constK)
		t.Log(cache.Len())
	}
}

func TestSetVe(t *testing.T) {
	cache := cachemap.NewMap[string, int](100)

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

func TestDelete(t *testing.T) {
	cache := cachemap.NewMap[string, int](100)

	for idx := 0; idx < 10; idx++ {
		t.Log("(", idx, ")")
		const constK = "abc"
		var res0 = rand.IntN(100)
		cache.Set(constK, res0)

		res1, err1, done := cache.Get(constK)
		t.Log(res1, err1)
		require.True(t, done)
		require.NoError(t, err1)
		require.Equal(t, res0, res1)

		cache.Delete(constK)

		res2, err2, done2 := cache.Get(constK)
		t.Log(res2, err2)
		require.False(t, done2)
		require.Error(t, err2)
		require.Equal(t, "not exist", err2.Error())
		require.Equal(t, 0, res2)

		t.Log(cache.Len())
	}
}

func TestLen(t *testing.T) {
	cache := cachemap.NewMap[string, int](100)

	for idx := 0; idx < 10; idx++ {
		t.Log("(", idx, ")")
		key := "key" + string(rune('A'+idx))
		var res0 = rand.IntN(100)
		cache.Set(key, res0)

		t.Log("Len:", cache.Len())
		require.Equal(t, idx+1, cache.Len())

		res1, err1, done := cache.Get(key)
		t.Log(res1, err1)
		require.True(t, done)
		require.NoError(t, err1)
		require.Equal(t, res0, res1)
	}

	for idx := 0; idx < 10; idx++ {
		key := "key" + string(rune('A'+idx))
		cache.Delete(key)
		t.Log("Len after delete:", cache.Len())
		require.Equal(t, 9-idx, cache.Len())
	}
}

func TestRange(t *testing.T) {
	cache := cachemap.NewMap[string, int](100)

	// Populate map with random values
	expected := make(map[string]int)
	for idx := 0; idx < 10; idx++ {
		key := "key" + string(rune('A'+idx))
		value := rand.IntN(100)
		cache.Set(key, value)
		expected[key] = value
	}

	// Test Range
	collected := make(map[string]int)
	cache.Range(func(k string, v int, err error, done bool) bool {
		t.Log("Range:", k, v, err, done)
		require.True(t, done)
		require.NoError(t, err)
		collected[k] = v
		return true
	})

	// Verify collected values
	for k, v := range expected {
		v2, ok := collected[k]
		require.True(t, ok)
		require.Equal(t, v, v2)
	}
	require.Equal(t, len(expected), len(collected))

	// Test Range with early termination
	count := 0
	cache.Range(func(k string, v int, err error, done bool) bool {
		count++
		return count < 5 // Stop after 5 iterations
	})
	require.Equal(t, 5, count)
}
