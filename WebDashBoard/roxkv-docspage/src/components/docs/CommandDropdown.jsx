import { useState, useRef, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { ChevronDown } from 'lucide-react';
import { commands } from '../../data/commands';
import { motion, AnimatePresence } from 'framer-motion';

export default function CommandDropdown() {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef(null);
  const navigate = useNavigate();
  const location = useLocation();

  const currentCommandPath = location.pathname.split('/').pop();
  const activeCommand = commands.find(c => c.name.toLowerCase() === currentCommandPath?.toLowerCase());

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleSelect = (cmdName) => {
    setIsOpen(false);
    navigate(`/docs/roxkv/commands/${cmdName.toLowerCase()}`);
  };

  return (
    <div className="relative w-full mb-6" ref={dropdownRef}>
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-full bg-[#16161f] border border-[#1e1e2e] hover:border-[rgba(124,58,237,0.4)] text-[#f1f1f4] rounded-lg px-4 py-2 flex items-center justify-between transition-colors"
      >
        <span className="font-mono text-sm truncate">
          {activeCommand ? activeCommand.name : 'Select a command...'}
        </span>
        <ChevronDown className={`w-4 h-4 text-[#9898ab] transition-transform ${isOpen ? 'rotate-180' : ''}`} />
      </button>

      <AnimatePresence>
        {isOpen && (
          <motion.div
            initial={{ opacity: 0, y: -5 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -5 }}
            transition={{ duration: 0.15 }}
            className="absolute z-50 w-full mt-2 bg-[#16161f] border border-[#1e1e2e] rounded-lg shadow-xl max-h-60 overflow-y-auto scrollbar-thin"
          >
            {commands.map((cmd) => (
              <button
                key={cmd.name}
                onClick={() => handleSelect(cmd.name)}
                className={`w-full text-left px-4 py-2 text-sm font-mono transition-colors ${
                  activeCommand?.name === cmd.name
                    ? 'bg-[rgba(124,58,237,0.1)] text-[#a78bfa]'
                    : 'text-[#9898ab] hover:bg-[#1e1e2e] hover:text-[#f1f1f4]'
                }`}
              >
                {cmd.name}
              </button>
            ))}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
