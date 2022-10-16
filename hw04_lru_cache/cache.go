package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	item, ok := l.items[key]
	if ok {
		l.queue.MoveToFront(item)
		cItem := item.Value.(cacheItem)
		cItem.value = value
		item.Value = cItem

		return true
	}

	l.items[key] = l.queue.PushFront(cacheItem{
		key:   key,
		value: value,
	})

	if l.queue.Len() > l.capacity {
		backItem := l.queue.Back()
		cBackItem := backItem.Value.(cacheItem)
		delete(l.items, cBackItem.key)
		l.queue.Remove(backItem)
	}

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	item, ok := l.items[key]
	if !ok {
		return nil, false
	}

	l.queue.MoveToFront(item)
	item = l.queue.Front()
	cItem := item.Value.(cacheItem)

	return cItem.value, true
}

func (l *lruCache) Clear() {
	l.items = make(map[Key]*ListItem, l.capacity)
	l.queue = NewList()
}

type cacheItem struct {
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
