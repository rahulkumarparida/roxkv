import CodeBlock from '../../../components/docs/CodeBlock';
import FlowDiagram from '../../../components/docs/FlowDiagram';
import InfoTable from '../../../components/docs/InfoTable';

export default function RoxaiArchitecture() {
  const flowSteps = [
    { label: 'User Query', highlight: false, sub: 'Received via TCP or HTTP' },
    { label: 'Session Retrieval', highlight: false, sub: 'Keyed by masteragent:<user_id>' },
    { label: 'Semantic Tool Router', highlight: true, sub: 'Selects top 5 relevant tools' },
    { label: 'Intent LLM Call', highlight: true, sub: 'Determines required tool calls' },
    { label: 'Tool Execution', highlight: false, sub: 'Switch-based dispatch to Go modules' },
    { label: 'Result Formatting', highlight: false, sub: 'Formats tools output as JSON/text' },
    { label: 'Synthesis LLM Call', highlight: true, sub: 'Converts raw results into natural language' },
    { label: 'Response', highlight: false, sub: 'Streamed back to the client' },
  ];

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">RoxAI Architecture</h1>
        <p className="text-xl text-[#9898ab]">Internal design and request flow of the AI orchestrator</p>
      </div>
      <div className="section-divider" />

      <h2 id="request-flow" className="text-2xl font-semibold text-white mb-4">Request Flow</h2>
      <p className="text-[#9898ab] leading-relaxed mb-6">
        When a user submits a query to RoxAI, it goes through a sophisticated pipeline orchestrated primarily by the Master Agent.
      </p>
      
      <div className="mb-8 p-6 bg-[#0a0a0f] rounded-xl border border-[#1e1e2e] flex justify-center">
        <FlowDiagram steps={flowSteps} direction="vertical" />
      </div>

      <h2 id="two-phase-llm" className="text-2xl font-semibold text-white mb-4">Two-Phase LLM Call Pattern</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        RoxAI implements a robust two-phase execution pattern for handling complex operations:
      </p>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-[#a78bfa] mb-2">Phase 1: Intent Recognition</h3>
          <p className="text-[#5e5e73]">The Master Agent provides the user's query and a dynamically routed list of tools to the LLM. The LLM's goal is to determine the intent and output the exact tool calls with required parameters to fulfill the request.</p>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-[#a78bfa] mb-2">Phase 2: Synthesis</h3>
          <p className="text-[#5e5e73]">After executing the tools internally, the results (often raw JSON or structural data) are injected back into the context. A second LLM call synthesizes this raw data into a coherent, human-readable response.</p>
        </div>
      </div>

      <h2 id="component-architecture" className="text-2xl font-semibold text-white mb-4">Component Architecture</h2>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">abstractor/</h3>
          <p className="text-[#5e5e73]">Provides a unified interface for multiple LLM providers, normalizing responses and handling failovers.</p>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">router/</h3>
          <p className="text-[#5e5e73]">Implements semantic tool routing to select the most relevant tools for a given query, saving context window space.</p>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">registry/ & Master/</h3>
          <p className="text-[#5e5e73]">The tool registry maintains definitions, while the Master Agent handles tool dispatch and session coordination.</p>
        </div>
      </div>
    </div>
  );
}
