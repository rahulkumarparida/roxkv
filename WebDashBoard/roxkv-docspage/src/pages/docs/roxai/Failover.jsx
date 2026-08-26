import FlowDiagram from '../../../components/docs/FlowDiagram';

export default function Failover() {
  const failoverSteps = [
    { label: 'Primary Provider', highlight: true, sub: 'Attempt 1' },
    { label: 'Retry 1', highlight: false, sub: 'Wait 1s' },
    { label: 'Retry 2', highlight: false, sub: 'Wait 2s' },
    { label: 'Retry 3', highlight: false, sub: 'Wait 2s' },
    { label: 'Fallback Provider 1', highlight: true, sub: 'Attempt 1' },
    { label: 'Retry 1', highlight: false, sub: 'Wait 1s' },
    { label: 'Fallback Provider 2', highlight: true, sub: 'Attempt 1' },
  ];

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Provider Failover</h1>
        <p className="text-xl text-[#9898ab]">Ensuring high availability across LLM APIs</p>
      </div>
      <div className="section-divider" />

      <h2 id="overview" className="text-2xl font-semibold text-white mb-4">Automatic Failover Mechanism</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        Network issues, rate limits, and provider outages happen. RoxAI includes an intelligent <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">ExecuteWithFailover</code> mechanism built into the abstractor layer to handle these scenarios gracefully without dropping user requests.
      </p>

      <h2 id="error-classification" className="text-2xl font-semibold text-white mb-4">Error Classification</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">Retryable Errors</h3>
          <p className="text-[#5e5e73] mb-2">Errors identified by <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">IsRetryableError()</code>:</p>
          <ul className="list-disc list-inside text-[#5e5e73]">
            <li>HTTP 502 / 503 / 504</li>
            <li>Connection resets</li>
            <li>Timeouts</li>
          </ul>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">Non-Retryable Errors</h3>
          <p className="text-[#5e5e73] mb-2">Errors that trigger an immediate switch to a fallback provider:</p>
          <ul className="list-disc list-inside text-[#5e5e73]">
            <li>401 Unauthorized (Bad API Key)</li>
            <li>400 Bad Request</li>
          </ul>
        </div>
      </div>

      <h2 id="retry-policy" className="text-2xl font-semibold text-white mb-4">Retry Policy</h2>
      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#9898ab]">For retryable errors, the active provider is given up to 3 retries (4 total attempts). The backoff delay increases incrementally (e.g., 1s, 2s) to avoid overwhelming the endpoint.</p>
      </div>

      <h2 id="fallback-sequence" className="text-2xl font-semibold text-white mb-4">Fallback Sequence</h2>
      <p className="text-[#9898ab] leading-relaxed mb-6">
        If the primary provider exhausts its retry attempts or encounters a non-retryable error, the failover system retrieves a list of other authenticated providers via <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">GetConfiguredProviders()</code>. It then iterates through these fallback providers, applying the same retry policy, until a successful response is received or all options are exhausted.
      </p>

      <div className="mb-8 p-6 bg-[#0a0a0f] rounded-xl border border-[#1e1e2e]">
        <FlowDiagram steps={failoverSteps} direction="horizontal" />
      </div>
    </div>
  );
}
