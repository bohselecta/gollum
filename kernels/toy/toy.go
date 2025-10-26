package toy

import (
	"math/rand"
	"strings"
	"time"

	"github.com/haydenlabs/gollum/engine"
)

type ToyBackend struct{}

func NewToyBackend() *ToyBackend { return &ToyBackend{} }

func (t *ToyBackend) Prefill(batch *engine.Batch) error {
	// Toy implementation: just simulate some processing time
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (t *ToyBackend) Decode(step *engine.Step) (int, error) {
	// Toy implementation: simulate processing time
	time.Sleep(5 * time.Millisecond)
	return step.BatchSize, nil
}

func (t *ToyBackend) PredictNext(ctxs []engine.DecodeCtx) []string {
	results := make([]string, len(ctxs))
	for i, ctx := range ctxs {
		results[i] = t.NextToken(ctx.Prompt)
	}
	return results
}

// NextToken emits a tiny "continuation" that looks language-ish.
func (t *ToyBackend) NextToken(prefix string) string {
	rand.Seed(time.Now().UnixNano())
	w := []string{" llama", " on", " the", " high", " plain", ".", " gentle", ",", " wind", " hums", " softly"}
	if strings.HasSuffix(prefix, ".") {
		return " "
	}
	return w[rand.Intn(len(w))]
}
