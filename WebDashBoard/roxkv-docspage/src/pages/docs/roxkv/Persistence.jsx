import React from 'react';

export default function Persistence() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Persistence</h1>
        <p className="text-xl text-[#9898ab]">Data durability through JSON snapshots.</p>
      </div>
      
      <div className="section-divider" />

      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#f1f1f4]"><strong>Important:</strong> RoxKV persistence uses JSON-based snapshots. It is NOT equivalent to Redis RDB (binary format) or AOF (append-only file), and it does NOT use Go's gob encoding.</p>
      </div>

      <h2 id="snapshot-worker" className="text-2xl font-semibold text-white mb-4 mt-8">Snapshot Worker</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        The <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">SnapshotWorker</code> runs as a background goroutine on a 10-minute ticker. It periodically dumps the entire memory state to disk.
      </p>

      <h2 id="manual-trigger" className="text-2xl font-semibold text-white mb-4 mt-8">Manual Trigger (SAVE)</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        You can manually trigger a snapshot using the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">SAVE</code> command.
        <br/>
        <em>Warning:</em> This is a synchronous operation and will block other store operations during the dump.
      </p>

      <h2 id="storage-format" className="text-2xl font-semibold text-white mb-4 mt-8">Storage Format and Location</h2>
      <ul className="list-disc list-inside text-[#9898ab] space-y-2 mb-4">
        <li><strong>Format:</strong> Data is serialized using <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">json.MarshalIndent</code>. It includes all key-values and associated metadata.</li>
        <li><strong>Location:</strong> Files are saved in the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">.roxkv/</code> directory.</li>
        <li><strong>Filename:</strong> Based on the timestamp, e.g., <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">"2026-08-26 15:04:05.json"</code>.</li>
      </ul>

      <h2 id="startup-loading" className="text-2xl font-semibold text-white mb-4 mt-8">Loading on Startup</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        Upon starting, RoxKV scans the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">.roxkv/</code> directory, sorts the snapshot files by time, and loads the most recent valid JSON snapshot into memory.
      </p>

      <h2 id="lastsave-command" className="text-2xl font-semibold text-white mb-4 mt-8">Checking Snapshot Time</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        The <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">LASTSAVE</code> command returns the Unix timestamp of the last successful snapshot.
      </p>

      <div className="terminal-window">
        <div className="terminal-header">
          <div className="terminal-dot bg-red-500"></div>
          <div className="terminal-dot bg-yellow-500"></div>
          <div className="terminal-dot bg-green-500"></div>
        </div>
        <div className="terminal-body font-mono text-sm p-4 text-[#f1f1f4]">
          <div><span className="prompt text-[#7c3aed]">&gt;</span> <span className="command text-[#a78bfa]">SAVE</span></div>
          <div className="output text-[#9898ab] mb-2">OK</div>
          <div><span className="prompt text-[#7c3aed]">&gt;</span> <span className="command text-[#a78bfa]">LASTSAVE</span></div>
          <div className="output text-[#9898ab]">(integer) 1787729605</div>
        </div>
      </div>
    </div>
  );
}
