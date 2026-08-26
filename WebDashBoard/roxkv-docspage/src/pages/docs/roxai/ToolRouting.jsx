import InfoTable from '../../../components/docs/InfoTable';

export default function ToolRouting() {
  const scoringRules = [
    { criteria: 'Exact keyword match', points: '+3 × weight', description: 'Highest confidence match' },
    { criteria: 'Synonym match', points: '+2 × weight', description: 'Matches related terms' },
    { criteria: 'Example phrase overlap', points: '+2 × weight per word', description: 'Matches words in example phrases' },
    { criteria: 'Substring/partial match', points: '+1', description: 'Partial token matching' },
    { criteria: 'Category match', points: '+1', description: 'Matches tool category' },
  ];

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Semantic Tool Routing</h1>
        <p className="text-xl text-[#9898ab]">Optimizing LLM context windows through intelligent tool selection</p>
      </div>
      <div className="section-divider" />

      <h2 id="why-it-exists" className="text-2xl font-semibold text-white mb-4">Why Semantic Routing?</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        RoxAI has over 40 registered tools. Sending the full JSON schema of all 42 tools to the LLM on every request consumes a massive amount of context window space. This is especially problematic for local models like Ollama, leading to slower inference times and higher memory usage. Semantic routing solves this by only providing the most relevant tools for a given query.
      </p>
      
      <div className="bg-[rgba(124,58,237,0.05)] border-l-4 border-[#7c3aed] p-4 my-6">
        <p className="text-[#9898ab]"><strong>Concept:</strong> The <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">RouteTools</code> function in <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">agents/router/router.go</code> tokenizes the query, removes stop words, and scores all available tools.</p>
      </div>

      <h2 id="scoring-algorithm" className="text-2xl font-semibold text-white mb-4">Scoring Algorithm</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        Each tool starts with a score of 0. The router analyzes the tokens and applies the following scoring rules:
      </p>
      <div className="mb-8">
        <InfoTable 
          headers={['Match Criteria', 'Score Impact', 'Description']} 
          rows={scoringRules.map(r => [r.criteria, r.points, r.description])} 
        />
      </div>

      <h2 id="filtering-and-capping" className="text-2xl font-semibold text-white mb-4">Filtering and Capping</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        After scoring, the tools are sorted in descending order based on their total score. The router then caps the result set, returning only the top <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">MaxTools</code> (default is 5). This guarantees that the LLM receives a highly focused, optimized subset of tools.
      </p>

      <h2 id="example" className="text-2xl font-semibold text-white mb-4">Example Scenario</h2>
      <div className="bg-[#16161f] border border-[#1e1e2e] p-6 rounded-xl">
        <p className="text-[#9898ab] mb-2"><strong>User Query:</strong> "How much RAM is being used?"</p>
        <p className="text-[#9898ab] mb-2"><strong>Tokens:</strong> <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">["ram", "used"]</code></p>
        <p className="text-[#9898ab] mb-2"><strong>Resulting Top Tools:</strong></p>
        <ul className="list-disc list-inside text-[#5e5e73] space-y-1 ml-4">
          <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">get_ram_usage</code> (Highest score due to exact match on "ram")</li>
          <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">get_computer_usage</code></li>
          <li><code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">get_runtime_stats</code></li>
        </ul>
      </div>
    </div>
  );
}
