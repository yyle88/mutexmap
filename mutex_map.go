package mutexmap

import "sync"

type Map[K comparable, V any] struct {
	mp    map[K]V
	mutex *sync.RWMutex
}

func NewMap[K comparable, V any](cap int) *Map[K, V] {
	return &Map[K, V]{
		mp:    make(map[K]V, cap),
		mutex: &sync.RWMutex{},
	}
}

func (a *Map[K, V]) Get(k K) (V, bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	if v, ok := a.mp[k]; ok {
		return v, ok
	} else {
		return v, false
	}
}

func (a *Map[K, V]) Set(k K, v V) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.mp[k] = v
}

func (a *Map[K, V]) Delete(k K) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	delete(a.mp, k)
}

func (a *Map[K, V]) Len() int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return len(a.mp)
}

func (a *Map[K, V]) Range(run func(k K, v V) bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	for k, v := range a.mp {
		if ok := run(k, v); !ok {
			return
		}
	}
}

// GetOrzSet get a value, if value is not exist, then create an object and set into map.
// return value, and created(true). if exist return not created(false).
// so when already exist do not change the value, return the old value.
// Orz means "or" but "or" is too ugly, more ugly than ugly. So I use Orz instead of it.
func (a *Map[K, V]) GetOrzSet(k K, newValue func() V) (v V, created bool) {
	if v, ok := a.Get(k); ok {
		return v, false
	}
	//读锁释放，启动写锁，但很明显假设有两个线程同时读不到，就都会同时去占用写锁
	a.mutex.Lock()
	defer a.mutex.Unlock()
	//因此增加读锁以后二次确认内容是否在map里面，这样第二次占用写锁的，就不会创建新对象啦
	if v, ok := a.mp[k]; ok {
		return v, false
	}
	//当内容确实是不在map里时，即首次占用写锁时，这才创建新对象，设置到map里
	v = newValue() // this func maybe always cost a lot of time 假如你非常介意这个的话这个包里还有另一个cache能用
	a.mp[k] = v
	return v, true
}
