import React from 'react';
import { roadmap } from '../data/content';
import { motion } from 'framer-motion';
import SectionWrapper from './SectionWrapper';

export default function Roadmap() {
  return (
    <SectionWrapper id="roadmap">
      <h2 className="text-3xl sm:text-4xl font-bold gradient-text text-center">
        Future Roadmap
      </h2>
      
      <div className="max-w-2xl mx-auto mt-16 relative">
        <div className="absolute left-[15px] md:left-1/2 top-0 bottom-0 w-0.5 bg-gradient-to-b from-purple-500/40 via-purple-500/20 to-transparent md:transform md:-translate-x-1/2" />
        
        {roadmap.map((item, index) => (
          <motion.div
            key={index}
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: index * 0.1 }}
            className={`relative flex items-start gap-6 mb-10 md:mb-12 ${index % 2 === 0 ? 'md:flex-row-reverse' : ''}`}
          >
            <div className="relative z-10 w-8 h-8 rounded-full bg-purple-500/20 border-2 border-purple-500/40 flex items-center justify-center flex-shrink-0 md:absolute md:left-1/2 md:-translate-x-1/2">
              <div className="w-2 h-2 rounded-full bg-purple-400 pulse-dot" />
            </div>
            <div className="glass rounded-xl p-5 glow-border flex-1 md:max-w-[calc(50%-2rem)]">
              <h3 className="font-semibold text-white">{item.title}</h3>
              <p className="text-sm text-[#9898ab] mt-1">{item.description}</p>
            </div>
          </motion.div>
        ))}
      </div>
    </SectionWrapper>
  );
}
