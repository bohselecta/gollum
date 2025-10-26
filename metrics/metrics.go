package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// Duration histograms
	TTFTMs = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "gollum_ttft_ms",
		Help:    "Time to first token (ms)",
		Buckets: []float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1000},
	}, []string{"model"})

	TPOTMs = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "gollum_tpot_ms",
		Help:    "Time to produce all tokens (ms)",
		Buckets: []float64{4, 8, 16, 32, 64, 128, 256, 512, 1000, 2000, 5000},
	}, []string{"model"})

	BatchSize = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "gollum_batch_size",
		Help:    "Decode batch size",
		Buckets: []float64{1, 2, 4, 8, 16, 32, 64},
	}, []string{"model"})

	// Cache events: hits/misses by type and model
	CacheEvents = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "gollum_cache_events_total",
		Help: "Cache events (hits/misses) by type and model",
	}, []string{"cache_type", "hit_kind", "model"}) // cache_type: prompt|prefix, hit_kind: hit|miss

	// KV events: alloc/pin/unpin/evict by model
	KVEvents = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "gollum_kv_events_total",
		Help: "KV pager events by action and model",
	}, []string{"action", "model"}) // action: alloc|pin|unpin|evict

	DecodeSteps = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "gollum_decode_steps_total",
		Help: "Decode steps executed",
	}, []string{"model"})
)

func MustRegister() {
	prometheus.MustRegister(
		TTFTMs, TPOTMs, BatchSize,
		CacheEvents, KVEvents, DecodeSteps,
	)
}
