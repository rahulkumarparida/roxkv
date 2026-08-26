import { Outlet, useLocation } from 'react-router-dom';
import Navbar from '../Navbar';
import DocsSidebar from './DocsSidebar';
import DocsRightSidebar from './DocsRightSidebar';
import DocsCommandSidebar from './DocsCommandSidebar';
import SearchPalette from './SearchPalette';
import SectionSwitcher from './SectionSwitcher';

export default function DocsLayout() {
  const location = useLocation();
  const isCommandRoute = location.pathname.includes('/commands/') || location.pathname.endsWith('/commands');
  const isCommandDetailRoute = location.pathname.match(/\/commands\/[^/]+$/);
  
  return (
    <div className="min-h-screen bg-[#0a0a0f] text-[#f1f1f4] flex flex-col">
      <Navbar />
      <div className="flex-1 flex pt-16 max-w-[1536px] mx-auto w-full">
        {/* Main Left Sidebar */}
        <DocsSidebar />
        
        {/* Conditional Command Sidebar */}
        {isCommandDetailRoute && <DocsCommandSidebar />}
        
        {/* Main Content Area */}
        <main className="flex-1 min-w-0 py-10 px-6 sm:px-10 lg:pl-12 lg:pr-8">
          <div className="w-full max-w-4xl mx-auto prose prose-invert prose-purple max-w-none">
            <Outlet key={location.pathname} />
          </div>
        </main>
        
        {/* Right Sidebar */}
        <DocsRightSidebar key={location.pathname + '-sidebar'} />
      </div>
      
      {/* Search Launcher & Palette */}
      <SearchPalette />
    </div>
  );
}
