package main

import (
	"context"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/robfig/cron/v3"
	"go-actor-model/configs"
	"go-actor-model/data"
	"go-actor-model/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"log"
	"math/rand"
	"sync"
	"time"
)

// Actor pool sizes
const (
	transactionActorPoolSize = 100
	statusCheckActorPoolSize = 100

	numTransactions = 100_000 // process 100K transactions
)

type TransactionWorkerService struct {
	transactionProcessor models.ITransactionProcessor
	statusCheckProcessor models.IStatusCheckProcessor
}

func NewTransactionWorkerService(transactionProcessor models.ITransactionProcessor, statusCheckProcessor models.IStatusCheckProcessor) *TransactionWorkerService {
	return &TransactionWorkerService{
		transactionProcessor: transactionProcessor,
		statusCheckProcessor: statusCheckProcessor,
	}
}

// StartWorkerService is a function triggers a cron job to process transactions every 5 minutes.
func (tws *TransactionWorkerService) StartWorkerService(ctx context.Context) func() {
	kron := cron.New()

	if _, err := kron.AddFunc("@every 5m", tws.processTransactions); err != nil {
		log.Fatalf("Failed to start cron job: %v", err)
	}

	// initialize the tracer
	tp := tws.initTracer()

	// start initial processing
	go tws.processTransactions()

	// start the cron job
	kron.Start()

	// Ensure all the spans are flushed before the application exits
	return func() {
		<-ctx.Done()
		kron.Stop()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}
}

// processTransactions is a method that simulates processing of transactions using actors.
func (tws *TransactionWorkerService) processTransactions() {
	// create the actor system
	system := actor.NewActorSystem(
	//actor.WithMetricProviders(),
	)
	defer system.Shutdown()

	// create pools of Transaction and StatusCheck Actors
	transactionActorPIDs := make([]*actor.PID, transactionActorPoolSize)
	statusCheckActorPIDs := make([]*actor.PID, statusCheckActorPoolSize)

	// initialize the actor pools
	for i := 0; i < transactionActorPoolSize; i++ {
		pid := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return tws.transactionProcessor
		}))
		transactionActorPIDs[i] = pid
	}
	for i := 0; i < statusCheckActorPoolSize; i++ {
		pid := system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return tws.statusCheckProcessor
		}))
		statusCheckActorPIDs[i] = pid
	}

	var wg sync.WaitGroup
	// simulate processing of transactions
	for i := 0; i < numTransactions; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// create a new transaction
			rand.New(rand.NewSource(time.Now().UnixNano()))
			amount := rand.Float64() * 100
			transaction := data.NewTransaction(amount)
			system.Root.Send(transactionActorPIDs[i%transactionActorPoolSize], transaction)

			// create a new status check
			system.Root.Send(statusCheckActorPIDs[i%statusCheckActorPoolSize], transaction.ID)
		}(i)
	}

	wg.Wait()

	// get all the transaction results
	total, processed, failed := tws.statusCheckProcessor.GetResults()
	log.Printf("Total transactions: %d, Processed: %d, Failed: %d\n", total, processed, failed)
}

func (*TransactionWorkerService) initTracer() *sdktrace.TracerProvider {
	ctx := context.Background()

	// configure zipkin exporter
	zexp, err := zipkin.New("http://localhost:9411/api/v2/spans")
	if err != nil {
		log.Fatalf("failed to initialize Zipkin exporter: %v", err)
	}

	// Configure the OTLP exporter to send data to the OpenTelemetry Collector
	exporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			otlptracegrpc.WithEndpoint("localhost:4317"),
			otlptracegrpc.WithInsecure(), // Use WithTLSCredentials if your collector is setup with TLS
		),
	)
	if err != nil {
		log.Fatalf("failed to initialize OTLP exporter: %v", err)
	}

	// Create a tracer provider with the exporter and resource detector
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithBatcher(zexp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("GoActorModel"),
		)),
	)

	// Set the global trace provider
	otel.SetTracerProvider(tp)

	// Set the global propagator to tracecontext
	configs.OtelTracer = tp.Tracer("GoActorModel")

	return tp
}
