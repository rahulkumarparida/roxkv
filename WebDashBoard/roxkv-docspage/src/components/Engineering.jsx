import { motion } from 'framer-motion';
import { engineeringDecisions } from '../data/content';
import SectionWrapper from './SectionWrapper';

export default function Engineering() {
  return (
    <SectionWrapper id="engineering">
      <h2 className="text-3xl sm:text-4xl font-bold gradient-text text-center">
        Engineering Decisions
      </h2>
      <p className="text-[#9898ab] text-lg text-center mt-4">
        What makes this project architecturally unique.
      </p>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-16">
        {engineeringDecisions.map((decision, index) => (
          <motion.div
            key={decision.title}
            className="glass rounded-2xl p-7 glow-border"
            initial={{ opacity: 0, y: 30 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: index * 0.08 }}
            whileHover={{ y: -4 }}
          >
            <div className="w-8 h-8 rounded-lg bg-purple-500/20 text-purple-400 flex items-center justify-center text-sm font-bold mb-4">
              {(index + 1).toString().padStart(2, '0')}
            </div>
            <h3 className="text-lg font-semibold text-white">{decision.title}</h3>
            <p className="text-sm text-[#9898ab] mt-2 leading-relaxed">
              {decision.description}
            </p>
          </motion.div>
        ))}
      </div>
    </SectionWrapper>
  );
}
