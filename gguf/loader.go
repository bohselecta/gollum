package gguf

import (
	"fmt"
	"os"
	"path/filepath"
)

// Model represents a loaded GGUF model
type Model struct {
	Path       string
	GGUF       *GGUF
	EmbedDim   int
	VocabSize  int
	ContextLen int
	NumLayers  int
	NumHeads   int
	HeadDim    int
}

// LoadModel loads a GGUF file from the models directory
func LoadModel(path string) (*Model, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open GGUF file: %w", err)
	}
	defer f.Close()

	gguf, err := Parse(f)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GGUF: %w", err)
	}

	m := &Model{
		Path: path,
		GGUF: gguf,
	}

	// Extract metadata from tensors (simplified)
	if tokEmb, ok := gguf.Tensors["tok_embeddings.weight"]; ok {
		if len(tokEmb.Dims) >= 2 {
			m.VocabSize = int(tokEmb.Dims[0])
			m.EmbedDim = int(tokEmb.Dims[1])
		}
	}

	// Try to infer architecture from tensor names
	m.NumLayers = 0
	for name := range gguf.Tensors {
		if len(name) > 8 && name[0:8] == "layers." {
			layerIdx := 0
			if _, err := fmt.Sscanf(name, "layers.%d", &layerIdx); err == nil {
				if layerIdx >= m.NumLayers {
					m.NumLayers = layerIdx + 1
				}
			}
		}
	}

	// Default values if not found
	if m.EmbedDim == 0 {
		m.EmbedDim = 4096 // Llama default
	}
	if m.VocabSize == 0 {
		m.VocabSize = 32000 // Llama default
	}
	if m.NumLayers == 0 {
		m.NumLayers = 32 // Default
	}

	m.NumHeads = 32 // Common default
	m.HeadDim = m.EmbedDim / m.NumHeads
	m.ContextLen = 2048 // Default context

	return m, nil
}

// FindModels scans the models directory for GGUF files
func FindModels(basePath string) ([]string, error) {
	modelsDir := filepath.Join(basePath, "models")

	files, err := os.ReadDir(modelsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read models directory: %w", err)
	}

	var models []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".gguf" {
			models = append(models, filepath.Join(modelsDir, file.Name()))
		}
	}

	return models, nil
}
