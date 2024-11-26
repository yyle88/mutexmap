package mutexmap

import (
	"testing"
)

func TestZero(t *testing.T) {
	res := Zero[int]()
	t.Log(res)
}
