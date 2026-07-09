import React, { useState } from 'react';
import { gallery } from '../data/content';
import { X, Maximize2 } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import SectionWrapper from './SectionWrapper';

export default function Gallery() {
  const [selected, setSelected] = useState(null);

  return (
    <SectionWrapper id="gallery">
      <h2 className="text-3xl sm:text-4xl font-bold gradient-text text-center">
        Project Gallery
      </h2>
      
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 mt-12">
        {gallery.map((item, index) => (
          <motion.div
            key={item.id || index}
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: index * 0.1 }}
            whileHover={{ y: -4, scale: 1.02 }}
            onClick={() => setSelected(item)}
            className="glass rounded-xl overflow-hidden cursor-pointer glow-border aspect-video"
          >
            <div className="flex flex-col items-center justify-center h-full bg-gradient-to-br from-purple-900/20 to-[#0a0a0f]">
              <Maximize2 className="w-6 h-6 text-purple-500/30" />
              <span className="text-sm font-medium text-[#5e5e73] mt-2">{item.label}</span>
            </div>
          </motion.div>
        ))}
      </div>

      <AnimatePresence>
        {selected && (
          <div 
            className="modal-overlay" 
            onClick={() => setSelected(null)}
          >
            <motion.div 
              initial={{ opacity: 0, scale: 0.9 }} 
              animate={{ opacity: 1, scale: 1 }} 
              exit={{ opacity: 0, scale: 0.9 }} 
              onClick={e => e.stopPropagation()}
              className="max-w-4xl w-full mx-4 aspect-video glass rounded-2xl flex flex-col items-center justify-center relative bg-gradient-to-br from-purple-900/20 to-[#0a0a0f]"
            >
              <button 
                className="absolute top-4 right-4 w-8 h-8 rounded-full bg-white/5 flex items-center justify-center hover:bg-white/10 cursor-pointer text-white"
                onClick={() => setSelected(null)}
              >
                <X className="w-5 h-5" />
              </button>
              <span className="text-xl font-semibold text-[#5e5e73]">{selected.label}</span>
              <span className="text-sm text-[#3e3e53] mt-2">Screenshot placeholder</span>
            </motion.div>
          </div>
        )}
      </AnimatePresence>
    </SectionWrapper>
  );
}
