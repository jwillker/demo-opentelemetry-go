package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
)

func initMeter() *push.Controller {
	pusher, hf, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}
	http.HandleFunc("/metrics", hf)
	go func() {
		_ = http.ListenAndServe(":9080", nil)
	}()

	return pusher
}

func main() {
	defer initMeter().Stop()

	meter := global.Meter("transactions")

	transactions := metric.Must(meter).NewInt64Counter(
		"transactions.total",
		metric.WithKeys(key.New("status")),
	)
	transactionsLabels := meter.Labels(key.String("status", "pending"))

	ctx := context.Background()

	meter.RecordBatch(
		ctx,
		transactionsLabels,
		transactions.Measurement(1),
	)
	time.Sleep(5 * time.Second)

	meter.RecordBatch(
		ctx,
		meter.Labels(key.String("status", "completed")),
		transactions.Measurement(1),
	)
	time.Sleep(5 * time.Second)

	meter.RecordBatch(
		ctx,
		transactionsLabels,
		transactions.Measurement(1),
	)

	time.Sleep(1000 * time.Second)
}
