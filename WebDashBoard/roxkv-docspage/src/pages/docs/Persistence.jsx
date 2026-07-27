export default function Persistence() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Persistence & Expiration</h1>
        <p className="text-xl text-[#9898ab]">
          How RoxKV manages data durability and Time-To-Live (TTL).
        </p>
      </div>

      <div className="section-divider" />

      <h2 id="persistence" className="text-2xl font-bold text-white mt-12 mb-4">Data Persistence</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        Since RoxKV is primarily an in-memory data store, it needs mechanisms to write its state to disk so that data survives server restarts.
      </p>

      <h3 id="save-command" className="text-xl font-semibold text-[#a78bfa] mt-8 mb-2">The SAVE Command</h3>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        Invoking the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1 py-0.5 rounded">SAVE</code> command forces RoxKV to synchronously dump the entire memory state to disk.
        Since the operation runs synchronously on the main executor thread, it will block all other clients while the dump is occurring. 
      </p>

      <h3 id="lastsave-command" className="text-xl font-semibold text-[#a78bfa] mt-8 mb-2">The LASTSAVE Command</h3>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        Returns a UNIX timestamp representing the time of the last successful data dump to disk. You can use this to monitor whether automated background saves (if configured via RoxKV internals) are working as expected.
      </p>

      <h2 id="expiration" className="text-2xl font-bold text-white mt-12 mb-4">Time-To-Live (TTL)</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        RoxKV supports setting expiration times on individual keys.
      </p>

      <ul className="list-disc list-inside space-y-2 text-[#9898ab] ml-4 mb-6">
        <li>
          <strong className="text-white">EXPIRE:</strong> Sets a timeout on a key in seconds. Once the timeout expires, the key is automatically deleted.
        </li>
        <li>
          <strong className="text-white">TTL:</strong> Returns the remaining time to live for a key. RoxKV returns `-2` if the key exists but has no expiration set.
        </li>
        <li>
          <strong className="text-white">PERSIST:</strong> Removes the existing timeout on a key, turning it from a volatile key to a persistent key.
        </li>
      </ul>

      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <h4 className="text-white font-semibold mb-2">Note on Standard SET</h4>
        <p className="text-[#9898ab] text-sm">
          Unlike standard Redis, RoxKV's <code className="text-white">SET</code> command currently does not support <code className="text-white">EX</code> or <code className="text-white">PX</code> arguments directly. You must explicitly call <code className="text-white">EXPIRE</code> after setting the key if you wish to assign a TTL.
        </p>
      </div>
    </div>
  );
}
