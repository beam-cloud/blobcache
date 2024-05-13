# blobcache

## Overview
A very simple in-memory cache used as a content-addressed storage system. Exposes a GRPC server that can be embedded directly in a golang application. Persistence is backed by disk. The main use for blobcache is to store large blobs of content-addressed data for fast lookup by a distributed filesystem.