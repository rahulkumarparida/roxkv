import React from 'react';
import CodeBlock from '../../../components/docs/CodeBlock';

export default function CoreDatabase() {
  const structCode = `type MemoryAlloc struct {
    Data map[string]Item
    mu   sync.RWMutex
}

type Item struct {
    Key  string
    Val  []byte
    Meta Metadata
}

type Metadata struct {
    TTL            time.Time
    UpdatedAt      time.Time
    CreatedAt      time.Time
    LastAccessedBy string
    KeyAccessCount int
    Size           int
    Namespace      string
}`;

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Core Database Engine</h1>
        <p className="text-xl text-[#9898ab]">In-memory storage structures and concurrency.</p>
      </div>
      
      <div className="section-divider" />

      <h2 id="storage-structures" className="text-2xl font-semibold text-white mb-4">Storage Structures</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        At the heart of RoxKV is the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">MemoryAlloc</code> struct, which wraps a Go map and a read-write mutex for thread-safe access.
      </p>
      <CodeBlock code={structCode} language="go" title="Data Structures" />

      <h2 id="binary-safety" className="text-2xl font-semibold text-white mb-4 mt-8">Binary Safety</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        Values are stored as raw <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">[]byte</code> arrays. This means RoxKV is binary-safe and can store arbitrary data, including images, protocol buffers, or compressed payloads, without assuming any string encoding like UTF-8.
      </p>

      <h2 id="concurrency-model" className="text-2xl font-semibold text-white mb-4 mt-8">Concurrency Model</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        Concurrency is managed using <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">sync.RWMutex</code>.
      </p>
      <ul className="list-disc list-inside text-[#9898ab] space-y-2 mb-4">
        <li><strong>Reads:</strong> Acquire <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">RLock()</code> allowing multiple concurrent readers.</li>
        <li><strong>Writes:</strong> Acquire <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">Lock()</code> ensuring exclusive access during modifications.</li>
      </ul>

      <h2 id="namespaces" className="text-2xl font-semibold text-white mb-4 mt-8">Namespace Derivation</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        The <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">Namespace</code> field in the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">Metadata</code> struct is derived automatically from the key name. 
      </p>
      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#f1f1f4]">Keys containing colons are split. For example, a key named <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">users:profile:123</code> is assigned the Namespace <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">users</code>.</p>
      </div>

    </div>
  );
}
