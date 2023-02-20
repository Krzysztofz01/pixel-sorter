package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryCacheShouldCreate(t *testing.T) {
	cache, err := CreateInMemoryCache[int, float64](25)

	assert.Nil(t, err)
	assert.NotNil(t, cache)
}

func TestInMemoryCacheShouldSetValue(t *testing.T) {
	cache, err := CreateInMemoryCache[int, float64](25)
	assert.Nil(t, err)

	expectedKey := 10
	expectedValue := 20.5

	ok := cache.SetValue(expectedKey, expectedValue)
	assert.True(t, ok)
}

func TestInMemoryCacheShouldGetValue(t *testing.T) {
	cache, err := CreateInMemoryCache[int, float64](25)
	assert.Nil(t, err)

	expectedKey := 10
	expectedValue := 20.5

	ok := cache.SetValue(expectedKey, expectedValue)
	assert.True(t, ok)

	actualValue, err := cache.GetValue(expectedKey)
	assert.Nil(t, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestInMemoryCacheShouldTellIfKeyExists(t *testing.T) {
	cache, err := CreateInMemoryCache[int, float64](25)
	assert.Nil(t, err)

	expectedKey := 10
	expectedValue := 20.5

	exists := cache.HasKey(expectedKey)
	assert.False(t, exists)

	ok := cache.SetValue(expectedKey, expectedValue)
	assert.True(t, ok)

	exists = cache.HasKey(expectedKey)
	assert.True(t, exists)
}

func TestInMemoryCacheShouldSetAndOverwriteValue(t *testing.T) {
	cache, err := CreateInMemoryCache[int, float64](25)
	assert.Nil(t, err)

	expectedKey := 10
	expectedValue := 20.5

	ok := cache.SetValue(expectedKey, expectedValue)
	assert.True(t, ok)

	actualValue, err := cache.GetValue(expectedKey)
	assert.Nil(t, err)
	assert.Equal(t, expectedValue, actualValue)

	expectedValueChanged := 30.5

	ok = cache.SetValue(expectedKey, expectedValueChanged)
	assert.True(t, ok)

	actualValue, err = cache.GetValue(expectedKey)
	assert.Nil(t, err)
	assert.Equal(t, expectedValueChanged, actualValue)
}

func TestInMemoryCacheShouldRespectSpecifiedCacheSize(t *testing.T) {
	cache, err := CreateInMemoryCache[int, float64](1)
	assert.Nil(t, err)

	ok := cache.SetValue(10, 20.5)
	assert.True(t, ok)

	ok = cache.SetValue(5, 20.5)
	assert.False(t, ok)

	ok = cache.SetValue(10, 30.5)
	assert.True(t, ok)
}
