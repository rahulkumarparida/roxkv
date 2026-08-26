import { Link, useLocation } from 'react-router-dom';
import { ChevronRight, Layers } from 'lucide-react';
import { roxkvNav, roxaiNav } from '../../data/navigation';
import SectionSwitcher from './SectionSwitcher';

export default function DocsSidebar() {
  const location = useLocation();
  const isRoxAI = location.pathname.startsWith('/docs/roxai');
  const isArchitecture = location.pathname === '/docs/architecture';
  
  const navItems = isRoxAI ? roxaiNav : roxkvNav;

  return (
    <aside className="w-64 flex-shrink-0 hidden md:block overflow-y-auto h-[calc(100vh-64px)] sticky top-16 border-r border-[#1e1e2e] bg-[#0a0a0f] p-6 scrollbar-thin">
      <div className="space-y-6">
        {/* Section Switcher */}
        <SectionSwitcher />
        
        {/* Architecture Link */}
        <Link
          to="/docs/architecture"
          className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm transition-colors ${
            isArchitecture
              ? 'bg-[rgba(124,58,237,0.1)] text-[#a78bfa] font-medium border border-[rgba(124,58,237,0.2)]'
              : 'text-[#9898ab] hover:text-[#f1f1f4] hover:bg-[#16161f]'
          }`}
        >
          <Layers className="w-3.5 h-3.5" />
          System Architecture
        </Link>
        
        <div className="h-px bg-[#1e1e2e]" />

        {/* Section Navigation */}
        {navItems.map((section) => (
          <div key={section.title}>
            <h4 className="font-semibold text-[#f1f1f4] mb-3 text-xs uppercase tracking-wider">
              {section.title}
            </h4>
            <ul className="space-y-1">
              {section.links.map((link) => {
                const isActive = location.pathname === link.href || location.pathname.startsWith(link.href + '/');
                return (
                  <li key={link.href}>
                    <Link
                      to={link.href}
                      className={`flex items-center text-sm transition-colors duration-200 px-3 py-1.5 rounded-lg ${
                        isActive
                          ? 'bg-[rgba(124,58,237,0.1)] text-[#a78bfa] font-medium border border-[rgba(124,58,237,0.2)]'
                          : 'text-[#9898ab] hover:text-[#f1f1f4] hover:bg-[#16161f] border border-transparent'
                      }`}
                    >
                      {link.icon && <link.icon className="w-3.5 h-3.5 mr-2 flex-shrink-0" />}
                      <span>{link.name}</span>
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
