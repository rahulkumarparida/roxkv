import { useState } from 'react';
import { Globe, Server, Users, Clock, Bot, ChevronDown, Cpu, MonitorCog } from 'lucide-react';
import GlassCard from '../Common/GlassCard';
import StatusBadge from '../Common/StatusBadge';
import LiveValue from '../Common/LiveValue';
import { formatUptime } from '../../utils/formatTime';

export default function SystemOverview() {
  const [showDetails, setShowDetails] = useState(false);

  return (
    <GlassCard className="p-4 flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <h2 className="text-sm font-bold tracking-wider text-text-primary uppercase flex items-center gap-2">
          <Server className="h-4 w-4 text-purple-400" /> System Overview
        </h2>
        <StatusBadge status="active" label="ACTIVE" />
      </div>
      <div className="grid grid-cols-2 gap-3">
        <div className="bg-white/5 rounded-xl p-3 border border-white/10">
          <div className="text-xs text-text-muted mb-1 flex items-center gap-1.5"><Globe className="h-3 w-3" /> Hostname</div>
          <div className="font-mono text-sm text-text-primary">
            <LiveValue channel="monitor" selector={d => d.machineinfo?.user} fallback="roxxd" />
          </div>
        </div>
        <div className="bg-white/5 rounded-xl p-3 border border-white/10">
          <div className="text-xs text-text-muted mb-1 flex items-center gap-1.5"><Users className="h-3 w-3" /> Connected</div>
          <div className="font-mono text-sm text-purple-400">
            <LiveValue channel="monitor" selector={d => d.machineinfo?.connectedclient} fallback={0} />
          </div>
        </div>
        <div className="bg-white/5 rounded-xl p-3 border border-white/10">
          <div className="text-xs text-text-muted mb-1 flex items-center gap-1.5"><Clock className="h-3 w-3" /> Uptime</div>
          <div className="font-mono text-sm text-text-primary">
            <LiveValue channel="monitor" selector={d => formatUptime(d.machineinfo?.serverUptime || 0)} fallback="0s" />
          </div>
        </div>
        <div className="bg-white/5 rounded-xl p-3 border border-white/10">
          <div className="text-xs text-text-muted mb-1 flex items-center gap-1.5"><Bot className="h-3 w-3" /> Master Agent</div>
          <div className="font-mono text-sm text-green-400">Online</div>
        </div>
      </div>

      <div className="rounded-xl border border-white/10 bg-white/5 overflow-hidden">
        <button
          type="button"
          onClick={() => setShowDetails((prev) => !prev)}
          className="w-full flex items-center justify-between px-3 py-2.5 text-left text-xs font-semibold tracking-wide text-text-secondary uppercase hover:bg-white/5 transition-colors"
        >
          <span className="flex items-center gap-2">
            <MonitorCog className="h-3.5 w-3.5 text-blue-400" />
            System Details
          </span>
          <ChevronDown className={`h-4 w-4 text-text-muted transition-transform ${showDetails ? 'rotate-180' : ''}`} />
        </button>

        {showDetails && (
          <div className="grid grid-cols-2 gap-3 px-3 pb-3 pt-1 border-t border-white/10">
            <div>
              <div className="text-[10px] text-text-muted mb-1 flex items-center gap-1.5">
                <Cpu className="h-3 w-3" /> CPUs
              </div>
              <div className="font-mono text-sm text-text-primary">
                <LiveValue channel="monitor" selector={(d) => d.machineinfo?.cpus} fallback="0" />
              </div>
            </div>
            <div>
              <div className="text-[10px] text-text-muted mb-1 flex items-center gap-1.5">
                <Server className="h-3 w-3" /> Arch
              </div>
              <div className="font-mono text-sm text-text-primary">
                <LiveValue channel="monitor" selector={(d) => d.machineinfo?.arch} fallback="unknown" />
              </div>
            </div>
          </div>
        )}
      </div>
    </GlassCard>
  );
}
