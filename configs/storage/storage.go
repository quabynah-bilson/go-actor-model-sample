package storage

import (
	"errors"
	"time"
)

var (
	errDataReadFailed   = errors.New("failed to read data from storage")
	errDataWriteFailed  = errors.New("failed to write data to storage")
	errDataDeleteFailed = errors.New("failed to delete data from storage")

	cacheExpiration = 30 * time.Second
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
