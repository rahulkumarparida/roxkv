import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { dashboardTabs } from '../data/content';
import SectionWrapper from './SectionWrapper';
import dashboardImage from '../assets/Dasboard.png';
import cliImage from '../assets/cli.png';
import monitoringImage from '../assets/Monitoring.png';
import architectureImage from '../assets/Architecture.png';

const tabImages = {
  dashboard: dashboardImage,
  cli: cliImage,
  monitoring: monitoringImage,
  architecture: architectureImage,
  chat: dashboardImage,
};

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
            className="h-full min-h-[260px] sm:min-h-[400px]"
          >
            <img
              src={tabImages[dashboardTabs[activeTab].id]}
              alt={`${dashboardTabs[activeTab].label} preview`}
              className="h-full w-full object-cover object-top"
            />
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
