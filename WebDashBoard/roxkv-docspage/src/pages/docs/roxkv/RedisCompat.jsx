export default function RedisCompat() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Redis Compatibility</h1>
        <p className="text-xl text-[#9898ab]">Supported commands and differences from standard Redis</p>
      </div>
      <div className="section-divider" />
      
      <h2 id="supported-commands" className="text-2xl font-bold text-white mt-12 mb-4">Supported Command Categories</h2>
      <p className="text-[#f1f1f4] mb-4">
        RoxKV supports a core subset of Redis commands designed to work out-of-the-box with existing Redis clients and standard connection flows.
      </p>
      <div className="overflow-x-auto rounded-xl border border-[#1e1e2e] bg-[#16161f] mb-8">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="border-b border-[#1e1e2e] bg-[rgba(124,58,237,0.05)]">
              <th className="py-3 px-4 text-sm font-semibold text-[#a78bfa]">Category</th>
              <th className="py-3 px-4 text-sm font-semibold text-[#a78bfa]">Commands</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-[#1e1e2e]">
            <tr className="hover:bg-[rgba(255,255,255,0.02)] transition-colors">
              <td className="py-3 px-4 text-sm font-bold text-white">Strings</td>
              <td className="py-3 px-4 text-sm font-mono text-[#f1f1f4]">SET, GET, MGET, STRLEN, ECHO, INCR, DECR, INCRBY, DECRBY</td>
            </tr>
            <tr className="hover:bg-[rgba(255,255,255,0.02)] transition-colors">
              <td className="py-3 px-4 text-sm font-bold text-white">Keys</td>
              <td className="py-3 px-4 text-sm font-mono text-[#f1f1f4]">DEL, EXISTS, KEYS, RENAME, RANDOMKEY, DBSIZE, FLUSHDB</td>
            </tr>
            <tr className="hover:bg-[rgba(255,255,255,0.02)] transition-colors">
              <td className="py-3 px-4 text-sm font-bold text-white">Expiration</td>
              <td className="py-3 px-4 text-sm font-mono text-[#f1f1f4]">EXPIRE, TTL, PERSIST</td>
            </tr>
            <tr className="hover:bg-[rgba(255,255,255,0.02)] transition-colors">
              <td className="py-3 px-4 text-sm font-bold text-white">Pub/Sub</td>
              <td className="py-3 px-4 text-sm font-mono text-[#f1f1f4]">SUBSCRIBE, UNSUBSCRIBE, PSUBSCRIBE, PUNSUBSCRIBE, PUBLISH</td>
            </tr>
            <tr className="hover:bg-[rgba(255,255,255,0.02)] transition-colors">
              <td className="py-3 px-4 text-sm font-bold text-white">Persistence</td>
              <td className="py-3 px-4 text-sm font-mono text-[#f1f1f4]">SAVE, LASTSAVE</td>
            </tr>
            <tr className="hover:bg-[rgba(255,255,255,0.02)] transition-colors">
              <td className="py-3 px-4 text-sm font-bold text-white">Connection</td>
              <td className="py-3 px-4 text-sm font-mono text-[#f1f1f4]">PING, QUIT, RESET, COMMAND, COMMAND DOCS, CLIENT, SELECT</td>
            </tr>
          </tbody>
        </table>
      </div>

      <h2 id="handshake" className="text-2xl font-bold text-white mt-12 mb-4">Client Handshakes</h2>
      <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl mb-4">
        <p className="text-[#f1f1f4] mb-3">
          To ensure compatibility with redis-cli and modern Redis SDKs, RoxKV actively implements:
        </p>
        <ul className="list-disc pl-5 space-y-2 text-[#9898ab]">
          <li><strong className="text-white">COMMAND and COMMAND DOCS:</strong> Properly handles capability discovery so clients don't crash on connection.</li>
          <li><strong className="text-white">CLIENT SETINFO/GETNAME:</strong> Safe responses to client metadata tracking.</li>
          <li><strong className="text-white">SELECT:</strong> The <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">SELECT</code> command is mocked to return <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">+OK</code>. RoxKV operates on a single database space.</li>
        </ul>
      </div>

      <h2 id="limitations" className="text-2xl font-bold text-white mt-12 mb-4">Limitations</h2>
      <p className="text-[#f1f1f4] mb-4">
        RoxKV is heavily inspired by Redis but is not a 1:1 clone. It has the following limitations:
      </p>
      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <ul className="list-disc pl-5 space-y-2 text-[#9898ab]">
          <li><strong>No Clustering:</strong> Distributed cluster mechanisms are unsupported.</li>
          <li><strong>No Lua Scripting or Streams:</strong> EVAL and stream-based data patterns are absent.</li>
          <li><strong>No Advanced Data Structures:</strong> Does not currently support Sorted Sets, Hashes, Lists, or Sets natively.</li>
          <li><strong>Inline Expiration:</strong> <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">SET EX</code> and <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">SET PX</code> modifiers are not supported. Use <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">EXPIRE</code> separately instead.</li>
          <li><strong>Simplified Wildcards:</strong> The <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">KEYS</code> command relies on <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">strings.Contains</code> rather than a full regex glob parser.</li>
          <li><strong>Pub/Sub Patterns:</strong> <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">PSUBSCRIBE</code> uses namespace wildcards, not standard Redis glob syntax.</li>
        </ul>
      </div>
    </div>
  );
}
