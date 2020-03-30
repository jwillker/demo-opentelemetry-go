package main

import (
	"context"
	"log"
	"net/http"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"

	"github.com/go-chi/chi"
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

	meter := global.Meter("transactionss")

	transactions := metric.Must(meter).NewInt64Counter(
		"transactions.volume",
		metric.WithKeys(key.New("status")),
	)
	transactionsLabels := meter.Labels(key.String("status", "performed"))

	ctx := context.Background()

	//rest
	r := chi.NewRouter()
	r.Get("/transactions", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
		// add count
		meter.RecordBatch(
			ctx,
			transactionsLabels,
			transactions.Measurement(1),
		)
	})

	http.ListenAndServe(":3000", r)
}
