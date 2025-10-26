package main

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
