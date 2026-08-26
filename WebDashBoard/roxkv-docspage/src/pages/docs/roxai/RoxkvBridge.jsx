export default function RoxkvBridge() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">RoxAI ↔ RoxKV Bridge</h1>
        <p className="text-xl text-[#9898ab]">Zero-overhead integration between AI and Data Layer</p>
      </div>
      <div className="section-divider" />

      <h2 id="direct-invocation" className="text-2xl font-semibold text-white mb-4">Direct Function Invocation</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        The defining architectural feature of RoxAI is its relationship with RoxKV. Unlike external AI tools that connect to databases via network protocols (like RESP or HTTP), RoxAI tools call RoxKV internals <strong>directly</strong> via Go function calls.
      </p>
      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#9898ab]"><strong>Same Process Architecture:</strong> Both the RoxAI orchestrator and the RoxKV data store run compiled together within the exact same Go binary.</p>
      </div>

      <h2 id="internal-call-mappings" className="text-2xl font-semibold text-white mb-4">Internal Call Mappings</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        Tool execution maps directly to underlying package functions:
      </p>
      <ul className="list-disc list-inside text-[#5e5e73] space-y-2 mb-8 ml-4">
        <li><strong>KV Tools:</strong> Call methods like <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">store.GetKv()</code>, <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">store.Set()</code>, <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">store.DelKv()</code>, <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">store.Keys()</code>.</li>
        <li><strong>Monitoring Tools:</strong> Call methods in the metrics package, e.g., <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">metrics.GetComputerUsage()</code>.</li>
        <li><strong>PubSub Tools:</strong> Call core pubsub management functions for topics and subscribers.</li>
        <li><strong>Persistence Tools:</strong> Can directly trigger or read snapshot operations in the storage layer.</li>
      </ul>

      <h2 id="advantages-and-limitations" className="text-2xl font-semibold text-white mb-4">Advantages & Limitations</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">Advantages</h3>
          <ul className="list-disc list-inside text-[#5e5e73]">
            <li>Zero network overhead or latency.</li>
            <li>Direct memory access for extreme performance.</li>
            <li>Strict Go type safety between agent and data store.</li>
          </ul>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">Limitations</h3>
          <ul className="list-disc list-inside text-[#5e5e73]">
            <li>Tightly coupled architecture.</li>
            <li>RoxAI cannot currently run as a standalone service separate from RoxKV nodes.</li>
          </ul>
        </div>
      </div>
    </div>
  );
}
