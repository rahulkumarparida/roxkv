import { Link } from 'react-router-dom';
import CodeBlock from '../../components/docs/CodeBlock';
import FlowDiagram from '../../components/docs/FlowDiagram';
import InfoTable from '../../components/docs/InfoTable';

export default function HighLevelArchitecture() {
  const stackSteps = [
    { label: 'RoxAI', highlight: true, sub: 'AI Orchestrator' },
    { label: 'Master Agent', highlight: false, sub: 'Coordination' },
    { label: 'Tool/LLM Layer', highlight: false, sub: 'Abstractions' },
    { label: 'RoxKV', highlight: true, sub: 'Core Data Engine' },
    { label: 'Storage / TCP / RESP', highlight: false, sub: 'Internal Branches' },
    { label: 'Redis Clients', highlight: true, sub: 'External Access' },
  ];

  const ports = [
    { port: '6969', service: 'Native CLI', description: 'RoxKV native command-line interface' },
    { port: '6970', service: 'RoxAI TCP', description: 'AI agent CLI access (roxai> prompt)' },
    { port: '6971', service: 'Events HTTP', description: 'Server-Sent Events (SSE) stream' },
    { port: '6972', service: 'RoxAI HTTP', description: 'AI Chat REST API & Web Dashboard' },
    { port: '6973', service: 'RESP Server', description: 'Redis-compatible protocol (use with redis-cli)' },
  ];

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">High-Level Architecture</h1>
        <p className="text-xl text-[#9898ab]">How RoxAI and RoxKV integrate into a unified system</p>
      </div>
      <div className="section-divider" />

      <h2 id="system-stack" className="text-2xl font-semibold text-white mb-4">Unified System Stack</h2>
      <p className="text-[#9898ab] leading-relaxed mb-6">
        RoxKV and RoxAI are built together as a single cohesive Go binary. RoxAI sits on top of RoxKV, interacting with its internals via direct function calls, while RoxKV continues to serve external data clients simultaneously.
      </p>

      <div className="mb-8 p-6 bg-[#0a0a0f] rounded-xl border border-[#1e1e2e] flex justify-center">
        <FlowDiagram steps={stackSteps} direction="vertical" />
      </div>

      <h2 id="core-capabilities" className="text-2xl font-semibold text-white mb-4">Core Capabilities</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-4">RoxKV (Data Layer)</h3>
          <ul className="list-disc list-inside text-[#5e5e73] space-y-2">
            <li>High-performance in-memory key-value store</li>
            <li>Full RESP protocol compatibility</li>
            <li>Native Pub/Sub messaging system</li>
            <li>Snapshot-based persistence</li>
            <li>TTL support and background cleanup</li>
            <li>28+ supported database commands</li>
          </ul>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-[#a78bfa] mb-4">RoxAI (Intelligence Layer)</h3>
          <ul className="list-disc list-inside text-[#5e5e73] space-y-2">
            <li>Master Agent intent orchestration</li>
            <li>42+ native tools querying internal state</li>
            <li>Support for 5 major LLM providers</li>
            <li>Semantic tool routing for context optimization</li>
            <li>Automatic provider failover and retries</li>
            <li>Natural language to DB operation translation</li>
          </ul>
        </div>
      </div>

      <h2 id="process-communication" className="text-2xl font-semibold text-white mb-4">Zero-Overhead Communication</h2>
      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#9898ab]">Communication between RoxAI and RoxKV happens via direct Go function calls in the same process memory space. There are no HTTP overheads, no serialization bottlenecks, and no network latency when the AI queries or mutates the database state.</p>
      </div>

      <h2 id="port-overview" className="text-2xl font-semibold text-white mb-4">Network Topology & Ports</h2>
      <div className="mb-8">
        <InfoTable 
          headers={['Port', 'Service', 'Description']} 
          rows={ports.map(p => [p.port, p.service, p.description])} 
        />
      </div>
      
      <div className="flex gap-4">
        <Link to="/docs/roxkv" className="text-[#a78bfa] hover:text-white underline">Explore RoxKV Docs</Link>
        <Link to="/docs/roxai" className="text-[#a78bfa] hover:text-white underline">Explore RoxAI Docs</Link>
      </div>
    </div>
  );
}
