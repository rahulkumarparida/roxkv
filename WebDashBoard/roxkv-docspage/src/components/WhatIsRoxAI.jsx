import { motion } from 'framer-motion';
import { whatIsRoxAI } from '../data/content';
import SectionWrapper from './SectionWrapper';

export default function WhatIsRoxAI() {
  return (
    <SectionWrapper id="what-is-roxai">
      <div className="lg:grid lg:grid-cols-2 gap-16 items-center">
        {/* Left column – heading & bullet points */}
        <div>
          <motion.h2
            className="text-3xl sm:text-4xl font-bold gradient-text"
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.5 }}
          >
            What is RoxAI?
          </motion.h2>

          <div className="mt-8 space-y-1">
            {whatIsRoxAI.points.map((point, index) => (
              <motion.div
                key={index}
                className="flex items-start gap-3 mb-3"
                initial={{ opacity: 0, x: -20 }}
                whileInView={{ opacity: 1, x: 0 }}
                viewport={{ once: true }}
                transition={{ duration: 0.4, delay: index * 0.1 }}
              >
                <span className="w-2 h-2 rounded-full bg-purple-500 mt-2 flex-shrink-0" />
                <span className="text-[#9898ab]">{point}</span>
              </motion.div>
            ))}
          </div>
        </div>

        {/* Right column – vertical workflow diagram */}
        <div className="mt-12 lg:mt-0">
          {whatIsRoxAI.workflow.map((step, index) => (
            <div key={index}>
              <motion.div
                className="glass rounded-xl p-4 flex items-center gap-4 glow-border"
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ duration: 0.4, delay: 0.15 * index }}
              >
                <step.icon className="w-10 h-10 text-purple-400" />
                <span className="font-semibold text-white">{step.label}</span>
              </motion.div>

              {/* Connector line between steps */}
              {index < whatIsRoxAI.workflow.length - 1 && (
                <div className="w-0.5 h-8 mx-auto bg-gradient-to-b from-purple-500/50 to-purple-500/10" />
              )}
            </div>
          ))}
        </div>
      </div>
    </SectionWrapper>
  );
}
