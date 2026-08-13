import React, { useState } from 'react';
import { runLocally } from '../data/content';
import { Copy, Check } from 'lucide-react';
import SectionWrapper from './SectionWrapper';

export default function RunLocally() {
  const [copiedIndex, setCopiedIndex] = useState(null);

  const handleCopy = (steps, index) => {
    const commands = steps.map(step => step.command).join('\n');
    navigator.clipboard.writeText(commands);
    setCopiedIndex(index);
    setTimeout(() => setCopiedIndex(null), 2000);
  };

  return (
    <SectionWrapper id="run-locally">
      <h2 className="gradient-text text-3xl sm:text-4xl font-bold text-center">
        {runLocally.heading}
      </h2>
      <p className="text-[#9898ab] text-lg text-center max-w-2xl mx-auto mt-4">
        {runLocally.description}
      </p>
      
      <div className="max-w-4xl mx-auto mt-12 grid grid-cols-1 gap-8">
        {runLocally.methods.map((method, mIndex) => (
          <div key={mIndex} className="bg-[#13131a] rounded-xl border border-white/5 overflow-hidden">
            <div className="px-6 py-4 border-b border-white/5 bg-[#1a1a24]">
              <h3 className="text-xl font-semibold text-white">{method.name}</h3>
            </div>
            <div className="terminal-window m-6 relative">
              <div className="terminal-header">
                <div className="terminal-dot bg-[#ef4444]" />
                <div className="terminal-dot bg-[#f59e0b]" />
                <div className="terminal-dot bg-[#10b981]" />
                <span className="ml-3 text-xs text-[#5e5e73]">terminal ~ roxkv</span>
              </div>
              
              <div className="terminal-body relative">
                <button 
                  onClick={() => handleCopy(method.steps, mIndex)}
                  className="absolute top-2 right-2 p-2 rounded-lg text-[#5e5e73] hover:text-purple-400 transition cursor-pointer"
                  aria-label="Copy commands"
                >
                  {copiedIndex === mIndex ? <Check className="w-4 h-4" /> : <Copy className="w-4 h-4" />}
                </button>
                
                {method.steps.map((step, index) => (
                  <React.Fragment key={index}>
                    <div><span className="comment"># {step.label}</span></div>
                    <div><span className="prompt">$ </span><span className="command">{step.command}</span></div>
                    {index < method.steps.length - 1 && <div className="h-3"></div>}
                  </React.Fragment>
                ))}
              </div>
            </div>
          </div>
        ))}
      </div>
      
      <p className="text-sm text-[#5e5e73] text-center mt-8 italic">
        {runLocally.note}
      </p>
    </SectionWrapper>
  );
}
