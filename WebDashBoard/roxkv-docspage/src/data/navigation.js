import {
  BookOpen, Layers, Database, Server, FileCode, Terminal as TerminalIcon,
  Radio, Save, Clock, Settings, Download, GitCompare, Code,
  Brain, Wrench, GitBranch, Cpu, Globe, Shield, MessageSquare,
  LayoutDashboard, Link2, Zap
} from 'lucide-react';

export const roxkvNav = [
  {
    title: 'Getting Started',
    links: [
      { name: 'Overview', href: '/docs/roxkv/overview', icon: BookOpen, description: 'Introduction to RoxKV database' },
      { name: 'Installation', href: '/docs/roxkv/installation', icon: Download, description: 'Docker, manual build, and binary setup' },
      { name: 'Architecture', href: '/docs/roxkv/architecture', icon: Layers, description: 'System architecture and data flow' },
    ],
  },
  {
    title: 'Core',
    links: [
      { name: 'Core Database', href: '/docs/roxkv/core', icon: Database, description: 'In-memory store, concurrency, and binary safety' },
      { name: 'TCP Server', href: '/docs/roxkv/tcp-server', icon: Server, description: 'Connection handling and client management' },
      { name: 'RESP Protocol', href: '/docs/roxkv/resp', icon: FileCode, description: 'Redis Serialization Protocol implementation' },
    ],
  },
  {
    title: 'Features',
    links: [
      { name: 'Commands', href: '/docs/roxkv/commands', icon: TerminalIcon, description: 'Supported Redis commands reference' },
      { name: 'Pub/Sub', href: '/docs/roxkv/pubsub', icon: Radio, description: 'Publish/Subscribe and pattern matching' },
      { name: 'Persistence', href: '/docs/roxkv/persistence', icon: Save, description: 'JSON snapshot persistence' },
      { name: 'TTL & Expiration', href: '/docs/roxkv/ttl', icon: Clock, description: 'Key expiration and background cleanup' },
    ],
  },
  {
    title: 'Reference',
    links: [
      { name: 'Configuration', href: '/docs/roxkv/configuration', icon: Settings, description: 'Ports, environment variables, and Docker' },
      { name: 'Redis Compatibility', href: '/docs/roxkv/redis-compat', icon: GitCompare, description: 'Supported features and limitations' },
      { name: 'Usage Examples', href: '/docs/roxkv/usage', icon: Code, description: 'Client integration and code examples' },
    ],
  },
];

export const roxaiNav = [
  {
    title: 'Getting Started',
    links: [
      { name: 'Overview', href: '/docs/roxai/overview', icon: BookOpen, description: 'Introduction to RoxAI agent system' },
      { name: 'Architecture', href: '/docs/roxai/architecture', icon: Layers, description: 'Request flow and system design' },
    ],
  },
  {
    title: 'Agent System',
    links: [
      { name: 'Master Agent', href: '/docs/roxai/master-agent', icon: Brain, description: 'Central orchestration layer' },
      { name: 'Tool Registry', href: '/docs/roxai/tool-registry', icon: Wrench, description: '40+ tools across 5 categories' },
      { name: 'Tool Routing', href: '/docs/roxai/tool-routing', icon: GitBranch, description: 'Semantic scoring and selection' },
      { name: 'Tool Execution', href: '/docs/roxai/tool-execution', icon: Zap, description: 'Dispatch, execution, and result synthesis' },
    ],
  },
  {
    title: 'LLM Integration',
    links: [
      { name: 'LLM Abstraction', href: '/docs/roxai/llm-layer', icon: Cpu, description: 'Provider-agnostic interface' },
      { name: 'Providers', href: '/docs/roxai/providers', icon: Globe, description: 'Ollama, OpenAI, Anthropic, Gemini, Groq' },
      { name: 'Failover', href: '/docs/roxai/failover', icon: Shield, description: 'Retry logic and provider fallback' },
    ],
  },
  {
    title: 'Integration',
    links: [
      { name: 'Sessions', href: '/docs/roxai/sessions', icon: MessageSquare, description: 'Conversation and context management' },
      { name: 'Configuration', href: '/docs/roxai/configuration', icon: Settings, description: 'model.json and runtime settings' },
      { name: 'Dashboard', href: '/docs/roxai/dashboard', icon: LayoutDashboard, description: 'Web UI and REST API endpoints' },
      { name: 'RoxKV Bridge', href: '/docs/roxai/roxkv-bridge', icon: Link2, description: 'How RoxAI interacts with the database' },
    ],
  },
];

// Flat list for search indexing
export const allDocPages = [
  // High-level
  { name: 'System Architecture', href: '/docs/architecture', section: 'Overview', description: 'How RoxKV and RoxAI fit together' },
  // RoxKV
  ...roxkvNav.flatMap(group => group.links.map(link => ({ ...link, section: 'RoxKV' }))),
  // RoxAI
  ...roxaiNav.flatMap(group => group.links.map(link => ({ ...link, section: 'RoxAI' }))),
];
