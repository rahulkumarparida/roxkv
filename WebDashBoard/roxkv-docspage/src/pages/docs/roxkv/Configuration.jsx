export default function Configuration() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Configuration</h1>
        <p className="text-xl text-[#9898ab]">Configure RoxKV ports, environment, and persistence</p>
      </div>
      <div className="section-divider" />
      
      <h2 id="port-mapping" className="text-2xl font-bold text-white mt-12 mb-4">Port Mapping</h2>
      <p className="text-[#f1f1f4] mb-4">
        RoxKV exposes multiple services on different ports to handle native CLI, AI interactions, events, and Redis compatibility.
      </p>
      <div className="overflow-x-auto rounded-xl border border-[#1e1e2e] bg-[#16161f]">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="border-b border-[#1e1e2e] bg-[rgba(124,58,237,0.05)]">
              <th className="py-3 px-4 text-sm font-semibold text-[#a78bfa]">Port</th>
              <th className="py-3 px-4 text-sm font-semibold text-[#a78bfa]">Service</th>
              <th className="py-3 px-4 text-sm font-semibold text-[#a78bfa]">Description</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-[#1e1e2e]">
            <tr className="hover:bg-[rgba(255,255,255,0.02)] transition-colors">
              <td className="py-3 px-4 text-sm font-mono text-white">6969</td>
              <td className="py-3 px-4 text-sm text-[#f1f1f4]">Native CLI</td>
              <td className="py-3 px-4 text-sm text-[#9898ab]">Default port for RoxKV's native CLI connection.</td>
            </tr>
            <tr className="hover:bg-[rgba(255,255,255,0.02)] transition-colors">
              <td className="py-3 px-4 text-sm font-mono text-white">6970</td>
              <td className="py-3 px-4 text-sm text-[#f1f1f4]">AI TCP Server</td>
              <td className="py-3 px-4 text-sm text-[#9898ab]">Accepts TCP connections for the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">roxai&gt;</code> prompt.</td>
            </tr>
            <tr className="hover:bg-[rgba(255,255,255,0.02)] transition-colors">
              <td className="py-3 px-4 text-sm font-mono text-white">6971</td>
              <td className="py-3 px-4 text-sm text-[#f1f1f4]">Events HTTP</td>
              <td className="py-3 px-4 text-sm text-[#9898ab]">HTTP server for Server-Sent Events (SSE).</td>
            </tr>
            <tr className="hover:bg-[rgba(255,255,255,0.02)] transition-colors">
              <td className="py-3 px-4 text-sm font-mono text-white">6972</td>
              <td className="py-3 px-4 text-sm text-[#f1f1f4]">AI Chat HTTP</td>
              <td className="py-3 px-4 text-sm text-[#9898ab]">REST API endpoints for AI chat interactions.</td>
            </tr>
            <tr className="hover:bg-[rgba(255,255,255,0.02)] transition-colors">
              <td className="py-3 px-4 text-sm font-mono text-white">6973</td>
              <td className="py-3 px-4 text-sm text-[#f1f1f4]">RESP Server</td>
              <td className="py-3 px-4 text-sm text-[#9898ab]">Redis-compatible server (use this for redis-cli).</td>
            </tr>
          </tbody>
        </table>
      </div>

      <h2 id="environment" className="text-2xl font-bold text-white mt-12 mb-4">Environment Variables</h2>
      <p className="text-[#f1f1f4] mb-4">
        You can configure environment variables using the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">.env</code> file. Create a copy of <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">.env.example</code> to get started.
      </p>

      <h2 id="docker" className="text-2xl font-bold text-white mt-12 mb-4">Docker Configuration</h2>
      <p className="text-[#f1f1f4] mb-4">
        The <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">docker-compose.yml</code> file orchestrates the services. It maps all required ports and mounts volumes for persistence.
      </p>
      <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl mb-4">
        <ul className="list-disc pl-5 space-y-2 text-[#9898ab]">
          <li><strong className="text-white">Services:</strong> The core RoxKV server and an integrated Ollama container.</li>
          <li><strong className="text-white">Volumes:</strong> Maps the local <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">.roxkv</code> directory to ensure snapshots and data survive container restarts.</li>
          <li><strong className="text-white">AI Provider:</strong> Uses <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">model.json</code> to determine the AI provider configuration.</li>
        </ul>
      </div>

      <h2 id="data-persistence" className="text-2xl font-bold text-white mt-12 mb-4">Data Directory</h2>
      <p className="text-[#f1f1f4] mb-4">
        RoxKV stores its persistent data inside the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">.roxkv/</code> directory. This directory contains snapshots of your key-value store, enabling the system to recover its state upon restarting.
      </p>
      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#9898ab] text-sm">
          <strong>Note:</strong> Ensure the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">.roxkv</code> directory is included in your backups if you want to preserve your database snapshots securely.
        </p>
      </div>
    </div>
  );
}
