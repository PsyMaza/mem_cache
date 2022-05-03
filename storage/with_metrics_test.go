package storage

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

const parallelFactor = 10_000

func BenchmarkBalancedLoad(b *testing.B) {
	cache := NewMemCache()
	for i := 0; i < b.N; i++ {
		benchEmulateLoad(cache, parallelFactor)
	}
}

func BenchmarkReadIntensiveLoad(b *testing.B) {
	cache := NewMemCache()
	for i := 0; i < b.N; i++ {
		benchEmulateReadIntensiveLoad(cache, parallelFactor)
	}
}

func benchEmulateLoad(c Cache, parallelFactor int) {
	wg := sync.WaitGroup{}

	for i := 0; i < parallelFactor; i++ {
		key := fmt.Sprintf("#{%v}-key", i)
		value := fmt.Sprintf("#{%v}-value", i)

		wg.Add(1)
		go func(key, value string) {
			err := c.Set(key, value)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(key, value)

		wg.Add(1)
		go func(key, value string) {
			_, err := c.Get(key)

			if err != nil && !errors.Is(err, ErrNotFound) {
				panic(err)
			}
			wg.Done()
		}(key, value)

		wg.Add(1)
		go func(key string) {
			err := c.Delete(key)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(key)
	}
}

func benchEmulateReadIntensiveLoad(c Cache, parallelFactor int) {
	wg := sync.WaitGroup{}

	for i := 0; i < parallelFactor/10; i++ {
		key := fmt.Sprintf("#{%v}-key", i)
		value := fmt.Sprintf("#{%v}-value", i)

		wg.Add(1)
		go func(key, value string) {
			err := c.Set(key, value)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(key, value)

		wg.Add(1)
		go func(key string) {
			err := c.Delete(key)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(key)
	}

	for i := 0; i < parallelFactor; i++ {
		key := fmt.Sprintf("#{%v}-key", i)
		value := fmt.Sprintf("#{%v}-value", i)

		wg.Add(1)
		go func(key, value string) {
			_, err := c.Get(key)

			if err != nil && !errors.Is(err, ErrNotFound) {
				panic(err)
			}
			wg.Done()
		}(key, value)
	}
}
