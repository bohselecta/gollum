package impl
import ("context"
	"github.com/haydenlabs/gollum/engine")
type goEngine struct{ sched *engine.Scheduler }
func NewEngine() engine.Engine { ops := &metalOps{}; s := engine.NewScheduler(ops); ge := &goEngine{sched:s}; go s.Run(context.Background()); return ge }
func (e *goEngine) Generate(ctx context.Context, req *engine.GenRequest)(<-chan engine.Token,*engine.Trace,error){
    if req.MaxTokens<=0 { req.MaxTokens=64 }
    ch, trace := e.sched.Enqueue(ctx, req)
    return ch, trace, nil
}
func (e *goEngine) Embeddings(ctx context.Context, input string)([]float32,error){
    out:=make([]float32,128); for i := range out { out[i] = float32((i*7)%13)/13.0 }; return out, nil
}
