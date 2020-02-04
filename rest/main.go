package main

import (
	"context"
	"log"
	"net/http"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporter/metric/prometheus"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"

	"github.com/go-chi/chi"
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
		"transaction.volume",
		metric.WithKeys(key.New("status")),
	)
	transactionLabels := meter.Labels(key.String("status", "pending"))

	ctx := context.Background()

	//rest
	r := chi.NewRouter()
	r.Get("/transaction", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
		// add count
		meter.RecordBatch(
			ctx,
			transactionLabels,
			transaction.Measurement(1),
		)
	})

	http.ListenAndServe(":3000", r)
}
