import { motion } from 'framer-motion';
import { techStack } from '../data/content';
import SectionWrapper from './SectionWrapper';

export default function TechStack() {
  return (
    <SectionWrapper id="tech-stack">
      <h2 className="text-3xl sm:text-4xl font-bold gradient-text text-center">
        Technology Stack
      </h2>
      <p className="text-[#9898ab] text-lg text-center mt-4">
        Built with modern, battle-tested technologies.
      </p>

      <div className="flex flex-wrap justify-center gap-4 mt-16 max-w-4xl mx-auto">
        {techStack.map((item, index) => (
          <motion.div
            key={item.name}
            className="glass rounded-xl px-6 py-4 glow-border flex items-center gap-3 cursor-default"
            initial={{ opacity: 0, scale: 0.8 }}
            whileInView={{ opacity: 1, scale: 1 }}
            viewport={{ once: true }}
            transition={{ delay: index * 0.03, duration: 0.2 }}
            whileHover={{ y: -4, scale: 1.05 }}
          >
            <div
              className="w-3 h-3 rounded-full"
              style={{ backgroundColor: item.color }}
            />
            <span className="font-medium text-white text-sm">{item.name}</span>
          </motion.div>
        ))}
      </div>
    </SectionWrapper>
  );
}
