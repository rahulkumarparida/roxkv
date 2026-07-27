export default function Architecture() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Architecture</h1>
        <p className="text-xl text-[#9898ab]">
          Understanding how RoxKV handles TCP connections, parses RESP, and manages state.
        </p>
      </div>

      <div className="section-divider" />

      <h2 id="core-components" className="text-2xl font-bold text-white mt-12 mb-4">Core Components</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        RoxKV's architecture is broken down into modular packages to isolate concerns.
      </p>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 my-8">
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl hover:border-[rgba(124,58,237,0.3)] transition-colors">
          <h3 className="text-lg font-semibold text-white mb-2 text-[#a78bfa]">TCP Server (`server/`)</h3>
          <p className="text-[#9898ab] text-sm">
            Listens for incoming TCP connections. Spawns a lightweight goroutine for every active client connection to read from the socket.
          </p>
        </div>
        
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl hover:border-[rgba(124,58,237,0.3)] transition-colors">
          <h3 className="text-lg font-semibold text-white mb-2 text-[#a78bfa]">Protocol Parser (`resp-parser/`)</h3>
          <p className="text-[#9898ab] text-sm">
            Decodes incoming bytes into RESP data types (Arrays, Bulk Strings, Integers) and encodes server responses back into the wire format.
          </p>
        </div>

        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl hover:border-[rgba(124,58,237,0.3)] transition-colors">
          <h3 className="text-lg font-semibold text-white mb-2 text-[#a78bfa]">Key-Value Store (`store/`)</h3>
          <p className="text-[#9898ab] text-sm">
            The core thread-safe dictionary that stores all keys, values, and metadata (like TTL). Values are internally managed.
          </p>
        </div>

        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl hover:border-[rgba(124,58,237,0.3)] transition-colors">
          <h3 className="text-lg font-semibold text-white mb-2 text-[#a78bfa]">Command Executor (`commands/`)</h3>
          <p className="text-[#9898ab] text-sm">
            Evaluates the parsed AST (Abstract Syntax Tree) and executes the requested database operations against the store.
          </p>
        </div>
      </div>

      <h2 id="client-state" className="text-2xl font-bold text-white mt-12 mb-4">Client State & Modes</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        RoxKV tracks the state of every connected client. When a client executes a command like `SUBSCRIBE` or `PSUBSCRIBE`, RoxKV flips the client into <strong>Subscriber Mode</strong>.
      </p>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        While in Subscriber Mode, standard database operations (like `GET` or `SET`) are strictly blocked. The client may only execute `SUBSCRIBE`, `UNSUBSCRIBE`, `PSUBSCRIBE`, `PUNSUBSCRIBE`, `PING`, `QUIT`, and `RESET`. Issuing an invalid command results in a specific protocol error.
      </p>

    </div>
  );
}
