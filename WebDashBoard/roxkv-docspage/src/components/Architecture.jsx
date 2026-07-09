import { motion } from 'framer-motion';
import { architectureNodes, architectureLeaves } from '../data/content';
import SectionWrapper from './SectionWrapper';

export default function Architecture() {
  return (
    <SectionWrapper id="architecture">
      <h2 className="text-3xl sm:text-4xl font-bold gradient-text text-center">
        System Architecture
      </h2>
      <p className="text-[#9898ab] text-lg text-center mt-4">
        How every query flows from user to database and back.
      </p>

      <div className="max-w-2xl mx-auto mt-16">
        {architectureNodes.map((node, index) => (
          <div key={node.label}>
            <motion.div
              className="glass rounded-xl p-5 flex items-center gap-4 glow-border"
              initial={{ opacity: 0, x: -20 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true }}
              transition={{ delay: index * 0.15 }}
            >
              <node.icon className="w-8 h-8 text-purple-400" />
              <span className="font-semibold text-white">{node.label}</span>
            </motion.div>

            {index < architectureNodes.length - 1 && (
              <div className="w-0.5 h-10 mx-auto bg-gradient-to-b from-purple-500/40 to-purple-500/10" />
            )}
          </div>
        ))}

        {/* Vertical line from last node */}
        <div className="w-0.5 h-6 mx-auto bg-purple-500/30" />

        {/* Horizontal branch line */}
        <div className="h-0.5 w-full max-w-md mx-auto bg-gradient-to-r from-transparent via-purple-500/30 to-transparent" />

        {/* Leaf nodes */}
        <div className="flex flex-wrap justify-center gap-4 mt-4">
          {architectureLeaves.map((leaf, index) => (
            <motion.div
              key={leaf.label}
              className="glass rounded-xl px-4 py-3 flex items-center gap-2 glow-border text-sm"
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: architectureNodes.length * 0.15 + index * 0.1 }}
            >
              <leaf.icon className="w-5 h-5 text-purple-400" />
              <span className="text-white">{leaf.label}</span>
            </motion.div>
          ))}
        </div>
      </div>
    </SectionWrapper>
  );
}
