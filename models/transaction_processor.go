package models

import (
	"context"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"go-actor-model/configs"
	"go-actor-model/data"
	"math/rand"
)

type ITransactionProcessor interface {
	Receive(actor.Context)
}

type TransactionActor struct {
	storage configs.IStorage
	ITransactionProcessor
}

func NewTransactionProcessor(storage configs.IStorage) ITransactionProcessor {
	return &TransactionActor{storage: storage}
}

func (a *TransactionActor) Receive(ctx actor.Context) {
	// start a new span
	if _, span := configs.OtelTracer.Start(context.Background(), "TransactionActor.Receive"); span.IsRecording() {
		defer span.End()
	}

	switch transaction := ctx.Message().(type) {
	case *actor.Started:
	case *actor.Stopping:
		fmt.Printf("Stopping, actor is about shut down: %v\n", ctx.Self().Id)
	case *actor.Stopped:
		fmt.Printf("Stopped, actor and its children are stopped: %v\n", ctx.Self().Id)
	case *actor.Restarting:
		fmt.Printf("Restarting, actor is about restart: %v\n", ctx.Self().Id)
	case *data.Transaction:
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
	}

	// ignore all other inputs
}
