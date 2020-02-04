package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporter/metric/prometheus"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
)

func initMeter() *push.Controller {
	pusher, hf, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}
	http.HandleFunc("/", hf)
	go func() {
		_ = http.ListenAndServe(":9080", nil)
	}()

	return pusher
}

func main() {
	defer initMeter().Stop()

	meter := global.MeterProvider().Meter("ex.com/basic")

	transaction := meter.NewInt64Counter(
		"transaction.total",
		metric.WithKeys(key.New("status")),
	)
	transactionLabels := meter.Labels(key.String("status", "pending"))

	ctx := context.Background()

	meter.RecordBatch(
		ctx,
		transactionLabels,
		transaction.Measurement(1),
	)
	time.Sleep(5 * time.Second)

	meter.RecordBatch(
		ctx,
		meter.Labels(key.String("status", "completed")),
		transaction.Measurement(1),
	)
	time.Sleep(5 * time.Second)

	meter.RecordBatch(
		ctx,
		transactionLabels,
		transaction.Measurement(1),
	)
	time.Sleep(5 * time.Second)
}
