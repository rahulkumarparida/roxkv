import React from 'react';

export default function TcpServer() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">TCP Server</h1>
        <p className="text-xl text-[#9898ab]">Network layer handling connections and fragmentation.</p>
      </div>
      
      <div className="section-divider" />

      <h2 id="connection-lifecycle" className="text-2xl font-semibold text-white mb-4">Connection Lifecycle</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        The RoxKV server binds to port <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">6973</code>. The lifecycle follows a standard Go networking pattern:
      </p>
      <ol className="list-decimal list-inside text-[#9898ab] space-y-2 mb-6">
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">net.Listen("tcp", ":6973")</code> initializes the listener.</li>
        <li>Incoming connections are accepted via <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">Accept()</code>.</li>
        <li>Each connection is wrapped in a <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">Client</code> struct using <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">utils.NewClient</code>.</li>
        <li>A new goroutine <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">go HandleRespConnections</code> is spawned to handle the client.</li>
      </ol>

      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#f1f1f4]">A global mutex tracks the active connection count to limit maximum concurrent connections.</p>
      </div>

      <h2 id="client-struct" className="text-2xl font-semibold text-white mb-4 mt-8">The Client Struct</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        The internal <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">Client</code> struct maintains state for each connection:
      </p>
      <ul className="list-disc list-inside text-[#9898ab] space-y-2 mb-6">
        <li><strong>Conn:</strong> The underlying <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">net.Conn</code>.</li>
        <li><strong>Buffer:</strong> A <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">[]byte</code> slice for buffering partial reads (TCP fragmentation).</li>
        <li><strong>Mu:</strong> A <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">sync.RWMutex</code> to synchronize state changes.</li>
        <li><strong>Mode:</strong> Tracks if the client has entered subscriber mode.</li>
      </ul>

      <h2 id="tcp-fragmentation" className="text-2xl font-semibold text-white mb-4 mt-8">Handling TCP Fragmentation</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        Because TCP is a stream protocol, RESP commands might arrive in multiple fragments. The <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">ReadAndHandleConnection</code> loop continuously reads from the socket into the client's <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">Buffer</code>. It then calls <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">HandleMultipleCommands</code> which attempts to parse full RESP frames. Incomplete frames are left in the buffer until more data arrives.
      </p>

      <h2 id="subscriber-mode" className="text-2xl font-semibold text-white mb-4 mt-8">Subscriber Mode Enforcement</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        Once a client executes a <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">SUBSCRIBE</code> or <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">PSUBSCRIBE</code> command, its <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">Mode</code> is switched. In subscriber mode, the server enforces strict command validation—standard storage commands like <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">GET</code> and <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">SET</code> are blocked until the client unsubscribes from all channels.
      </p>
    </div>
  );
}
