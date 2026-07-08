import { motion } from 'framer-motion';

/**
 * Bouncing dots typing indicator for agent responses.
 */
export default function TypingIndicator() {
  return (
    <div className="flex items-center gap-1.5 py-2 px-1">
      {[0, 1, 2].map((i) => (
        <motion.span
          key={i}
          className="h-1.5 w-1.5 rounded-full bg-purple-400"
          animate={{ scale: [1, 1.4, 1], opacity: [0.4, 1, 0.4] }}
          transition={{ duration: 1, repeat: Infinity, delay: i * 0.2, ease: 'easeInOut' }}
        />
      ))}
      <span className="text-xs text-text-secondary ml-2">Agent is thinking...</span>
    </div>
  );
}
