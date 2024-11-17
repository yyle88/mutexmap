# MutexMap - A Thread-Safe Map with Mutex Locking

A thread-safe map implementation for Go, using a `sync.RWMutex` to protect concurrent access. This map ensures that read operations do not block each other, while write operations are properly synchronized to avoid data races.

## README
[中文说明](README.zh.md)

## Overview

Go’s standard `map` type is not safe for concurrent use. This package wraps a `map[K]V` with a read-write mutex (`sync.RWMutex`) to provide safe concurrent access. Read operations (Get and Range) are optimized with a read lock, allowing multiple readers to access the map concurrently. Write operations (Set, Delete, and Getset) acquire a write lock to ensure that only one writer modifies the map at a time.

The main advantage of this approach is that it allows efficient concurrent access, with multiple goroutines able to read from the map simultaneously while ensuring that write operations are safely coordinated.

## Installation

```bash
go get github.com/yyle88/mutexmap
```

## Usage

### Example 1: Basic Map Operations

```go
package main

import (
	"fmt"
	"github.com/yyle88/mutexmap"
)

func main() {
	// Create a new MutexMap with string keys and int values
	mp := mutexmap.NewMap[string, int](0)
	mp.Set("a", 10)
	mp.Set("b", 20)
	mp.Set("c", 30)

	// Get a value
	if v, ok := mp.Get("b"); ok {
		fmt.Println("Key 'b' has value:", v)
	} else {
		fmt.Println("Key 'b' not found")
	}

	// Range through the map
	mp.Range(func(k string, v int) bool {
		fmt.Println(k, v)
		return true
	})
}
```

### Example 2: Using the `Getset` Method

```go
package main

import (
	"fmt"
	"github.com/yyle88/mutexmap"
)

// Example of using the Getset method with a map that stores strings
func main() {
	mp := mutexmap.NewMap[string, string](0) 

	// Using Get a value
	value, created := mp.Getset("newKey", func() string {
		return "This is a newly created value"
	})
	if created {
		fmt.Println("New value created:", value)
	} else {
		fmt.Println("Existing value:", value)
	}

	// Attempt to Getset again - this should not create a new value
	value, created = mp.Getset("newKey", func() string {
		return "This should not be created again"
	})
	if created {
		fmt.Println("New value created:", value)
	} else {
		fmt.Println("Existing value:", value)
	}
}
```

## Features

- **Thread-Safety**: Uses a `sync.RWMutex` to ensure safe concurrent access to the map.
    - Read operations (`Get`, `Range`) are lock-free for other readers.
    - Write operations (`Set`, `Delete`, `Getset`) are serialized using a write lock.

- **Optimized for Concurrency**: Multiple readers can access the map simultaneously, while write operations are exclusive.

- **`Getset` Method**: A special method that retrieves the value for a given key, and if the key doesn't exist, it creates and stores a new value. This method helps avoid redundant object creation when multiple goroutines might attempt to create a value at the same time.

## Key Methods

- **`NewMap[K comparable, V any](cap int) *Map[K, V]`**: Initializes a new `Map` with the specified initial capacity.

- **`Get(k K) (V, bool)`**: Retrieves the value associated with key `k`. Returns `false` if the key is not found.

- **`Set(k K, v V)`**: Sets the value for the given key `k`.

- **`Delete(k K)`**: Deletes the key-value pair for the given key `k`.

- **`Len() int`**: Returns the number of key-value pairs in the map.

- **`Range(run func(k K, v V) bool)`**: Iterates over all key-value pairs in the map and applies the `run` function. Iteration stops if the function returns `false`.

- **`Getset(k K, createNewValue func() V) (v V, created bool)`**: Retrieves the value for the given key `k`. If the key doesn’t exist, it creates and stores a new value using the `createNewValue` function. Returns the value and a `created` flag indicating whether a new value was created.

## Why Use This Package?

This package is ideal for situations where you need a thread-safe map but still want to allow concurrent reads. By using a `sync.RWMutex`, you get the benefit of multiple readers accessing the map without blocking each other, while ensuring that writes are properly synchronized.

The `Getset` method is particularly useful in scenarios where you need to ensure that a value is created only once (even when multiple goroutines try to access the same key simultaneously), reducing unnecessary computation and locking.

## Performance Considerations

- **Concurrent Reads**: Multiple goroutines can perform read operations (`Get` and `Range`) simultaneously without blocking each other.

- **Writes**: Write operations (`Set`, `Delete`, and `Getset`) are serialized, so only one goroutine can perform a write at any given time. However, read operations do not block writes, and vice versa.

## When to Use

This package is useful in concurrent programs where multiple goroutines need to read from and write to a shared map. It is especially beneficial when:

- **You need thread-safety** for a shared map.
- **You want concurrent read access** but want to serialize writes to avoid data races.
- **You need to ensure the value is created only once** for a given key, using the `Getset` method.

## Example Test

Check out the test file for a full working example: [mutexmap_test.go](mutexmap_test.go).

## Conclusion

The `mutexmap` package provides a simple yet powerful solution for concurrent maps in Go. By using a read-write mutex (`sync.RWMutex`), it allows multiple readers while serializing writes, ensuring thread safety and high concurrency performance. With the addition of the `Getset` method, this package offers an elegant solution for safely creating and caching values in a concurrent environment.

If you find this package helpful, please consider giving it a star! Thank you for using it!

## Thank You

If you find this package valuable, give it a star on GitHub! Thank you!!!
