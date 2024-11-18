package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptionsCachePrefixValue(t *testing.T) {
	cacheOptions := ApplyOptions(WithPrefix("my-prefix"))

	// When - Then
	assert.Equal(t, "my-prefix", cacheOptions.CachePrefix)
}

func Test_applyOption(t *testing.T) {
	// Given
	options := &Options{}

	// When
	options = ApplyOptions(WithPrefix("my-prefix"))

	// Then
	assert.Equal(t, "my-prefix", options.CachePrefix)
}

func Test_applyOptionsWithDefault(t *testing.T) {
	// Given
	defaultOptions := &Options{
		CachePrefix: "my-default-prefix",
	}

	// When
	options := ApplyOptionsWithDefault(defaultOptions, WithPrefix("my-prefix"))

	// Then
	assert.Equal(t, "my-prefix", options.CachePrefix)
}
