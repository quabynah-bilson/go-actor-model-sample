package models

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"go-actor-model/configs/storage"
	"go-actor-model/data"
	"log"
	"math/rand"
	"time"
)

type ITransactionProcessor interface {
	Receive(actor.Context)
	SetStatusCheckActorPID(*actor.PID)
}

type TransactionActor struct {
	storage                storage.IStorage
	statusCheckProcessorID *actor.PID
	ITransactionProcessor
}

func NewTransactionProcessor(storage storage.IStorage) ITransactionProcessor {
	return &TransactionActor{storage: storage}
}

func (a *TransactionActor) SetStatusCheckActorPID(pid *actor.PID) {
	a.statusCheckProcessorID = pid
}

func (a *TransactionActor) Receive(ctx actor.Context) {
	// start a new span
	//if _, span := configs.OtelTracer.Start(context.Background(), "TransactionActor.Receive"); span.IsRecording() {
	//	defer span.End()
	//}

	switch transaction := ctx.Message().(type) {
	case *actor.Started:
	case *actor.Stopping:
		fmt.Printf("Stopping, actor is about shut down: %v\n", ctx.Self().Id)
	case *actor.Stopped:
		fmt.Printf("Stopped, actor and its children are stopped: %v\n", ctx.Self().Id)
	case *actor.Restarting:
		fmt.Printf("Restarting, actor is about restart: %v\n", ctx.Self().Id)
	case *data.Transaction:
		// Update transaction with current time
		transaction.CreatedAt = time.Now()

		rand.New(rand.NewSource(time.Now().UnixNano()))
		// Simulating random processing errors
		if rand.Float32() < 0.05 {
			transaction.UpdateStatus(data.Failed)
		} else {
			transaction.UpdateStatus(data.Processed)
		}

		// Store transaction as processed
		if err := a.storage.Save(transaction.ID, transaction.Stringify()); err != nil {
			fmt.Println(err)
			return
		}

		// send response to the status checker actor
		if a.statusCheckProcessorID != nil {
			ctx.Send(a.statusCheckProcessorID, transaction.ID)
		}
		log.Printf("Transaction %s processed after %v\n", transaction.ID, time.Since(transaction.CreatedAt))
	}
}
