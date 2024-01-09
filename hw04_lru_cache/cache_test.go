package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("capacity logic", func(t *testing.T) {
		c := NewCache(4)
		cache := c.(*lruCache)

		c.Set("k1", "v1")
		c.Set("k2", "v2")
		c.Set("k3", "v3")
		c.Set("k4", "v4")

		require.Equal(t, 4, cache.queue.Len())
		require.Equal(t, 4, len(cache.items))

		// ОР: элемент k5 вытеснит элемент k1
		c.Set("k5", "v5")

		require.Equal(t, 4, cache.queue.Len())
		require.Equal(t, 4, len(cache.items))

		value, ok := c.Get("k1")
		require.Nil(t, value)
		require.False(t, ok)

		value, ok = c.Get("k2")
		require.Equal(t, "v2", value)
		require.True(t, ok)

		// ОР: элемент k6 вытеснит элемент k3, т.к. k2 перемещен в начало списка
		c.Set("k6", "v6")

		value, ok = c.Get("k3")
		require.Nil(t, value)
		require.False(t, ok)

		value, ok = c.Get("k2")
		require.Equal(t, "v2", value)
		require.True(t, ok)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(4)
		cache := c.(*lruCache)

		c.Set("k1", "v1")
		c.Set("k2", "v2")
		c.Set("k3", "v3")
		c.Set("k4", "v4")

		value, ok := c.Get("k1")
		require.Equal(t, "v1", value)
		require.True(t, ok)

		require.Equal(t, 4, cache.queue.Len())
		require.Equal(t, 4, len(cache.items))

		c.Clear()
		require.Equal(t, 0, cache.queue.Len())
		require.Equal(t, 0, len(cache.items))

		value, ok = c.Get("k1")
		require.Nil(t, value)
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
