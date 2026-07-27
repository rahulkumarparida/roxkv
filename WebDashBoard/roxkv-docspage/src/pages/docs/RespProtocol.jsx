export default function RespProtocol() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">RESP Protocol</h1>
        <p className="text-xl text-[#9898ab]">
          How RoxKV communicates with clients over TCP.
        </p>
      </div>

      <div className="section-divider" />

      <h2 id="overview" className="text-2xl font-bold text-white mt-12 mb-4">Overview</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        RoxKV implements the REdis Serialization Protocol (RESP). It is a human-readable, binary-safe text protocol. All communication between the client (like `redis-cli`) and RoxKV happens via RESP frames over a TCP socket.
      </p>

      <div className="bg-[#16161f] border border-[#1e1e2e] rounded-xl p-8 my-8 overflow scrollbar-thin">
        <div className="flex items-center min-w-[700px] justify-between text-center gap-2">
          
          <div className="flex flex-col items-center">
            <div className="bg-[#0a0a0f] border border-purple-500/30 text-white rounded-lg px-4 py-3 shadow-lg w-32 shadow-purple-900/20">
              <span className="font-semibold block mb-1">Client</span>
              <span className="text-xs text-[#9898ab] font-mono">redis-cli</span>
            </div>
          </div>

          <div className="flex flex-col items-center justify-center flex-1 min-w-[40px] relative">
            <div className="h-px bg-gradient-to-r from-purple-500/20 via-purple-500/50 to-purple-500/20 w-full mb-1"></div>
            <span className="text-[10px] uppercase tracking-wider text-purple-400/80 font-mono absolute -top-4">RESP Array</span>
          </div>

          <div className="flex flex-col items-center">
            <div className="bg-[#1e1e2e] border border-[#2a2a3e] text-white rounded-lg px-4 py-3 shadow-lg w-32">
              <span className="font-semibold block mb-1">Decoder</span>
              <span className="text-xs text-[#9898ab] font-mono">Parser</span>
            </div>
          </div>

          <div className="flex flex-col items-center justify-center flex-1 min-w-[40px]">
            <div className="h-px bg-[#2a2a3e] w-full"></div>
          </div>

          <div className="flex flex-col items-center">
            <div className="bg-[#1e1e2e] border border-[#2a2a3e] text-white rounded-lg px-4 py-3 shadow-lg w-32">
              <span className="font-semibold block mb-1">Dispatcher</span>
              <span className="text-xs text-[#9898ab] font-mono">Handler</span>
            </div>
          </div>

          <div className="flex flex-col items-center justify-center flex-1 min-w-[40px]">
            <div className="h-px bg-[#2a2a3e] w-full"></div>
          </div>

          <div className="flex flex-col items-center">
            <div className="bg-[#1e1e2e] border border-[#2a2a3e] text-white rounded-lg px-4 py-3 shadow-lg w-32">
              <span className="font-semibold block mb-1">Encoder</span>
              <span className="text-xs text-[#9898ab] font-mono">Response</span>
            </div>
          </div>

          <div className="flex flex-col items-center justify-center flex-1 min-w-[40px] relative">
            <div className="h-px bg-gradient-to-l from-emerald-500/20 via-emerald-500/50 to-emerald-500/20 w-full mb-1"></div>
            <span className="text-[10px] uppercase tracking-wider text-emerald-400/80 font-mono absolute -top-4">RESP String</span>
          </div>

          <div className="flex flex-col items-center">
            <div className="bg-[#0a0a0f] border border-emerald-500/30 text-white rounded-lg px-4 py-3 shadow-lg w-32 shadow-emerald-900/20">
              <span className="font-semibold block mb-1">Client</span>
              <span className="text-xs text-[#9898ab] font-mono">TCP Read</span>
            </div>
          </div>
        </div>
      </div>

      <h2 id="supported-types" className="text-2xl font-bold text-white mt-12 mb-4">Supported RESP Types</h2>
      
      <div className="space-y-6">
        <div>
          <h3 id="decoding" className="text-xl font-semibold text-[#a78bfa] mb-2">Decoding (Client to Server)</h3>
          <p className="text-[#f1f1f4] mb-3">RoxKV can parse the following data types sent by clients:</p>
          <ul className="list-disc list-inside space-y-1 text-[#9898ab] ml-4">
            <li><code className="text-[#a78bfa]">+</code> Simple Strings</li>
            <li><code className="text-[#a78bfa]">-</code> Simple Errors</li>
            <li><code className="text-[#a78bfa]">:</code> Integers</li>
            <li><code className="text-[#a78bfa]">$</code> Bulk Strings</li>
            <li><code className="text-[#a78bfa]">*</code> Arrays</li>
            <li><code className="text-[#a78bfa]">#</code> Booleans</li>
            <li><code className="text-[#a78bfa]">,</code> Doubles</li>
          </ul>
        </div>

        <div>
          <h3 id="encoding" className="text-xl font-semibold text-[#a78bfa] mb-2">Encoding (Server to Client)</h3>
          <p className="text-[#f1f1f4] mb-3">RoxKV responds to clients using these data types:</p>
          <ul className="list-disc list-inside space-y-1 text-[#9898ab] ml-4">
            <li><code className="text-[#a78bfa]">+</code> Simple Strings (e.g. `+OK\r\n`)</li>
            <li><code className="text-[#a78bfa]">-</code> Simple Errors (e.g. `-ERR invalid command\r\n`)</li>
            <li><code className="text-[#a78bfa]">:</code> Integers</li>
            <li><code className="text-[#a78bfa]">$</code> Bulk Strings</li>
            <li><code className="text-[#a78bfa]">$-1</code> Null Values (Returned for non-existent keys)</li>
            <li><code className="text-[#a78bfa]">*</code> Arrays (Returned for commands like MGET or KEYS)</li>
            <li><code className="text-[#a78bfa]">#</code> Booleans</li>
          </ul>
        </div>
      </div>

    </div>
  );
}
