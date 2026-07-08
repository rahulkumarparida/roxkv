import { Settings, Globe, Zap, Bot } from 'lucide-react';
import { APP_META } from '../utils/constants';
import bgImage from '../assets/4.png';

/**
 * DashboardLayout provides the persistent shell with background image,
 * dark overlay, and status bar.
 */
export default function DashboardLayout({ children }) {
  return (
    <div className="relative h-screen w-screen overflow-hidden">
      {/* Background Image */}
      <div
        className="absolute inset-0 bg-cover bg-center bg-no-repeat"
        style={{ backgroundImage: `url(${bgImage})` }}
      />

      {/* Dark Overlay */}
      <div className="absolute inset-0 bg-black/80" />

      {/* Main Content */}
      <div className="relative z-10 flex flex-col h-full">
        {/* Dashboard Content */}
        <main className="flex-1 min-h-0">
          {children}
        </main>

        {/* Bottom Status Bar */}
        <footer className="shrink-0 px-4 py-2 border-t border-white/5 bg-black/40 backdrop-blur-sm">
          <div className="flex items-center justify-between text-[10px]">
            {/* Left: App Info */}
            <div className="flex items-center gap-3">
              <span className="font-bold text-text-primary tracking-wider">{APP_META.name}</span>
              <span className="text-purple-400 font-mono">{APP_META.version}</span>
              <span className="text-text-muted">{APP_META.tagline}</span>
            </div>

            {/* Center: Quick Icons */}
            <div className="flex items-center gap-3 text-text-muted">
              <button className="hover:text-purple-400 transition-colors p-1"><Settings className="h-3 w-3" /></button>
              <button className="hover:text-purple-400 transition-colors p-1"><Globe className="h-3 w-3" /></button>
              <button className="hover:text-purple-400 transition-colors p-1"><Zap className="h-3 w-3" /></button>
              <button className="hover:text-purple-400 transition-colors p-1"><Bot className="h-3 w-3" /></button>
            </div>

            {/* Right: Connection Info */}
            <div className="flex items-center gap-4">
              <span className="text-text-muted">
                Connection: <span className="text-text-secondary font-mono">{APP_META.connectionHost}</span>
              </span>
              <span className="text-green-400 flex items-center gap-1">
                <span className="h-1.5 w-1.5 rounded-full bg-green-400 animate-pulse" />
                Auto-refresh: ON
              </span>
            </div>
          </div>
        </footer>
      </div>

      {/* Master Agent Button */}
      <div className="fixed top-12 left-1/2 -translate-x-1/2 z-20">
        <button className="flex items-center gap-2 px-4 py-1.5 rounded-full bg-purple-500/20 border border-purple-500/30 text-xs text-purple-300 hover:bg-purple-500/30 hover:shadow-[0_0_15px_rgba(168,85,247,0.4)] transition-all duration-300 active:scale-95 backdrop-blur-md">
          <Bot className="h-3.5 w-3.5" />
          Master Agent
        </button>
      </div>
    </div>
  );
}
