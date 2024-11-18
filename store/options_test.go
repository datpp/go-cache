package store

import (
	gocache_store "github.com/eko/gocache/lib/v4/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptionsForceRefreshValueTrue(t *testing.T) {
	storeOptions := gocache_store.ApplyOptions(WithForceRefresh(true))

	// When - Then
	assert.Contains(t, storeOptions.Tags, TAG_FORCE_REFRESH)
}

func TestOptionsForceRefreshValueFalse(t *testing.T) {
	storeOptions := gocache_store.ApplyOptions(WithForceRefresh(false))

	// When - Then
	assert.NotContains(t, storeOptions.Tags, TAG_FORCE_REFRESH)
}

func TestOptionsIgnoreErrorValueTrue(t *testing.T) {
	storeOptions := gocache_store.ApplyOptions(WithIgnoreCacheError(true))

	// When - Then
	assert.Contains(t, storeOptions.Tags, TAG_IGNORE_ERROR)
}

func TestOptionsIgnoreErrorValueFalse(t *testing.T) {
	storeOptions := gocache_store.ApplyOptions(WithIgnoreCacheError(false))

	// When - Then
	assert.NotContains(t, storeOptions.Tags, TAG_IGNORE_ERROR)
}
