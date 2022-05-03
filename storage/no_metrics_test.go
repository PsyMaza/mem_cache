package storage

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func Test_Cache(t *testing.T) {
	t.Parallel()

	var testCache Cache
	testCache = NewMemCache()

	t.Run("correctly stored value", func(t *testing.T) {
		t.Parallel()

		key := "key"
		value := "value"

		err := testCache.Set(key, value)
		assert.NoError(t, err)
		storedValue, err := testCache.Get(key)
		assert.NoError(t, err)

		assert.Equal(t, value, storedValue)
	})

	t.Run("correctly update value", func(t *testing.T) {
		t.Parallel()

		key := "key"
		value := "value"

		err := testCache.Set(key, value)
		assert.NoError(t, err)
		storedValue, err := testCache.Get(key)
		assert.NoError(t, err)

		assert.Equal(t, value, storedValue)

		newValue := "value2"

		err = testCache.Set(key, newValue)
		assert.NoError(t, err)
		newStoredValue, err := testCache.Get(key)
		assert.NoError(t, err)

		assert.Equal(t, newValue, newStoredValue)
	})

	t.Run("correctly delete value", func(t *testing.T) {
		t.Parallel()

		key := "key"
		value := "value"

		err := testCache.Set(key, value)
		assert.NoError(t, err)
		err = testCache.Delete(key)
		assert.NoError(t, err)
		_, err = testCache.Get(key)
		assert.EqualError(t, err, "value not found")
	})

	t.Run("no data races", func(t *testing.T) {
		t.Parallel()

		var testCache Cache
		testCache = NewMemCache()

		parallelFactor := 100_000_0

		emulateLoad(t, testCache, parallelFactor)
	})
}

func emulateLoad(t *testing.T, c Cache, parallelFactor int) {
	wg := sync.WaitGroup{}

	for i := 0; i < parallelFactor; i++ {
		key := fmt.Sprintf("#{%v}-key", i)
		value := fmt.Sprintf("#{%v}-value", i)

		wg.Add(1)
		go func(key, value string) {
			err := c.Set(key, value)
			assert.NoError(t, err)
			wg.Done()
		}(key, value)

		wg.Add(1)
		go func(key, value string) {
			v, err := c.Get(key)
			if !errors.Is(err, ErrNotFound) {
				assert.Equal(t, value, v)
			}
			wg.Done()
		}(key, value)

		wg.Add(1)
		go func(key string) {
			err := c.Delete(key)
			assert.NoError(t, err)
			wg.Done()
		}(key)
	}

	wg.Wait()
}
