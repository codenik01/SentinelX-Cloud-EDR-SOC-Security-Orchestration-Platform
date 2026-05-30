import { useState, useEffect } from 'react';
import { Cpu, Terminal, Clock, RefreshCw, AlertTriangle, MonitorPlay } from 'lucide-react';

interface Host {
  host_id: string;
  hostname: string;
  os: string;
  status: string; // "online", "offline"
  last_heartbeat: string;
  created_at: string;
}

function Hosts() {
  const [hosts, setHosts] = useState<Host[]>([]);
  const [loading, setLoading] = useState(true);
  const [isRefreshing, setIsRefreshing] = useState(false);

  const fetchHosts = () => {
    setIsRefreshing(true);
    fetch('http://localhost:8080/api/v1/hosts')
      .then((res) => res.json())
      .then((data) => {
        setHosts(data);
        setLoading(false);
        setIsRefreshing(false);
      })
      .catch(() => {
        setLoading(false);
        setIsRefreshing(false);
      });
  };

  useEffect(() => {
    fetchHosts();
    const interval = setInterval(fetchHosts, 5000);
    return () => clearInterval(interval);
  }, []);

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center h-96 space-y-4">
        <RefreshCw className="w-8 h-8 text-blue-500 animate-spin" />
        <span className="text-slate-400 font-mono text-sm">Scanning endpoint nodes...</span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Overview stats */}
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-md font-bold heading-cyber text-slate-200">Registered Endpoints Fleet</h3>
          <p className="text-xs text-slate-400">Total systems reporting telemetry inside the EDR ecosystem</p>
        </div>
        <button
          onClick={fetchHosts}
          className="flex items-center space-x-2 px-3 py-1.5 text-xs bg-[#1E293B] border border-slate-800 rounded-md hover:bg-slate-700/60 transition-colors text-slate-200"
        >
          <RefreshCw className={`w-3.5 h-3.5 ${isRefreshing ? 'animate-spin' : ''}`} />
          <span>Refresh Fleet</span>
        </button>
      </div>

      {hosts.length === 0 ? (
        <div className="glass-panel p-12 rounded-xl border border-slate-800/80 text-center max-w-lg">
          <AlertTriangle className="w-12 h-12 text-amber-500 mx-auto mb-4" />
          <h4 className="text-md font-bold text-slate-200 mb-2">No Registered Endpoints Found</h4>
          <p className="text-xs text-slate-400 leading-relaxed mb-6">
            Make sure your Go Endpoint Agent daemon is built and currently running. Once launched, it will automatically register heartbeats.
          </p>
          <div className="p-3 bg-[#0A0F1D] border border-slate-800 rounded-lg text-[10px] font-mono text-slate-300 text-left">
            cd agent && go run .
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {hosts.map((host) => {
            const isOnline = host.status === 'online';
            return (
              <div key={host.host_id} className="glass-panel rounded-xl border border-[rgba(40,50,80,0.4)] overflow-hidden flex flex-col justify-between hover:border-blue-500/50 transition-colors">
                {/* Host Title */}
                <div className="p-6 border-b border-[rgba(40,50,80,0.2)] bg-slate-900/35">
                  <div className="flex items-center justify-between mb-4">
                    <span className="text-[10px] text-blue-400 font-mono tracking-wider bg-blue-600/10 border border-blue-500/20 px-2 py-0.5 rounded uppercase">
                      {host.os}
                    </span>
                    <span className="flex items-center space-x-1.5">
                      <span className={`w-2 h-2 rounded-full ${isOnline ? 'bg-emerald-500 shadow-[0_0_6px_rgba(16,185,129,0.5)] blink-indicator' : 'bg-slate-500'}`} />
                      <span className={`text-[10px] font-mono font-bold ${isOnline ? 'text-emerald-400' : 'text-slate-400'}`}>
                        {isOnline ? 'ONLINE' : 'OFFLINE'}
                      </span>
                    </span>
                  </div>
                  <h4 className="text-lg font-bold text-slate-200 heading-cyber truncate">{host.hostname}</h4>
                  <p className="text-[10px] text-slate-400 font-mono mt-1">ID: {host.host_id}</p>
                </div>

                {/* Info Panel */}
                <div className="p-6 space-y-4 text-xs font-mono text-slate-300 flex-1">
                  <div className="flex items-center justify-between">
                    <span className="text-slate-500 flex items-center space-x-1">
                      <Terminal className="w-3.5 h-3.5" />
                      <span>Node IP:</span>
                    </span>
                    <span>192.168.1.15</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-slate-500 flex items-center space-x-1">
                      <Clock className="w-3.5 h-3.5" />
                      <span>Last Seen:</span>
                    </span>
                    <span>{new Date(host.last_heartbeat).toLocaleTimeString()}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-slate-500 flex items-center space-x-1">
                      <Cpu className="w-3.5 h-3.5" />
                      <span>Reporting Type:</span>
                    </span>
                    <span className="text-blue-400 uppercase">{host.os === 'linux' ? 'Audit Engine' : 'Simulation Mode'}</span>
                  </div>
                </div>

                {/* Operations Footer */}
                <div className="p-4 bg-slate-900/60 border-t border-[rgba(40,50,80,0.2)] flex items-center justify-between">
                  <span className="text-[10px] text-slate-400 font-mono">Agent Active</span>
                  <div className="flex items-center space-x-2 text-[10px] text-emerald-400 bg-emerald-500/10 border border-emerald-500/20 px-2.5 py-1 rounded">
                    <MonitorPlay className="w-3.5 h-3.5" />
                    <span>Telemetry Syncing</span>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}

export default Hosts;
