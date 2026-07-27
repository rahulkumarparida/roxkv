import { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { commands } from '../../data/commands';

export default function CommandIndex() {
  const [selectedCategory, setSelectedCategory] = useState('All');

  const categories = ['All', ...new Set(commands.map(cmd => cmd.category))];

  const filteredCommands = useMemo(() => {
    return commands.filter(cmd => {
      return selectedCategory === 'All' || cmd.category === selectedCategory;
    });
  }, [selectedCategory]);

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Command Reference</h1>
        <p className="text-xl text-[#9898ab]">
          Explore the complete list of commands supported by RoxKV.
        </p>
      </div>

      <div className="section-divider" />

      {/* Filters */}
     <div className="flex gap-2 mt-8 mb-6 overflow-x-auto pb-2 scrollbar-none md:pb-0">
  {categories.map(cat => (
    <button
      key={cat}
      onClick={() => setSelectedCategory(cat)}
      className={`px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-colors ${
        selectedCategory === cat 
          ? 'bg-[#7c3aed] text-white' 
          : 'bg-[#16161f] text-[#9898ab] border border-[#1e1e2e] hover:border-[#7c3aed]'
      }`}
    >
      {cat}
    </button>
  ))}
</div>
      {/* Command Grid */}
      {filteredCommands.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {filteredCommands.map(cmd => (
            <Link 
              key={cmd.name}
              to={`/docs/commands/${cmd.name.toLowerCase()}`}
              className="group bg-[#16161f] border border-[#1e1e2e] p-5 rounded-xl hover:border-[rgba(124,58,237,0.5)] hover:bg-[rgba(124,58,237,0.02)] transition-all block"
            >
              <div className="flex justify-between items-start mb-2">
                <h3 className="text-lg font-bold text-[#a78bfa] group-hover:text-[#c4b5fd] transition-colors">{cmd.name}</h3>
                <span className="text-xs px-2 py-1 rounded bg-[#1e1e2e] text-[#9898ab]">{cmd.category}</span>
              </div>
              <p className="text-sm text-[#9898ab] line-clamp-2">{cmd.description}</p>
            </Link>
          ))}
        </div>
      ) : (
        <div className="text-center py-12 bg-[#16161f] border border-[#1e1e2e] rounded-xl">
          <p className="text-[#9898ab]">No commands found matching your criteria.</p>
        </div>
      )}
    </div>
  );
}
