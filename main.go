package main

import (
	cache "cache/pkg"
	"log"
)

func main() {
	s, err := cache.NewCacheService(cache.CacheServiceConfig.PersistencePath, cache.CacheServiceConfig.CacheSize, cache.CacheServiceConfig.PageSize)
	if err != nil {
		log.Fatal(err)
	}

	s.StartServer(cache.CacheServiceConfig.Address)
}
