import { useState, useEffect } from 'react';
import { Terminal, Globe, Key, FileText, Activity, AlertTriangle, ShieldCheck } from 'lucide-react';

interface EventItem {
  id: number;
  host_id: string;
  timestamp: string;
  type: string; // "process", "network", "login", "file"
  payload: string;
}

function Timeline() {
  const [events, setEvents] = useState<EventItem[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchTimeline = () => {
    fetch('http://localhost:8080/api/v1/timeline?limit=50')
      .then((res) => res.json())
      .then((data) => {
        setEvents(data);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  };

  useEffect(() => {
    fetchTimeline();
    const interval = setInterval(fetchTimeline, 3000);
    return () => clearInterval(interval);
  }, []);

  const getEventMeta = (evt: EventItem) => {
    let parsed: any = {};
    try {
      parsed = JSON.parse(evt.payload || '{}');
    } catch {
      // Ignored
    }

    switch (evt.type) {
      case 'process':
        return {
          icon: Terminal,
          color: 'text-indigo-400 border-indigo-500/30 bg-indigo-500/10',
          title: 'Process Spawned',
          body: `Command: "${parsed.command}" executed by "${parsed.user}" (PID: ${parsed.pid})`,
        };
      case 'network':
        return {
          icon: Globe,
          color: 'text-sky-400 border-sky-500/30 bg-sky-500/10',
          title: 'Network Socket Outbound',
          body: `Connected to ${parsed.destination} (${parsed.ip}:${parsed.port}) via protocol ${parsed.protocol} (Process: ${parsed.process})`,
        };
      case 'login':
        const isFailed = parsed.event === 'failed_login';
        return {
          icon: Key,
          color: isFailed ? 'text-red-400 border-red-500/30 bg-red-500/10' : 'text-emerald-400 border-emerald-500/30 bg-emerald-500/10',
          title: isFailed ? 'Failed Authentication Attempt' : 'Successful System Authentication',
          body: isFailed
            ? `Failed SSH password login for user "${parsed.user}" from source IP: ${parsed.source_ip}`
            : `User "${parsed.user}" logged in successfully from source IP: ${parsed.source_ip}`,
        };
      case 'file':
        return {
          icon: FileText,
          color: 'text-teal-400 border-teal-500/30 bg-teal-500/10',
          title: 'Sensitive Registry Write',
          body: `File "${parsed.path}" experienced operation "${parsed.action}" by user "${parsed.user}"`,
        };
      default:
        return {
          icon: Activity,
          color: 'text-slate-400 border-slate-500/30 bg-slate-500/10',
          title: 'Agent Event Ingestion',
          body: 'Low-level telemetry registered into audit pipeline.',
        };
    }
  };

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center h-96 space-y-4">
        <Activity className="w-8 h-8 text-blue-500 animate-spin" />
        <span className="text-slate-400 font-mono text-sm">Streaming operational logs timeline...</span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-md font-bold heading-cyber text-slate-200">Chronological Telemetry Stream</h3>
        <p className="text-xs text-slate-400">Aggregated raw events transmitted by monitored Linux & macOS agents</p>
      </div>

      {events.length === 0 ? (
        <div className="glass-panel p-16 rounded-xl border border-slate-800/80 text-center max-w-lg mx-auto flex flex-col items-center justify-center">
          <Activity className="w-12 h-12 text-slate-600 mb-4" />
          <h4 className="text-md font-bold text-slate-300 mb-1">Queue is empty</h4>
          <p className="text-xs text-slate-500 max-w-xs leading-relaxed">
            Waiting for agent connection. Run the Go agent to begin logging low-level endpoint telemetry.
          </p>
        </div>
      ) : (
        <div className="relative border-l border-slate-800 ml-3 pl-8 space-y-6 font-mono text-xs text-slate-300">
          {events.map((evt) => {
            const meta = getEventMeta(evt);
            const Icon = meta.icon;
            return (
              <div key={evt.id} className="relative transition-transform duration-150 hover:translate-x-1">
                {/* Visual node */}
                <div className={`absolute -left-[45px] top-0 p-2 border rounded-lg ${meta.color}`}>
                  <Icon className="w-4 h-4" />
                </div>

                {/* Event Card */}
                <div className="glass-panel p-4 rounded-xl border border-[rgba(40,50,80,0.4)] space-y-2">
                  <div className="flex items-center justify-between text-[10px] text-slate-400 border-b border-slate-800/60 pb-2">
                    <div className="flex items-center space-x-2">
                      <span className="text-blue-400 font-bold">{evt.host_id}</span>
                      <span>•</span>
                      <span className="font-semibold text-slate-200 uppercase">{evt.type}</span>
                    </div>
                    <span>{new Date(evt.timestamp).toLocaleTimeString()}</span>
                  </div>
                  <h4 className="text-sm font-semibold text-slate-200">{meta.title}</h4>
                  <p className="text-slate-400 leading-relaxed">{meta.body}</p>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}

export default Timeline;
