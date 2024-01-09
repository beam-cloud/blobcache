package cache

type cacheServiceConfig struct {
	Address         string
	PersistencePath string
	PageSize        int64
	CacheSize       int64
}

var CacheServiceConfig cacheServiceConfig = cacheServiceConfig{
	Address:         "0.0.0.0:2049",
	PersistencePath: "/cache",
	PageSize:        1 << 21, // 2Mb
	CacheSize:       100000000,
}
