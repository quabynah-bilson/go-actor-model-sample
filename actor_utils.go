package main

import "github.com/asynkron/protoactor-go/actor"

// Actor pool sizes
const (
	transactionActorPoolSize = 10
	statusCheckActorPoolSize = 10

	numTransactions = 100_000 // process 100K transactions
)

var (
	transactionActorPIDs []*actor.PID
	statusCheckActorPIDs []*actor.PID
)
