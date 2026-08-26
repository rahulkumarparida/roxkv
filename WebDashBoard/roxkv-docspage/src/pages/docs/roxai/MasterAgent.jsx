import CodeBlock from '../../../components/docs/CodeBlock';

export default function MasterAgent() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Master Agent</h1>
        <p className="text-xl text-[#9898ab]">The central orchestration component of RoxAI</p>
      </div>
      <div className="section-divider" />

      <h2 id="overview" className="text-2xl font-semibold text-white mb-4">Overview</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        The Master Agent, defined in <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">agents/Master/Masteragent.go</code>, is the core orchestrator of the RoxAI system. It handles session management, manages the two-phase LLM execution cycle, and dispatches tool calls directly to internal Go functions.
      </p>

      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#9898ab]"><strong>Key Responsibility:</strong> The Master Agent acts as the bridge between the natural language reasoning of the LLM and the low-level, high-performance operations of the RoxKV database.</p>
      </div>

      <h2 id="session-management" className="text-2xl font-semibold text-white mb-4">Session Management</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        Every interaction is stateful. The Master Agent retrieves an <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">AgentSession</code> from the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">sessionRegistry</code> using a key formatted as <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">masteragent:&lt;user_id&gt;</code>. This ensures continuity and context awareness across multiple queries from the same user.
      </p>

      <h2 id="two-phase-execution" className="text-2xl font-semibold text-white mb-4">Two-Phase Execution</h2>
      <div className="space-y-6 mb-8">
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">Phase 1: Intent</h3>
          <p className="text-[#5e5e73]">
            The agent passes the user query along with the top 5 tools (filtered by the semantic router) to the LLM. The LLM determines the required actions and requests specific tool calls.
          </p>
        </div>
        <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
          <h3 className="text-lg font-medium text-white mb-2">Phase 2: Synthesis</h3>
          <p className="text-[#5e5e73]">
            After executing the requested tools, the results are appended to the context. A secondary LLM call is made to synthesize these results into a clear, natural language response.
          </p>
        </div>
      </div>

      <h2 id="tool-dispatch" className="text-2xl font-semibold text-white mb-4">Tool Dispatch Mechanism</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        Tool execution is highly efficient. Instead of making network requests or shelling out to external scripts, the Master Agent uses a large switch statement mapping tool names directly to Go functions within the RoxKV binary.
      </p>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        For example, a tool call to <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">get_computer_usage</code> routes directly to <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">metrics.GetComputerUsage()</code>.
      </p>

      <h2 id="subscriber-awareness" className="text-2xl font-semibold text-white mb-4">Subscriber Mode Awareness</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        The Master Agent is fully aware of client state. It understands if a connected client is in "subscriber mode" (listening to PubSub topics) and can adapt its behavior and tool access accordingly.
      </p>
    </div>
  );
}
