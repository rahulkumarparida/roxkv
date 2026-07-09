import { motion } from 'framer-motion';
import { Play, ExternalLink } from 'lucide-react';
import { demo } from '../data/content';
import SectionWrapper from './SectionWrapper';

export default function Demo() {
  return (
    <SectionWrapper id="demo">
      <h2 className="gradient-text text-3xl sm:text-4xl font-bold text-center">
        {demo.heading}
      </h2>
      <p className="text-[#9898ab] text-lg text-center max-w-2xl mx-auto mt-4">
        {demo.description}
      </p>

      {/* Video placeholder */}
      <motion.div
        className="glass rounded-2xl overflow-hidden aspect-video max-w-4xl mx-auto mt-12 glow-border flex flex-col items-center justify-center bg-gradient-to-br from-purple-900/10 to-[#0a0a0f]"
        initial={{ opacity: 0, y: 30 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.5 }}
      >
        <div className="w-20 h-20 rounded-full bg-purple-500/10 border border-purple-500/20 flex items-center justify-center">
          <Play className="w-10 h-10 text-purple-400" />
        </div>
        <span className="text-[#5e5e73] mt-4 text-sm">Click to play demo</span>
      </motion.div>

      {/* CTA Button */}
      <a
        href={demo.ctaHref}
        className="btn-primary mt-8 mx-auto flex w-fit"
        target="_blank"
        rel="noopener noreferrer"
      >
        <ExternalLink className="w-5 h-5" />
        {demo.ctaLabel}
      </a>
    </SectionWrapper>
  );
}
