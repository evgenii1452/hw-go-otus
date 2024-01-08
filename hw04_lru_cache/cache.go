package hw04lrucache

import (
	"fmt"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type listItemValueWrapper struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	wrappedValue := listItemValueWrapper{
		key:   key,
		value: value,
	}

	if listItem, ok := c.items[key]; ok {
		listItem.Value = wrappedValue
		c.queue.MoveToFront(listItem)
		return true
	}

	if c.capacity == c.queue.Len() {
		lastItem, ok := c.queue.Back().Value.(listItemValueWrapper)
		if !ok {
			panic("Cache value is not listItemValueWrapper")
		}

		delete(c.items, lastItem.key)
		c.queue.Remove(c.queue.Back())
	}

	c.items[key] = c.queue.PushFront(wrappedValue)

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	fmt.Println(item.Value)

	valueWrapper, ok := item.Value.(listItemValueWrapper)
	if !ok {
		panic("Cache value is not listItemValueWrapper")
	}

	c.queue.MoveToFront(item)

	return valueWrapper.value, true
}

func (c *lruCache) Clear() {
	c.queue = new(list)
	c.items = make(map[Key]*ListItem, c.capacity)
}
