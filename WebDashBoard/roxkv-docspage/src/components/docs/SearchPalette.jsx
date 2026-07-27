import { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { Search, Command, X } from 'lucide-react';
import { commands } from '../../data/commands';

export default function SearchPalette() {
  const [isOpen, setIsOpen] = useState(false);
  const [query, setQuery] = useState('');
  const [selectedIndex, setSelectedIndex] = useState(0);
  const inputRef = useRef(null);
  const navigate = useNavigate();

  const filteredCommands = commands.filter(cmd =>
    cmd.name.toLowerCase().includes(query.toLowerCase()) ||
    cmd.description.toLowerCase().includes(query.toLowerCase())
  );

  useEffect(() => {
    const handleKeyDown = (e) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault();
        setIsOpen((prev) => !prev);
      }
      
      if (e.key === 'Escape') {
        setIsOpen(false);
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, []);

  useEffect(() => {
    if (isOpen) {
      setTimeout(() => inputRef.current?.focus(), 50);
      setQuery('');
      setSelectedIndex(0);
    }
  }, [isOpen]);

  useEffect(() => {
    setSelectedIndex(0);
  }, [query]);

  const handleModalKeyDown = (e) => {
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      setSelectedIndex((prev) => (prev + 1) % filteredCommands.length);
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      setSelectedIndex((prev) => (prev - 1 + filteredCommands.length) % filteredCommands.length);
    } else if (e.key === 'Enter') {
      e.preventDefault();
      if (filteredCommands[selectedIndex]) {
        navigate(`/docs/commands/${filteredCommands[selectedIndex].name.toLowerCase()}`);
        setIsOpen(false);
      }
    }
  };

  return (
    <>
      {/* Floating Launcher */}
      <button
        onClick={() => setIsOpen(true)}
        className="fixed bottom-6 left-1/2 transform -translate-x-1/2 z-40 bg-[#16161f] border border-[#1e1e2e] hover:border-[rgba(124,58,237,0.5)] text-[#9898ab] px-4 py-2 rounded-full shadow-2xl flex items-center gap-3 transition-all hover:scale-105"
      >
        <Search className="w-4 h-4" />
        <span className="text-sm font-medium">Search Documentation</span>
        <div className="flex gap-1 ml-2">
          <kbd className="bg-[#0a0a0f] border border-[#1e1e2e] rounded px-1.5 py-0.5 text-xs text-[#f1f1f4]">
            <Command className="w-3 h-3 inline" />
          </kbd>
          <kbd className="bg-[#0a0a0f] border border-[#1e1e2e] rounded px-1.5 py-0.5 text-xs text-[#f1f1f4]">K</kbd>
        </div>
      </button>

      {/* Modal Overlay */}
      {isOpen && (
        <div className="fixed inset-0 z-50 flex items-start justify-center pt-[15vh] px-4 sm:px-0">
          <div 
            className="absolute inset-0 bg-black/60 backdrop-blur-sm"
            onClick={() => setIsOpen(false)}
          />
          <div 
            className="relative w-full max-w-xl bg-[#16161f] border border-[#1e1e2e] rounded-xl shadow-2xl overflow-hidden flex flex-col max-h-[60vh]"
            role="dialog"
            aria-modal="true"
          >
            {/* Search Input */}
            <div className="flex items-center px-4 py-4 border-b border-[#1e1e2e]">
              <Search className="w-5 h-5 text-[#9898ab] mr-3" />
              <input
                ref={inputRef}
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                onKeyDown={handleModalKeyDown}
                className="flex-1 bg-transparent border-none text-[#f1f1f4] placeholder-[#5e5e73] focus:outline-none focus:ring-0 text-lg"
                placeholder="Search commands..."
              />
              <button onClick={() => setIsOpen(false)} className="text-[#9898ab] hover:text-white p-1">
                <X className="w-5 h-5" />
              </button>
            </div>

            {/* Results */}
            <div className="overflow-y-auto p-2 scrollbar-thin">
              {filteredCommands.length > 0 ? (
                filteredCommands.map((cmd, index) => (
                  <button
                    key={cmd.name}
                    onClick={() => {
                      navigate(`/docs/commands/${cmd.name.toLowerCase()}`);
                      setIsOpen(false);
                    }}
                    onMouseEnter={() => setSelectedIndex(index)}
                    className={`w-full text-left flex items-center justify-between px-4 py-3 rounded-lg transition-colors ${
                      index === selectedIndex 
                        ? 'bg-[rgba(124,58,237,0.15)] border border-[rgba(124,58,237,0.3)]' 
                        : 'border border-transparent hover:bg-[#1e1e2e]'
                    }`}
                  >
                    <div>
                      <h4 className={`font-mono text-base ${index === selectedIndex ? 'text-[#a78bfa]' : 'text-[#f1f1f4]'}`}>
                        {cmd.name}
                      </h4>
                      <p className="text-[#9898ab] text-xs mt-1 truncate max-w-[300px] sm:max-w-[400px]">
                        {cmd.description}
                      </p>
                    </div>
                    <span className="text-xs px-2 py-1 bg-[#0a0a0f] rounded text-[#5e5e73]">
                      {cmd.category}
                    </span>
                  </button>
                ))
              ) : (
                <div className="py-12 text-center text-[#9898ab]">
                  No commands found matching "{query}"
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </>
  );
}
