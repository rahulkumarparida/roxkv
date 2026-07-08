import { History, Save } from 'lucide-react';
import GlassCard from '../Common/GlassCard';
import LiveValue from '../Common/LiveValue';
import { useDatabaseHealth } from '../../hooks/useDashboard';
import { formatBytes } from '../../utils/formatBytes';
import { formatCountdown } from '../../utils/formatTime';

export default function SnapshotInfo() {
  const { totalSnapshots, totalSnapshotSizeBytes, latestSnapshotTimestamp } = useDatabaseHealth();

  return (
    <GlassCard className="p-4" hover>
      <div className="flex items-center gap-2 mb-4">
        <History className="h-4 w-4 text-purple-400" />
        <h3 className="text-xs font-semibold tracking-wider text-text-primary uppercase">Snapshot Status</h3>
      </div>
      
      <div className="space-y-4">
        <div className="flex justify-between items-end">
          <div>
            <div className="text-[10px] text-text-muted uppercase tracking-wider mb-1">Next Auto-Snapshot</div>
            <div className="text-2xl font-mono text-text-primary">
              <LiveValue 
                channel="storage" 
                selector={d => formatCountdown(Math.max(0, Math.floor(Number(d.nextSnap) || 0)))}
                fallback="00:00:00" 
              />
            </div>
          </div>
          <button className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-purple-500/20 text-purple-300 text-xs hover:bg-purple-500/30 transition-colors border border-purple-500/30">
            <Save className="h-3 w-3" /> Force
          </button>
        </div>

        <div className="grid grid-cols-2 gap-2 pt-3 border-t border-white/5">
          <div>
            <div className="text-[10px] text-text-muted mb-0.5">Total Snapshots</div>
            <div className="text-xs font-mono text-text-secondary">{totalSnapshots}</div>
          </div>
          <div>
            <div className="text-[10px] text-text-muted mb-0.5">Total Size</div>
            <div className="text-xs font-mono text-text-secondary">{formatBytes(totalSnapshotSizeBytes)}</div>
          </div>
        </div>
        
        {latestSnapshotTimestamp && (
          <div className="pt-2">
            <div className="text-[10px] text-text-muted mb-0.5">Latest Snapshot</div>
            <div className="text-[10px] font-mono text-text-secondary opacity-70">
              {new Date(latestSnapshotTimestamp).toLocaleString()}
            </div>
          </div>
        )}
      </div>
    </GlassCard>
  );
}
