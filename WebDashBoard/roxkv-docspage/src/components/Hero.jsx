import { motion } from 'framer-motion';
import { Play, GitBranch, Bot, Terminal } from 'lucide-react';
import { hero } from '../data/content';
import logo from '../assets/roxkvLogo.png';

const ease = [0.25, 0.4, 0.25, 1];

export default function Hero() {
  return (
    <section className="relative min-h-screen flex items-center justify-center overflow-hidden bg-[#0a0a0f]">
      {/* Background Layers */}
      <div className="absolute inset-0 pointer-events-none">
        <div className="hero-grid absolute inset-0 opacity-40" />
        <div className="orb orb-1" />
        <div className="orb orb-2" />
        <div className="orb orb-3" />
        <div className="absolute bottom-0 left-0 right-0 h-40 bg-gradient-to-t from-[#0a0a0f] via-transparent to-transparent" />
      </div>

      {/* Content */}
      <div className="relative z-10 text-center max-w-4xl mx-auto px-4 pt-24">
        {/* Badge */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.7, delay: 0, ease }}
          className="inline-flex items-center gap-2 px-4 py-2 rounded-full border border-purple-500/20 bg-purple-500/5 text-sm text-purple-300"
        >
          <Bot className="w-4 h-4" />
          {hero.badge}
        </motion.div>

        {/* Title Line 1 */}
        <motion.h1
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.7, delay: 0.15, ease }}
          className="text-5xl sm:text-6xl lg:text-7xl font-extrabold leading-[1.1] tracking-tight text-white mt-8"
        >
          {hero.titleLine1}
        </motion.h1>

        {/* Title Line 2 */}
        <motion.span
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.7, delay: 0.3, ease }}
          className="text-5xl sm:text-6xl lg:text-7xl font-extrabold leading-[1.1] tracking-tight gradient-text-purple block mt-2"
        >
          {hero.titleLine2}
        </motion.span>

        {/* Subtitle */}
        <motion.p
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.7, delay: 0.45, ease }}
          className="text-lg sm:text-xl text-[#9898ab] max-w-2xl mx-auto mt-6"
        >
          {hero.subtitle}
        </motion.p>

        {/* CTA Buttons */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.7, delay: 0.6, ease }}
          className="flex flex-col sm:flex-row gap-4 justify-center mt-10"
        >
          <a href={hero.primaryCta.href} className="btn-primary inline-flex items-center justify-center gap-2">
            <Play className="w-5 h-5" />
            {hero.primaryCta.label}
          </a>
          <a
            href={hero.secondaryCta.href}
            target="_blank"
            rel="noopener noreferrer"
            className="btn-secondary inline-flex items-center justify-center gap-2"
          >
            <GitBranch className="w-5 h-5" />
            {hero.secondaryCta.label}
          </a>
        </motion.div>

        {/* Terminal Decoration */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.7, delay: 0.8, ease }}
          className="max-w-lg mx-auto mt-16"
        >
          <div className="terminal-window">
            <div className="terminal-header">
              <span className="terminal-dot" style={{ background: '#ef4444' }} />
              <span className="terminal-dot" style={{ background: '#f59e0b' }} />
              <span className="terminal-dot" style={{ background: '#10b981' }} />
              <span className="ml-3 text-xs text-[#5e5e73]">roxkv-chat</span>
            </div>
            <div className="terminal-body font-mono text-sm text-left">
              <div>
                <span className="prompt">$ </span>
                <span className="command">roxai query "What is the health of my database?"</span>
              </div>
              <br />
              <div>
                <span className="comment">✓ Analyzing 45 tools...</span>
              </div>
              <div>
                <span className="comment">✓ Routing: get_system_health, get_cpu_usage, get_ram_usage</span>
              </div>
              <div>
                <span className="output">✓ Database is healthy. CPU: 12%, RAM: 340MB, 156 keys active.</span>
              </div>
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  );
}
