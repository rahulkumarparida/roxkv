import React, { useState } from 'react';
import { gallery } from '../data/content';
import { X, Maximize2 } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import SectionWrapper from './SectionWrapper';
import dashboardImage from '../assets/Dasboard.png';
import cliImage from '../assets/cli.png';
import architectureImage from '../assets/Architecture.png';
import monitoringImage from '../assets/Monitoring.png';

const galleryImages = {
  dashboard: dashboardImage,
  cli: cliImage,
  architecture: architectureImage,
  monitoring: monitoringImage,
};

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
            className="group glass rounded-xl overflow-hidden cursor-pointer glow-border aspect-video"
          >
            <div className="relative h-full">
              <img
                src={galleryImages[item.id]}
                alt={item.label}
                className="h-full w-full object-cover object-top transition-transform duration-300 group-hover:scale-[1.02]"
              />
              <div className="absolute inset-0 bg-gradient-to-t from-[#0a0a0f]/90 via-transparent to-transparent" />
              <div className="absolute bottom-0 left-0 right-0 flex items-end justify-between p-4">
                <div>
                  <span className="block text-sm font-semibold text-white">{item.label}</span>
                  <span className="block text-xs text-[#c3c3d3] mt-1">{item.description}</span>
                </div>
                <Maximize2 className="w-5 h-5 text-white/70 shrink-0" />
              </div>
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
              className="max-w-5xl w-full mx-4 glass rounded-2xl relative overflow-hidden"
            >
              <button 
                className="absolute top-4 right-4 w-8 h-8 rounded-full bg-white/5 flex items-center justify-center hover:bg-white/10 cursor-pointer text-white"
                onClick={() => setSelected(null)}
              >
                <X className="w-5 h-5" />
              </button>
              <img
                src={galleryImages[selected.id]}
                alt={selected.label}
                className="w-full max-h-[80vh] object-contain bg-[#050508]"
              />
            </motion.div>
          </div>
        )}
      </AnimatePresence>
    </SectionWrapper>
  );
}
