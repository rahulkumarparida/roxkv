import { Link, useLocation } from 'react-router-dom';
import { commands } from '../../data/commands';
import CommandDropdown from './CommandDropdown';

export default function DocsCommandSidebar() {
  const location = useLocation();

  const groupedCommands = commands.reduce((acc, cmd) => {
    if (!acc[cmd.category]) acc[cmd.category] = [];
    acc[cmd.category].push(cmd);
    return acc;
  }, {});

  return (
    <aside className="w-full lg:w-56 flex-shrink-0 lg:h-[calc(100vh-64px)] lg:sticky lg:top-16 lg:border-r border-[#1e1e2e] bg-[#0a0a0f] p-6 lg:overflow-y-auto scrollbar-none">
      
      {/* Dropdown for RESP Client Command Selector */}
      <CommandDropdown />

      <div className="space-y-6 hidden lg:block">
        {Object.entries(groupedCommands).map(([category, cmds]) => (
          <div key={category}>
            <h4 className="font-semibold text-[#f1f1f4] mb-2 text-xs uppercase tracking-widest opacity-60">
              {category}
            </h4>
            <ul className="space-y-1">
              {cmds.map((cmd) => {
                const cmdPath = `/docs/commands/${cmd.name.toLowerCase()}`;
                const isActive = location.pathname === cmdPath;
                
                return (
                  <li key={cmd.name}>
                    <Link
                      to={cmdPath}
                      className={`block px-3 py-1.5 rounded-lg text-sm transition-colors duration-200 ${
                        isActive
                          ? 'bg-[rgba(124,58,237,0.1)] text-[#a78bfa] font-medium border border-[rgba(124,58,237,0.2)]'
                          : 'text-[#9898ab] hover:text-[#f1f1f4] hover:bg-[#16161f]'
                      }`}
                    >
                      {cmd.name}
                    </Link>
                  </li>
                );
              })}
            </ul>
          </div>
        ))}
      </div>
    </aside>
  );
}
