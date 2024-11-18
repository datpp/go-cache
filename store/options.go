package store

import "github.com/eko/gocache/lib/v4/store"

var TAG_FORCE_REFRESH = "force_refresh"
var TAG_IGNORE_ERROR = "ignore_error"

func WithForceRefresh(v bool) store.Option {
	return func(o *store.Options) {
		if v {
			o.Tags = append(o.Tags, TAG_FORCE_REFRESH)
		}
	}
}

func WithIgnoreCacheError(v bool) store.Option {
	return func(o *store.Options) {
		if v {
			o.Tags = append(o.Tags, TAG_IGNORE_ERROR)
		}
	}
}
