package cache

import (
	"errors"
	"sync"
)

type CacheKeyConstraint interface {
	comparable
}

type CacheValueConstraint interface {
	any
}

// Utility used to cache values using key-value mapping
type Cache[TKey CacheKeyConstraint, TValue CacheValueConstraint] interface {
	// Check if the cache contains a value associated to a given key
	HasKey(k TKey) bool

	// Retrieve a value associated to a given key
	GetValue(k TKey) (TValue, error)

	// Add or overwrite a cached value associated to a given key. The return value indicates if the operation was successful
	SetValue(k TKey, v TValue) bool
}

type inMemoryCache[TKey CacheKeyConstraint, TValue CacheValueConstraint] struct {
	cacheMap map[TKey]TValue
	mutex    sync.RWMutex
	maxSize  int
}

// Create a new instance of a in-memory cache with a given cache size. If the size is set to 0 the cache will have not size limit
func CreateInMemoryCache[TKey CacheKeyConstraint, TValue CacheValueConstraint](maxSize int) (Cache[TKey, TValue], error) {
	if maxSize < 0 {
		return nil, errors.New("cache: invalid cache size specified")
	}

	cache := new(inMemoryCache[TKey, TValue])
	cache.cacheMap = make(map[TKey]TValue, maxSize)
	cache.mutex = sync.RWMutex{}
	cache.maxSize = maxSize
	return cache, nil
}

func (cache *inMemoryCache[TKey, TValue]) HasKey(k TKey) bool {
	cache.mutex.RLock()
	_, ok := cache.cacheMap[k]
	cache.mutex.RUnlock()

	return ok
}

func (cache *inMemoryCache[TKey, TValue]) GetValue(k TKey) (TValue, error) {
	cache.mutex.RLock()
	value, ok := cache.cacheMap[k]
	cache.mutex.RUnlock()

	if !ok {
		return *new(TValue), errors.New("cache: the cache does not contain a value associated to the given key")
	}

	return value, nil
}

func (cache *inMemoryCache[TKey, TValue]) SetValue(k TKey, v TValue) bool {
	cache.mutex.Lock()

	_, ok := cache.cacheMap[k]
	if ok {
		cache.cacheMap[k] = v
		cache.mutex.Unlock()
		return true
	}

	if cache.maxSize != 0 && len(cache.cacheMap) == cache.maxSize {
		cache.mutex.Unlock()
		return false
	}

	cache.cacheMap[k] = v
	cache.mutex.Unlock()
	return true
}
