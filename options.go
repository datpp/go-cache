package cache

type Option func(o *Options)

type Options struct {
	CachePrefix string
}

func (o *Options) IsEmpty() bool {
	return o.CachePrefix == ""
}

func ApplyOptions(opts ...Option) *Options {
	o := &Options{}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func ApplyOptionsWithDefault(defaultOptions *Options, opts ...Option) *Options {
	returnedOptions := &Options{}
	*returnedOptions = *defaultOptions

	for _, opt := range opts {
		opt(returnedOptions)
	}

	return returnedOptions
}

// WithPrefix allows setting default prefix for all cache keys.
// it should be set as Name of Service in most cases. (to follow naming convention)
func WithPrefix(prefix string) Option {
	return func(o *Options) {
		o.CachePrefix = prefix
	}
}
