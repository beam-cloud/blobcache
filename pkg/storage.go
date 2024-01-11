package blobcache

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/dgraph-io/ristretto"
)

type ContentAddressableStorage struct {
	dir       string
	inMemory  *ristretto.Cache
	diskCache string
	chunkSize int64
	mu        sync.RWMutex
}

func NewContentAddressableStorage(dir string, size int64, chunkSize int64) (*ContentAddressableStorage, error) {
	if size <= 0 || chunkSize <= 0 {
		return nil, errors.New("invalid cache configuration")
	}

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: size * 100,
		MaxCost:     size * 1000,
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}

	return &ContentAddressableStorage{
		dir:       dir,
		inMemory:  cache,
		diskCache: dir,
		chunkSize: chunkSize,
	}, nil
}

func (cas *ContentAddressableStorage) Add(content []byte) (string, error) {
	hash := sha256.Sum256(content)
	hashStr := hex.EncodeToString(hash[:])

	cas.mu.Lock()
	defer cas.mu.Unlock()

	dirPath := filepath.Join(cas.diskCache, hashStr)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Break content into chunks and store on disk
	for offset := int64(0); offset < int64(len(content)); offset += cas.chunkSize {
		chunkIdx := offset / cas.chunkSize
		end := offset + cas.chunkSize
		if end > int64(len(content)) {
			end = int64(len(content))
		}

		chunk := content[offset:end]
		filePath := filepath.Join(dirPath, fmt.Sprintf("%d", chunkIdx))
		if err := os.WriteFile(filePath, chunk, 0644); err != nil {
			return "", fmt.Errorf("failed to write file: %w", err)
		}

		chunkKey := fmt.Sprintf("%s-%d", hashStr, chunkIdx)
		cas.inMemory.Set(chunkKey, chunk, int64(len(chunk)))
	}

	return hashStr, nil
}

func (cas *ContentAddressableStorage) Get(hash string, offset, length int64) ([]byte, error) {
	cas.mu.RLock()
	defer cas.mu.RUnlock()

	var result []byte
	remainingLength := length

	o := offset
	for remainingLength > 0 {
		var chunkBytes []byte
		chunkIdx := o / cas.chunkSize
		chunkKey := fmt.Sprintf("%s-%d", hash, chunkIdx)

		// Check in-memory cache first
		if chunk, found := cas.inMemory.Get(chunkKey); found {
			chunkBytes = chunk.([]byte)
		} else {
			// Check disk cache
			chunkPath := filepath.Join(cas.diskCache, hash, fmt.Sprintf("%d", chunkIdx))
			chunkBytesRead, err := os.ReadFile(chunkPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read from disk cache: %w", err)
			}
			chunkBytes = chunkBytesRead
			cas.inMemory.Set(chunkKey, chunkBytes, int64(len(chunkBytes)))
		}

		start := o % cas.chunkSize
		chunkRemaining := int64(len(chunkBytes)) - start
		if chunkRemaining <= 0 {
			// No more data in this chunk, break out of the loop
			break
		}

		readLength := min(remainingLength, chunkRemaining)
		end := start + readLength

		// Validate start and end positions
		if start < 0 || end <= start || end > int64(len(chunkBytes)) {
			return nil, fmt.Errorf("invalid chunk boundaries: start %d, end %d, chunk size %d", start, end, len(chunkBytes))
		}

		// Read only the required portion of the chunk
		requiredChunkBytes := chunkBytes[start:end]

		// Append the required bytes to the result
		result = append(result, requiredChunkBytes...)

		// Update the remaining length and current offset for the next iteration
		remainingLength -= readLength
		o += readLength
	}

	return result, nil
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
