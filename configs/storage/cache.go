package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

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
