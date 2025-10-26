<p align="center">
  <img src="assets/mascot_gollum.svg" width="160" alt="GoLLuM mascot"/>
</p>

# GoLLuM (recreated)


## Prompt cache
- SHA-256 keyed by (prompt, model, temperature, max_tokens), size-limited LRU (default 512 entries).
- Scheduler checks the cache on enqueue; full hits replay instantly. Misses/partials run normally and store on completion.
- This is a scaffolding for a real prefix-cache tied to KV blocks (coming in the KV integration step).


## KV + Prefix reuse (scaffold)
This integrates a `PrefixCache` mapping `(model, prefix)` to `KVRef{BlockID, Tokens}`.
- Prefill checks for prefix hits, pins blocks for reuse.
- Completion stores prefix→block and unpins for LRU reuse.
This enables prefix-KV block reuse in future attention kernels.
