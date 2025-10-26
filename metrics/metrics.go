package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	TTFTMs = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "gollum_ttft_ms", Help: "Time to first token (ms)",
		Buckets: []float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1000},
	})
	TPOTMs = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "gollum_tpot_ms", Help: "Time to produce all tokens (ms)",
		Buckets: []float64{4, 8, 16, 32, 64, 128, 256, 512, 1000, 2000, 5000},
	})

	BatchSize = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "gollum_batch_size", Help: "Decode batch size",
		Buckets: []float64{1, 2, 4, 8, 16, 32, 64},
	})

	PromptCacheHit  = prometheus.NewCounter(prometheus.CounterOpts{Name: "gollum_promptcache_hit_total"})
	PromptCacheMiss = prometheus.NewCounter(prometheus.CounterOpts{Name: "gollum_promptcache_miss_total"})

	PrefixHit  = prometheus.NewCounter(prometheus.CounterOpts{Name: "gollum_prefix_hit_total"})
	PrefixMiss = prometheus.NewCounter(prometheus.CounterOpts{Name: "gollum_prefix_miss_total"})

	KVAlloc     = prometheus.NewCounter(prometheus.CounterOpts{Name: "gollum_kv_alloc_total"})
	KVEvict     = prometheus.NewCounter(prometheus.CounterOpts{Name: "gollum_kv_evict_total"})
	KVPin       = prometheus.NewCounter(prometheus.CounterOpts{Name: "gollum_kv_pin_total"})
	KVUnpin     = prometheus.NewCounter(prometheus.CounterOpts{Name: "gollum_kv_unpin_total"})
	DecodeSteps = prometheus.NewCounter(prometheus.CounterOpts{Name: "gollum_decode_steps_total"})
)

func MustRegister() {
	prometheus.MustRegister(
		TTFTMs, TPOTMs, BatchSize,
		PromptCacheHit, PromptCacheMiss,
		PrefixHit, PrefixMiss,
		KVAlloc, KVEvict, KVPin, KVUnpin,
		DecodeSteps,
	)
}
