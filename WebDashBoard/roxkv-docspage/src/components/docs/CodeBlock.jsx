import { useState } from 'react';
import { Copy, Check } from 'lucide-react';

export default function CodeBlock({ code, language = 'bash', title = '' }) {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="terminal-window my-4">
      <div className="terminal-header flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="flex gap-2">
            <div className="terminal-dot bg-red-500"></div>
            <div className="terminal-dot bg-yellow-500"></div>
            <div className="terminal-dot bg-green-500"></div>
          </div>
          {title && <span className="ml-4 text-xs text-[#9898ab]">{title}</span>}
          {!title && language && <span className="ml-4 text-xs text-[#9898ab]">{language}</span>}
        </div>
        <button
          onClick={handleCopy}
          className="text-[#5e5e73] hover:text-[#a78bfa] transition-colors p-1"
          title="Copy to clipboard"
        >
          {copied ? <Check className="w-3.5 h-3.5 text-emerald-400" /> : <Copy className="w-3.5 h-3.5" />}
        </button>
      </div>
      <div className="terminal-body bg-[#16161f] p-4 overflow-x-auto">
        <pre className="font-mono text-sm leading-relaxed text-[#f1f1f4] whitespace-pre">{code}</pre>
      </div>
    </div>
  );
}
