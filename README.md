[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yyle88/mutexmap/release.yml?branch=main&label=BUILD)](https://github.com/yyle88/mutexmap/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yyle88/mutexmap)](https://pkg.go.dev/github.com/yyle88/mutexmap)
[![Coverage Status](https://img.shields.io/coveralls/github/yyle88/mutexmap/master.svg)](https://coveralls.io/github/yyle88/mutexmap?branch=main)
![Supported Go Versions](https://img.shields.io/badge/Go-1.22%2C%201.23-lightgrey.svg)
[![GitHub Release](https://img.shields.io/github/release/yyle88/mutexmap.svg)](https://github.com/yyle88/mutexmap/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yyle88/mutexmap)](https://goreportcard.com/report/github.com/yyle88/mutexmap)

# MutexMap - A Thread-Safe Map for Go

A thread-safe map implementation for Go, using `sync.RWMutex` to synchronize access. This package is optimized for scenarios involving concurrent reads and writes, providing efficient and reliable operations for multi-thread applications.

## README

[中文说明](README.zh.md)

## Overview

Go’s standard `map` is not safe for concurrent access. This package wraps a `map[K]V` with `sync.RWMutex` to provide thread safety. With `RWMutex`, multiple readers can access the map simultaneously, while writes are synchronized to prevent race conditions.

Key highlights:
- **Thread-safety**: Prevents data races in concurrent environments.
- **Efficient Reads**: Read operations (`Get`, `Range`) are lock-free for other readers.
- **Synchronous Write**: Write operations (`Set`, `Delete`, `Getset`) are synchronous.

This package is suitable for use cases requiring frequent reads and occasional writes with a shared map.

## Installation

```bash  
go get github.com/yyle88/mutexmap  
```  

## Example Usage

### Basic Operations

```go  
package main  

import (  
	"fmt"  
	"github.com/yyle88/mutexmap"  
)  

func main() {  
	mp := mutexmap.NewMap   

	mp.Set("key1", 100)  
	mp.Set("key2", 200)  

	if value, found := mp.Get("key1"); found {  
		fmt.Println("Key1 Value:", value)  
	}  

	mp.Range(func(key string, value int) bool {  
		fmt.Println(key, value)  
		return true  
	})  
}  
```  

### Using `Getset` for Cached Initialization

```go  
package main  

import (  
	"fmt"  
	"github.com/yyle88/mutexmap"  
)  

func main() {  
	mp := mutexmap.NewMap   

	value, created := mp.Getset("exampleKey", func() string {  
		return "This is a computed value"  
	})  
	fmt.Println("Created:", created, "Value:", value)  

	value, created = mp.Getset("exampleKey", func() string {  
		return "Another computed value"  
	})  
	fmt.Println("Created:", created, "Value:", value)  
}  
```  

## Features

- **Concurrent Access**: Allows multiple goroutines to safely read and write to the map.
- **Optimized Reads**: Supports simultaneous reads for better performance.
- **Custom Initialization**: Use `Getset` to initialize values only if they don’t already exist.

---

## Method Summary

| Method                                 | Description                                                                                      |  
|----------------------------------------|--------------------------------------------------------------------------------------------------|  
| `NewMap[K comparable, V any](cap int)` | Creates a new `Map` with an optional initial capacity.                                           |  
| `Get(k K) (V, bool)`                   | Retrieves the value for the given key. Returns `false` if the key is not found.                  |  
| `Set(k K, v V)`                        | Sets a value for the given key. Overwrites if the key already exists.                            |  
| `Delete(k K)`                          | Removes the key-value pair from the map.                                                         |  
| `Len() int`                            | Returns the number of elements in the map.                                                       |  
| `Range(func(k K, v V) bool)`           | Iterates over all key-value pairs. Stops if the callback returns `false`.                        |  
| `Getset(k K, func() V) (V, bool)`      | Gets a value or creates it if it doesn’t exist, ensuring the creation is atomic and thread-safe. |  

---

## Why Use MutexMap?

1. **Thread Safety**: Essential for shared maps in multi-threaded environments.
2. **Efficient Reads**: Read lock (`RLock`) ensures non-blocking reads for other readers.
3. **Write Synchronization**: Write lock (`Lock`) ensures data integrity during modifications.
4. **Flexible Initialization**: The `Getset` method prevents redundant computations.

---

## License

`mutexmap` is open-source and released under the MIT License. See the [LICENSE](LICENSE) file for more information.

---

## Support

Welcome to contribute to this project by submitting pull requests or reporting issues.

If you find this package helpful, give it a star on GitHub!

**Thank you for your support!**

**Happy Coding with `mutexmap`!** 🎉

Give me stars. Thank you!!!

## GitHub Stars

[![starring](https://starchart.cc/yyle88/mutexmap.svg?variant=adaptive)](https://starchart.cc/yyle88/mutexmap)
