import CodeBlock from '../../../components/docs/CodeBlock';

export default function RoxaiConfiguration() {
  const envExample = `OPENAI_API_KEY=your_key_here
ANTHROPIC_API_KEY=your_key_here
GEMINI_API_KEY=your_key_here
GROQ_API_KEY=your_key_here`;

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">RoxAI Configuration</h1>
        <p className="text-xl text-[#9898ab]">Setting up and configuring the AI engine</p>
      </div>
      <div className="section-divider" />

      <h2 id="model-json" className="text-2xl font-semibold text-white mb-4">model.json Configuration</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        The primary configuration file for RoxAI is <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">model.json</code>. The location of this file is managed by <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">abstractor.ConfigPath()</code>. This file persistently stores the active configuration state.
      </p>
      <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl mb-8">
        <p className="text-[#5e5e73] mb-2">The structure includes:</p>
        <ul className="list-disc list-inside text-[#5e5e73]">
          <li>Provider Name</li>
          <li>Model Name</li>
          <li>API Endpoint URL</li>
          <li>API Key</li>
        </ul>
      </div>

      <h2 id="environment-variables" className="text-2xl font-semibold text-white mb-4">Environment Variables</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        For cloud providers, API keys are typically loaded from the environment variables, usually defined in a <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">.env</code> file (see <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">.env.example</code>).
      </p>
      <div className="mb-8">
        <CodeBlock code={envExample} language="env" title=".env" />
      </div>

      <h2 id="port-assignments" className="text-2xl font-semibold text-white mb-4">Network & Port Assignments</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        RoxAI binds to specific ports independently of the core RoxKV database ports:
      </p>
      <ul className="space-y-2 text-[#9898ab] mb-8">
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">6970</code> - TCP Interface for raw socket connections (AI Chat CLI).</li>
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">6972</code> - HTTP Web Server serving REST API endpoints for the dashboard.</li>
      </ul>

      <h2 id="ollama-setup" className="text-2xl font-semibold text-white mb-4">Ollama Local Setup</h2>
      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#9898ab]">By default, RoxAI expects Ollama to be running locally or via Docker to serve the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">llama3.2:3b</code> model. Ensure Ollama is running and accessible on port 11434 if you intend to use the default local configuration.</p>
      </div>
    </div>
  );
}
