import React from 'react';

export default function PubSub() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Pub/Sub</h1>
        <p className="text-xl text-[#9898ab]">Message brokering and pattern subscriptions.</p>
      </div>
      
      <div className="section-divider" />

      <h2 id="overview" className="text-2xl font-semibold text-white mb-4">Overview</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        RoxKV provides publish/subscribe messaging functionality. The <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">AllChannels</code> broker manages a map of <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">SubrChannel</code> instances. Each channel tracks its subscribers, last publisher, total messages, and byte size, keeping a short message history ring.
      </p>

      <h2 id="subscriber-mode" className="text-2xl font-semibold text-white mb-4 mt-8">Subscriber Mode</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        When a client issues a <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">SUBSCRIBE</code> or <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">PSUBSCRIBE</code> command, it is placed into <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">ModeSubscriber</code>.
      </p>
      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#f1f1f4]"><strong>Restriction:</strong> In subscriber mode, only the following commands are allowed: <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">SUBSCRIBE</code>, <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">UNSUBSCRIBE</code>, <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">PSUBSCRIBE</code>, <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">PUNSUBSCRIBE</code>, <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">PING</code>, <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">QUIT</code>, and <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">RESET</code>.</p>
      </div>

      <h2 id="publishing" className="text-2xl font-semibold text-white mb-4 mt-8">Publishing Messages</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        When a <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">PUBLISH</code> command is received, the Broker concurrently writes to all subscriber TCP sockets via <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">go DeliverMessage</code>, using a <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">sync.WaitGroup</code> to coordinate delivery without blocking the publisher.
      </p>

      <h2 id="pattern-matching" className="text-2xl font-semibold text-white mb-4 mt-8">Pattern Matching</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">PSUBSCRIBE</code> uses namespace-based wildcard matching. It splits channels by <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">:</code> and <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">.</code> separators. (Note: This is NOT standard glob pattern matching).
      </p>

      <h2 id="examples" className="text-2xl font-semibold text-white mb-4 mt-8">Examples</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="terminal-window">
          <div className="terminal-header">
            <div className="terminal-dot bg-red-500"></div>
            <div className="terminal-dot bg-yellow-500"></div>
            <div className="terminal-dot bg-green-500"></div>
            <span className="ml-2 text-xs text-[#5e5e73]">Terminal 1 (Subscriber)</span>
          </div>
          <div className="terminal-body font-mono text-sm p-4 text-[#f1f1f4]">
            <div><span className="prompt text-[#7c3aed]">&gt;</span> <span className="command text-[#a78bfa]">SUBSCRIBE</span> news sports</div>
            <div className="output text-[#9898ab]">1) "subscribe"</div>
            <div className="output text-[#9898ab]">2) "news"</div>
            <div className="output text-[#9898ab]">3) (integer) 1</div>
            <div className="output text-[#9898ab]">1) "subscribe"</div>
            <div className="output text-[#9898ab]">2) "sports"</div>
            <div className="output text-[#9898ab] mb-4">3) (integer) 2</div>
            
            <div className="comment text-[#5e5e73]"># Waiting for messages...</div>
            <div className="output text-[#9898ab] mt-2">1) "message"</div>
            <div className="output text-[#9898ab]">2) "news"</div>
            <div className="output text-[#9898ab]">3) "Hello World"</div>
          </div>
        </div>

        <div className="terminal-window">
          <div className="terminal-header">
            <div className="terminal-dot bg-red-500"></div>
            <div className="terminal-dot bg-yellow-500"></div>
            <div className="terminal-dot bg-green-500"></div>
            <span className="ml-2 text-xs text-[#5e5e73]">Terminal 2 (Publisher)</span>
          </div>
          <div className="terminal-body font-mono text-sm p-4 text-[#f1f1f4]">
            <div><span className="prompt text-[#7c3aed]">&gt;</span> <span className="command text-[#a78bfa]">PUBLISH</span> news "Hello World"</div>
            <div className="output text-[#9898ab]">(integer) 1</div>
          </div>
        </div>
      </div>
    </div>
  );
}
