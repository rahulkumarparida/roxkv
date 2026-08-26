import React from 'react';

export default function TtlExpiration() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">TTL & Expiration</h1>
        <p className="text-xl text-[#9898ab]">Key lifetime management and background cleanup.</p>
      </div>
      
      <div className="section-divider" />

      <h2 id="expiry-worker" className="text-2xl font-semibold text-white mb-4">Background Expiry Worker</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        Active expiration is handled by the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">ExpiryWorker</code> in <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">internal/worker/worker.go</code>. It runs on a 1-second ticker, performing a full scan of the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">MemoryAlloc.Data</code> map. If it encounters a key where <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">TTL.Before(time.Now())</code> is true, it removes the key via <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">store.DelKv()</code>.
      </p>

      <h2 id="lazy-deletion" className="text-2xl font-semibold text-white mb-4 mt-8">Lazy Deletion</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        In addition to the background worker, lazy deletion is enforced. For example, during operations like <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">store.Keys()</code>, if an expired key is encountered, it is immediately deleted and omitted from the results.
      </p>

      <h2 id="metrics" className="text-2xl font-semibold text-white mb-4 mt-8">Metrics Tracking</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        The worker updates the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">TTLMetricsContainer</code>, incrementing the count of expired keys for observability.
      </p>

      <h2 id="commands" className="text-2xl font-semibold text-white mb-4 mt-8">TTL Commands</h2>
      
      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#f1f1f4]"><strong>Important Note:</strong> The <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">SET</code> command in RoxKV currently does NOT support inline <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">EX</code> or <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">PX</code> arguments. You must call <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">EXPIRE</code> separately.</p>
      </div>

      <div className="space-y-4">
        <div className="terminal-window">
          <div className="terminal-header">
            <div className="terminal-dot bg-red-500"></div>
            <div className="terminal-dot bg-yellow-500"></div>
            <div className="terminal-dot bg-green-500"></div>
          </div>
          <div className="terminal-body font-mono text-sm p-4 text-[#f1f1f4]">
            <div className="comment text-[#5e5e73] mb-1"># Set a key and apply a 60 second TTL</div>
            <div><span className="prompt text-[#7c3aed]">&gt;</span> <span className="command text-[#a78bfa]">SET</span> mykey "hello"</div>
            <div className="output text-[#9898ab] mb-2">OK</div>
            <div><span className="prompt text-[#7c3aed]">&gt;</span> <span className="command text-[#a78bfa]">EXPIRE</span> mykey 60</div>
            <div className="output text-[#9898ab] mb-2">(integer) 1</div>
            
            <div className="comment text-[#5e5e73] mt-4 mb-1"># Check remaining time-to-live</div>
            <div><span className="prompt text-[#7c3aed]">&gt;</span> <span className="command text-[#a78bfa]">TTL</span> mykey</div>
            <div className="output text-[#9898ab] mb-2">(integer) 58</div>
            
            <div className="comment text-[#5e5e73] mt-4 mb-1"># Remove the expiration (make it persistent)</div>
            <div><span className="prompt text-[#7c3aed]">&gt;</span> <span className="command text-[#a78bfa]">PERSIST</span> mykey</div>
            <div className="output text-[#9898ab]">(integer) 1</div>
          </div>
        </div>
      </div>
    </div>
  );
}
