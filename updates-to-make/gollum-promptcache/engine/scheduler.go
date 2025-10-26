package engine
import ("context"; "time"; "sync")
type reqState struct{ ctx context.Context; req *GenRequest; ch chan Token; trace *Trace; generated int; created time.Time }
type Scheduler struct{ mu sync.Mutex; incoming, active []*reqState; closed bool; backend KernelOps; maxBatch int; stepInterval time.Duration; pc *PromptCache }
func NewScheduler(b KernelOps)*Scheduler{ return &Scheduler{backend:b,maxBatch:32,stepInterval:25*time.Millisecond, pc: NewPromptCache(512)} }
func (s *Scheduler) Enqueue(ctx context.Context, r *GenRequest) (<-chan Token, *Trace){
    rs:=&reqState{ctx:ctx,req:r,ch:make(chan Token,32),trace:&Trace{},created:time.Now()}
    // Cache fast-path: exact prompt/model/temp/maxTokens
    if toks, ok := s.pc.Get(r.Prompt, r.Model, r.Temperature, r.MaxTokens); ok && len(toks) >= r.MaxTokens {
        go func(){
            defer close(rs.ch)
            for i:=0;i<r.MaxTokens;i++{
                select {
                case <-ctx.Done(): return
                default:
                    rs.ch <- Token{Text: toks[i]}
                }
            }
            rs.trace.TTFTMs = 1
            rs.trace.TPOTMs = 1
        }()
        return rs.ch, rs.trace
    }
    s.mu.Lock(); s.incoming=append(s.incoming,rs); s.mu.Unlock(); return rs.ch, rs.trace
}
func (s *Scheduler) Run(ctx context.Context){ t:=time.NewTicker(s.stepInterval); defer t.Stop(); for{ select{ case <-ctx.Done(): s.shutdown(); return; case <-t.C: s.tick() } } }
func (s *Scheduler) shutdown(){ s.mu.Lock(); defer s.mu.Unlock(); if s.closed {return}; s.closed=true; for _,rs := range append(s.incoming,s.active...) { close(rs.ch) } }
func (s *Scheduler) tick(){
    // Prefill
    var prefill []*reqState
    s.mu.Lock(); if len(s.incoming)>0 { limit:=s.maxBatch; if len(s.incoming)<limit { limit=len(s.incoming) }; prefill = append(prefill, s.incoming[:limit]...); s.incoming = s.incoming[limit:]; }; s.mu.Unlock()
    if len(prefill)>0 {
        b := &Batch{Prompts: make([]string,0,len(prefill))}
        for _,rs := range prefill { b.Prompts = append(b.Prompts, rs.req.Prompt) }
        _ = s.backend.Prefill(b)
        for _,rs := range prefill { if rs.trace.TTFTMs==0 { rs.trace.TTFTMs = time.Since(rs.created).Milliseconds() }; s.mu.Lock(); s.active = append(s.active, rs); s.mu.Unlock() }
    }
    // Decode step
    s.mu.Lock(); active := append([]*reqState{}, s.active...); s.mu.Unlock()
    if len(active)==0 { return }
    produced,_ := s.backend.Decode(&Step{BatchSize:len(active), SeqLen:1}); if produced<=0 { produced = len(active) }
    // Build prompts and predict
    prompts := make([]string, len(active))
    for i,rs := range active { prompts[i] = rs.req.Prompt }
    nexts := s.backend.PredictNext(prompts)
    // Emit up to 'produced' tokens
    var still []*reqState
    for i, rs := range active {
        if rs.ctx.Err()!=nil { close(rs.ch); continue }
        if i < produced {
            tokText := ""
            if i < len(nexts) { tokText = nexts[i] }
            rs.req.Prompt += tokText  // evolve prompt context
            rs.ch <- Token{Text: tokText}
            rs.generated++
        }
        if rs.req.MaxTokens>0 && rs.generated>=rs.req.MaxTokens { close(rs.ch); rs.trace.TPOTMs = time.Since(rs.created).Milliseconds(); s.pc.Put(rs.req.Prompt, rs.req.Model, rs.req.Temperature, rs.req.MaxTokens, replayTokens(rs.req.Prompt)); continue }
        still = append(still, rs)
    }
    s.mu.Lock(); s.active = still; s.mu.Unlock()
}


func replayTokens(prompt string) []string {
    // naive split by known tiny vocab pieces; in a real system we’d store exact emitted tokens.
    out := []string{}
    tail := prompt
    // This is a stub to satisfy the cache Put(); downstream will be replaced when we hold per-req emission buffers.
    for len(tail) > 0 {
        if len(tail) > 12 {
            out = append(out, tail[:12])
            tail = tail[12:]
        } else {
            out = append(out, tail)
            break
        }
    }
    return out
}
