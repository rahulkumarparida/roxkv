import { Cpu, HardDrive, Download, Upload, Activity } from 'lucide-react';
import SystemOverview from '../components/Sidebar/SystemOverview';
import DatabaseHealth from '../components/Sidebar/DatabaseHealth';
import LLMManager from '../components/Sidebar/LLMManager';
import ChatBox from '../components/Chat/ChatBox';
import MetricCard from '../components/Monitor/MetricCard';
import ActivityFeed from '../components/Activity/ActivityFeed';
import GlassCard from '../components/Common/GlassCard';
import StatusBadge from '../components/Common/StatusBadge';

/**
 * System Monitor section in the right sidebar.
 * Since MetricCard now manages its own state via useMetric, 
 * this parent component renders exactly once.
 */
function SystemMonitor() {
  const metrics = [
    { icon: <Cpu className="h-4 w-4" />, label: 'CPU Usage', dataKey: (data) => data.cpuusage, unit: '%', color: '#f97316' },
    { icon: <Activity className="h-4 w-4" />, label: 'RAM Usage', dataKey: (data) => data.ramusage?.usedPercentge, unit: '%', color: '#3b82f6' },
    {
      icon: <HardDrive className="h-4 w-4" />,
      label: 'Disk Usage',
      dataKey: (data) => {
        const total = Number(data.diskusage?.total || 0);
        const free = Number(data.diskusage?.free || 0);
        if (!total) return 0;
        return ((total - free) / total) * 100;
      },
      unit: '%',
      color: '#22c55e'
    },
    { icon: <Download className="h-4 w-4" />, label: 'Network In', dataKey: (data) => data.networkstats?.downloadspeed, unit: 'Kbps', color: '#a855f7' },
    { icon: <Upload className="h-4 w-4" />, label: 'Network Up', dataKey: (data) => data.networkstats?.uploadspeed, unit: 'Kbps', color: '#ec4899' },
  ];

  return (
    <GlassCard className="p-4" hover>
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-xs font-semibold tracking-wider text-text-primary uppercase">System Monitor</h3>
        <StatusBadge status="live" label="REAL-TIME" />
      </div>
      <div className="space-y-2">
        {metrics.map((m) => (
          <MetricCard key={m.label} {...m} />
        ))}
      </div>
    </GlassCard>
  );
}

/**
 * Main Dashboard page composing all three columns.
 */
export default function Dashboard() {
  return (
    <div className="grid grid-cols-[310px_1fr_320px] gap-4 h-full p-4 overflow-hidden">
      {/* LEFT SIDEBAR — Metrics & Health */}
      <div className="flex flex-col gap-4 overflow-y-auto pr-1">
        <SystemOverview />
        <DatabaseHealth />
        <LLMManager />
      </div>

      {/* CENTER PANEL — Chat */}
      <div className="flex flex-col min-h-0">
        <ChatBox />
      </div>

      {/* RIGHT SIDEBAR — Telemetry & Logs */}
      <div className="flex flex-col gap-4 overflow-y-auto pl-1">
        <SystemMonitor />
        <ActivityFeed />
      </div>
    </div>
  );
}
