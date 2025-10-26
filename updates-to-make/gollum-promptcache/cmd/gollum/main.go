package main
import (
	"fmt"
    "log"
    "net/http"
    _ "net/http/pprof"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/haydenlabs/gollum/apis/openai"
    "github.com/haydenlabs/gollum/engine/impl"
    "github.com/haydenlabs/gollum/obs"
)
func main(){
    port := os.Getenv("PORT")
    if port=="" { port="8080" }
    engine := impl.NewEngine()
    r := gin.Default()
    obs.Install(r)

	// Static assets and favicon
	r.Static("/assets", "./assets")
	r.StaticFile("/favicon.ico", "./assets/icons/gollum_icon_256.png")
    r.GET("/metrics", gin.WrapH(promhttp.Handler()))
    api := openai.NewAPI(engine)
    api.Register(r)

	selfTestMetal()

	// Minimal landing page
	r.GET("/", func(c *gin.Context){
		c.Header("Content-Type","text/html; charset=utf-8")
		c.String(200, landingHTML)
	})
    log.Printf("GoLLuM listening on :%s", port)
    go func(){ log.Println(http.ListenAndServe(":6060", nil)) }()
    if err := r.Run(":"+port); err!=nil { log.Fatal(err) }
}


const landingHTML = `<!doctype html>
  <html lang="en">
    <head>
      <meta charset="utf-8" />
      <meta name="viewport" content="width=device-width, initial-scale=1" />
      <title>GoLLuM • Go-first LLM runtime</title>
      <link rel="icon" href="/favicon.ico">
      <style>
        :root { --ink:#0f172a; --muted:#475569; --bg:#0b1220; --card:#0f1b2d; --ring:#38bdf8; }
        *{box-sizing:border-box} body{margin:0;font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Inter,Roboto,Helvetica,Arial,sans-serif;background:linear-gradient(180deg,#08101c,#0b1626);color:#e2e8f0}
        .wrap{max-width:960px;margin:48px auto;padding:24px}
        .card{background:linear-gradient(180deg,var(--card),#0b1525);border:1px solid rgba(56,189,248,.2);border-radius:18px;padding:24px;box-shadow:0 10px 30px rgba(8,12,23,.35)}
        .row{display:flex;gap:24px;align-items:center;flex-wrap:wrap}
        .row img{width:128px;height:128px}
        code,pre{background:#0a1422;border:1px solid rgba(148,163,184,.2);border-radius:12px;padding:.5rem 0.75rem;color:#e5e7eb}
        pre{overflow:auto}
        h1{margin:0 0 8px 0;font-size:28px} p{margin:8px 0 0 0;color:#cbd5e1}
        a{color:#7dd3fc}
      </style>
    </head>
    <body>
      <div class="wrap">
        <div class="card">
          <div class="row">
            <img src="/assets/icons/gollum_icon_256.png" alt="GoLLuM icon"/>
            <div>
              <h1>GoLLuM</h1>
              <p>Go-first LLM inference runtime • OpenAI-compatible API • Continuous batching • Mac-ready (Metal)</p>
            </div>
          </div>
          <hr style="border:none;border-top:1px solid rgba(148,163,184,.2);margin:16px 0"/>
          <p>Server is up. Try a quick request:</p>
          <pre><code>curl -N -H "Content-Type: application/json" \
-d '{"model":"toy-1","messages":[{"role":"user","content":"Say hi!"}]}' \
http://localhost:8080/v1/chat/completions</code></pre>
          <p>OpenAI client examples are in <code>examples/</code>.</p>
        </div>
      </div>
    </body>
  </html>`


func selfTestMetal() {
	// Only prints if Metal is present (darwin build)
	defer func(){ recover() }()
	name := metalDeviceName()
	if name == "" { return }
	fmt.Printf("[GoLLuM] Metal device: %s\n", name)
	// 2x3 * 3x2 => 2x2
	A := []float32{1,2,3, 4,5,6}
	B := []float32{7,8, 9,10, 11,12}
	C, ok := metalMatMul(A,B,2,2,3)
	if !ok { fmt.Println("[GoLLuM] Metal matmul failed"); return }
	fmt.Printf("[GoLLuM] Metal matmul ok: %.1f %.1f | %.1f %.1f\n", C[0], C[1], C[2], C[3])
}


func metalDeviceName() string {
	return metal.DeviceName()
}
func metalMatMul(A,B []float32, M,N,K int)([]float32,bool){
	return metal.MatMul(A,B,M,N,K)
}
