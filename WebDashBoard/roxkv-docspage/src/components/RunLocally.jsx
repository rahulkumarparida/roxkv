import React, { useState } from 'react';
import { runLocally } from '../data/content';
import { Copy, Check } from 'lucide-react';
import SectionWrapper from './SectionWrapper';

export default function RunLocally() {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    const commands = runLocally.steps.map(step => step.command).join('\n');
    navigator.clipboard.writeText(commands);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <SectionWrapper id="run-locally">
      <h2 className="gradient-text text-3xl sm:text-4xl font-bold text-center">
        {runLocally.heading}
      </h2>
      <p className="text-[#9898ab] text-lg text-center max-w-2xl mx-auto mt-4">
        {runLocally.description}
      </p>
      
      <div className="max-w-3xl mx-auto mt-12 terminal-window relative">
        <div className="terminal-header">
          <div className="terminal-dot bg-[#ef4444]" />
          <div className="terminal-dot bg-[#f59e0b]" />
          <div className="terminal-dot bg-[#10b981]" />
          <span className="ml-3 text-xs text-[#5e5e73]">terminal ~ roxkv</span>
        </div>
        
        <div className="terminal-body relative">
          <button 
            onClick={handleCopy}
            className="absolute top-2 right-2 p-2 rounded-lg text-[#5e5e73] hover:text-purple-400 transition cursor-pointer"
            aria-label="Copy commands"
          >
            {copied ? <Check className="w-4 h-4" /> : <Copy className="w-4 h-4" />}
          </button>
          
          {runLocally.steps.map((step, index) => (
            <React.Fragment key={index}>
              <div><span className="comment"># {step.label}</span></div>
              <div><span className="prompt">$ </span><span className="command">{step.command}</span></div>
              {index < runLocally.steps.length - 1 && <div className="h-3"></div>}
            </React.Fragment>
          ))}
        </div>
      </div>
      
      <p className="text-sm text-[#5e5e73] text-center mt-6 italic">
        {runLocally.note}
      </p>
    </SectionWrapper>
  );
}
