import { motion, AnimatePresence } from 'framer-motion';
import GlassCard from '../Common/GlassCard';
import StatusBadge from '../Common/StatusBadge';
import { useActivityFeed } from '../../hooks/useDashboard';
import { cn } from '../../utils/classNames';

function ActivityItem({ event }) {
  const eventType = event.type?.toUpperCase() || 'INFO';
  let colorClass = 'text-blue-400';
  let badgeClass = 'text-blue-300 bg-blue-500/10';

  if (eventType === 'SUCCESS') {
    colorClass = 'text-green-400';
    badgeClass = 'text-green-300 bg-green-500/10';
  } else if (eventType === 'ERROR') {
    colorClass = 'text-red-400';
    badgeClass = 'text-red-300 bg-red-500/10';
  } else if (eventType === 'INFO') {
    colorClass = 'text-blue-400';
    badgeClass = 'text-blue-300 bg-blue-500/10';
  }

  return (
    <motion.div
      layout
      initial={{ opacity: 0, y: -10, height: 0 }}
      animate={{ opacity: 1, y: 0, height: 'auto' }}
      exit={{ opacity: 0, height: 0 }}
      transition={{ duration: 0.25, ease: 'easeOut' }}
      className="flex items-center gap-2 py-1.5 border-b border-white/5 last:border-0 text-xs"
    >
      <span className="text-text-muted font-mono shrink-0 w-[52px]">
        {event.timestamp?.split(', ')[1] || event.timestamp}
      </span>
      <span className="text-text-muted">{'\u203a'}</span>
      <span className={cn('flex-1 truncate', colorClass)}>{event.description}</span>
      <span className={cn(
        'shrink-0 font-mono text-[10px] px-1.5 py-0.5 rounded',
        badgeClass
      )}>
        {eventType}
      </span>
    </motion.div>
  );
}

export default function ActivityFeed() {
  const activityFeed = useActivityFeed();

  return (
    <GlassCard className="p-4 flex flex-col" hover>
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-xs font-semibold tracking-wider text-text-primary uppercase">Real-Time Activity</h3>
        <StatusBadge status="live" label="LIVE FEED" />
      </div>

      <div className="flex-1 overflow-y-auto max-h-[280px]">
        <AnimatePresence initial={false}>
          {activityFeed.map((event) => (
            <ActivityItem key={event.id} event={event} />
          ))}
        </AnimatePresence>
        {activityFeed.length === 0 && (
          <div className="text-center text-text-muted text-xs py-8">Waiting for events...</div>
        )}
      </div>
    </GlassCard>
  );
}
