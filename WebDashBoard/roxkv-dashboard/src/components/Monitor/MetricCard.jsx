import { memo } from 'react';
import MonitorGraph from './MonitorGraph';
import { useMetric } from '../../hooks/useDashboard';

const DEFAULT_FORMATTER = (value) => {
  if (!Number.isFinite(value)) return '0';
  if (Math.abs(value) >= 100) return Math.round(value).toString();
  return value.toFixed(2).replace(/\.00$/, '');
};

const MetricCard = memo(function MetricCard({ icon, label, dataKey, color, unit = '', formatter = DEFAULT_FORMATTER }) {
  // Each MetricCard subscribes independently to its specific dataKey.
  // This means if 'allocatedMemMB' changes, ONLY this specific card re-renders.
  const { current, history } = useMetric(dataKey);

  return (
    <div className="bg-white/5 border border-white/10 rounded-xl p-3 relative overflow-hidden group hover:border-white/20 transition-colors">
      <div className="flex justify-between items-start mb-2 relative z-10">
        <div className="flex items-center gap-2">
          <div className="p-1.5 rounded-lg bg-black/20" style={{ color }}>{icon}</div>
          <span className="text-xs text-text-secondary font-medium">{label}</span>
        </div>
        <div className="text-right">
          <div className="text-lg font-mono text-text-primary leading-none">
            {formatter(current)}
            {unit && <span className="text-xs text-text-muted ml-0.5">{unit}</span>}
          </div>
        </div>
      </div>
      <div className="h-10 mt-2 -mx-1">
        <MonitorGraph data={history} dataKey="value" color={color} />
      </div>
      {/* Decorative gradient glow */}
      <div 
        className="absolute -bottom-4 -right-4 w-16 h-16 rounded-full blur-2xl opacity-10 group-hover:opacity-20 transition-opacity"
        style={{ backgroundColor: color }}
      />
    </div>
  );
});

export default MetricCard;
