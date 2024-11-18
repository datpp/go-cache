package cache

import (
	"github.com/eko/gocache/lib/v4/codec"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	store := store.NewMockStoreInterface(ctrl)

	// When
	cache := New[any](store)

	// Then
	assert.IsType(t, new(Cache[any]), cache)
	assert.IsType(t, new(codec.Codec), cache.Cache.GetCodec())

	assert.Equal(t, store, cache.Cache.GetCodec().GetStore())
}
