package service

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	prometheus.MustRegister(rpcRequestCount)
	prometheus.MustRegister(rpcRequestDuration)
}

var (
	rpcRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rpc_request_count",
			Help: "Number of rpc request count",
		},
		[]string{"method"},
	)

	rpcRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "rpc_request_duration_milliseconds",
			Help:    "rpc request duration distribution",
			Buckets: []float64{10, 20, 30, 40, 50, 100, 150, 200, 250, 300, 350, 400, 450, 500},
		},
		[]string{"method", "date"},
	)
)

func PrometheusServer(port int) {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}
