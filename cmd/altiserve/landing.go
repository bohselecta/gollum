package main

const landingHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>GoLLuM • LLM Chat</title>
  <link rel="icon" href="/favicon.ico">
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    
    :root {
      --bg: #0b1220;
      --card: #0f1b2d;
      --border: rgba(56, 189, 248, 0.2);
      --text: #e2e8f0;
      --muted: #94a3b8;
      --primary: #38bdf8;
      --input-bg: #0a1422;
      --hover: rgba(56, 189, 248, 0.1);
    }
    
    body {
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Inter, Roboto, sans-serif;
      background: linear-gradient(180deg, #08101c, #0b1626);
      color: var(--text);
      height: 100vh;
      overflow: hidden;
    }
    
    #app {
      display: flex;
      flex-direction: column;
      height: 100vh;
      max-width: 960px;
      margin: 0 auto;
      padding: 16px;
    }
    
    .header {
      display: flex;
      align-items: center;
      gap: 16px;
      padding: 16px;
      background: var(--card);
      border: 1px solid var(--border);
      border-radius: 12px;
      margin-bottom: 16px;
      box-shadow: 0 4px 12px rgba(8, 12, 23, 0.3);
    }
    
    .logo {
      width: 48px;
      height: 48px;
      border-radius: 8px;
      filter: drop-shadow(0 2px 4px rgba(0,0,0,0.3));
    }
    
    .wordmark {
      height: 32px;
      filter: drop-shadow(0 2px 4px rgba(0,0,0,0.3));
    }
    
    .header-info {
      flex: 1;
    }
    
    .header p {
      font-size: 14px;
      color: var(--muted);
      margin-top: 4px;
    }
    
    .messages {
      flex: 1;
      overflow-y: auto;
      padding: 16px;
      background: var(--card);
      border: 1px solid var(--border);
      border-radius: 12px;
      margin-bottom: 16px;
      display: flex;
      flex-direction: column;
      gap: 12px;
    }
    
    .message {
      display: flex;
      gap: 12px;
      padding: 12px;
      background: rgba(10, 20, 34, 0.4);
      border-radius: 8px;
      animation: fadeIn 0.3s ease-in;
    }
    
    @keyframes fadeIn {
      from { opacity: 0; transform: translateY(4px); }
      to { opacity: 1; transform: translateY(0); }
    }
    
    .message.user { background: rgba(56, 189, 248, 0.1); }
    .message.assistant { background: rgba(10, 20, 34, 0.4); }
    
    .message-icon {
      width: 32px;
      height: 32px;
      border-radius: 50%;
      background: var(--primary);
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: bold;
      flex-shrink: 0;
      font-size: 14px;
      color: #0b1220;
    }
    
    .message.assistant .message-icon {
      background: linear-gradient(135deg, #667eea, #764ba2);
      color: white;
    }
    
    .message-content {
      flex: 1;
      line-height: 1.6;
      white-space: pre-wrap;
      word-wrap: break-word;
    }
    
    .input-area {
      display: flex;
      gap: 8px;
      padding: 16px;
      background: var(--card);
      border: 1px solid var(--border);
      border-radius: 12px;
    }
    
    #input {
      flex: 1;
      background: var(--input-bg);
      border: 1px solid var(--border);
      border-radius: 8px;
      padding: 12px 16px;
      color: var(--text);
      font-family: inherit;
      font-size: 14px;
      resize: none;
      min-height: 48px;
      max-height: 120px;
    }
    
    #input:focus {
      outline: none;
      border-color: var(--primary);
      box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.1);
    }
    
    #send {
      background: var(--primary);
      color: #0b1220;
      border: none;
      border-radius: 8px;
      padding: 12px 24px;
      font-weight: 600;
      cursor: pointer;
      transition: all 0.2s;
    }
    
    #send:hover {
      background: #2dd4bf;
      transform: translateY(-1px);
    }
    
    #send:active {
      transform: translateY(0);
    }
    
    #send:disabled {
      background: var(--muted);
      cursor: not-allowed;
      transform: none;
    }
    
    .typing-indicator {
      display: inline-block;
      animation: blink 1.4s infinite;
    }
    
    @keyframes blink {
      0%, 50% { opacity: 1; }
      51%, 100% { opacity: 0.3; }
    }
    
    ::-webkit-scrollbar {
      width: 8px;
    }
    
    ::-webkit-scrollbar-track {
      background: var(--input-bg);
      border-radius: 4px;
    }
    
    ::-webkit-scrollbar-thumb {
      background: var(--muted);
      border-radius: 4px;
    }
    
    ::-webkit-scrollbar-thumb:hover {
      background: var(--primary);
    }
  </style>
</head>
<body>
  <div id="app">
    <div class="header">
      <img src="/assets/gollum-logo.svg" alt="" class="logo">
      <img src="/assets/word-mark-gollum-text.svg" alt="GoLLuM" class="wordmark">
      <div class="header-info">
        <p>Go-first LLM inference runtime with advanced caching</p>
      </div>
    </div>
    
    <div class="messages" id="messages"></div>
    
    <div class="input-area">
      <textarea id="input" placeholder="Ask GoLLuM anything..." rows="1"></textarea>
      <button id="send">Send</button>
    </div>
  </div>
  
  <script>
    const messagesEl = document.getElementById('messages');
    const inputEl = document.getElementById('input');
    const sendBtn = document.getElementById('send');
    
    let currentAssistant = null;
    
    function addMessage(role, content) {
      const div = document.createElement('div');
      div.className = 'message ' + role;
      
      const icon = document.createElement('div');
      icon.className = 'message-icon';
      icon.textContent = role === 'user' ? 'You' : 'G';
      
      const contentEl = document.createElement('div');
      contentEl.className = 'message-content';
      contentEl.textContent = content;
      
      div.appendChild(icon);
      div.appendChild(contentEl);
      messagesEl.appendChild(div);
      
      messagesEl.scrollTop = messagesEl.scrollHeight;
      
      return { div, contentEl };
    }
    
    function updateAssistant(content) {
      if (!currentAssistant) {
        const msg = addMessage('assistant', content);
        currentAssistant = msg;
      } else {
        currentAssistant.contentEl.textContent = content + '▊';
        messagesEl.scrollTop = messagesEl.scrollHeight;
      }
    }
    
    function finalizeAssistant() {
      if (currentAssistant) {
        currentAssistant.contentEl.textContent = currentAssistant.contentEl.textContent.replace('▊', '');
        currentAssistant = null;
      }
    }
    
    async function sendMessage() {
      const text = inputEl.value.trim();
      if (!text) return;
      
      inputEl.value = '';
      sendBtn.disabled = true;
      
      addMessage('user', text);
      
      try {
        const response = await fetch('/v1/chat/completions', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            model: 'toy-1',
            messages: [{ role: 'user', content: text }],
            stream: true,
            max_tokens: 128
          })
        });
        
        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        let buffer = '';
        let fullContent = '';
        
        while (true) {
          const { done, value } = await reader.read();
          if (done) {
            finalizeAssistant();
            break;
          }
          
          buffer += decoder.decode(value, { stream: true });
          const lines = buffer.split('\n\n');
          buffer = lines.pop() || '';
          
          for (const line of lines) {
            if (line.startsWith('data: ') && !line.includes('[DONE]')) {
              try {
                const data = JSON.parse(line.slice(6));
                const content = data.choices?.[0]?.delta?.content || '';
                if (content) {
                  fullContent += content;
                  updateAssistant(fullContent);
                }
              } catch (e) {}
            }
          }
        }
      } catch (e) {
        addMessage('assistant', 'Error: ' + e.message);
      } finally {
        sendBtn.disabled = false;
      }
    }
    
    sendBtn.addEventListener('click', sendMessage);
    inputEl.addEventListener('keydown', (e) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        sendMessage();
      }
    });
    
    // Auto-resize textarea
    inputEl.addEventListener('input', function() {
      this.style.height = 'auto';
      this.style.height = Math.min(this.scrollHeight, 120) + 'px';
    });
  </script>
</body>
</html>`
