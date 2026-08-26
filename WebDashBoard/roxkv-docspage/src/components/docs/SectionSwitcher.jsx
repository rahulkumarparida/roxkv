import { Link, useLocation } from 'react-router-dom';
import { Database, Brain } from 'lucide-react';

export default function SectionSwitcher() {
  const location = useLocation();
  const isRoxAI = location.pathname.startsWith('/docs/roxai');
  const isArchitecture = location.pathname === '/docs/architecture';

  return (
    <div className="flex gap-1  bg-[#0a0a0f] rounded-lg border border-[#1e1e2e]">
      <Link
        to="/docs/roxkv/overview"
        className={`flex items-center gap-2 px-4 py-2 rounded-md text-sm font-medium transition-all ${
          !isRoxAI && !isArchitecture
            ? 'bg-[rgba(124,58,237,0.15)] text-[#a78bfa] border border-[rgba(124,58,237,0.2)]'
            : 'text-[#9898ab] hover:text-white hover:bg-[#16161f] border border-transparent'
        }`}
      >
        <Database className="w-3.5 h-3.5" />
        RoxKV
      </Link>
      <Link
        to="/docs/roxai/overview"
        className={`flex items-center gap-2 px-4 py-2 rounded-md text-sm font-medium transition-all ${
          isRoxAI
            ? 'bg-[rgba(124,58,237,0.15)] text-[#a78bfa] border border-[rgba(124,58,237,0.2)]'
            : 'text-[#9898ab] hover:text-white hover:bg-[#16161f] border border-transparent'
        }`}
      >
        <Brain className="w-3.5 h-3.5" />
        RoxAI
      </Link>
    </div>
  );
}
