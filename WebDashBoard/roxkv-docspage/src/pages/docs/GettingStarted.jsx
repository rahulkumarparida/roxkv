import { Terminal } from 'lucide-react';

export default function GettingStarted() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Getting Started with RoxKV</h1>
        <p className="text-xl text-[#9898ab]">
          A high-performance, Redis-compatible, agent-integrated key-value store built in Go.
        </p>
      </div>

      <div className="section-divider" />

      <h2 id="introduction" className="text-2xl font-bold text-white mt-12 mb-4">Introduction</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        RoxKV is an in-memory database that implements a subset of the Redis Serialization Protocol (RESP).
        It is designed to be fast, reliable, and integrated closely with agentic workflows (RoxAI).
        You can interact with RoxKV using standard Redis clients, such as `redis-cli`, or any programming language Redis driver.
      </p>

      <h2 id="running-roxkv" className="text-2xl font-bold text-white mt-12 mb-4">Running RoxKV</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        You can run RoxKV locally using Docker. We provide a convenient docker-compose file that spins up the RoxKV server along with the dashboard.
      </p>

      <div className="terminal-window my-6">
        <div className="terminal-header">
          <div className="flex gap-2">
            <div className="terminal-dot bg-red-500"></div>
            <div className="terminal-dot bg-yellow-500"></div>
            <div className="terminal-dot bg-green-500"></div>
          </div>
          <span className="ml-4 text-xs text-[#9898ab]">bash</span>
        </div>
        <div className="terminal-body bg-[#16161f] p-4 font-mono text-sm">
          <div><span className="prompt">$</span> <span className="command">git clone https://github.com/rahulkumarparida/roxkv.git</span></div>
          <div><span className="prompt">$</span> <span className="command">cd roxkv</span></div>
          <div><span className="prompt">$</span> <span className="command">docker compose up --build</span></div>
        </div>
      </div>

      <h2 id="connecting" className="text-2xl font-bold text-white mt-12 mb-4">Connecting</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        By default, RoxKV listens on port `6379`, matching standard Redis. This means that out of the box, existing tooling and libraries just work.
      </p>
      
      <div className="bg-[rgba(124,58,237,0.1)] border border-[rgba(124,58,237,0.2)] rounded-xl p-6 my-6 flex items-start gap-4">
        <Terminal className="text-[#a78bfa] flex-shrink-0 mt-1" />
        <div>
          <h4 className="text-white font-semibold mb-1">Using redis-cli</h4>
          <p className="text-[#9898ab] text-sm mb-3">
            If you have redis-cli installed, you can simply connect to localhost:
          </p>
          <code className="bg-[#0a0a0f] px-3 py-1.5 rounded-md text-sm border border-[#1e1e2e] text-[#a78bfa]">
            redis-cli -h 127.0.0.1 -p 6379
          </code>
        </div>
      </div>

      <p className="text-[#f1f1f4] leading-relaxed">
        Once connected, try running <code className="bg-[#1e1e2e] px-1.5 py-0.5 rounded text-[#a78bfa]">PING</code>. The server should reply with <code className="bg-[#1e1e2e] px-1.5 py-0.5 rounded text-white">PONG</code>.
      </p>
    </div>
  );
}
