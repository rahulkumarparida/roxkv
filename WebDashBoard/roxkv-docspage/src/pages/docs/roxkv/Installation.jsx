export default function Installation() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Installation</h1>
        <p className="text-xl text-[#9898ab]">Three different ways to get RoxKV up and running</p>
      </div>
      <div className="section-divider" />
      
      <h2 id="docker" className="text-2xl font-bold text-white mt-12 mb-4">1. Docker (Recommended)</h2>
      <p className="text-[#f1f1f4] mb-4">
        Using Docker is the easiest way to run RoxKV, as it bundles the core services along with an Ollama container for AI capabilities.
      </p>
      <div className="terminal-window my-4">
        <div className="terminal-header">
          <div className="flex gap-2">
            <div className="terminal-dot bg-red-500"></div>
            <div className="terminal-dot bg-yellow-500"></div>
            <div className="terminal-dot bg-green-500"></div>
          </div>
          <span className="ml-4 text-xs text-[#9898ab]">bash</span>
        </div>
        <div className="terminal-body bg-[#16161f] p-4 font-mono text-sm">
          <div><span className="prompt">$</span> <span className="command">git clone https://github.com/rahulkumarparida/roxkv.git && cd roxkv</span></div>
          <div><span className="prompt">$</span> <span className="command">cp .env.example .env</span></div>
          <div><span className="prompt">$</span> <span className="command">docker compose up -d --build</span></div>
        </div>
      </div>
      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#9898ab] text-sm">
          <strong>Note:</strong> This automatically exposes all necessary ports and handles Ollama integration for AI out-of-the-box.
        </p>
      </div>

      <h2 id="manual-build" className="text-2xl font-bold text-white mt-12 mb-4">2. Manual Build</h2>
      <p className="text-[#f1f1f4] mb-4">
        If you prefer to compile from source, you'll need Go (1.21+) and Node.js installed on your system.
      </p>
      <div className="terminal-window my-4">
        <div className="terminal-header">
          <div className="flex gap-2">
            <div className="terminal-dot bg-red-500"></div>
            <div className="terminal-dot bg-yellow-500"></div>
            <div className="terminal-dot bg-green-500"></div>
          </div>
          <span className="ml-4 text-xs text-[#9898ab]">bash</span>
        </div>
        <div className="terminal-body bg-[#16161f] p-4 font-mono text-sm">
          <div><span className="comment"># Clone and run the server</span></div>
          <div><span className="prompt">$</span> <span className="command">git clone https://github.com/rahulkumarparida/roxkv.git && cd roxkv</span></div>
          <div><span className="prompt">$</span> <span className="command">go run ./cmd/roxkv roxkv-tcp</span></div>
          <br/>
          <div><span className="comment"># Run the dashboard in a separate terminal</span></div>
          <div><span className="prompt">$</span> <span className="command">cd WebDashBoard/roxkv-dashboard</span></div>
          <div><span className="prompt">$</span> <span className="command">npm install && npm run dev</span></div>
        </div>
      </div>
      <p className="text-[#9898ab] text-sm mb-4">
        If you plan to use AI features with a manual build, you must install and run Ollama separately.
      </p>

      <h2 id="pre-built" className="text-2xl font-bold text-white mt-12 mb-4">3. Pre-built Binaries</h2>
      <p className="text-[#f1f1f4] mb-4">
        You can also run RoxKV directly using the pre-compiled binaries located in the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">executable/</code> directory.
      </p>
      <div className="terminal-window my-4">
        <div className="terminal-header">
          <div className="flex gap-2">
            <div className="terminal-dot bg-red-500"></div>
            <div className="terminal-dot bg-yellow-500"></div>
            <div className="terminal-dot bg-green-500"></div>
          </div>
          <span className="ml-4 text-xs text-[#9898ab]">bash</span>
        </div>
        <div className="terminal-body bg-[#16161f] p-4 font-mono text-sm">
          <div><span className="prompt">$</span> <span className="command">chmod +x ./executable/roxkv</span></div>
          <div><span className="prompt">$</span> <span className="command">./executable/roxkv roxkv-tcp</span></div>
        </div>
      </div>

      <h2 id="connecting" className="text-2xl font-bold text-white mt-12 mb-4">Connecting & Verifying</h2>
      <p className="text-[#f1f1f4] mb-4">
        Once your server is running, you can connect using the standard <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">redis-cli</code> on the Redis-compatible port <strong>6973</strong>.
      </p>
      <div className="terminal-window my-4">
        <div className="terminal-header">
          <div className="flex gap-2">
            <div className="terminal-dot bg-red-500"></div>
            <div className="terminal-dot bg-yellow-500"></div>
            <div className="terminal-dot bg-green-500"></div>
          </div>
          <span className="ml-4 text-xs text-[#9898ab]">redis-cli</span>
        </div>
        <div className="terminal-body bg-[#16161f] p-4 font-mono text-sm">
          <div><span className="prompt">$</span> <span className="command">redis-cli -h 127.0.0.1 -p 6973</span></div>
          <div><span className="prompt">127.0.0.1:6973&gt;</span> <span className="command">PING</span></div>
          <div><span className="output">PONG</span></div>
        </div>
      </div>
    </div>
  );
}
