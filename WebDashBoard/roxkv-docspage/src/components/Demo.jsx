import { motion } from 'framer-motion';
import { Wrench, ShieldCheck, Activity } from 'lucide-react';
import { demo } from '../data/content';
import SectionWrapper from './SectionWrapper';

export default function Demo() {
  return (
    <SectionWrapper id="demo">
      <h2 className="gradient-text text-3xl sm:text-4xl font-bold text-center">
        {demo.heading}
      </h2>
      <p className="text-[#9898ab] text-lg text-center max-w-3xl mx-auto mt-4 leading-relaxed">
        {demo.description}
      </p>

      <motion.div
        className="glass rounded-2xl overflow-hidden max-w-5xl mx-auto mt-12 glow-border bg-gradient-to-br from-purple-900/10 to-[#0a0a0f] p-8 sm:p-10"
        initial={{ opacity: 0, y: 30 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.5 }}
      >
        <div className="flex flex-col lg:flex-row gap-8 lg:items-start">
          <div className="lg:w-1/3">
            <div className="w-16 h-16 rounded-2xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center">
              <Wrench className="w-8 h-8 text-purple-400" />
            </div>
            <h3 className="text-white text-xl font-semibold mt-5">
              Continuous work is in progress
            </h3>
            <p className="text-[#9898ab] text-sm leading-7 mt-3">
              RoxAI is being iterated on continuously to make the system more robust before the full demo experience is presented here.
            </p>
          </div>

          <div className="lg:w-2/3 grid gap-4">
            {demo.details.map((detail, index) => {
              const Icon = index === 0 ? Activity : index === 1 ? ShieldCheck : Wrench;

              return (
                <div
                  key={detail}
                  className="rounded-2xl border border-purple-500/10 bg-black/20 p-5"
                >
                  <div className="flex items-start gap-3">
                    <Icon className="w-5 h-5 text-purple-300 mt-0.5 shrink-0" />
                    <p className="text-sm sm:text-base text-[#c7c7d4] leading-7">{detail}</p>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </motion.div>
    </SectionWrapper>
  );
}
