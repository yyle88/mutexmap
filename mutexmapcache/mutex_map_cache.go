package mutexmapcache

import (
	"sync"

	"github.com/pkg/errors"
)

type Map[K comparable, V any] struct {
	mp    map[K]*vaItem[V]
	mutex *sync.RWMutex
}

type vaItem[V any] struct {
	mu   *sync.RWMutex
	done bool //这里的done仅仅表示已经执行过，而不表示结果成功，有可能结果是失败的，但也标记为"done"
	res  V
	erx  error //即使是计算错误，也保存下来，再次遇到时直接返回错误
}

func NewMap[K comparable, V any](size int) *Map[K, V] {
	return &Map[K, V]{
		mp:    make(map[K]*vaItem[V], size),
		mutex: &sync.RWMutex{},
	}
}

func (a *Map[K, V]) Get(k K) (res V, erx error, done bool) {
	a.mutex.RLock() //首先用读锁去试探性的读数据
	defer a.mutex.RUnlock()
	vx, ok := a.mp[k]
	if !ok {
		return res, errors.New("not exist"), false
	}
	if !vx.done {
		return res, errors.New("not acted"), false
	}
	return vx.res, vx.erx, true
}

func (a *Map[K, V]) Set(k K, v V) {
	a.SetVE(k, v, nil)
}

func (a *Map[K, V]) SetVE(k K, v V, erx error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.mp[k] = &vaItem[V]{
		mu:   nil, //这里暂时不需要，因为目前所有用的地方会优先判断是否done，只要done的就不需要使用锁，以后需要时再说吧
		done: true,
		res:  v,
		erx:  erx,
	}
}

// Getset get a value, if value is not exist, then create an object and set into map.
// when not exist call the newValue to create new value and put it in to map as cache.
// not lock all the map during call newValue even newValue is very slow and waste time.
// so when already exist do not change the value, return the old value.
// Orz means "or" but "or" is too ugly, more ugly than ugly. So I use Orz instead of it.
func (a *Map[K, V]) Getset(k K, newValue func() (V, error)) (V, error) {
	if res, erx, done := a.Get(k); done {
		return res, erx
	}
	//在写锁内执行-因此在占用前有可能会执行别的操作-比如设置新值/删除值-设计时要考虑时序问题
	vx, rwMutex := a.setVaOnce(k)
	//释放mp的写锁-因此在释放后有可能会执行别的操作-比如设置新值/删除值-设计时要考虑时序问题
	if rwMutex != nil {
		defer rwMutex.Unlock()

		func() {
			defer func() {
				if ove := recover(); ove != nil {
					if erp, ok := ove.(error); ok {
						vx.erx = erp
					} else {
						vx.erx = errors.Errorf("panic occurred. reason: %v", ove)
					}
				}
			}()
			vx.res, vx.erx = newValue() //接着去计算数据，最终释放壳子的写锁，确保一个资源只被计算一次
		}()
		vx.done = true
		return vx.res, vx.erx
	} else {
		if vx.done {
			return vx.res, vx.erx
		} else {
			vx.mu.RLock() //抢占壳子的读锁，当能抢到的时候说明壳子已经有数据，抢不到说明还没计算完
			defer vx.mu.RUnlock()
			//当抢到读锁时，说明前面的执行已经完成（done），结果已经被赋值到壳子里，就可以使用和返回
			return vx.res, vx.erx
		}
	}
}

func (a *Map[K, V]) setVaOnce(k K) (*vaItem[V], *sync.RWMutex) {
	a.mutex.Lock()         //假如读不到数据再用写锁锁住，再去尝试读数据，假如确定还是读不到，就创建数据
	defer a.mutex.Unlock() //把最外面的大锁给释放掉，这把大锁的作用是避免重复创建空壳

	vx, exist := a.mp[k]
	if !exist {
		vx = &vaItem[V]{ //创建个空壳数据
			mu:   &sync.RWMutex{},
			done: false,
			//res:  zero[V](), // no need to set zero value
			//erx:  nil,       // no need to set none value
		}
		a.mp[k] = vx
		vx.mu.Lock() //把这个壳给锁住，再开始计算数据，这时map的锁即可释放，而元素的锁将生效以确保需要同样元素的请求都得等待后面运算结束，即等待元素的释放锁
		return vx, vx.mu
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

func (a *Map[K, V]) Range(run func(k K, v V, erx error, done bool) bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	for k, v := range a.mp {
		if ok := run(k, v.res, v.erx, v.done); !ok {
			return
		}
	}
}
