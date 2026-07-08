import { motion } from 'framer-motion';
import { cn } from '../../utils/classNames';

/**
 * Reusable glassmorphism card component.
 * @param {Object} props
 * @param {React.ReactNode} props.children
 * @param {string} [props.className] - Additional classes
 * @param {boolean} [props.hover=false] - Enable hover lift effect
 * @param {boolean} [props.pulse=false] - Enable border pulse animation
 */
export default function GlassCard({ children, className, hover = false, pulse = false, ...rest }) {
  return (
    <motion.div
      className={cn(
        'bg-white/5 backdrop-blur-lg border border-white/10 rounded-2xl',
        'transition-all duration-300 ease-out',
        hover && 'hover:-translate-y-0.5 hover:shadow-[0_0_15px_rgba(168,85,247,0.3)] hover:border-purple-500/30',
        pulse && 'animate-border-pulse',
        className
      )}
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4, ease: 'easeOut' }}
      {...rest}
    >
      {children}
    </motion.div>
  );
}
