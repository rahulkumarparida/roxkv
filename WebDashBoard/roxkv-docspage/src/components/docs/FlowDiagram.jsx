export default function FlowDiagram({ steps, direction = 'horizontal' }) {
  if (direction === 'vertical') {
    return (
      <div className="bg-[#16161f] border border-[#1e1e2e] rounded-xl p-6 my-6">
        <div className="flex flex-col items-center gap-0">
          {steps.map((step, idx) => (
            <div key={idx} className="flex flex-col items-center">
              <div className={`px-5 py-3 rounded-lg border text-center min-w-[200px] ${
                step.highlight
                  ? 'bg-[#0a0a0f] border-purple-500/30 shadow-lg shadow-purple-900/10'
                  : 'bg-[#1e1e2e] border-[#2a2a3e]'
              }`}>
                <span className="font-semibold text-white text-sm block">{step.label}</span>
                {step.sub && <span className="text-xs text-[#9898ab] font-mono mt-0.5 block">{step.sub}</span>}
              </div>
              {idx < steps.length - 1 && (
                <div className="flex flex-col items-center py-1">
                  <div className="w-px h-5 bg-purple-500/30"></div>
                  <div className="text-purple-400/60 text-xs">▼</div>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="bg-[#16161f] border border-[#1e1e2e] rounded-xl p-6 sm:p-8 my-6 overflow-x-auto">
      <div className="flex items-center gap-2 min-w-[600px] justify-center">
        {steps.map((step, idx) => (
          <div key={idx} className="flex items-center gap-2">
            <div className={`px-4 py-3 rounded-lg border text-center min-w-[100px] ${
              step.highlight
                ? 'bg-[#0a0a0f] border-purple-500/30 shadow-lg shadow-purple-900/10'
                : 'bg-[#1e1e2e] border-[#2a2a3e]'
            }`}>
              <span className="font-semibold text-white text-sm block">{step.label}</span>
              {step.sub && <span className="text-xs text-[#9898ab] font-mono mt-0.5 block">{step.sub}</span>}
            </div>
            {idx < steps.length - 1 && (
              <div className="flex items-center text-purple-500/40 text-xs px-1">→</div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
