import CodeBlock from '../../../components/docs/CodeBlock';
import InfoTable from '../../../components/docs/InfoTable';

export default function Providers() {
  const providers = [
    { name: 'Ollama', type: 'Local', defaultModel: 'llama3.2:3b', description: 'Default local provider, no API key required.' },
    { name: 'OpenAI', type: 'Cloud', defaultModel: 'gpt-4o-mini', description: 'Requires OPENAI_API_KEY.' },
    { name: 'Anthropic', type: 'Cloud', defaultModel: 'claude-3-5-sonnet', description: 'Requires ANTHROPIC_API_KEY.' },
    { name: 'Gemini', type: 'Cloud', defaultModel: 'gemini-1.5-flash', description: 'Requires GEMINI_API_KEY.' },
    { name: 'Groq', type: 'Cloud', defaultModel: 'llama-3.1-70b', description: 'Requires GROQ_API_KEY for fast inference.' },
  ];

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Providers & Configuration</h1>
        <p className="text-xl text-[#9898ab]">Managing LLM endpoints and runtime settings</p>
      </div>
      <div className="section-divider" />

      <h2 id="supported-providers" className="text-2xl font-semibold text-white mb-4">Supported Providers</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        RoxAI supports multiple LLM providers out of the box, allowing you to choose between local privacy and cloud performance. By default, it is configured to use <strong>Ollama</strong> running locally with <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">llama3.2:3b</code> at <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">http://localhost:11434</code>.
      </p>
      <div className="mb-8">
        <InfoTable 
          headers={['Provider', 'Type', 'Default Model', 'Notes']} 
          rows={providers.map(p => [p.name, p.type, p.defaultModel, p.description])} 
        />
      </div>

      <h2 id="configuration-storage" className="text-2xl font-semibold text-white mb-4">Configuration Storage</h2>
      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#9898ab]">The active provider configuration is stored persistently in a <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">model.json</code> file, managed by the abstractor package. This file holds the selected provider name, model name, API endpoint, and required API keys.</p>
      </div>

      <h2 id="rest-api-management" className="text-2xl font-semibold text-white mb-4">REST API Management</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        RoxAI allows dynamic, runtime provider switching without restarting the server via its REST API endpoints:
      </p>
      <ul className="space-y-3 text-[#9898ab] mb-8">
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">GET /api/chat/config</code> - Fetch current active configuration</li>
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">POST /api/chat/config</code> - Update active configuration</li>
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">GET /providers</code> - List all available/supported providers</li>
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">GET /providers/&#123;provider&#125;/models</code> - List available models for a specific provider</li>
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">GET /providers/&#123;provider&#125;/config</code> - Get credentials for a specific provider</li>
      </ul>
    </div>
  );
}
