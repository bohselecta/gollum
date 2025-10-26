package impl

import (
	"math"
	"strings"

	"github.com/haydenlabs/gollum/engine"
)

const Kdim = 64 // embedding dimension
const tinyVocabSize = 16

var tinyVocab = []string{" llama", " on", " the", " high", " plain", ".", " gentle", ",", " wind", " hums", " softly", " high", " plain", " wind", " gentle", " softly"}

// Simple projection matrix (random init for demo)
var projW []float32

func init() {
	projW = make([]float32, Kdim*tinyVocabSize)
	for i := range projW {
		projW[i] = float32(math.Sin(float64(i)) * 0.1)
	}
}

type stubOps struct{}

func NewMetalOps() *stubOps {
	return &stubOps{}
}

func (m *stubOps) Prefill(batch *engine.Batch) error {
	// Toy implementation: just simulate some processing time
	return nil
}

func (m *stubOps) Decode(step *engine.Step) (int, error) {
	// Toy implementation: simulate processing time
	return step.BatchSize, nil
}

// PredictNext with persistent KV: per item, create KV handle if 0, append (K=V=Q), attend over full sequence, project.
func (m *stubOps) PredictNext(ctxs []engine.DecodeCtx) []string {
	M := len(ctxs)
	if M == 0 {
		return nil
	}
	ctxVec := make([]float32, M*Kdim)
	// Build Q from prompt, ensure KV handle exists, append, attend
	for i := 0; i < M; i++ {
		p := ctxs[i].Prompt
		q := promptVec(p, Kdim)
		h := ctxs[i].KVHandle
		if h == 0 {
			h = i + 1 // Simple stub handle
			// Persist handle back into bound block when scheduler sees it next tick:
			// (scheduler reads rs.kv.KVHandle to pass it here; we set it via side-effect below)
		}
		// Append token K/V (demo: K=V=Q) - stub implementation
		// In a real implementation, this would call metal.KVAppend(h, q, q, Kdim)
		
		// Simple attention: just use the query vector for now
		copy(ctxVec[i*Kdim:(i+1)*Kdim], q)
		
		// Store handle into the KV block if available (idempotent). We expose a tiny hook:
		// In this minimal patch, scheduler carries KVHandle via DecodeCtx and expects us to reuse 'h' next call.
		ctxs[i].KVHandle = h
	}
	V := len(tinyVocab)
	// Project ctx -> logits
	logits := cpuMatMul(ctxVec, M, projW, Kdim, V)
	res := make([]string, M)
	for i := 0; i < M; i++ {
		best := 0
		bestv := logits[i*V]
		for j := 1; j < V; j++ {
			if logits[i*V+j] > bestv {
				best = j
				bestv = logits[i*V+j]
			}
		}
		res[i] = tinyVocab[best]
		if strings.HasSuffix(ctxs[i].Prompt, res[i]) {
			res[i] = "."
		}
	}
	return res
}

func promptVec(prompt string, dim int) []float32 {
	// Simple hash-based embedding
	vec := make([]float32, dim)
	hash := 0
	for _, c := range prompt {
		hash = hash*31 + int(c)
	}
	for i := 0; i < dim; i++ {
		vec[i] = float32(math.Sin(float64(hash+i)) * 0.1)
	}
	return vec
}

func cpuMatMul(A []float32, M int, B []float32, K, N int) []float32 {
	C := make([]float32, M*N)
	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			sum := float32(0)
			for k := 0; k < K; k++ {
				sum += A[i*K+k] * B[k*N+j]
			}
			C[i*N+j] = sum
		}
	}
	return C
}
