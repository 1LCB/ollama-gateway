package utils

import (
	"sync"
	"time"
)

type TemporaryMap[K comparable, V any] struct {
	data       map[K]item[V]
	mu         sync.RWMutex
	defaultTTL time.Duration
}

type item[V any] struct {
	value  V
	expiry time.Time
}

func NewTemporaryMap[K comparable, V any](defaultTTL time.Duration) *TemporaryMap[K, V] {
	tempMap := &TemporaryMap[K, V]{
		data:       make(map[K]item[V]),
		defaultTTL: defaultTTL,
	}

	go tempMap.Clean(15 * time.Minute)
	return tempMap
}

func (tm *TemporaryMap[K, V]) Set(key K, value V, ttl time.Duration) {
	if ttl == 0 {
		ttl = tm.defaultTTL
	}
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.data[key] = item[V]{
		value:  value,
		expiry: time.Now().Add(ttl),
	}

}

func (tm *TemporaryMap[K, V]) Get(key K) (value V, exists bool) {
	tm.mu.RLock()
	it, exists := tm.data[key]
	tm.mu.RUnlock()

	if !exists || it.expiry.Before(time.Now()) {
		var zero V
		return zero, false
	}

	return it.value, true
}

func (tm *TemporaryMap[K, V]) Delete(key K) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	delete(tm.data, key)
}

func (tm *TemporaryMap[K, V]) Clean(delay time.Duration) {
	for range time.Tick(delay) {
		now := time.Now()
		var keysToBeDeleted []K

		tm.mu.RLock()
		for key, it := range tm.data {
			if it.expiry.Before(now) {
				keysToBeDeleted = append(keysToBeDeleted, key)
			}
		}
		tm.mu.RUnlock()

		tm.mu.Lock()
		for _, k := range keysToBeDeleted {
			delete(tm.data, k)
		}
		tm.mu.Unlock()
	}
}
