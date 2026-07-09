import { useRef, useState, useEffect } from 'react';
import { motion, useInView } from 'framer-motion';
import { stats } from '../data/content';
import SectionWrapper from './SectionWrapper';

function CountUp({ target, suffix }) {
  const [count, setCount] = useState(0);
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true });

  useEffect(() => {
    if (!isInView) return;
    let start = 0;
    const duration = 2000;
    const startTime = performance.now();

    function animate(currentTime) {
      const elapsed = currentTime - startTime;
      const progress = Math.min(elapsed / duration, 1);
      // Ease out cubic
      const eased = 1 - Math.pow(1 - progress, 3);
      setCount(Math.floor(eased * target));
      if (progress < 1) requestAnimationFrame(animate);
    }
    requestAnimationFrame(animate);
  }, [isInView, target]);

  return <span ref={ref}>{count}{suffix}</span>;
}

export default function Stats() {
  return (
    <SectionWrapper id="stats">
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 sm:gap-6">
        {stats.map((stat, index) => (
          <motion.div
            key={stat.label}
            initial={{ opacity: 0, y: 30 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{
              duration: 0.7,
              delay: index * 0.1,
              ease: [0.25, 0.4, 0.25, 1],
            }}
            className="glass rounded-2xl p-6 text-center glow-border"
          >
            <div className="text-3xl sm:text-4xl font-bold gradient-text-purple">
              <CountUp target={stat.value} suffix={stat.suffix} />
            </div>
            <p className="text-sm text-[#9898ab] mt-2">{stat.label}</p>
          </motion.div>
        ))}
      </div>
    </SectionWrapper>
  );
}
