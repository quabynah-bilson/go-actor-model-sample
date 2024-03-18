package data

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// TransactionStatus is a custom type for the status of a transaction.
type TransactionStatus string
type ReadStatus string

const (
	Pending   TransactionStatus = "Pending"
	Processed                   = "Processed"
	Failed                      = "Failed"

	ReadPending   ReadStatus = "Pending"
	ReadProcessed            = "Processed"
)

// Transaction is a struct that represents a transaction.
type Transaction struct {
	ID              string            `json:"id"`     // The unique identifier for the transaction.
	Amount          float64           `json:"amount"` // The amount of the transaction.
	Status          TransactionStatus `json:"status"` // The status of the transaction (Pending, Completed, or Failed).
	CheckReadStatus ReadStatus        `json:"read_status"`
	CreatedAt       time.Time         `json:"created_at"`
}

// NewTransaction is a function that creates
// returns a new Transaction instance with the specified amount.
func NewTransaction(amount float64) *Transaction {
	return &Transaction{
		ID:              generateTransactionID(),
		Amount:          amount,
		Status:          Pending,
		CheckReadStatus: ReadPending,
	}
}

// UpdateStatus is a method that updates the status of the transaction.
func (t *Transaction) UpdateStatus(status TransactionStatus) {
	t.Status = status
}

func (t *Transaction) Stringify() string {
	data, err := json.Marshal(t)
	if err != nil {
		log.Printf("Unable to marshal transaction data with ID: %s", t.ID)
		return ""
	}

	return string(data)
}

func (t *Transaction) Parse(data string) error {
	return json.Unmarshal([]byte(data), t)
}

// generateTransactionID is a function that generates a new unique transaction ID.
func generateTransactionID() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("T%05d", rand.Intn(10000000))
}
