package cachemap

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/yyle88/mutexmap/internal/utils"
)

type Map[K comparable, V any] struct {
	mp    map[K]*valueBottle[V]
	mutex *sync.RWMutex
}

// New creates a new thread-safe map.
// New 创建一个线程安全 map。
func New[K comparable, V any]() *Map[K, V] {
	return NewMap[K, V](8)
}

type valueBottle[V any] struct {
	mutex *sync.RWMutex
	done  bool //这里的done仅仅表示已经执行过，而不表示结果成功，有可能结果是失败的，但也标记为"done"
	res   V
	err   error //即使是计算错误，也保存下来，再次遇到时直接返回错误
}

// NewMap creates a new thread-safe map.
// NewMap 创建一个线程安全 map。
func NewMap[K comparable, V any](cap int) *Map[K, V] {
	return &Map[K, V]{
		mp:    make(map[K]*valueBottle[V], cap),
		mutex: &sync.RWMutex{},
	}
}

// Get retrieves a value by key, returning the value, error, and whether it was computed.
// Get 根据键获取值，返回值、错误以及是否已计算的状态。
func (a *Map[K, V]) Get(k K) (res V, err error, done bool) {
	a.mutex.RLock() //首先用读锁去试探性的读数据
	defer a.mutex.RUnlock()
	vx, ok := a.mp[k]
	if !ok {
		return res, errors.New("not exist"), false
	}
	if !vx.done {
		return res, errors.New("not acted"), false
	}
	return vx.res, vx.err, true
}

// Set stores a value for a key with no associated error.
// Set 为键存储一个值，不带错误信息。
func (a *Map[K, V]) Set(k K, v V) {
	a.SetVe(k, v, nil)
}

// SetKv stores a value for a key with no associated error (alias for Set).
// SetKv 为键存储一个值，不带错误信息（Set 的别名）。
func (a *Map[K, V]) SetKv(k K, v V) {
	a.SetVe(k, v, nil)
}

// SetKe stores an error for a key with a zero value.
// SetKe 为键存储一个错误，使用零值作为值。
func (a *Map[K, V]) SetKe(k K, err error) {
	a.SetVe(k, utils.Zero[V](), err)
}

// SetVe stores a value and an optional error for a key, overwriting any existing item.
// SetVe 为键存储值和可选错误，覆盖现有条目。
func (a *Map[K, V]) SetVe(k K, v V, err error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.mp[k] = &valueBottle[V]{
		mutex: &sync.RWMutex{}, //其实当设置done以后就不会再用锁啦，但没有关系的，依然写在这里避免疑惑和出错
		done:  true,
		res:   v,
		err:   err,
	}
}

// Getset retrieves a value by key or computes it using newValue if it doesn't exist, caching the result.
// Getset 根据键获取值，若不存在则使用 newValue 计算并缓存结果。
// Getset get a value, if value is not exist, then create an object and set into map.
// when not exist, call the newValue to create new value and put it in to map as cache.
// not lock all the map during call newValue even newValue is very slow and waste time.
// so when already exist do not change the value, return the old value.
// Orz means "or" but "or" is too ugly, more ugly than ugly. So I use Orz instead of it.
func (a *Map[K, V]) Getset(k K, newValue func() (V, error)) (V, error) {
	if res, err, done := a.Get(k); done {
		return res, err
	}
	//在写锁内执行-因此在占用前有可能会执行别的操作-比如设置新值/删除值-设计时要考虑时序问题
	vb, mu := a.setBottle(k)
	//释放mp的写锁-因此在释放后有可能会执行别的操作-比如设置新值/删除值-设计时要考虑时序问题
	if mu != nil {
		defer mu.Unlock()

		func() {
			defer func() {
				if ove := recover(); ove != nil {
					if erp, ok := ove.(error); ok {
						vb.err = erp
					} else {
						vb.err = errors.Errorf("panic occurred. reason: %v", ove)
					}
				}
			}()
			vb.res, vb.err = newValue() //接着去计算数据，最终释放壳子的写锁，确保一个资源只被计算一次
		}()
		vb.done = true
		return vb.res, vb.err
	} else {
		if vb.done {
			return vb.res, vb.err
		} else {
			vb.mutex.RLock() //抢占壳子的读锁，当能抢到的时候说明壳子已经有数据，抢不到说明还没计算完
			defer vb.mutex.RUnlock()
			//当抢到读锁时，说明前面的执行已经完成（done），结果已经被赋值到壳子里，就可以使用和返回
			return vb.res, vb.err
		}
	}
}

// setBottle creates or retrieves a valueBottle for a key, returning the bottle and its mutex if new.
// setBottle 为键创建或获取一个 valueBottle，若为新创建则返回 bottle 及其互斥锁。
func (a *Map[K, V]) setBottle(k K) (*valueBottle[V], *sync.RWMutex) {
	a.mutex.Lock()         //假如读不到数据再用写锁锁住，再去尝试读数据，假如确定还是读不到，就创建数据
	defer a.mutex.Unlock() //把最外面的大锁给释放掉，这把大锁的作用是避免重复创建空壳

	vx, exist := a.mp[k]
	if !exist {
		vx = &valueBottle[V]{ //创建个空壳数据
			mutex: &sync.RWMutex{},
			done:  false,
			res:   utils.Zero[V](), // no need to set zero value
			err:   nil,             // no need to set none value
		}
		a.mp[k] = vx
		vx.mutex.Lock() //把这个壳给锁住，再开始计算数据，这时map的锁即可释放，而元素的锁将生效以确保需要同样元素的请求都得等待后面运算结束，即等待元素的释放锁
		return vx, vx.mutex
	}

	return vx, nil
}

// Delete removes a key and its associated value from the map.
// Delete 从 map 中删除一个键及其关联的值。
func (a *Map[K, V]) Delete(k K) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	delete(a.mp, k)
}

// Len returns the number of key-value pairs in the map.
// Len 返回 map 中键值对的数量。
func (a *Map[K, V]) Len() int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return len(a.mp)
}

// Range iterates over the map, calling the provided function for each key-value pair until it returns false.
// Range 遍历 map，对每个键值对调用提供的函数，直到函数返回 false。
func (a *Map[K, V]) Range(run func(k K, v V, err error, done bool) bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	for k, v := range a.mp {
		if ok := run(k, v.res, v.err, v.done); !ok {
			return
		}
	}
}
