package impl

import (
	"context"

	"github.com/haydenlabs/gollum/engine"
	"github.com/haydenlabs/gollum/kernels/toy"
)

type goEngine struct {
	scheduler *engine.Scheduler
}

func NewEngine() engine.Engine {
	backend := toy.NewToyBackend()
	scheduler := engine.NewScheduler(backend)
	// Start the scheduler in the background
	go scheduler.Run(context.Background())
	return &goEngine{scheduler: scheduler}
}

func (e *goEngine) Generate(ctx context.Context, req *engine.GenRequest) (<-chan engine.Token, *engine.Trace, error) {
	if req.MaxTokens <= 0 {
		req.MaxTokens = 64
	}
	ch, trace := e.scheduler.Enqueue(ctx, req)
	return ch, trace, nil
}

func (e *goEngine) Embeddings(ctx context.Context, input string) ([]float32, error) {
	// stub: fixed-size vector
	out := make([]float32, 128)
	for i := range out {
		out[i] = float32((i*7)%13) / 13.0
	}
	return out, nil
}
