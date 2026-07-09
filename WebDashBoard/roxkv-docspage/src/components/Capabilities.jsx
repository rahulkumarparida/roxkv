import { motion } from 'framer-motion';
import { capabilities } from '../data/content';
import SectionWrapper from './SectionWrapper';

export default function Capabilities() {
  return (
    <SectionWrapper id="capabilities">
      <motion.h2
        className="text-3xl sm:text-4xl font-bold gradient-text text-center"
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.5 }}
      >
        Capabilities
      </motion.h2>

      <motion.p
        className="text-[#9898ab] text-lg text-center mt-4 max-w-2xl mx-auto"
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.5, delay: 0.1 }}
      >
        Everything you need to manage, monitor, and interact with your database.
      </motion.p>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5 mt-16">
        {capabilities.map((cap, index) => (
          <motion.div
            key={index}
            className="glass rounded-2xl p-6 glow-border"
            whileHover={{ y: -6 }}
            transition={{ duration: 0.3 }}
            initial={{ opacity: 0, y: 30 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
          >
            <cap.icon className="w-10 h-10 text-purple-400 mb-4" />
            <h3 className="text-lg font-semibold text-white">{cap.title}</h3>
            <p className="text-sm text-[#9898ab] mt-2 leading-relaxed">
              {cap.description}
            </p>
          </motion.div>
        ))}
      </div>
    </SectionWrapper>
  );
}
