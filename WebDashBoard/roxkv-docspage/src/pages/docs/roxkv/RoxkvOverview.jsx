import React from 'react';
import CodeBlock from '../../../components/docs/CodeBlock';
import FlowDiagram from '../../../components/docs/FlowDiagram';
import InfoTable from '../../../components/docs/InfoTable';

export default function RoxkvOverview() {
  const flowSteps = [
    { label: 'Client' },
    { label: 'TCP Server', sub: 'Port 6973' },
    { label: 'RESP Decoder' },
    { label: 'Command Parser' },
    { label: 'Executor' },
    { label: 'Store', sub: 'In-Memory', highlight: true },
    { label: 'RESP Encoder' },
    { label: 'Client' }
  ];

  const capabilitiesHeaders = ['Feature', 'Description'];
  const capabilitiesRows = [
    ['Commands', '28+ Redis commands implemented'],
    ['Pub/Sub', 'Channel and Pattern-based message brokering'],
    ['TTL & Expiration', 'Key expiration with background cleanup worker'],
    ['Persistence', 'JSON snapshot saving and loading'],
    ['Network', 'TCP server on port 6973'],
    ['Binary Safety', 'Values stored as raw []byte']
  ];

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">RoxKV Overview</h1>
        <p className="text-xl text-[#9898ab]">High-performance in-memory key-value store with RESP compatibility.</p>
      </div>
      
      <div className="section-divider" />

      <h2 id="introduction" className="text-2xl font-semibold text-white mb-4">Introduction</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        RoxKV is an in-memory Key-Value store written in Go. It implements the Redis Serialization Protocol (RESP), making it compatible with existing Redis clients, while focusing on simplicity and specific use cases.
      </p>

      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#f1f1f4]"><strong>Note:</strong> RoxKV is NOT a full Redis replacement. It is a subset implementation designed for specific application needs.</p>
      </div>

      <h2 id="core-capabilities" className="text-2xl font-semibold text-white mb-4 mt-8">Core Capabilities</h2>
      <InfoTable headers={capabilitiesHeaders} rows={capabilitiesRows} />

      <h2 id="architecture-overview" className="text-2xl font-semibold text-white mb-4 mt-8">Architecture Overview</h2>
      <p className="text-[#9898ab] leading-relaxed mb-6">
        RoxKV uses a goroutine-per-connection model and <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">sync.RWMutex</code> for concurrency control.
      </p>
      
      <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
        <h3 className="text-lg font-medium text-white mb-6">Request Lifecycle</h3>
        <FlowDiagram steps={flowSteps} direction="horizontal" />
      </div>

    </div>
  );
}
