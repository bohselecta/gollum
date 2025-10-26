package openai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/haydenlabs/gollum/engine"
)

type API struct {
	eng engine.Engine
}

func NewAPI(eng engine.Engine) *API { return &API{eng: eng} }

func (a *API) Register(r *gin.Engine) {
	v1 := r.Group("/v1")
	v1.GET("/models", a.Models)
	v1.POST("/chat/completions", a.ChatCompletions)
	v1.POST("/embeddings", a.Embeddings)
}

func (a *API) Models(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"object": "list",
		"data":   []gin.H{{"id": "toy-1", "object": "model", "created": time.Now().Unix(), "owned_by": "gollum"}},
	})
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Stream      bool          `json:"stream"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float32       `json:"temperature"`
}

func (a *API) ChatCompletions(c *gin.Context) {
	var req ChatCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Model == "" {
		req.Model = "toy-1"
	}
	id := "chatcmpl_" + uuid.New().String()

	// Collapse messages to a prompt (very naive for toy backend)
	prompt := ""
	for _, m := range req.Messages {
		if m.Role == "user" || m.Role == "system" {
			prompt += m.Content + "\n"
		}
	}

	ctx := c.Request.Context()
	stream := c.Writer
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Status(http.StatusOK)

	ch, trace, err := a.eng.Generate(ctx, &engine.GenRequest{
		Model: req.Model, Prompt: prompt, MaxTokens: req.MaxTokens, Temperature: req.Temperature,
	})
	if err != nil {
		fmt.Fprintf(stream, "data: %s\n\n", toJSON(gin.H{"error": err.Error()}))
		return
	}
	// Send an initial event like OpenAI
	fmt.Fprintf(stream, "data: %s\n\n", toJSON(gin.H{
		"id": id, "object": "chat.completion.chunk", "created": time.Now().Unix(),
		"model": req.Model, "choices": []gin.H{{"index": 0, "delta": gin.H{"role": "assistant", "content": ""}, "finish_reason": nil}},
	}))

	for tok := range ch {
		fmt.Fprintf(stream, "data: %s\n\n", toJSON(gin.H{
			"id": id, "object": "chat.completion.chunk", "created": time.Now().Unix(),
			"model": req.Model, "choices": []gin.H{{"index": 0, "delta": gin.H{"content": tok.Text}, "finish_reason": nil}},
		}))
		stream.(http.Flusher).Flush()
	}
	// done
	fmt.Fprintf(stream, "data: %s\n\n", toJSON(gin.H{
		"id": id, "object": "chat.completion.chunk", "created": time.Now().Unix(),
		"model": req.Model, "choices": []gin.H{{"index": 0, "delta": gin.H{}, "finish_reason": "stop"}},
	}))
	fmt.Fprint(stream, "data: [DONE]\n\n")

	_ = trace // currently unused; wire into logs later
}

func (a *API) Embeddings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"object": "list",
		"data":   []gin.H{},
	})
}

func toJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
