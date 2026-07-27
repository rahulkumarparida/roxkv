import { Link, useLocation } from 'react-router-dom';
import { ChevronRight } from 'lucide-react';

const navItems = [
  {
    title: 'Overview',
    links: [
      { name: 'Introduction', href: '/docs/overview' },
    ],
  },
  {
    title: 'Core Concepts',
    links: [
      { name: 'Architecture', href: '/docs/architecture' },
      { name: 'Persistence', href: '/docs/persistence' },
    ],
  },
  {
    title: 'Protocol & API',
    links: [
      { name: 'RESP Protocol', href: '/docs/resp' },
      { name: 'Commands Index', href: '/docs/commands' },
      { name: 'Pub/Sub', href: '/docs/pubsub' },
    ],
  },
];

export default function DocsSidebar() {
  const location = useLocation();

  return (
    <aside className="w-64 flex-shrink-0 hidden md:block overflow-y-auto h-[calc(100vh-64px)] sticky top-16 border-r border-[#1e1e2e] bg-[#0a0a0f] p-6 scrollbar-thin">
      <div className="space-y-8">
        {navItems.map((section) => (
          <div key={section.title}>
            <h4 className="font-semibold text-[#f1f1f4] mb-3 text-sm uppercase tracking-wider">
              {section.title}
            </h4>
            <ul className="space-y-2">
              {section.links.map((link) => {
                const isActive = location.pathname === link.href || location.pathname.startsWith(link.href + '/');
                return (
                  <li key={link.href}>
                    <Link
                      to={link.href}
                      className={`flex items-center text-sm transition-colors duration-200 ${
                        isActive
                          ? 'text-[#a78bfa] font-medium'
                          : 'text-[#9898ab] hover:text-[#f1f1f4]'
                      }`}
                    >
                      {isActive && <ChevronRight className="w-4 h-4 mr-1 text-[#7c3aed]" />}
                      <span className={isActive ? '' : 'ml-5'}>{link.name}</span>
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
