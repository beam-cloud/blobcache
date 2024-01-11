package blobcache

type blobCacheServiceConfig struct {
	Address         string
	PersistencePath string
	PageSize        int64
	CacheSize       int64
}

var BlobCacheServiceConfig blobCacheServiceConfig = blobCacheServiceConfig{
	Address:         "0.0.0.0:2049",
	PersistencePath: "/cache",
	PageSize:        1 << 21, // 2Mb
	CacheSize:       100000000,
}
