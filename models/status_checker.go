package models

import (
	"github.com/asynkron/protoactor-go/actor"
	"go-actor-model/configs/storage"
	"go-actor-model/data"
	"log"
)

type IStatusCheckProcessor interface {
	Receive(actor.Context)
	GetResults() (int32, int32, int32)
}

type StatusCheckActor struct {
	storage storage.IStorage
	IStatusCheckProcessor
}

func NewStatusCheckProcessor(storage storage.IStorage) IStatusCheckProcessor {
	return &StatusCheckActor{storage: storage}
}

func (s *StatusCheckActor) Receive(ctx actor.Context) {
	// start a new span
	//if _, span := configs.OtelTracer.Start(context.Background(), "StatusCheckActor.Receive"); span.IsRecording() {
	//	defer span.End()
	//}

	switch transactionID := ctx.Message().(type) {
	case *actor.Started:
	case *actor.Stopping:
		log.Printf("Stopping, actor is about shut down: %v\n", ctx.Self().Id)
	case *actor.Stopped:
		log.Printf("Stopped, actor and its children are stopped: %v\n", ctx.Self().Id)
	case *actor.Restarting:
		log.Printf("Restarting, actor is about restart: %v\n", ctx.Self().Id)
	case string:
		transactionStr, err := s.storage.Get(transactionID)
		if err != nil {
			log.Printf("Status for transaction ID: %s is unknown\n", transactionID)
			return
		}

		// parse transaction JSON string to struct
		transaction := data.Transaction{}
		if err = transaction.Parse(*transactionStr); err != nil {
			transaction.Status = data.Failed
		}
		transaction.CheckReadStatus = data.ReadProcessed

		if err := s.storage.Save(transaction.ID, transaction.Stringify()); err != nil {
			log.Printf("Failed to update transaction status in storage %v\n", err)
			return
		}

		log.Printf("Status check for transaction %s completed with status: %v\n", transaction.ID, transaction.Status)
	}

	// ignore all other inputs
}

func (s *StatusCheckActor) GetResults() (int32, int32, int32) {
	return s.storage.GetStats()
}
