import { useState, useEffect } from 'react';
import { Shield, AlertTriangle, Cpu, List, FileText, Database, Activity, RefreshCw } from 'lucide-react';

function App() {
  const [activeTab, setActiveTab] = useState<'overview' | 'alerts' | 'hosts' | 'timeline' | 'threatintel'>('overview');
  const [apiStatus, setApiStatus] = useState<'connected' | 'disconnected'>('disconnected');
  const [isSyncing, setIsSyncing] = useState(false);

  // Check backend server health status
  useEffect(() => {
    const checkHealth = () => {
      fetch('http://localhost:8080/health')
        .then((res) => (res.ok ? setApiStatus('connected') : setApiStatus('disconnected')))
        .catch(() => setApiStatus('disconnected'));
    };

    checkHealth();
    const interval = setInterval(checkHealth, 5000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="min-h-screen bg-[#0A0F1D] text-slate-100 font-sans flex overflow-hidden">
      {/* 1. Sidebar Panel */}
      <aside className="w-64 bg-[#0D1527] border-r border-[rgba(40,50,80,0.4)] flex flex-col z-20">
        {/* Brand Header */}
        <div className="p-6 border-b border-[rgba(40,50,80,0.4)] flex items-center space-x-3">
          <div className="p-2 bg-blue-600/20 border border-blue-500/50 rounded-lg text-blue-400 animate-pulse-slow">
            <Shield className="w-6 h-6" />
          </div>
          <div>
            <h1 className="text-lg font-bold tracking-wider text-slate-100 heading-cyber">SENTINELX</h1>
            <p className="text-[10px] text-blue-400 font-mono tracking-widest uppercase">EDR + SOC Command</p>
          </div>
        </div>

        {/* Navigation Items */}
        <nav className="flex-1 px-4 py-6 space-y-2">
          {[
            { id: 'overview', name: 'SOC Overview', icon: Activity },
            { id: 'alerts', name: 'Active Alerts', icon: AlertTriangle },
            { id: 'hosts', name: 'Endpoints Fleet', icon: Cpu },
            { id: 'timeline', name: 'System Logs', icon: List },
            { id: 'threatintel', name: 'Threat Intel Feed', icon: Database },
          ].map((item) => {
            const Icon = item.icon;
            const isActive = activeTab === item.id;
            return (
              <button
                key={item.id}
                onClick={() => setActiveTab(item.id as any)}
                className={`w-full flex items-center space-x-3 px-4 py-3 rounded-lg text-sm transition-all duration-200 text-left ${
                  isActive
                    ? 'bg-blue-600/15 border-l-4 border-blue-500 text-blue-400 font-medium'
                    : 'text-slate-400 hover:bg-slate-800/40 hover:text-slate-200'
                }`}
              >
                <Icon className="w-5 h-5" />
                <span>{item.name}</span>
              </button>
            );
          })}
        </nav>

        {/* Footer Status Indicators */}
        <div className="p-4 border-t border-[rgba(40,50,80,0.4)] bg-[#0A0F1D]/50 text-xs font-mono space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-slate-400">Agent API:</span>
            <span className="flex items-center space-x-1.5">
              <span className={`w-2.5 h-2.5 rounded-full ${apiStatus === 'connected' ? 'bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)] blink-indicator' : 'bg-rose-500 shadow-[0_0_8px_rgba(244,63,94,0.5)]'}`} />
              <span className={apiStatus === 'connected' ? 'text-emerald-400' : 'text-rose-400'}>
                {apiStatus === 'connected' ? 'ONLINE' : 'OFFLINE'}
              </span>
            </span>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-slate-400">System Time:</span>
            <span className="text-slate-300">GMT+5:30</span>
          </div>
        </div>
      </aside>

      {/* 2. Main Content Viewport */}
      <main className="flex-1 flex flex-col overflow-hidden relative">
        {/* Top Header Grid */}
        <header className="h-16 border-b border-[rgba(40,50,80,0.4)] px-8 flex items-center justify-between bg-[#0A0F1D]/80 backdrop-blur-md z-10">
          <div>
            <h2 className="text-lg font-bold text-slate-100 heading-cyber capitalize">
              {activeTab === 'threatintel' ? 'Threat Intelligence' : activeTab}
            </h2>
            <p className="text-xs text-slate-400">Real-time endpoint security status & logs feed</p>
          </div>

          {/* Sync Trigger Action */}
          <div className="flex items-center space-x-4">
            <button
              onClick={() => {
                setIsSyncing(true);
                setTimeout(() => setIsSyncing(false), 800);
              }}
              className="flex items-center space-x-2 px-3 py-1.5 text-xs bg-[#1E293B] border border-[rgba(40,50,80,0.4)] rounded-md hover:bg-slate-700/60 transition-colors"
            >
              <RefreshCw className={`w-3.5 h-3.5 ${isSyncing ? 'animate-spin text-blue-400' : ''}`} />
              <span>Force Synchronize</span>
            </button>

            <div className="flex items-center space-x-2 text-xs border-l border-slate-700 pl-4">
              <span className="text-slate-400">Analyst:</span>
              <span className="font-semibold text-blue-400">codenik</span>
            </div>
          </div>
        </header>

        {/* View Switcher viewport */}
        <div className="flex-1 overflow-y-auto p-8 relative">
          <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-blue-900/5 via-transparent to-transparent pointer-events-none" />
          
          {activeTab === 'overview' && (
            <div className="space-y-6">
              <div className="glass-panel p-8 rounded-xl border border-[rgba(40,50,80,0.4)] max-w-2xl">
                <h3 className="text-xl font-bold heading-cyber text-blue-400 mb-2">SentinelX SOC Terminal</h3>
                <p className="text-sm text-slate-300 leading-relaxed mb-4">
                  Theme and layout loaded. Dashboard foundations initialized. Component models will link database pipelines on the next commit cycles.
                </p>
                <div className="flex items-center space-x-3 text-xs font-mono bg-[#0A0F1D]/80 p-3 rounded-lg border border-slate-800 text-slate-400">
                  <span className="text-blue-400 font-bold">$</span>
                  <span>sentinelxctl --status --verbose</span>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'alerts' && <div className="text-slate-400 font-mono">Alerts component rendering initialized...</div>}
          {activeTab === 'hosts' && <div className="text-slate-400 font-mono">Host network tracking pending backend connection...</div>}
          {activeTab === 'timeline' && <div className="text-slate-400 font-mono">Loading telemetry stream logs...</div>}
          {activeTab === 'threatintel' && <div className="text-slate-400 font-mono">Threat feed query pool mapping...</div>}
        </div>
      </main>
    </div>
  );
}

export default App;
