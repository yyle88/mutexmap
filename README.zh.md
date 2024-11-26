# MutexMap - 基于 Mutex 实现的线程安全 Map

MutexMap 是一个基于 `sync.RWMutex` 和标准 `map[K]V` 的线程安全实现，专为需要高效并发读写的场景设计。通过 `RWMutex` 实现读写锁，支持多线程安全访问，适合多协程环境下的共享数据操作。

## 目录

[ENGLISH README](README.md)

---

## 项目概览

Go 标准库的 `map` 在并发访问场景下并不安全。本库对 `map[K]V` 进行了封装，借助 `sync.RWMutex` 提供线程安全性，支持高效的并发操作。

### 核心特点

- **线程安全**：通过读写锁避免数据竞争。
- **高效读操作**：允许多个读操作同时进行，不互相阻塞。
- **同步写操作**：使用写锁控制写操作，确保数据一致性。

### 适用场景

适合以下需求的应用场景：
- 需要共享存储的数据结构。
- 多读少写的并发访问场景。
- 需要通过函数安全初始化值的场景（`Getset` 方法）。

---

## 安装

使用 `go get` 安装本库：

```bash
go get github.com/yyle88/mutexmap
```

---

## 示例代码

### 基础操作示例

```go
package main

import (
	"fmt"
	"github.com/yyle88/mutexmap"
)

func main() {
	mp := mutexmap.NewMap.Set("key1", 100)
	mp.Set("key2", 200)

	// 获取值
	if value, found := mp.Get("key1"); found {
		fmt.Println("Key1 的值:", value)
	} else {
		fmt.Println("Key1 不存在")
	}

	// 遍历所有键值对
	mp.Range(func(key string, value int) bool {
		fmt.Println(key, value)
		return true
	})
}
```

### 使用 `Getset` 方法实现值的缓存初始化

```go
package main

import (
	"fmt"
	"github.com/yyle88/mutexmap"
)

func main() {
	mp := mutexmap.NewMap 

	// 如果键不存在，则, created := mp.Getset("exampleKey", func() string {
		return "这是一段计算得到的值"
	})
	fmt.Println("值是否新创建:", created, "值:", value)

	// 再次调用 Getset 方法，值不会被重新创建
	value, created = mp.Getset("exampleKey", func() string {
		return "这段值不会被创建"
	})
	fmt.Println("值是否新创建:", created, "值:", value)
}
```

---

## 功能概览

### 方法一览

| 方法                                     | 描述                                            |
|----------------------------------------|-----------------------------------------------|
| `NewMap[K comparable, V any](cap int)` | 创建一个新的线程安全 Map，支持设置初始容量。                      |
| `Get(k K) (V, bool)`                   | 获取指定键的值，如果键不存在，返回 `false`。                    |
| `Set(k K, v V)`                        | 设置指定键的值，如果键已存在则覆盖旧值。                          |
| `Delete(k K)`                          | 删除指定键的键值对。                                    |
| `Len() int`                            | 返回当前 Map 中的键值对数量。                             |
| `Range(func(k K, v V) bool)`           | 遍历所有键值对，并对每个键值对执行传入的回调函数。如果回调返回 `false`，终止遍历。 |
| `Getset(k K, func() V) (V, bool)`      | 获取指定键的值。如果键不存在，使用回调函数计算值并存储，返回值和是否新创建的标志。     |

---

## MutexMap 的优势

1. **线程安全**  
   借助读写锁实现线程安全，避免并发操作导致的数据竞争问题。

2. **高效的读操作**  
   使用读锁 (`RLock`) 支持多个读操作同时进行，提升并发性能。

3. **写操作的完整性**  
   写锁 (`Lock`) 确保写操作的同步性，防止数据不一致。

4. **灵活的值初始化**  
   `Getset` 方法确保每个键的值只被安全创建一次，避免重复计算或初始化。

---

## 贡献与支持

非常欢迎提交问题、建议或功能请求!!!

如果觉得本库对您有帮助，请在 GitHub 上给个 ⭐，感谢支持！！！
