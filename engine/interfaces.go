package engine

import "context"

type Token struct {
	ID   int
	Text string
}

type Trace struct {
	TTFTMs int64
	TPOTMs int64
}

type GenRequest struct {
	Model       string
	Prompt      string
	MaxTokens   int
	Temperature float32
	Priority    int
}

type Engine interface {
	Generate(ctx context.Context, req *GenRequest) (<-chan Token, *Trace, error)
	Embeddings(ctx context.Context, input string) ([]float32, error)
}

type KernelOps interface {
	Prefill(batch *Batch) error
	Decode(step *Step) (int, error)
	PredictNext(prompts []string) []string
}

type Batch struct {
	Prompts []string
	// future: tensors, kv handles, dtypes
}

type Step struct {
	BatchSize int
	SeqLen    int
	// future: kv ptrs, positions
}
