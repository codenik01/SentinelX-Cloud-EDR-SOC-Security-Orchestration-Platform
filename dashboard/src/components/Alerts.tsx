import { useState, useEffect } from 'react';
import { AlertTriangle, Clock, Server, Check, X, ShieldAlert, FileCode } from 'lucide-react';

interface Alert {
  id: number;
  host_id: string;
  host: { hostname: string; os: string };
  title: string;
  description: string;
  severity: string;
  matched_rule: string;
  payload: string;
  status: string;
  timestamp: string;
}

function Alerts() {
  const [alerts, setAlerts] = useState<Alert[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedAlert, setSelectedAlert] = useState<Alert | null>(null);
  const [filterSeverity, setFilterSeverity] = useState<string>('');
  const [filterStatus, setFilterStatus] = useState<string>('unresolved');

  const fetchAlerts = () => {
    let url = 'http://localhost:8080/api/v1/alerts?limit=50';
    if (filterSeverity) url += `&severity=${filterSeverity}`;
    if (filterStatus) url += `&status=${filterStatus}`;

    fetch(url)
      .then((res) => res.json())
      .then((data) => {
        setAlerts(data);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  };

  useEffect(() => {
    fetchAlerts();
    const interval = setInterval(fetchAlerts, 2000);
    return () => clearInterval(interval);
  }, [filterSeverity, filterStatus]);

  const handleResolveAlert = (id: number, status: 'resolved' | 'false_positive') => {
    fetch(`http://localhost:8080/api/v1/alerts/${id}/resolve`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ status }),
    })
      .then((res) => res.json())
      .then(() => {
        fetchAlerts();
        setSelectedAlert(null);
      });
  };

  const getSeverityStyle = (sev: string) => {
    switch (sev.toLowerCase()) {
      case 'critical':
        return 'bg-red-500/10 text-red-400 border border-red-500/30';
      case 'high':
        return 'bg-orange-500/10 text-orange-400 border border-orange-500/30';
      case 'medium':
        return 'bg-amber-500/10 text-amber-400 border border-amber-500/30';
      default:
        return 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/30';
    }
  };

  return (
    <div className="space-y-6 flex flex-col h-full">
      {/* Search & Filter Header Grid */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        {/* Filters Group */}
        <div className="flex flex-wrap items-center gap-3 text-xs font-mono">
          <select
            value={filterStatus}
            onChange={(e) => setFilterStatus(e.target.value)}
            className="bg-[#1E293B] border border-slate-800 text-slate-200 px-3 py-1.5 rounded-md focus:border-blue-500 outline-none"
          >
            <option value="unresolved">Status: Active</option>
            <option value="resolved">Status: Resolved</option>
            <option value="false_positive">Status: False Positive</option>
            <option value="">Status: All Tickets</option>
          </select>

          <select
            value={filterSeverity}
            onChange={(e) => setFilterSeverity(e.target.value)}
            className="bg-[#1E293B] border border-slate-800 text-slate-200 px-3 py-1.5 rounded-md focus:border-blue-500 outline-none"
          >
            <option value="">Severity: All Levels</option>
            <option value="critical">Severity: Critical</option>
            <option value="high">Severity: High</option>
            <option value="medium">Severity: Medium</option>
            <option value="low">Severity: Low</option>
          </select>
        </div>

        <span className="text-[10px] text-slate-400 font-mono">
          Displaying {alerts.length} historical records
        </span>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-20 flex-1">
          <Clock className="w-8 h-8 text-blue-500 animate-spin" />
        </div>
      ) : alerts.length === 0 ? (
        <div className="glass-panel p-16 rounded-xl border border-slate-800/80 text-center flex-1 max-w-lg mx-auto flex flex-col items-center justify-center">
          <ShieldAlert className="w-12 h-12 text-slate-500 mb-4" />
          <h4 className="text-md font-bold text-slate-300 mb-1">Queue is fully clear</h4>
          <p className="text-xs text-slate-500 max-w-xs">
            Excellent! No unresolved EDR threat warnings currently require manual SOC analysis.
          </p>
        </div>
      ) : (
        <div className="overflow-x-auto border border-slate-800/80 rounded-xl bg-slate-950/20">
          <table className="w-full text-left border-collapse text-xs font-mono">
            <thead>
              <tr className="border-b border-slate-800 bg-[#0E1527] text-slate-400">
                <th className="p-4 font-semibold">SEVERITY</th>
                <th className="p-4 font-semibold">SECURITY INCIDENT</th>
                <th className="p-4 font-semibold">TARGET HOST</th>
                <th className="p-4 font-semibold">CORRELATED RULE</th>
                <th className="p-4 font-semibold">TIME TRIGGERED</th>
                <th className="p-4 font-semibold">STATUS</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800 text-slate-300">
              {alerts.map((alert) => (
                <tr
                  key={alert.id}
                  onClick={() => setSelectedAlert(alert)}
                  className="hover:bg-slate-900/50 cursor-pointer transition-colors"
                >
                  <td className="p-4">
                    <span className={`px-2 py-0.5 rounded text-[10px] font-bold uppercase ${getSeverityStyle(alert.severity)}`}>
                      {alert.severity}
                    </span>
                  </td>
                  <td className="p-4">
                    <div className="font-semibold text-slate-100">{alert.title}</div>
                    <div className="text-[10px] text-slate-500 mt-0.5 truncate max-w-xs">{alert.description}</div>
                  </td>
                  <td className="p-4">
                    <div className="flex items-center space-x-1.5">
                      <Server className="w-3.5 h-3.5 text-slate-500" />
                      <span>{alert.host ? alert.host.hostname : alert.host_id}</span>
                    </div>
                  </td>
                  <td className="p-4 text-slate-400 font-semibold">{alert.matched_rule}</td>
                  <td className="p-4 text-slate-400">
                    {new Date(alert.timestamp).toLocaleTimeString()} ({new Date(alert.timestamp).toLocaleDateString()})
                  </td>
                  <td className="p-4">
                    <span className={`px-2 py-0.5 rounded-[3px] text-[10px] uppercase ${alert.status === 'unresolved' ? 'bg-red-500/10 text-red-400 border border-red-500/20 font-bold' : 'bg-slate-800 text-slate-400'}`}>
                      {alert.status === 'unresolved' ? 'unresolved' : alert.status}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* 4. Slide-Over Payload Panel */}
      {selectedAlert && (
        <div className="fixed inset-0 z-50 overflow-hidden flex justify-end bg-black/60 backdrop-blur-sm">
          <div className="w-full max-w-xl glass-panel-heavy h-full flex flex-col z-50 text-slate-300">
            {/* Slide Header */}
            <div className="p-6 border-b border-slate-800 flex items-center justify-between bg-[#0C1427]">
              <div>
                <span className={`px-2.5 py-0.5 rounded text-[10px] font-bold uppercase ${getSeverityStyle(selectedAlert.severity)}`}>
                  {selectedAlert.severity} ALERT
                </span>
                <h3 className="text-md font-bold heading-cyber text-slate-100 mt-2">{selectedAlert.title}</h3>
              </div>
              <button
                onClick={() => setSelectedAlert(null)}
                className="p-1.5 rounded-lg bg-slate-900 border border-slate-800 hover:text-white transition-colors text-slate-400"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            {/* Slide Body */}
            <div className="flex-1 overflow-y-auto p-6 space-y-6 text-xs font-mono">
              <div className="space-y-2">
                <span className="text-slate-500 uppercase tracking-widest text-[9px] block">Incident Description</span>
                <p className="bg-[#0A0F1D] border border-slate-900 p-4 rounded-lg text-slate-200 leading-relaxed">
                  {selectedAlert.description}
                </p>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="p-4 bg-slate-900/40 border border-slate-800 rounded-lg space-y-1.5">
                  <span className="text-slate-500 text-[9px] uppercase">ENDPOINT NAME</span>
                  <div className="text-slate-200 font-semibold">{selectedAlert.host ? selectedAlert.host.hostname : 'Unknown'}</div>
                </div>
                <div className="p-4 bg-slate-900/40 border border-slate-800 rounded-lg space-y-1.5">
                  <span className="text-slate-500 text-[9px] uppercase">CORRELATED EVENT ID</span>
                  <div className="text-slate-200 font-semibold">EVT-{selectedAlert.id * 13 + 120}</div>
                </div>
              </div>

              <div className="space-y-2">
                <span className="text-slate-500 uppercase tracking-widest text-[9px] flex items-center space-x-1.5">
                  <FileCode className="w-3.5 h-3.5 text-blue-400" />
                  <span>Low-Level Telemetry Event payload (JSON)</span>
                </span>
                <pre className="bg-[#0A0F1D] border border-slate-900 p-4 rounded-lg text-blue-300 overflow-x-auto text-[10px] leading-relaxed">
                  {JSON.stringify(JSON.parse(selectedAlert.payload || '{}'), null, 2)}
                </pre>
              </div>
            </div>

            {/* Slide Actions Footer */}
            {selectedAlert.status === 'unresolved' && (
              <div className="p-6 border-t border-slate-800 bg-[#0C1427] flex space-x-3 justify-end">
                <button
                  onClick={() => handleResolveAlert(selectedAlert.id, 'false_positive')}
                  className="px-4 py-2 border border-slate-800 bg-[#0A0F1D] rounded-lg text-xs hover:bg-slate-800 hover:text-white transition-colors"
                >
                  Mark False Positive
                </button>
                <button
                  onClick={() => handleResolveAlert(selectedAlert.id, 'resolved')}
                  className="flex items-center space-x-1.5 px-4 py-2 bg-blue-600 border border-blue-500 rounded-lg text-xs text-white hover:bg-blue-500 transition-colors font-semibold"
                >
                  <Check className="w-4 h-4" />
                  <span>Resolve Incident</span>
                </button>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}

export default Alerts;
