package main

import (
	"log"

	blobcache "github.com/beam-cloud/blobcache/pkg"
	"github.com/checkpoint-restore/go-criu/v7"
)

func main() {
	s, err := blobcache.NewCacheService(blobcache.BlobCacheServiceConfig.PersistencePath, blobcache.BlobCacheServiceConfig.CacheSize, blobcache.BlobCacheServiceConfig.PageSize)
	if err != nil {
		log.Fatal(err)
	}

	c := criu.MakeCriu()
	result, err := c.IsCriuAtLeast(31100)
	if err != nil {
		log.Println("err: ", err)
	}

	log.Println("criu: ", result)

	s.StartServer(blobcache.BlobCacheServiceConfig.Address)
}
