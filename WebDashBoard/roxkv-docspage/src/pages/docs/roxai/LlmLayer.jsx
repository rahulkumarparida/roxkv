import CodeBlock from '../../../components/docs/CodeBlock';

export default function LlmLayer() {
  const providerInterface = `type Provider interface {
    Chat(ctx context.Context, messages []GenericMessage, tools []GenericToolDefinition, config ProviderConfig) (*ProviderResponse, error)
    Name() string
}`;

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">LLM Abstraction Layer</h1>
        <p className="text-xl text-[#9898ab]">SDK-agnostic interfaces for multi-provider support</p>
      </div>
      <div className="section-divider" />

      <h2 id="why-an-abstraction-layer" className="text-2xl font-semibold text-white mb-4">Why an Abstraction Layer?</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        Located in <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">agents/abstractor/</code>, the LLM Abstraction Layer exists to decouple the core logic of RoxAI from the specific SDKs of various LLM providers (OpenAI, Anthropic, Ollama, etc.). This decoupling allows providers to be hot-swapped dynamically without altering the Master Agent's behavior or logic.
      </p>

      <h2 id="the-provider-interface" className="text-2xl font-semibold text-white mb-4">The Provider Interface</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        All integrated LLMs must implement a unified interface. This ensures standard behavior across the platform.
      </p>
      <div className="mb-8">
        <CodeBlock code={providerInterface} language="go" title="agents/abstractor/provider.go" />
      </div>

      <h2 id="generic-types" className="text-2xl font-semibold text-white mb-4">Generic Data Types</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        To prevent provider-specific structs from leaking into the core application, the abstractor defines generic types in <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">types.go</code>:
      </p>
      <ul className="list-disc list-inside text-[#9898ab] space-y-2 mb-6">
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">GenericMessage</code>: Represents a standard chat message (user, assistant, system, tool).</li>
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">GenericToolCall</code>: A standardized representation of a requested tool execution.</li>
        <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">GenericToolDefinition</code>: The schema format expected by the routing and LLM layers.</li>
      </ul>

      <h2 id="response-normalization" className="text-2xl font-semibold text-white mb-4">Response Normalization</h2>
      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#9898ab]">
          Every provider SDK returns data differently. OpenAI uses <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">ChatCompletion</code>, Anthropic uses <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">MessageResponse</code>. The abstraction layer's primary job is mapping these varying responses into a single, predictable <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">ProviderResponse</code> object that the Master Agent can reliably process.
        </p>
      </div>
    </div>
  );
}
