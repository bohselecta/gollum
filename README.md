<p align="center">
  <img src="assets/mascot_gollum.svg" width="160" alt="GoLLuM mascot"/>
</p>

# GoLLuM (recreated)


## Prompt cache
- SHA-256 keyed by (prompt, model, temperature, max_tokens), size-limited LRU (default 512 entries).
- Scheduler checks the cache on enqueue; full hits replay instantly. Misses/partials run normally and store on completion.

## KV + Prefix reuse
This integrates a `PrefixCache` mapping `(model, prefix)` to `KVRef{BlockID, Tokens}`.
- Prefill checks for prefix hits, pins blocks for reuse.
- Completion stores prefix→block and unpins for LRU reuse.
This enables prefix-KV block reuse in future attention kernels.


## Quick start
```bash
make run
# In another terminal:
curl http://localhost:8080/v1/models
curl -N -H "Content-Type: application/json" -d '{"model":"toy-1","messages":[{"role":"user","content":"Write a 7-word poem about llamas."}]}' http://localhost:8080/v1/chat/completions
```

## What works now
- OpenAI-style routes: `/v1/models`, `/v1/chat/completions` (SSE streaming), `/v1/embeddings` (stub)
- Continuous-batching-friendly engine interfaces
- Simple scheduler and paged KV cache **interfaces** (toy impl is stateless)
- Prometheus metrics at `/metrics`, pprof at `/debug/pprof/`
- Prompt caching with LRU eviction
- Beautiful landing page with GoLLuM branding

## Layout
```
/cmd/altiserve        # server main
/apis/openai          # HTTP routes (Gin)
/engine               # engine, scheduler, kv pager interfaces
/engine/impl          # minimal engine implementation with toy backend
/kernels/toy          # toy "kernel" (just fake decode loop)
/kernels/metal        # placeholder (Obj-C shim stubs)
/kernels/cuda         # placeholder (C shim stubs)
/tokenizer            # stub tokenizer + hooks for sentencepiece/tt later
/obs                  # metrics and tracing helpers
/assets               # static assets and icons
/examples             # client examples for Node.js and Python
```

## Extending to real backends
1. Implement fused ops in C/Obj-C and expose via cgo in `/kernels/metal` or `/kernels/cuda`.
2. Satisfy `KernelOps` in `/engine/impl/backends/` by calling those fused ops in batches.
3. Keep Go on the **serving/scheduling** hot path, but minimize cgo crossings (batch large calls).

## License
MIT (for now).


## macOS (Metal) notes
- Requires Xcode Command Line Tools: `xcode-select --install`
- Build/run as usual: `make run` (the Metal shim compiles only on macOS).
- On launch, you'll see logs include the detected Metal device via the shim.
- The current shim performs no real GPU math; it's a safe place to add **fused decode** later.

### Where to add real kernels next
- Implement fused attention/ffn in `kernels/metal/shim.mm` and expose via C functions.
- Add a proper `KernelOps` impl in `engine/impl/backends/metal.go` that batches requests and calls those functions once per step.
- Keep Go in charge of **scheduling** and **serving**; make a single cgo call per fused step.