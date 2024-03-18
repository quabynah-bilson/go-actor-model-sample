package main

import (
	"context"
	"fmt"
	"go-actor-model/configs"
	"go-actor-model/configs/storage"
	"go-actor-model/models"
)

// steps to start the worker
// 1. run `make start-services` to start the Redis server
// 2. run the go application using your IDE or `go run main.go`
func main() {
	fmt.Println("Starting the Go Actor Model Demo Application...")

	// create storage instance
	ctx, cancel := context.WithCancel(context.Background())

	// create storage instance
	rdc := configs.NewRedisClient(ctx)
	cacheStorage := storage.NewCacheStorage(ctx, rdc)

	// create processors
	transactionProcessor := models.NewTransactionProcessor(cacheStorage)
	statusCheckProcessor := models.NewStatusCheckProcessor(cacheStorage)

	// create worker service
	tws := NewTransactionWorkerService(transactionProcessor, statusCheckProcessor)
	cleanup := tws.StartWorkerService(ctx)
	defer func() {
		cleanup() // cleanup the worker service (stop the cron job and flush the spans)
		cancel()  // cancel the context
	}()

	// keep the application running
	select {}
}
