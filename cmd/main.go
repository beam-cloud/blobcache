package main

import (
	"log"

	blobcache "github.com/beam-cloud/blobcache/pkg"
	"github.com/checkpoint-restore/go-criu/v7"
	"github.com/checkpoint-restore/go-criu/v7/rpc"
)

func main() {
	s, err := blobcache.NewCacheService(blobcache.BlobCacheServiceConfig.PersistencePath, blobcache.BlobCacheServiceConfig.CacheSize, blobcache.BlobCacheServiceConfig.PageSize)
	if err != nil {
		log.Fatal(err)
	}

	c := criu.MakeCriu()
	_, err = c.IsCriuAtLeast(31100)
	if err != nil {
		log.Println("err??: ", err)
	}

	var imagesFd int32 = -1
	var imagesDir string = "/tmp/test"
	var shellJob bool = true
	err = c.Dump(&rpc.CriuOpts{
		ImagesDirFd: &imagesFd,
		ImagesDir:   &imagesDir,
		ShellJob:    &shellJob,
	}, criu.NoNotify{})

	log.Println("err dump: ", err)

	s.StartServer(blobcache.BlobCacheServiceConfig.Address)
}
