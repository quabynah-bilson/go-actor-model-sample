package configs

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

// NewRedisClient is a function that creates and returns a new Redis client instance.
// It takes a context as an argument which is used for the Redis Ping operation.
// The function creates a new Redis client with the specified options (address, password, and DB).
// It then pings the Redis server to check if the connection is successful.
// If the ping operation fails, it logs a fatal error and the program will exit.
// If the ping operation is successful, it returns the Redis client instance.
//
// Parameters:
//
//	ctx : context.Context : The context to use for the Redis Ping operation.
//
// Returns:
//
//	*redis.Client : A pointer to the new Redis client instance.
func NewRedisClient(ctx context.Context) *redis.Client {

	// Create a new Redis client with the specified options.
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // The address of the Redis server.
		Password: "",               // The password for the Redis server (empty for testing purposes).
		DB:       0,                // The DB number to use.
	})

	// Ping the Redis server to check if the connection is successful.
	if _, err := client.Ping(ctx).Result(); err != nil {
		// If the ping operation fails, log a fatal error and the program will exit.
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Return the Redis client instance.
	log.Println("Redis client connected successfully.")
	return client
}

func NewCockroachClient(ctx context.Context) *pgxpool.Pool {
	// Create a new CockroachDB connection pool with the specified options.

	return nil
}
