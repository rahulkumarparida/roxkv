export default function PubSub() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Pub/Sub</h1>
        <p className="text-xl text-[#9898ab]">
          Real-time message routing with exact and pattern matching.
        </p>
      </div>

      <div className="section-divider" />

      <h2 id="overview" className="text-2xl font-bold text-white mt-12 mb-4">Overview</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        RoxKV implements the Publish/Subscribe messaging paradigm. Senders (publishers) do not program messages to be sent directly to specific receivers (subscribers). Instead, messages are published to channels, without knowledge of what (if any) subscribers there may be.
      </p>

      <h2 id="subscribe" className="text-2xl font-bold text-white mt-12 mb-4">SUBSCRIBE & UNSUBSCRIBE</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        When a client issues a <code className="text-[#a78bfa] bg-[#1e1e2e] px-1 py-0.5 rounded">SUBSCRIBE channel</code> command, the server places that client into <strong>Subscriber Mode</strong>. The client remains blocked, waiting for messages to arrive on the channel.
      </p>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        To exit this mode, the client must either issue <code className="text-[#a78bfa] bg-[#1e1e2e] px-1 py-0.5 rounded">UNSUBSCRIBE</code> for all its channels, send a <code className="text-[#a78bfa] bg-[#1e1e2e] px-1 py-0.5 rounded">RESET</code> command, or close the connection via <code className="text-[#a78bfa] bg-[#1e1e2e] px-1 py-0.5 rounded">QUIT</code>.
      </p>

      <h2 id="publish" className="text-2xl font-bold text-white mt-12 mb-4">PUBLISH</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        Publishing a message pushes the payload to all currently connected exact-match subscribers, as well as all pattern-match subscribers whose pattern covers the channel name.
        RoxKV maintains channel history internally, tracking the last publisher, total messages, and byte size of the topic.
      </p>

      <h2 id="pattern-subscriptions" className="text-2xl font-bold text-white mt-12 mb-4">Pattern Subscriptions (PSUBSCRIBE)</h2>
      <p className="text-[#f1f1f4] leading-relaxed mb-4">
        RoxKV supports pattern matching via the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1 py-0.5 rounded">PSUBSCRIBE</code> command. However, the pattern matching rules are specific to RoxKV's implementation.
      </p>

      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <h4 className="text-white font-semibold mb-2">Pattern Matching Constraints</h4>
        <p className="text-[#9898ab] text-sm">
          RoxKV's pattern matcher currently evaluates wildcards (<code className="text-white">*</code>) based on namespaces separated by dots (<code className="text-white">.</code>) or colons (<code className="text-white">:</code>). It does not implement full glob-style character matching.
          <br /><br />
          For example, <code className="text-white">news.*.notification</code> will match <code className="text-white">news.tech</code> or <code className="text-white">news.sports</code>.
        </p>
      </div>

    </div>
  );
}
