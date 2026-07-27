export default function Overview() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Overview</h1>
        <p className="text-xl text-[#9898ab]">
          Technical architecture and operational flow of RoxKV.
        </p>
      </div>

      <div className="section-divider" />

      <h2 id="what-it-does" className="text-2xl font-bold text-white mt-12 mb-4">What RoxKV Does</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        RoxKV is a custom, in-memory key-value database written in Go. It implements a subset of the Redis Serialization Protocol (RESP), allowing standard Redis clients to interact with it seamlessly. At its core, RoxKV provides fast thread-safe storage, basic mathematical operations, key expiration (TTL), point-in-time persistence, and a Pub/Sub messaging system.
      </p>

      <h2 id="client-interaction" className="text-2xl font-bold text-white mt-12 mb-4">Client Interaction & Connection Handling</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        Clients communicate with RoxKV over standard TCP sockets. When a client connects (typically on port <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">6379</code>), the server spawns a dedicated goroutine. This connection is encapsulated in a <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">utils.NewClient</code> struct which maintains the client's state, such as its active Pub/Sub topic subscriptions and its current operating mode.
      </p>
      
      <h2 id="resp-data-structure" className="text-2xl font-bold text-white mt-12 mb-4">RESP Data Structure</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        All commands and responses are transmitted as RESP frames. According to the implementation in <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">decoder.go</code> and <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">encoder.go</code>, RoxKV handles:
      </p>
      <ul className="list-disc list-inside space-y-2 text-[#9898ab] ml-4 mb-6">
        <li><strong className="text-white">Simple Strings (+):</strong> Used for simple status responses like <code className="text-white">+OK\r\n</code>.</li>
        <li><strong className="text-white">Errors (-):</strong> Used to return operational errors to the client.</li>
        <li><strong className="text-white">Integers (:):</strong> Used for counts, lengths, and numerical command results.</li>
        <li><strong className="text-white">Bulk Strings ($):</strong> Binary-safe strings up to 512MB, used for actual key/value payloads.</li>
        <li><strong className="text-white">Arrays (*):</strong> Used by clients to send commands (e.g., <code className="text-white">["SET", "key", "value"]</code>).</li>
        <li><strong className="text-white">Nulls ($-1):</strong> Indicates a non-existent value.</li>
      </ul>

      <h2 id="command-lifecycle" className="text-2xl font-bold text-white mt-12 mb-4">Command Lifecycle</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        The journey of a command through the system is strictly defined:
      </p>
      <ol className="list-decimal list-inside space-y-3 text-[#f1f1f4] ml-4 mb-6">
        <li><strong className="text-[#a78bfa]">Read:</strong> The server reads the raw byte stream from the TCP socket.</li>
        <li><strong className="text-[#a78bfa]">Decode:</strong> The bytes are parsed into a string array using <code className="text-white bg-[#1e1e2e] px-1.5 py-0.5 rounded">DecodeArrayString</code>.</li>
        <li><strong className="text-[#a78bfa]">Structure:</strong> The array is mapped to a <code className="text-white bg-[#1e1e2e] px-1.5 py-0.5 rounded">RedisInput</code> struct containing the `Cmd` and `Args`.</li>
        <li><strong className="text-[#a78bfa]">Mode Validation:</strong> <code className="text-white bg-[#1e1e2e] px-1.5 py-0.5 rounded">parseCommand</code> checks the client's mode. If the client is in Subscriber mode, it rejects standard storage commands.</li>
        <li><strong className="text-[#a78bfa]">Dispatch & Execution:</strong> A `switch` statement routes the command to the appropriate handler in <code className="text-white bg-[#1e1e2e] px-1.5 py-0.5 rounded">executor.go</code>, which safely acquires a mutex on the <code className="text-white bg-[#1e1e2e] px-1.5 py-0.5 rounded">store.StoreHelper</code> to perform the operation.</li>
        <li><strong className="text-[#a78bfa]">Encode & Write:</strong> The result is RESP-encoded and written back to the client's TCP connection.</li>
      </ol>

      <h2 id="implementation-decisions" className="text-2xl font-bold text-white mt-12 mb-4">Important Implementation Decisions</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        To keep the system lightweight and maintainable, several explicit technical decisions were made:
      </p>
      <ul className="list-disc list-inside space-y-2 text-[#9898ab] ml-4 mb-6">
        <li><strong className="text-white">Values are String Lists:</strong> Internally, RoxKV treats all stored values as slices of strings. The <code className="text-white">GET</code> command simply joins them with a space.</li>
        <li><strong className="text-white">Simplified Pattern Matching:</strong> The <code className="text-white">PSUBSCRIBE</code> command uses a custom pattern matcher that evaluates namespaces delimited by dots or colons, rather than standard regex or globbing.</li>
        <li><strong className="text-white">Synchronous Persistence:</strong> The <code className="text-white">SAVE</code> command executes a blocking memory dump. It is designed for manual or low-frequency automated snapshots.</li>
      </ul>

      <h2 id="limitations" className="text-2xl font-bold text-white mt-12 mb-4">Limitations</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        As a subset implementation of Redis, certain behaviors are restricted:
        Standard <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">SET</code> options like <code className="text-white">EX</code> are ignored; users must explicitly call <code className="text-white">EXPIRE</code>.
        Advanced data structures like Hashes, Sets, and Sorted Sets are currently not implemented.
      </p>
    </div>
  );
}
