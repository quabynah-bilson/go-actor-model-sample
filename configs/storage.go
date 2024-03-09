package configs

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var (
	errDataReadFailed   = errors.New("failed to read data from storage")
	errDataWriteFailed  = errors.New("failed to write data to storage")
	errDataDeleteFailed = errors.New("failed to delete data from storage")

	cacheExpiration = 2 * time.Minute
)

// IStorage is an interface that defines the methods for a storage system.
// It includes methods for saving, retrieving, and deleting key-value pairs.
type IStorage interface {
	// Save is a method that takes a key and a value as strings.
	// It saves the key-value pair to the storage.
	// It returns an error if the save operation fails.
	Save(string, string) error

	// Get is a method that takes a key as a string.
	// It retrieves the value associated with the key from the storage.
	// It returns a pointer to the string value and an error.
	// If the key does not exist, the string pointer will be nil.
	// If the retrieval operation fails, it returns an error.
	Get(string) (*string, error)

	// Delete is a method that takes a key as a string.
	// It removes the key-value pair associated with the key from the storage.
	// It returns an error if the delete operation fails.
	Delete(string) error

	// GetStats is a method that returns the number of saved entries, processed entries, and failed entries.
	GetStats() (int32, int32, int32)
}

// CacheStorage is a struct that represents a cache storage system.
type CacheStorage struct {
	cache *redis.Client
	ctx   context.Context
	IStorage
}

func NewCacheStorage(ctx context.Context, cache *redis.Client) IStorage {
	return &CacheStorage{cache: cache, ctx: ctx}
}

func (cs *CacheStorage) Save(key, value string) error {
	if err := cs.cache.Set(cs.ctx, key, value, cacheExpiration).Err(); err != nil {
		log.Printf("Failed to save data to Redis: %v", err)
		return errDataWriteFailed
	}

	return nil
}

func (cs *CacheStorage) Get(key string) (*string, error) {
	data, err := cs.cache.Get(cs.ctx, key).Result()
	if err != nil {
		log.Printf("Failed to get data from Redis %s: %v", key, err)
		return nil, errDataReadFailed
	}
	return &data, nil
}

func (cs *CacheStorage) Delete(key string) error {
	if err := cs.cache.Del(cs.ctx, key).Err(); err != nil {
		log.Printf("Failed to delete data from Redis %s: %v", key, err)
		return errDataDeleteFailed
	}
	return nil
}

func (cs *CacheStorage) GetStats() (int32, int32, int32) {
	savedTransactionsCount, processedTransactionsCount, failedTransactionsCount := int32(0), int32(0), int32(0)
	if numSaved, err := cs.cache.DBSize(cs.ctx).Result(); err == nil {
		savedTransactionsCount = int32(numSaved)
	}

	pipeline := cs.cache.Pipeline()
	var processedCmd, failedCmd *redis.StringCmd

	// get the number of processed transactions
	processedCmd = pipeline.Get(cs.ctx, "processed")

	// get the number of failed transactions
	failedCmd = pipeline.Get(cs.ctx, "failed")

	if _, err := pipeline.Exec(cs.ctx); err == nil {
		if processed, err := processedCmd.Int64(); err == nil {
			processedTransactionsCount = int32(processed)
		}
		if failed, err := failedCmd.Int64(); err == nil {
			failedTransactionsCount = int32(failed)
		}
	}

	return savedTransactionsCount, processedTransactionsCount, failedTransactionsCount
}
