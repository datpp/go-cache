# Go Cache

This document extend from https://github.com/eko/gocache/blob/master/README.md.

As GoCache have some limitation on datatype when set and not support GetOnce method, so we have customized with our own library.

## Usage
### CacheOptions 
```go
cacheStore := cache.NewStore() // your storage here
cacheManager := cache.New[YourDataType](cacheStore, WithPrefix("your_service_name"))
...
```

### GetOnce
```go
cacheStore := cache.NewStore() // your storage here
cacheManager := cache.New[YourDataType](cacheStore)
cacheData, err := cacheManager.GetOnce(ctx, cacheKey, func() (YourDataType, error) {
    // your logic here
    return data, nil
}, WithForceRefresh(true))

if err != nill {
	fmt.Errorf("eff")
}

fmt.Println(cacheData)
```

- WithForceRefresh(true): force refresh cache 

### Delete with Wildcard
Now just support if you are using our own storage - `github.com/datpp/go-cache/store/redis`
```go
import "github.com/datpp/go-cache/store/redis"

cacheStore := redis.NewRedis(redisConn, options...)
cacheManager := cache.New[YourDataType](cacheStore)
cacheManager.Delete("your-key-*")
```

### Limitation 
- Some tech issue with storage Rueidis (https://github.com/eko/gocache/tree/master/store/rueidis) so for temporary we not support this storage. 