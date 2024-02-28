package main

import (
	"log"

	blobcache "github.com/beam-cloud/blobcache/pkg"
)

func main() {
	s, err := blobcache.NewCacheService(blobcache.BlobCacheServiceConfig.PersistencePath, blobcache.BlobCacheServiceConfig.CacheSize, blobcache.BlobCacheServiceConfig.PageSize)
	if err != nil {
		log.Fatal(err)
	}

	s.StartServer(blobcache.BlobCacheServiceConfig.Address)
}
