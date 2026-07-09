import { motion } from 'framer-motion';
import { Terminal } from 'lucide-react';
import { whyItExists } from '../data/content';
import SectionWrapper from './SectionWrapper';

export default function WhyItExists() {
  return (
    <SectionWrapper id="why-it-exists">
      <motion.h2
        className="text-3xl sm:text-4xl font-bold gradient-text text-center"
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.5 }}
      >
        {whyItExists.heading}
      </motion.h2>

      <motion.p
        className="text-[#9898ab] text-lg max-w-3xl mx-auto text-center mt-4"
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.5, delay: 0.1 }}
      >
        {whyItExists.description}
      </motion.p>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 mt-16">
        {whyItExists.examples.map((example, index) => (
          <motion.div
            key={index}
            className="glass rounded-2xl p-6 glow-border"
            whileHover={{ y: -4 }}
            initial={{ opacity: 0, y: 30 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.4, delay: index * 0.1 }}
          >
            <span className="text-3xl text-purple-500/30 font-serif leading-none">
              &ldquo;
            </span>
            <p className="text-lg font-medium text-white mt-1">
              {example.query}
            </p>
            <div className="flex items-center gap-2 mt-4">
              <Terminal className="w-3.5 h-3.5 text-purple-400/60" />
              <span className="text-xs font-mono text-purple-400/60">
                {example.tool}
              </span>
            </div>
          </motion.div>
        ))}
      </div>
    </SectionWrapper>
  );
}
