import React from 'react';
import FlowDiagram from '../../../components/docs/FlowDiagram';
import InfoTable from '../../../components/docs/InfoTable';

export default function RoxkvArchitecture() {
  const flowSteps = [
    { label: 'Client Connection', sub: 'TCP Port 6973' },
    { label: 'RESP Decoder', sub: 'DecodeArrayString' },
    { label: 'Command Parser', sub: 'RedisInput (Cmd + Args)' },
    { label: 'Mode Validation', sub: 'Check Subscriber State' },
    { label: 'Switch Dispatch', sub: 'Command Routing' },
    { label: 'Command Handler', sub: 'Execute Logic' },
    { label: 'Store Operation', sub: 'MemoryAlloc + sync.RWMutex', highlight: true },
    { label: 'RESP Encoder', sub: 'Format Response' },
    { label: 'TCP Write', sub: 'Send to Client' }
  ];

  const portHeaders = ['Port', 'Service'];
  const portRows = [
    ['6969', 'CLI'],
    ['6970', 'AI TCP'],
    ['6971', 'Events HTTP'],
    ['6972', 'AI Chat HTTP'],
    ['6973', 'RESP (RoxKV)']
  ];

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Architecture & Data Flow</h1>
        <p className="text-xl text-[#9898ab]">Internal structure and request lifecycle of RoxKV.</p>
      </div>
      
      <div className="section-divider" />

      <h2 id="data-flow" className="text-2xl font-semibold text-white mb-4">Request Data Flow</h2>
      <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl mb-8">
        <FlowDiagram steps={flowSteps} direction="vertical" />
      </div>

      <h2 id="module-layout" className="text-2xl font-semibold text-white mb-4 mt-8">Module Architecture</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">server/</h3>
          <p className="text-[#9898ab]">TCP server, connection acceptance, and goroutine spawning.</p>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">internal/resp-parser/</h3>
          <p className="text-[#9898ab]">decoder.go (RESP decoding), encoder.go (RESP encoding), parser.go (command dispatch), executor.go (handlers).</p>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">internal/store/</h3>
          <p className="text-[#9898ab]">MemoryAlloc struct managing the map[string]Item and sync.RWMutex.</p>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">internal/commands/</h3>
          <p className="text-[#9898ab]">Implementations of specific commands like SetCommand.</p>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">internal/pubsub/</h3>
          <p className="text-[#9898ab]">AllChannels broker, SubrChannel tracking, and concurrent message delivery.</p>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">internal/worker/</h3>
          <p className="text-[#9898ab]">Background tasks: ExpiryWorker (1s ticker) and SnapshotWorker (10m ticker).</p>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl md:col-span-2">
          <h3 className="text-lg font-medium text-white mb-2">internal/persistence/</h3>
          <p className="text-[#9898ab]">JSON snapshot saving and loading functionality.</p>
        </div>
      </div>

      <h2 id="port-mapping" className="text-2xl font-semibold text-white mb-4 mt-8">Service Port Mapping</h2>
      <InfoTable headers={portHeaders} rows={portRows} />
    </div>
  );
}
