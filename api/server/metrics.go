// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package server

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/MetalBlockchain/metalgo/utils/wrappers"
)

type metrics struct {
	numProcessing *prometheus.GaugeVec
	numCalls      *prometheus.CounterVec
	totalDuration *prometheus.GaugeVec
}

func newMetrics(namespace string, registerer prometheus.Registerer) (*metrics, error) {
	m := &metrics{
		numProcessing: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "calls_processing",
				Help:      "The number of calls this API is currently processing",
			},
			[]string{"base"},
		),
		numCalls: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "calls",
				Help:      "The number of calls this API has processed",
			},
			[]string{"base"},
		),
		totalDuration: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "calls_duration",
				Help:      "The total amount of time, in nanoseconds, spent handling API calls",
			},
			[]string{"base"},
		),
	}

	errs := wrappers.Errs{}
	errs.Add(
		registerer.Register(m.numProcessing),
		registerer.Register(m.numCalls),
		registerer.Register(m.totalDuration),
	)
	return m, errs.Err
}

func (m *metrics) wrapHandler(chainName string, handler http.Handler) http.Handler {
	numProcessing := m.numProcessing.WithLabelValues(chainName)
	numCalls := m.numCalls.WithLabelValues(chainName)
	totalDuration := m.totalDuration.WithLabelValues(chainName)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		numProcessing.Inc()

		defer func() {
			numProcessing.Dec()
			numCalls.Inc()
			totalDuration.Add(float64(time.Since(startTime)))
		}()

		handler.ServeHTTP(w, r)
	})
}
