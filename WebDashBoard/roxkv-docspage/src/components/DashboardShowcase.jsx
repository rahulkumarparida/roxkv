import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Monitor } from 'lucide-react';
import { dashboardTabs } from '../data/content';
import SectionWrapper from './SectionWrapper';

export default function DashboardShowcase() {
  const [activeTab, setActiveTab] = useState(0);

  return (
    <SectionWrapper id="dashboard">
      <motion.h2
        className="text-3xl sm:text-4xl font-bold gradient-text text-center"
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.5 }}
      >
        See It in Action
      </motion.h2>

      {/* Tab buttons */}
      <div className="flex flex-wrap justify-center gap-3 mt-8">
        {dashboardTabs.map((tab, index) => (
          <button
            key={index}
            onClick={() => setActiveTab(index)}
            className={`rounded-xl px-5 py-2.5 text-sm font-medium transition-all ${
              activeTab === index
                ? 'bg-purple-500/20 text-purple-300 border border-purple-500/30'
                : 'text-[#5e5e73] hover:text-[#9898ab] border border-transparent'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* Image / preview area */}
      <div className="mt-8 glass rounded-2xl overflow-hidden aspect-video max-w-5xl mx-auto relative">
        <AnimatePresence mode="wait">
          <motion.div
            key={activeTab}
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -20 }}
            transition={{ duration: 0.3 }}
            className="flex flex-col items-center justify-center h-full bg-gradient-to-br from-purple-900/20 to-[#0a0a0f] min-h-[400px]"
          >
            <div className="border-2 border-dashed border-[#1e1e2e] rounded-xl p-12 flex flex-col items-center">
              <Monitor className="w-12 h-12 text-purple-500/30" />
              <p className="text-lg font-medium text-[#5e5e73] mt-4">
                {dashboardTabs[activeTab].label}
              </p>
              <p className="text-sm text-[#3e3e53] mt-1">
                Screenshot placeholder
              </p>
            </div>
          </motion.div>
        </AnimatePresence>
      </div>

      {/* Description */}
      <p className="text-[#9898ab] text-center mt-6">
        {dashboardTabs[activeTab].description}
      </p>
    </SectionWrapper>
  );
}
