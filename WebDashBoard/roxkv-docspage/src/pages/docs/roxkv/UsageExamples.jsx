export default function UsageExamples() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Usage Examples</h1>
        <p className="text-xl text-[#9898ab]">Learn how to integrate RoxKV with various languages and tools</p>
      </div>
      <div className="section-divider" />
      
      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#f1f1f4] font-medium">
          Important: Remember to connect to port <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">6973</code> (the RoxKV RESP port) instead of the default Redis port (6379).
        </p>
      </div>

      <h2 id="basic" className="text-2xl font-bold text-white mt-12 mb-4">Basic Operations</h2>
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
          <div><span className="prompt">127.0.0.1:6973&gt;</span> <span className="command">SET mykey "hello world"</span></div>
          <div><span className="output">OK</span></div>
          <div><span className="prompt">127.0.0.1:6973&gt;</span> <span className="command">GET mykey</span></div>
          <div><span className="output">"hello world"</span></div>
          <div><span className="prompt">127.0.0.1:6973&gt;</span> <span className="command">EXISTS mykey</span></div>
          <div><span className="output">(integer) 1</span></div>
          <div><span className="prompt">127.0.0.1:6973&gt;</span> <span className="command">KEYS *</span></div>
          <div><span className="output">1) "mykey"</span></div>
          <div><span className="prompt">127.0.0.1:6973&gt;</span> <span className="command">DEL mykey</span></div>
          <div><span className="output">(integer) 1</span></div>
        </div>
      </div>

      <h2 id="arithmetic" className="text-2xl font-bold text-white mt-12 mb-4">Arithmetic Operations</h2>
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
          <div><span className="prompt">127.0.0.1:6973&gt;</span> <span className="command">SET counter 10</span></div>
          <div><span className="output">OK</span></div>
          <div><span className="prompt">127.0.0.1:6973&gt;</span> <span className="command">INCR counter</span></div>
          <div><span className="output">(integer) 11</span></div>
          <div><span className="prompt">127.0.0.1:6973&gt;</span> <span className="command">DECRBY counter 5</span></div>
          <div><span className="output">(integer) 6</span></div>
        </div>
      </div>

      <h2 id="ttl" className="text-2xl font-bold text-white mt-12 mb-4">TTL and Expiration</h2>
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
          <div><span className="prompt">127.0.0.1:6973&gt;</span> <span className="command">SET session abc123</span></div>
          <div><span className="output">OK</span></div>
          <div><span className="prompt">127.0.0.1:6973&gt;</span> <span className="command">EXPIRE session 300</span></div>
          <div><span className="output">(integer) 1</span></div>
          <div><span className="prompt">127.0.0.1:6973&gt;</span> <span className="command">TTL session</span></div>
          <div><span className="output">(integer) 298</span></div>
          <div><span className="prompt">127.0.0.1:6973&gt;</span> <span className="command">PERSIST session</span></div>
          <div><span className="output">(integer) 1</span></div>
        </div>
      </div>

      <h2 id="pubsub" className="text-2xl font-bold text-white mt-12 mb-4">Pub/Sub Implementation</h2>
      <p className="text-[#f1f1f4] mb-4">Requires two separate terminal sessions connecting to RoxKV.</p>
      
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div className="terminal-window">
          <div className="terminal-header">
            <div className="flex gap-2">
              <div className="terminal-dot bg-red-500"></div>
              <div className="terminal-dot bg-yellow-500"></div>
              <div className="terminal-dot bg-green-500"></div>
            </div>
            <span className="ml-4 text-xs text-[#9898ab]">Terminal 1 (Subscriber)</span>
          </div>
          <div className="terminal-body bg-[#16161f] p-4 font-mono text-sm">
            <div><span className="prompt">127.0.0.1:6973&gt;</span> <span className="command">SUBSCRIBE news</span></div>
            <div><span className="output">Reading messages... (press Ctrl-C to quit)</span></div>
            <div><span className="output">1) "subscribe"</span></div>
            <div><span className="output">2) "news"</span></div>
            <div><span className="output">3) (integer) 1</span></div>
            <br />
            <div><span className="comment"># After publisher sends message:</span></div>
            <div><span className="output">1) "message"</span></div>
            <div><span className="output">2) "news"</span></div>
            <div><span className="output">3) "Breaking news!"</span></div>
          </div>
        </div>
        
        <div className="terminal-window">
          <div className="terminal-header">
            <div className="flex gap-2">
              <div className="terminal-dot bg-red-500"></div>
              <div className="terminal-dot bg-yellow-500"></div>
              <div className="terminal-dot bg-green-500"></div>
            </div>
            <span className="ml-4 text-xs text-[#9898ab]">Terminal 2 (Publisher)</span>
          </div>
          <div className="terminal-body bg-[#16161f] p-4 font-mono text-sm">
            <div><span className="prompt">127.0.0.1:6973&gt;</span> <span className="command">PUBLISH news "Breaking news!"</span></div>
            <div><span className="output">(integer) 1</span></div>
          </div>
        </div>
      </div>

      <h2 id="clients" className="text-2xl font-bold text-white mt-12 mb-4">Client Examples</h2>
      
      <div className="space-y-6">
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-xl font-bold text-[#a78bfa] mb-4">Go (go-redis)</h3>
          <div className="terminal-body bg-[#0a0a0f] p-4 font-mono text-sm rounded">
            <div className="text-[#f1f1f4]">rdb := redis.NewClient(&amp;redis.Options{'{'}</div>
            <div className="text-[#f1f1f4] ml-4">Addr: <span className="text-green-400">"localhost:6973"</span>,</div>
            <div className="text-[#f1f1f4] ml-4">Password: <span className="text-green-400">""</span>, <span className="text-[#5e5e73]">// No password set</span></div>
            <div className="text-[#f1f1f4] ml-4">DB: <span className="text-yellow-400">0</span>,</div>
            <div className="text-[#f1f1f4]">{'}'})</div>
            <br/>
            <div className="text-[#f1f1f4]">err := rdb.Set(ctx, <span className="text-green-400">"key"</span>, <span className="text-green-400">"value"</span>, <span className="text-yellow-400">0</span>).Err()</div>
          </div>
        </div>

        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-xl font-bold text-[#a78bfa] mb-4">Python (redis-py)</h3>
          <div className="terminal-body bg-[#0a0a0f] p-4 font-mono text-sm rounded">
            <div className="text-[#f1f1f4]"><span className="text-[#7c3aed]">import</span> redis</div>
            <br/>
            <div className="text-[#f1f1f4]">r = redis.Redis(host=<span className="text-green-400">'localhost'</span>, port=<span className="text-yellow-400">6973</span>, decode_responses=<span className="text-[#7c3aed]">True</span>)</div>
            <div className="text-[#f1f1f4]">r.set(<span className="text-green-400">'foo'</span>, <span className="text-green-400">'bar'</span>)</div>
            <div className="text-[#f1f1f4]">val = r.get(<span className="text-green-400">'foo'</span>)</div>
          </div>
        </div>

        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-xl font-bold text-[#a78bfa] mb-4">Node.js (ioredis)</h3>
          <div className="terminal-body bg-[#0a0a0f] p-4 font-mono text-sm rounded">
            <div className="text-[#f1f1f4]"><span className="text-[#7c3aed]">const</span> Redis = require(<span className="text-green-400">"ioredis"</span>);</div>
            <br/>
            <div className="text-[#f1f1f4]"><span className="text-[#7c3aed]">const</span> redis = <span className="text-[#7c3aed]">new</span> Redis(<span className="text-yellow-400">6973</span>); <span className="text-[#5e5e73]">// Connect to 127.0.0.1:6973</span></div>
            <div className="text-[#f1f1f4]">await redis.set(<span className="text-green-400">"hello"</span>, <span className="text-green-400">"world"</span>);</div>
            <div className="text-[#f1f1f4]">console.log(await redis.get(<span className="text-green-400">"hello"</span>));</div>
          </div>
        </div>
      </div>
    </div>
  );
}
