export default function Dashboard() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Dashboard & UI Integration</h1>
        <p className="text-xl text-[#9898ab]">Interacting with RoxAI via TCP and Web Interfaces</p>
      </div>
      <div className="section-divider" />

      <h2 id="integration-points" className="text-2xl font-semibold text-white mb-4">Primary Integration Points</h2>
      <p className="text-[#9898ab] leading-relaxed mb-6">
        RoxAI exposes two main avenues for interaction, catering to both command-line users and web interfaces.
      </p>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">1. TCP CLI (Port 6970)</h3>
          <p className="text-[#5e5e73]">Managed by <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">ChatServer.go</code>. Provides a raw TCP socket connection featuring a shell-like <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">roxai&gt;</code> prompt. Ideal for direct terminal access and scripting.</p>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">2. HTTP API (Port 6972)</h3>
          <p className="text-[#5e5e73]">Managed by <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">ChatWebServer.go</code>. Provides RESTful endpoints designed specifically to power the React-based Web Dashboard.</p>
        </div>
      </div>

      <h2 id="rest-endpoints" className="text-2xl font-semibold text-white mb-4">REST API Endpoints</h2>
      <ul className="space-y-3 text-[#9898ab] mb-8">
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">POST /api/chat/</code> - Submits a query to the AI.</li>
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">GET /api/chat/config</code> - Retrieves current LLM config.</li>
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">POST /api/chat/config</code> - Saves new configuration settings.</li>
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">GET /providers</code> - Lists available providers.</li>
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">GET /providers/&#123;provider&#125;/models</code> - Lists models for a provider.</li>
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">GET /providers/&#123;provider&#125;/config</code> - Fetches provider settings.</li>
      </ul>

      <h2 id="connection-hijacking" className="text-2xl font-semibold text-white mb-4">HTTP Connection Hijacking</h2>
      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#9898ab]">
          For real-time streaming of chat responses over the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">POST /api/chat/</code> endpoint, RoxAI utilizes HTTP connection hijacking <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">w.(http.Hijacker)</code>. This takes over the raw underlying TCP connection from the Go HTTP server to stream plain text responses back to the dashboard efficiently, avoiding standard buffered HTTP request/response cycles.
        </p>
      </div>

      <h2 id="react-dashboard" className="text-2xl font-semibold text-white mb-4">The Web Dashboard</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        The official UI for RoxAI and RoxKV is a React application located at <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">WebDashBoard/roxkv-dashboard/</code>. It consumes the REST API on port 6972 to provide a visual chat interface and configuration management screens.
      </p>
    </div>
  );
}
