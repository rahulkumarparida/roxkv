import { Database } from 'lucide-react';
import GlassCard from '../Common/GlassCard';
import StatusBadge from '../Common/StatusBadge';
import { useDatabaseHealth } from '../../hooks/useDashboard';
import { formatBytes } from '../../utils/formatBytes';
import { formatDateDisplay } from '../../utils/formatTime';

function MetricRow({ label, children }) {
  return (
    <div className="flex items-center justify-between py-1.5 border-b border-white/5 last:border-0">
      <span className="text-sm text-text-secondary">{label}</span>
      <span className="text-sm font-mono text-text-primary">{children}</span>
    </div>
  );
}

export default function DatabaseHealth() {
  const {
    persistence,
    storageHealth,
    snapshotCount,
    databaseSizeBytes,
    permanentKeys,
    pubsubTopicsCount,
    latestSnapshotTimestamp,
  } = useDatabaseHealth();

  return (
    <GlassCard className="p-4" hover>
      <div className="flex items-center gap-2 mb-3">
        <Database className="h-4 w-4 text-purple-400" />
        <h3 className="text-xs font-semibold tracking-wider text-text-primary uppercase">Database Health</h3>
      </div>

      <div className="space-y-0">
        <MetricRow label="Persistence"><StatusBadge status={persistence?.toLowerCase()} label={persistence} dot={false} /></MetricRow>
        <MetricRow label="Snapshots">{snapshotCount}</MetricRow>
        <MetricRow label="Latest Snapshot">{latestSnapshotTimestamp ? formatDateDisplay(latestSnapshotTimestamp) : '--'}</MetricRow>
        <MetricRow label="Database Size">{formatBytes(databaseSizeBytes)}</MetricRow>
        <MetricRow label="Storage Health"><StatusBadge status={storageHealth?.toLowerCase()} label={storageHealth} dot={false} /></MetricRow>
        <MetricRow label="Permanent Keys"><span className="text-green-400">{permanentKeys}</span></MetricRow>
        <MetricRow label="Pub/Sub Topics"><span className="text-text-muted">{pubsubTopicsCount}</span></MetricRow>
      </div>
    </GlassCard>
  );
}
