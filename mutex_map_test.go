package mutexmap_test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/mutexmap"
)

func TestMap_Getset(t *testing.T) {
	a := mutexmap.NewMap[int, string](0)
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

func TestMutexMap_SetAndGet(t *testing.T) {
	m := mutexmap.New[int, string]()
	m.Set(1, "value1")

	// 测试正常获取
	value, ok := m.Get(1)
	require.True(t, ok)
	require.Equal(t, "value1", value)

	// 测试获取不存在的键
	_, ok = m.Get(2)
	require.False(t, ok)
}

func TestMutexMap_ConcurrencyRun(t *testing.T) {
	m := mutexmap.New[int, string]()
	var wg sync.WaitGroup

	// 模拟并发写入
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Set(i, "v"+strconv.Itoa(i))
		}(i)
	}

	wg.Wait()

	// 验证写入结果
	for i := 0; i < 100; i++ {
		value, ok := m.Get(i)
		require.True(t, ok)
		require.Equal(t, "v"+strconv.Itoa(i), value)
	}
}
