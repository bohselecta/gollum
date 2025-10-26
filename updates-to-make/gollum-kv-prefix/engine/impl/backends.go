package impl

import (
	"math"
	"strings"

	"github.com/haydenlabs/gollum/engine"
	"github.com/haydenlabs/gollum/kernels/metal"
)

// Vocab for tiny predictor (token pieces)
var tinyVocab = []string{" llama"," on"," the"," high"," plain","."," gentle",","," wind"," hums"}

// Fixed weights KxV for projection (deterministic small numbers)
const K = 8
var projW = []float32{
	// K rows stacked row-major (K x V)
	0.3, 0.1, 0.2, 0.1, 0.2, -0.1, 0.2, -0.2, 0.1, -0.1,
	0.1, 0.2, 0.1, 0.3, 0.0,  0.2, 0.1,  0.1, 0.2,  0.0,
	-0.2,0.1, 0.2, 0.2, -0.1, 0.3, -0.1, 0.2, -0.2, 0.1,
	0.2, 0.2, -0.1,0.1, 0.3,  0.1, 0.2,  0.0, 0.1,  0.2,
	0.0, 0.2, 0.3, -0.1,0.1,  0.2, 0.1,  0.2, 0.1,  0.2,
	0.1, 0.0, 0.2, 0.2, 0.2,  0.1, 0.3, -0.1,0.1,  0.1,
	0.2, 0.1, 0.0, 0.3, -0.1, 0.2, 0.1,  0.2, -0.1, 0.3,
	0.3, -0.1,0.2, 0.1, 0.2,  0.0, 0.1,  0.2, 0.1,  0.2,
}

// Hash prompt into an 8-dim float vector (stable)
func promptVec(p string) []float32 {
	v := make([]float32, K)
	for i, r := range p {
		idx := i % K
		v[idx] += float32((int(r)%97)-48) / 50.0
	}
	// tanh-like clamp and L2 normalize
	var sum float32
	for i:=0;i<K;i++{ if v[i]>1{v[i]=1}; if v[i]<-1{v[i]=-1}; sum += v[i]*v[i] }
	if sum>0{ inv := 1/float32(math.Sqrt(float64(sum))); for i:=0;i<K;i++{ v[i]*=inv } }
	return v
}

// CPU fallback matmul: (M x K) * (K x V) -> (M x V)
func cpuMatMul(A []float32, M int, B []float32, K int, V int) []float32 {
	C := make([]float32, M*V)
	for m:=0;m<M;m++{
		for v:=0; v<V; v++{
			var s float32
			for k:=0; k<K; k++{
				s += A[m*K+k] * B[k*V+v]
			}
			C[m*V+v] = s
		}
	}
	return C
}

type metalOps struct{}

func (m *metalOps) Prefill(b *engine.Batch) error { _ = metal.Prefill(len(b.Prompts)); return nil }
func (m *metalOps) Decode(s *engine.Step) (int,error){ return metal.DecodeStep(s.BatchSize, s.SeqLen), nil }

func (m *metalOps) PredictNext(prompts []string) []string {
	M := len(prompts)
	if M == 0 { return nil }
	// Build A (M x K) row-major
	A := make([]float32, M*K)
	for i,p := range prompts {
		v := promptVec(p)
		copy(A[i*K:(i+1)*K], v)
	}
	// Try Metal
	var logits []float32
	if out, ok := metal.MatMul(A, projW, M, len(tinyVocab), K); ok {
		logits = out
	} else {
		logits = cpuMatMul(A, M, projW, K, len(tinyVocab))
	}
	// Argmax per row -> token string
	res := make([]string, M)
	V := len(tinyVocab)
	for i:=0;i<M;i++{
		best := 0; bestv := logits[i*V]
		for j:=1;j<V;j++{
			if logits[i*V+j] > bestv { best=j; bestv=logits[i*V+j] }
		}
		res[i] = tinyVocab[best]
		// Tiny heuristic: avoid repeating last token if prompt already ends with it
		if strings.HasSuffix(prompts[i], res[i]) { res[i] = "." }
	}
	return res
}
