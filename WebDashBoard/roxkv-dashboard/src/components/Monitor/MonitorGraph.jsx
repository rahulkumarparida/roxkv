import { memo } from 'react';
import { AreaChart, Area, ResponsiveContainer } from 'recharts';

/**
 * Reusable sparkline graph component using Recharts.
 * Memoized to prevent unnecessary re-renders on parent updates.
 */
const MonitorGraph = memo(function MonitorGraph({ data, color, height = 40 }) {
  if (!data || data.length === 0) return null;

  const gradientId = `gradient-${color.replace('#', '')}-${Math.random().toString(36).slice(2, 6)}`;

  return (
    <ResponsiveContainer width="100%" height={height}>
      <AreaChart data={data} margin={{ top: 0, right: 0, bottom: 0, left: 0 }}>
        <defs>
          <linearGradient id={gradientId} x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor={color} stopOpacity={0.3} />
            <stop offset="100%" stopColor={color} stopOpacity={0.0} />
          </linearGradient>
        </defs>
        <Area
          type="monotone"
          dataKey="value"
          stroke={color}
          strokeWidth={1.5}
          fill={`url(#${gradientId})`}
          isAnimationActive={false}
          dot={false}
        />
      </AreaChart>
    </ResponsiveContainer>
  );
});

export default MonitorGraph;
