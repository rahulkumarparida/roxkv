import CodeBlock from '../../../components/docs/CodeBlock';
import FlowDiagram from '../../../components/docs/FlowDiagram';
import InfoTable from '../../../components/docs/InfoTable';

export default function RoxaiOverview() {
  const steps = [
    { label: 'User Query', highlight: false, sub: 'TCP or HTTP' },
    { label: 'Master Agent', highlight: true, sub: 'Orchestration' },
    { label: 'Semantic Router', highlight: false, sub: 'Selects tools' },
    { label: 'LLM', highlight: true, sub: 'Decides action' },
    { label: 'Tool Calls', highlight: false, sub: 'Direct Go funcs' },
    { label: 'RoxKV', highlight: true, sub: 'Database operations' },
    { label: 'Synthesis', highlight: false, sub: 'Natural language' },
    { label: 'Response', highlight: true, sub: 'Back to client' },
  ];

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">RoxAI Overview</h1>
        <p className="text-xl text-[#9898ab]">AI Orchestration Layer built on top of RoxKV</p>
      </div>
      <div className="section-divider" />

      <h2 id="what-is-roxai" className="text-2xl font-semibold text-white mb-4">What is RoxAI?</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        RoxAI is not a simple chatbot. It is a sophisticated tool-calling agent system designed to translate natural language into complex database operations. Built directly on top of RoxKV, it acts as an intelligent orchestration layer.
      </p>
      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#9898ab]"><strong>Note:</strong> RoxAI processes natural language queries and automatically determines the appropriate RoxKV operations to execute, making database management intuitive and conversational.</p>
      </div>

      <h2 id="capabilities" className="text-2xl font-semibold text-white mb-4">Capabilities</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">Master Agent Orchestration</h3>
          <p className="text-[#5e5e73]">Centralized agent managing the flow of intents, tool execution, and response synthesis.</p>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">40+ Tools</h3>
          <p className="text-[#5e5e73]">Extensive toolset across 5 categories to interact with all aspects of RoxKV.</p>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">Multi-LLM Support</h3>
          <p className="text-[#5e5e73]">Supports 5 major LLM providers: Ollama, OpenAI, Anthropic, Gemini, and Groq.</p>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">Semantic Routing & Failover</h3>
          <p className="text-[#5e5e73]">Intelligent tool routing and automatic provider failover for high availability.</p>
        </div>
      </div>

      <h2 id="high-level-flow" className="text-2xl font-semibold text-white mb-4">High-Level Flow</h2>
      <div className="mb-8 p-6 bg-[#0a0a0f] rounded-xl border border-[#1e1e2e]">
        <FlowDiagram steps={steps} direction="horizontal" />
      </div>

      <h2 id="interfaces-and-configuration" className="text-2xl font-semibold text-white mb-4">Interfaces & Configuration</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        RoxAI provides two primary interfaces for interaction:
      </p>
      <ul className="list-disc list-inside text-[#9898ab] space-y-2 mb-6">
        <li><strong>TCP:</strong> Available on port <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">6970</code>, providing a shell-like <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">roxai&gt;</code> prompt.</li>
        <li><strong>HTTP REST:</strong> Available on port <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">6972</code>, supporting integration with web dashboards.</li>
      </ul>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        By default, RoxAI is configured to use <strong>Ollama</strong> with the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">llama3.2:3b</code> model, connecting to <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">localhost:11434</code>.
      </p>
    </div>
  );
}
