export default function Sessions() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Sessions & Context</h1>
        <p className="text-xl text-[#9898ab]">Managing conversational state in RoxAI</p>
      </div>
      <div className="section-divider" />

      <h2 id="agent-session" className="text-2xl font-semibold text-white mb-4">The AgentSession Structure</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        Interactions with RoxAI are inherently stateful. The state is maintained within an <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">AgentSession</code> struct. This session object is responsible for holding the complete conversation context, including user prompts, LLM responses, and the raw tool execution outputs.
      </p>

      <h2 id="session-registry" className="text-2xl font-semibold text-white mb-4">Session Registry</h2>
      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#9898ab]">
          Sessions are stored globally in memory using a session registry. Every connected client or user is assigned a unique session key, formatted as <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">masteragent:&lt;user_id&gt;</code>. 
        </p>
      </div>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        When a request arrives, the Master Agent queries the registry. If a session exists for the user ID, it is retrieved. If not, a new session lifecycle begins. This design allows context to persist seamlessly across multiple interactions over time.
      </p>

      <h2 id="context-management" className="text-2xl font-semibold text-white mb-4">Context & Message Threading</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        The session context is critical for the LLM to understand the history of the conversation and the actions it has taken. Messages are threaded and appended in sequence:
      </p>
      <ol className="list-decimal list-inside text-[#5e5e73] space-y-2 mb-8 ml-4">
        <li><strong>User Message:</strong> Appended when a query is received.</li>
        <li><strong>Assistant Message (Tool Call):</strong> Appended when the Intent Phase LLM decides to call a tool.</li>
        <li><strong>Tool Response Message:</strong> Appended containing the raw JSON/text results after internal execution.</li>
        <li><strong>Assistant Message (Synthesis):</strong> Appended containing the final natural language answer provided to the user.</li>
      </ol>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        This continuous context is provided to the LLM during both the Intent and Synthesis phases, ensuring it always has the full picture of the ongoing interaction.
      </p>
    </div>
  );
}
