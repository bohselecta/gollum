package impl

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/haydenlabs/gollum/engine"
	"github.com/haydenlabs/gollum/gguf"
)

type goEngine struct {
	scheduler *engine.Scheduler
	models    map[string]*gguf.Model // model name -> loaded model
}

func NewEngine() engine.Engine {
	// Find and load models
	models := make(map[string]*gguf.Model)
	modelFiles, err := gguf.FindModels(".")
	if err != nil {
		log.Printf("Warning: failed to find models: %v", err)
	} else {
		for _, path := range modelFiles {
			log.Printf("Loading model: %s", path)
			model, err := gguf.LoadModel(path)
			if err != nil {
				log.Printf("Failed to load %s: %v", path, err)
				continue
			}
			// Use filename (without extension) as model ID
			modelName := filepath.Base(path)
			modelName = modelName[:len(modelName)-5] // remove .gguf
			models[modelName] = model
			log.Printf("Loaded model: %s (layers=%d, embed=%d, vocab=%d)", 
				modelName, model.NumLayers, model.EmbedDim, model.VocabSize)
		}
	}
	
	if len(models) == 0 {
		log.Printf("No GGUF models found, using toy backend")
	}
	
	backend := NewMetalOps()
	scheduler := engine.NewScheduler(backend)
	// Start the scheduler in the background
	go scheduler.Run(context.Background())
	return &goEngine{
		scheduler: scheduler,
		models:   models,
	}
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

func (e *goEngine) GetModel(name string) (*gguf.Model, error) {
	if model, ok := e.models[name]; ok {
		return model, nil
	}
	return nil, fmt.Errorf("model not found: %s", name)
}
