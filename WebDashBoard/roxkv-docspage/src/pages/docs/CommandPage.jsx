import { useParams, Link, Navigate } from 'react-router-dom';
import { commands } from '../../data/commands';
import { ArrowLeft } from 'lucide-react';

export default function CommandPage() {
  const { cmd } = useParams();
  
  const commandData = commands.find(c => c.name.toLowerCase() === cmd?.toLowerCase());

  if (!commandData) {
    return <Navigate to="/docs/commands" replace />;
  }

  return (
    <div className="space-y-8 pb-10">
      <div>
        <Link to="/docs/commands" className="inline-flex items-center text-sm text-[#9898ab] hover:text-[#a78bfa] transition-colors mb-6">
          <ArrowLeft className="w-4 h-4 mr-1" />
          Back to Commands
        </Link>
        <div className="flex items-center gap-4 mb-4">
          <h1 className="text-4xl font-bold text-white">{commandData.name}</h1>
          <span className="text-xs px-3 py-1 rounded-full bg-[#1e1e2e] text-[#a78bfa] border border-[rgba(124,58,237,0.2)]">
            {commandData.category}
          </span>
        </div>
        <p className="text-xl text-[#9898ab]">
          {commandData.description}
        </p>
      </div>

      <div className="section-divider" />

      <h2 id="syntax" className="text-2xl font-bold text-white mt-12 mb-4">Syntax</h2>
      <div className="bg-[#16161f] border border-[#1e1e2e] rounded-xl p-4 overflow-x-auto">
        <code className="text-lg font-mono text-[#f1f1f4]">{commandData.syntax}</code>
      </div>

      {commandData.arguments.length > 0 && (
        <>
          <h2 id="arguments" className="text-2xl font-bold text-white mt-12 mb-4">Arguments</h2>
          <div className="overflow-x-auto rounded-xl border border-[#1e1e2e] bg-[#16161f]">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="border-b border-[#1e1e2e] bg-[rgba(124,58,237,0.05)]">
                  <th className="py-3 px-4 text-sm font-semibold text-[#a78bfa]">Argument</th>
                  <th className="py-3 px-4 text-sm font-semibold text-[#a78bfa]">Description</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[#1e1e2e]">
                {commandData.arguments.map((arg, idx) => (
                  <tr key={idx} className="hover:bg-[rgba(255,255,255,0.02)] transition-colors">
                    <td className="py-3 px-4 text-sm font-mono text-white whitespace-nowrap">{arg.name}</td>
                    <td className="py-3 px-4 text-sm text-[#9898ab]">{arg.description}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </>
      )}

      <h2 id="behavior" className="text-2xl font-bold text-white mt-12 mb-4">Behavior</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        {commandData.behavior}
      </p>

      <h2 id="return-value" className="text-2xl font-bold text-white mt-12 mb-4">Return Value</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4" dangerouslySetInnerHTML={{ __html: commandData.returnValue.replace(/`([^`]+)`/g, '<code class="bg-[#1e1e2e] text-[#a78bfa] px-1.5 py-0.5 rounded">$1</code>') }} />

      <h2 id="examples" className="text-2xl font-bold text-white mt-12 mb-4">Examples</h2>
      <div className="space-y-4">
        {commandData.examples.map((example, idx) => (
          <div key={idx} className="terminal-window">
            <div className="terminal-header">
              <div className="flex gap-2">
                <div className="terminal-dot bg-red-500"></div>
                <div className="terminal-dot bg-yellow-500"></div>
                <div className="terminal-dot bg-green-500"></div>
              </div>
              <span className="ml-4 text-xs text-[#9898ab]">redis-cli</span>
            </div>
            <div className="terminal-body bg-[#16161f] p-4 font-mono text-sm">
              <div><span className="prompt">127.0.0.1:6379</span> <span className="command">{example}</span></div>
            </div>
          </div>
        ))}
      </div>

      {commandData.notes && (
        <>
          <h2 id="notes" className="text-2xl font-bold text-white mt-12 mb-4">Notes</h2>
          <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4">
            <p className="text-[#9898ab] text-sm">{commandData.notes}</p>
          </div>
        </>
      )}

      {commandData.related && commandData.related.length > 0 && (
        <>
          <h2 id="related-commands" className="text-2xl font-bold text-white mt-12 mb-4">Related Commands</h2>
          <div className="flex flex-wrap gap-2">
            {commandData.related.map(rel => (
              <Link 
                key={rel} 
                to={`/docs/commands/${rel.toLowerCase()}`}
                className="px-3 py-1.5 rounded-lg border border-[#1e1e2e] bg-[#16161f] text-sm text-[#a78bfa] hover:border-[#7c3aed] transition-colors"
              >
                {rel}
              </Link>
            ))}
          </div>
        </>
      )}
    </div>
  );
}
