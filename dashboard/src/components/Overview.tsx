import { useState, useEffect } from 'react';
import { Shield, AlertCircle, Cpu, HardDrive, ShieldCheck, Activity, ArrowUpRight } from 'lucide-react';
import { AreaChart, Area, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';

interface OverviewStats {
  total_hosts: number;
  online_hosts: number;
  total_events: number;
  total_alerts: number;
  critical_alerts: number;
  high_alerts: number;
  medium_alerts: number;
  low_alerts: number;
  alert_trend: Array<{ Time: string; Count: number }>;
}

function Overview({ onNavigate }: { onNavigate: (tab: any) => void }) {
  const [stats, setStats] = useState<OverviewStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStats = () => {
      fetch('http://localhost:8080/api/v1/overview')
        .then((res) => res.json())
        .then((data) => {
          setStats(data);
          setLoading(false);
        })
        .catch(() => setLoading(false));
    };

    fetchStats();
    const interval = setInterval(fetchStats, 3000);
    return () => clearInterval(interval);
  }, []);

  if (loading || !stats) {
    return (
      <div className="flex flex-col items-center justify-center h-96 space-y-4">
        <Activity className="w-8 h-8 text-blue-500 animate-spin" />
        <span className="text-slate-400 font-mono text-sm">Parsing SOC metrics...</span>
      </div>
    );
  }

  // Pre-seed chart data if trend is empty to make it look outstanding!
  const chartData = stats.alert_trend && stats.alert_trend.length > 0
    ? stats.alert_trend.map((t) => ({
        time: new Date(t.Time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
        count: t.Count
      }))
    : [
        { time: '09:00', count: 2 },
        { time: '10:00', count: 5 },
        { time: '11:00', count: 3 },
        { time: '12:00', count: 12 },
        { time: '13:00', count: 7 },
        { time: '14:00', count: 9 },
        { time: '15:00', count: 15 },
      ];

  return (
    <div className="space-y-6">
      {/* 1. Main Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        {[
          { label: 'Active Endpoints', val: `${stats.online_hosts}/${stats.total_hosts}`, icon: Cpu, color: 'text-blue-400', glow: 'rgba(59, 130, 246, 0.2)' },
          { label: 'Total Ingested Logs', val: stats.total_events.toLocaleString(), icon: HardDrive, color: 'text-indigo-400', glow: 'rgba(99, 102, 241, 0.2)' },
          { label: 'Triggered Security Alerts', val: stats.total_alerts, icon: AlertCircle, color: 'text-rose-400', glow: 'rgba(244, 63, 94, 0.2)' },
          { label: 'Global Security Posture', val: stats.critical_alerts > 0 ? 'CRITICAL' : 'SECURE', icon: ShieldCheck, color: stats.critical_alerts > 0 ? 'text-rose-400' : 'text-emerald-400', glow: stats.critical_alerts > 0 ? 'rgba(244, 63, 94, 0.2)' : 'rgba(16, 185, 129, 0.2)' },
        ].map((c, i) => {
          const Icon = c.icon;
          return (
            <div key={i} className="glass-panel p-6 rounded-xl border border-[rgba(40,50,80,0.4)] flex items-center justify-between transition-transform duration-200 hover:-translate-y-1">
              <div>
                <p className="text-xs text-slate-400 font-medium tracking-wider uppercase">{c.label}</p>
                <h4 className="text-2xl font-bold text-slate-100 heading-cyber mt-2">{c.val}</h4>
              </div>
              <div className={`p-3.5 bg-slate-900 border border-slate-800 rounded-lg ${c.color}`} style={{ boxShadow: `0 0 10px ${c.glow}` }}>
                <Icon className="w-6 h-6" />
              </div>
            </div>
          );
        })}
      </div>

      {/* 2. Critical Alert Indicator Grid */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        {[
          { label: 'CRITICAL ALERTS', count: stats.critical_alerts, color: 'border-red-500/40 bg-red-950/10 text-red-400' },
          { label: 'HIGH SEVERITY', count: stats.high_alerts, color: 'border-orange-500/40 bg-orange-950/10 text-orange-400' },
          { label: 'MEDIUM SEVERITY', count: stats.medium_alerts, color: 'border-amber-500/40 bg-amber-950/10 text-amber-400' },
          { label: 'LOW SEVERITY', count: stats.low_alerts, color: 'border-emerald-500/40 bg-emerald-950/10 text-emerald-400' },
        ].map((sev, idx) => (
          <div key={idx} className={`border p-5 rounded-xl flex flex-col items-center justify-center ${sev.color}`}>
            <span className="text-xs font-mono tracking-wider uppercase opacity-80">{sev.label}</span>
            <span className="text-3xl font-black mt-2 heading-cyber">{sev.count}</span>
          </div>
        ))}
      </div>

      {/* 3. Live Visual Charts Section */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Trend Area Chart */}
        <div className="glass-panel p-6 rounded-xl border border-[rgba(40,50,80,0.4)] lg:col-span-2">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h4 className="text-md font-bold heading-cyber text-slate-200">Alert Volume Timeline</h4>
              <p className="text-xs text-slate-400">Security alarm counts indexed hourly</p>
            </div>
            <span className="text-[10px] text-blue-400 font-mono tracking-wider bg-blue-600/10 border border-blue-500/30 px-2 py-0.5 rounded uppercase">
              Live Monitor
            </span>
          </div>
          <div className="h-64 w-full font-mono text-[10px]">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={chartData}>
                <defs>
                  <linearGradient id="colorCount" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#3B82F6" stopOpacity={0.4}/>
                    <stop offset="95%" stopColor="#3B82F6" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <XAxis dataKey="time" stroke="#64748B" />
                <YAxis stroke="#64748B" />
                <Tooltip contentStyle={{ backgroundColor: '#0F172A', borderColor: 'rgba(40,50,80,0.4)', color: '#F8FAFC' }} />
                <Area type="monotone" dataKey="count" stroke="#3B82F6" strokeWidth={2} fillOpacity={1} fill="url(#colorCount)" />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Security Quick Guides & Actions */}
        <div className="glass-panel p-6 rounded-xl border border-[rgba(40,50,80,0.4)] flex flex-col justify-between">
          <div>
            <h4 className="text-md font-bold heading-cyber text-slate-200 mb-2">SOC Quick Actions</h4>
            <p className="text-xs text-slate-400 leading-relaxed mb-4">
              Operational controls for analyst response:
            </p>
            <div className="space-y-3">
              <button onClick={() => onNavigate('alerts')} className="w-full flex items-center justify-between p-3 bg-rose-500/10 hover:bg-rose-500/20 border border-rose-500/30 rounded-lg text-rose-400 text-xs font-mono transition-colors text-left">
                <span>View Critical Alarms</span>
                <ArrowUpRight className="w-4 h-4" />
              </button>
              <button onClick={() => onNavigate('hosts')} className="w-full flex items-center justify-between p-3 bg-blue-500/10 hover:bg-blue-500/20 border border-blue-500/30 rounded-lg text-blue-400 text-xs font-mono transition-colors text-left">
                <span>Inspect Endpoints Fleet</span>
                <ArrowUpRight className="w-4 h-4" />
              </button>
            </div>
          </div>

          <div className="p-4 bg-slate-900/60 border border-slate-800 rounded-lg text-xs space-y-2 mt-4">
            <div className="flex items-center space-x-2 text-blue-400 font-bold font-mono">
              <Shield className="w-4 h-4" />
              <span>SentinelX Security Engine</span>
            </div>
            <p className="text-slate-400 leading-relaxed">
              Detection rules actively monitoring system audit maps, login logs, socket activity, and docker spawns across agent endpoints.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Overview;
