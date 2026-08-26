import React from 'react';
import InfoTable from '../../../components/docs/InfoTable';
import FlowDiagram from '../../../components/docs/FlowDiagram';

export default function RespProtocol() {
  const respHeaders = ['Type', 'Prefix', 'Example'];
  const respRows = [
    ['Simple String', '+', '+OK\\r\\n'],
    ['Error', '-', '-ERR message\\r\\n'],
    ['Integer', ':', ':42\\r\\n'],
    ['Bulk String', '$', '$5\\r\\nhello\\r\\n'],
    ['Array', '*', '*2\\r\\n$3\\r\\nGET\\r\\n$3\\r\\nkey\\r\\n'],
    ['Boolean', '#', ''],
    ['Double', ',', ''],
    ['Null', '$-1', '$-1\\r\\n']
  ];

  const flowSteps = [
    { label: 'Client' },
    { label: 'Decoder' },
    { label: 'Dispatcher' },
    { label: 'Encoder' },
    { label: 'Client' }
  ];

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">RESP Protocol</h1>
        <p className="text-xl text-[#9898ab]">Redis Serialization Protocol implementation details.</p>
      </div>
      
      <div className="section-divider" />

      <h2 id="overview" className="text-2xl font-semibold text-white mb-4">Overview</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        RoxKV communicates using the standard Redis Serialization Protocol (RESP). This allows standard Redis clients to connect and interact with RoxKV natively.
      </p>

      <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl my-6">
        <FlowDiagram steps={flowSteps} direction="horizontal" />
      </div>

      <h2 id="supported-types" className="text-2xl font-semibold text-white mb-4 mt-8">Supported Types</h2>
      <InfoTable headers={respHeaders} rows={respRows} />

      <h2 id="decoding-flow" className="text-2xl font-semibold text-white mb-4 mt-8">Decoding Flow</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        Incoming bytes are processed by <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">DecodeArrayString</code>, parsing the raw RESP payload into a structured <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">[][]byte</code> representing the command and its arguments.
      </p>

      <h2 id="encoding-flow" className="text-2xl font-semibold text-white mb-4 mt-8">Encoding Flow</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        Responses are formatted using dedicated encoding functions:
      </p>
      <ul className="list-disc list-inside text-[#9898ab] space-y-2 mb-4">
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">EncodeBulkBytes</code>: Prepends <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">$length\r\n</code> to raw byte arrays.</li>
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">EncodeBulkString</code></li>
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">EncodeSimpleString</code></li>
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">EncodeError</code></li>
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">EncodeInt</code></li>
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">EncodeArray</code></li>
      </ul>

      <h2 id="binary-safety" className="text-2xl font-semibold text-white mb-4 mt-8">Binary Safety</h2>
      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#f1f1f4]">Values are parsed and returned as raw bytes. The encoder and decoder make no UTF-8 assumptions, ensuring full binary safety for arbitrary payloads.</p>
      </div>
    </div>
  );
}
