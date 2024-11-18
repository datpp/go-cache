package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	store_cache "github.com/datpp/go-cache/store"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/mitchellh/hashstructure/v2"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/singleflight"
	"time"
)

var cacheProcessGroup singleflight.Group

type Cache[T any] struct {
	cache.Cache[any]
	options *Options
}

func New[T any](store store.StoreInterface, options ...Option) *Cache[T] {
	cacheOptions := ApplyOptions(options...)

	return &Cache[T]{
		Cache:   *cache.New[any](store),
		options: cacheOptions,
	}
}

func (s *Cache[T]) Get(ctx context.Context, key any) (T, error) {
	cKey, err := s.getCacheKey(key)
	if err != nil {
		return *new(T), fmt.Errorf("get / fail create cache key: %w", err)
	}

	data, err := s.Cache.Get(ctx, cKey)
	if err != nil {
		return *new(T), fmt.Errorf("cache get: %w", err)
	}

	rs, err := s.unmarshal(data)
	if err != nil {
		s.Delete(ctx, cKey) // usecase: different cache data version maybe make loop forever. So, we need to delete cache if we can't unmarshal it
		return *new(T), fmt.Errorf("cache get unmarshal: %w", err)
	}

	return rs, nil
}

func (s *Cache[T]) GetWithTTL(ctx context.Context, key any) (T, time.Duration, error) {
	cKey, err := s.getCacheKey(key)
	if err != nil {
		return *new(T), 0, fmt.Errorf("getwithttl / fail create cache key: %w", err)
	}

	data, ttl, err := s.Cache.GetWithTTL(ctx, cKey)
	if err != nil {
		return *new(T), 0, fmt.Errorf("cache get: %w", err)
	}

	rs, err := s.unmarshal(data)
	if err != nil {
		s.Delete(ctx, cKey) // usecase: different cache data version maybe make loop forever. So, we need to delete cache if we can't unmarshal it
		return *new(T), 0, fmt.Errorf("cache get unmarshal: %w", err)
	}

	return rs, ttl, nil
}

func (s *Cache[T]) Set(ctx context.Context, key any, value T, options ...store.Option) error {
	cKey, err := s.getCacheKey(key)
	if err != nil {
		return fmt.Errorf("set / fail create cache key: %w", err)
	}

	data, err := s.marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return s.Cache.Set(ctx, cKey, data, options...)
}

func (c *Cache[T]) Delete(ctx context.Context, key any) error {
	cKey, err := c.getCacheKey(key)
	if err != nil {
		return fmt.Errorf("delete / fail create cache key: %w", err)
	}

	return c.Cache.Delete(ctx, cKey)
}

// GetOnce will get the value from cache, if it doesn't exist, it will call the getFn function to get the value and set it to cache
// support extra option WithForceRefresh to force refresh the cache
// see more at README.md
func (s *Cache[T]) GetOnce(ctx context.Context, key any, getFn func() (T, error), cacheSetOptions ...store.Option) (T, error) {
	cKey, err := s.getCacheKey(key)
	if err != nil {
		return *new(T), fmt.Errorf("getonce / fail create cache key: %w", err)
	}

	loadAndSet := func() (T, error) {
		rs, err, _ := cacheProcessGroup.Do(cKey, func() (interface{}, error) {
			cbResponse, err := getFn()
			if err != nil {
				return nil, err
			}

			// if you can't set cache, don't block the request just return it and next time try it again
			_ = s.Set(ctx, key, cbResponse, cacheSetOptions...)

			return cbResponse, nil
		})

		if err != nil {
			return *new(T), fmt.Errorf("singleflight / cache get once: %w", err)
		}

		if v, ok := rs.(T); ok {
			return v, nil
		}

		return *new(T), fmt.Errorf("singleflight / get once convert data fail: %w", err)
	}

	// for force refresh
	cOptions := store.ApplyOptions(cacheSetOptions...)
	forceRefresh := slices.Contains(cOptions.Tags, store_cache.TAG_FORCE_REFRESH)
	ignoreError := slices.Contains(cOptions.Tags, store_cache.TAG_IGNORE_ERROR)

	if forceRefresh {
		return loadAndSet()
	}

	data, err := s.Get(ctx, key)
	if err != nil {
		var NotFoundErr *store.NotFound
		if ignoreError || errors.As(err, &NotFoundErr) {
			return loadAndSet()
		}

		return *new(T), fmt.Errorf("singleflight / get data fail: %w", err)
	}

	return data, nil
}

func (s *Cache[T]) getCacheKey(key any) (string, error) {
	var cKey string
	var err error

	switch v := key.(type) {
	case string:
		cKey = v
	default:
		cKey, err = checksum(key)
		if err != nil {
			return "", fmt.Errorf("failed to create cache key: %w", err)
		}
	}

	if s.options.CachePrefix == "" {
		return cKey, nil
	}

	return s.options.CachePrefix + ":" + cKey, nil
}

func (s *Cache[T]) marshal(value T) ([]byte, error) {
	bData, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize data: %w", err)
	}

	return bData, nil
}

func (s *Cache[T]) unmarshal(data any) (T, error) {
	var result T
	var err error

	switch v := data.(type) {
	case []byte:
		err = json.Unmarshal(v, &result)
	case string:
		err = json.Unmarshal([]byte(v), &result)
	}

	if err != nil {
		return *new(T), fmt.Errorf("failed to deserialize data: %w", err)
	}

	return result, nil
}

// checksum hashes a given object into a string
func checksum(object any) (string, error) {
	hash, err := hashstructure.Hash(object, hashstructure.FormatV2, nil)
	if err != nil {
		return "", fmt.Errorf("checksum / failed to create checksum: %w", err)
	}

	return fmt.Sprintf("%d", hash), nil
}
