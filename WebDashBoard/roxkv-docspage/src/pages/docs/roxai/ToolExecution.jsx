import CodeBlock from '../../../components/docs/CodeBlock';

export default function ToolExecution() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Tool Execution</h1>
        <p className="text-xl text-[#9898ab]">How LLM intent becomes direct database action</p>
      </div>
      <div className="section-divider" />

      <h2 id="dispatch-mechanism" className="text-2xl font-semibold text-white mb-4">Dispatch Mechanism</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        Once the LLM decides which tool to call in Phase 1, the Master Agent receives a tool call request with specific arguments. The Master Agent implements a massive switch statement that maps the tool name (e.g., <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">get_keys</code>) directly to the corresponding Go function.
      </p>

      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#9898ab]"><strong>Crucial Architecture Detail:</strong> Tool executions are NOT network calls. They are direct, zero-overhead Go function calls within the same process. For example, <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">store.GetKv()</code> or <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">metrics.GetCpuUsage()</code>.</p>
      </div>

      <h2 id="result-handling" className="text-2xl font-semibold text-white mb-4">Result Handling & Injection</h2>
      <div className="space-y-4 mb-8">
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">1. Formatting</h3>
          <p className="text-[#5e5e73]">
            The internal Go functions return structured data (structs, maps, slices). The execution layer serializes this output into JSON or formatted text strings.
          </p>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">2. Injection</h3>
          <p className="text-[#5e5e73]">
            The formatted result string is packaged into a tool response message and appended directly to the LLM's conversation context.
          </p>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">3. Synthesis</h3>
          <p className="text-[#5e5e73]">
            The secondary LLM call (Phase 2) reads this raw JSON/text output in the context and synthesizes it into a human-readable response for the user.
          </p>
        </div>
      </div>

      <h2 id="error-handling" className="text-2xl font-semibold text-white mb-4">Error Handling</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        If a tool execution fails (e.g., missing key, invalid parameters), the error is caught by the execution layer. Instead of crashing, the error message itself is formatted as the tool output and sent back to the LLM. This allows the LLM to understand the failure and potentially try an alternative approach or explain the error to the user gracefully.
      </p>
    </div>
  );
}
