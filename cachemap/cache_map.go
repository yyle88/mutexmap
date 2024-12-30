package cachemap

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/yyle88/mutexmap"
)

type Map[K comparable, V any] struct {
	mp    map[K]*valueBottle[V]
	mutex *sync.RWMutex
}

type valueBottle[V any] struct {
	mutex *sync.RWMutex
	done  bool //这里的done仅仅表示已经执行过，而不表示结果成功，有可能结果是失败的，但也标记为"done"
	res   V
	err   error //即使是计算错误，也保存下来，再次遇到时直接返回错误
}

func NewMap[K comparable, V any](size int) *Map[K, V] {
	return &Map[K, V]{
		mp:    make(map[K]*valueBottle[V], size),
		mutex: &sync.RWMutex{},
	}
}

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

func (a *Map[K, V]) Set(k K, v V) {
	a.SetVe(k, v, nil)
}

func (a *Map[K, V]) SetKv(k K, v V) {
	a.SetVe(k, v, nil)
}

func (a *Map[K, V]) SetKe(k K, err error) {
	a.SetVe(k, mutexmap.Zero[V](), err)
}

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

// Getset get a value, if value is not exist, then create an object and set into map.
// when not exist call the newValue to create new value and put it in to map as cache.
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

func (a *Map[K, V]) setBottle(k K) (*valueBottle[V], *sync.RWMutex) {
	a.mutex.Lock()         //假如读不到数据再用写锁锁住，再去尝试读数据，假如确定还是读不到，就创建数据
	defer a.mutex.Unlock() //把最外面的大锁给释放掉，这把大锁的作用是避免重复创建空壳

	vx, exist := a.mp[k]
	if !exist {
		vx = &valueBottle[V]{ //创建个空壳数据
			mutex: &sync.RWMutex{},
			done:  false,
			res:   mutexmap.Zero[V](), // no need to set zero value
			err:   nil,                // no need to set none value
		}
		a.mp[k] = vx
		vx.mutex.Lock() //把这个壳给锁住，再开始计算数据，这时map的锁即可释放，而元素的锁将生效以确保需要同样元素的请求都得等待后面运算结束，即等待元素的释放锁
		return vx, vx.mutex
	}

	return vx, nil
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

func (a *Map[K, V]) Range(run func(k K, v V, err error, done bool) bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	for k, v := range a.mp {
		if ok := run(k, v.res, v.err, v.done); !ok {
			return
		}
	}
}
