package helper

import (
	unsafeRandom "math/rand"
	"sync"
	"time"
)

var randomSourcePool sync.Pool

func init() {
	randomSourcePool.New = func() any {
		return unsafeRandom.New(unsafeRandom.NewSource(time.Now().UnixNano())) //nolint:gosec
	}
}

func UnsafeInt() int {
	randomSource := randomSourcePool.Get().(*unsafeRandom.Rand)
	defer randomSourcePool.Put(randomSource)

	return randomSource.Int()
}

func UnsafeIntn(n int) int {
	randomSource := randomSourcePool.Get().(*unsafeRandom.Rand)
	defer randomSourcePool.Put(randomSource)

	return randomSource.Intn(n)
}
