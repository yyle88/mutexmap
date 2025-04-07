package mutexmap

import (
	"sync"

	"github.com/yyle88/mutexmap/internal/utils"
)

// Map provides a thread-safe map implementation using a sync.RWMutex.
// Map 提供了一个使用 sync.RWMutex 的线程安全 map 实现。
type Map[K comparable, V any] struct {
	mp    map[K]V       // The underlying map. 内部 map。
	mutex *sync.RWMutex // Mutex for synchronizing access. 用于同步访问的互斥锁。
}

// NewMap creates a new thread-safe map.
// NewMap 创建线程安全 map。
func NewMap[K comparable, V any](cap int) *Map[K, V] {
	return &Map[K, V]{
		mp:    make(map[K]V, cap), // Initialize the internal map. 初始化内部 map。
		mutex: &sync.RWMutex{},    // Initialize the RWMutex. 初始化读写锁。
	}
}

// Get retrieves the value associated with the given key. It returns the value and a boolean indicating whether the key exists.
// Get 获取与指定键关联的值，它返回值以及一个布尔值，指示键是否存在。
func (a *Map[K, V]) Get(k K) (V, bool) {
	a.mutex.RLock()         // Acquire read lock. 获取读锁。
	defer a.mutex.RUnlock() // Ensure lock is released. 确保锁被释放。
	if v, ok := a.mp[k]; ok {
		return v, ok
	} else {
		return utils.Zero[V](), false // Explicitly return zero value if not found. 未找到时显式返回零值。
	}
}

// Set inserts or updates the value for the given key into map.
// Set 插入或更新指定键的值。
func (a *Map[K, V]) Set(k K, v V) {
	a.mutex.Lock()         // Acquire write lock. 获取写锁。
	defer a.mutex.Unlock() // Ensure lock is released. 确保锁被释放。
	a.mp[k] = v
}

// Delete removes the key-value.
// Delete 移除与指定键关联的键值对。
func (a *Map[K, V]) Delete(k K) {
	a.mutex.Lock()         // Acquire write lock. 获取写锁。
	defer a.mutex.Unlock() // Ensure lock is released. 确保锁被释放。
	delete(a.mp, k)
}

// Len returns the number of key-value pairs in the map.
// Len 返回 map 中键值对的数量。
func (a *Map[K, V]) Len() int {
	a.mutex.RLock()         // Acquire read lock. 获取读锁。
	defer a.mutex.RUnlock() // Ensure lock is released. 确保锁被释放。
	return len(a.mp)
}

// Range iterates over all key-value pairs in the map, applying the given function. If the function returns false, the iteration stops.
// Range 遍历 map 中的所有键值对，并应用给定的函数。 如果函数返回 false，迭代将停止。
func (a *Map[K, V]) Range(run func(k K, v V) bool) {
	a.mutex.RLock()         // Acquire read lock. 获取读锁。
	defer a.mutex.RUnlock() // Ensure lock is released. 确保锁被释放。
	for k, v := range a.mp {
		if ok := run(k, v); !ok { // Stop iteration if the callback returns false. 如果回调返回 false，则停止迭代。
			return
		}
	}
}

// Getset retrieves the value associated with the key, or computes and stores a new value if the key does not exist.
// It returns the value and a boolean indicating whether a new value was created.
// Getset 获取与键关联的值，如果键不存在，则计算并存储新值。
// 它返回值以及一个布尔值，指示是否创建了新值。
func (a *Map[K, V]) Getset(k K, calculate func() V) (v V, created bool) {
	if v, ok := a.Get(k); ok { // Attempt to read the value with a read lock. 尝试用读锁读取值。
		return v, false // Return existing value if found. 如果找到，返回已存在的值。
	}
	// 读锁释放，启动写锁，但假设有两个线程同时读不到，就都会同时占用写锁。
	a.mutex.Lock()
	defer a.mutex.Unlock()
	// 增加读锁以后二次确认内容是否在 map 里面，这样第二次占用写锁的线程就不会创建新对象。
	if v, ok := a.mp[k]; ok { // Double-check under write lock. 在写锁下再次检查。
		return v, false // Return the existing value if found. 如果找到，返回已有值。
	}
	// 当内容确实不在 map 里时，即首次占用写锁时，这才创建新对象，设置到 map 里。
	v = calculate() // Calculate the value. 计算值。
	// This function might be expensive. If this is a concern, consider alternative caching solutions.
	// 这个函数可能耗时较长，如果对此介意，可以考虑使用其他缓存方案。
	a.mp[k] = v
	return v, true // Return the new value and indicate it was created. 返回新值并指示是新创建的。
}
