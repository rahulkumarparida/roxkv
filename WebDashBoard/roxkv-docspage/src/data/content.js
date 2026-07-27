import {
  Database, Activity, HardDrive, Camera, Clock, Radio,
  Save, BarChart3, Users, ScrollText, LayoutDashboard,
  Bot, Server, Cpu, Globe, Terminal, Container,
  Zap, GitBranch, Lock, Shield, TrendingUp, Brain, Layers
} from 'lucide-react';

// ─── Hero ────────────────────────────────────────────────────────────────────
export const hero = {
  badge: 'AI-Powered Database Operations',
  titleLine1: 'Talk to Your Database.',
  titleLine2: 'Not Its Commands.',
  subtitle:
    'RoxAI is an intelligent orchestration layer built on top of RoxKV that enables natural language interaction with a custom in-memory database using AI-powered tool execution.',
  primaryCta: { label: 'Watch Demo', href: 'https://youtu.be/pscPMhkcZQg?si=7TXA0KQUJfyAoS_T' },
  secondaryCta: { label: 'View GitHub', href: 'https://github.com/rahulkumarparida/roxkv' },
};

// ─── Statistics ──────────────────────────────────────────────────────────────
export const stats = [
  { value: 5800, suffix: '+', label: 'Lines of Code' },
  { value: 15, suffix: '', label: 'Go Packages' },
  { value: 45, suffix: '', label: 'AI Tools' },
  { value: 50, suffix: '+', label: 'Supported Operations' },
  { value: 12, suffix: '', label: 'Monitoring Metrics' },
  { value: 5, suffix: '', label: 'Data Types' },
  { value: 100, suffix: '+', label: 'Concurrent Clients' },
  { value: 4, suffix: '', label: 'Real-time Streams' },
];

// ─── What is RoxAI ──────────────────────────────────────────────────────────
export const whatIsRoxAI = {
  heading: 'What is RoxAI?',
  points: [
    'RoxAI is not just a chatbot.',
    'It is an orchestration layer.',
    'It understands user intent.',
    'Selects relevant tools from 45 registered AI tools.',
    'Executes operations on the RoxKV database.',
    'Combines results from multiple tools.',
    'Produces intelligent summaries streamed via SSE.',
  ],
  workflow: [
    { label: 'User', icon: Users },
    { label: 'Master Agent', icon: Brain },
    { label: 'Semantic Tool Routing', icon: GitBranch },
    { label: '45 AI Tools', icon: Zap },
    { label: 'RoxKV Store', icon: Database },
    { label: 'Intelligent Response', icon: Bot },
  ],
};

// ─── Why It Exists ──────────────────────────────────────────────────────────
export const whyItExists = {
  heading: 'Why It Exists',
  description:
    'Traditional database management requires remembering commands, monitoring multiple dashboards, and manually collecting information. RoxAI allows users to simply ask questions in natural language.',
  examples: [
    { query: 'What is the health of my database?', tool: 'get_system_health' },
    { query: 'Show active clients.', tool: 'get_server_info' },
    { query: 'How much RAM is being used?', tool: 'get_memory_usage' },
    { query: 'List expired keys.', tool: 'get_expired_keys_count' },
    { query: 'What keys are stored?', tool: 'get_all_keys' },
    { query: 'Create a snapshot now.', tool: 'create_snapshot' },
  ],
};

// ─── Capabilities ───────────────────────────────────────────────────────────
export const capabilities = [
  { icon: Database, title: 'Database Operations', description: 'Full CRUD with strings, lists, hashes, sets, and sorted sets — GET, SET, DEL, KEYS, and more.' },
  { icon: Activity, title: 'Monitoring', description: 'Real-time CPU, memory, goroutine tracking with periodic metric sampling.' },
  { icon: HardDrive, title: 'Storage Analysis', description: 'Inspect memory usage, key counts, data types, and storage distribution.' },
  { icon: Camera, title: 'Snapshots', description: 'Point-in-time snapshots with GOB encoding. Create, list, and restore.' },
  { icon: Clock, title: 'TTL Analysis', description: 'Manage key expiration with TTL inspection, persistence, and expired key tracking.' },
  { icon: Radio, title: 'Pub/Sub', description: 'Channel-based publish/subscribe with real-time SSE streaming to subscribers.' },
  { icon: Save, title: 'Persistence', description: 'Dual strategy: Append-Only File for durability + snapshots for point-in-time recovery.' },
  { icon: BarChart3, title: 'Runtime Statistics', description: 'Server info, uptime, active goroutines, commands processed, and connection stats.' },
  { icon: Users, title: 'Client Monitoring', description: 'Track connected clients, concurrent connections, and per-client activity.' },
  { icon: ScrollText, title: 'Activity Logging', description: 'Ring buffer of recent operations with timestamps, client info, and command details.' },
  { icon: LayoutDashboard, title: 'Real-time Dashboard', description: 'React dashboard with 8 feature panels: keys, monitoring, pub/sub, snapshots, chat, terminal.' },
  { icon: Bot, title: 'AI Tool Calling', description: '45 registered tools across 5 categories orchestrated by a Master Agent via metadata-driven scoring.' },
];

// ─── Dashboard Showcase ─────────────────────────────────────────────────────
export const dashboardTabs = [
  { id: 'dashboard', label: 'Dashboard', description: 'Overview with metrics, key counts, and system health at a glance.' },
  { id: 'cli', label: 'CLI', description: 'Interactive command-line interface with syntax highlighting and tab completion.' },
  { id: 'monitoring', label: 'Live Monitoring', description: 'Real-time charts for CPU, memory, goroutines, and command throughput.' },
  { id: 'chat', label: 'AI Chat', description: 'Natural language interface with SSE-streamed AI responses.' },
  { id: 'architecture', label: 'Architecture', description: 'System architecture showing all interconnected components.' },
];

// ─── Architecture ───────────────────────────────────────────────────────────
export const architectureNodes = [
  { id: 'user', label: 'User', icon: Users, level: 0 },
  { id: 'dashboard', label: 'React Dashboard', icon: LayoutDashboard, level: 1 },
  { id: 'api', label: 'HTTP / SSE API', icon: Globe, level: 2 },
  { id: 'agent', label: 'Master Agent (Ollama)', icon: Brain, level: 3 },
  { id: 'router', label: 'Semantic Tool Router', icon: GitBranch, level: 4 },
  { id: 'tools', label: 'Tool Layer (45 Tools)', icon: Zap, level: 5 },
  { id: 'roxkv', label: 'RoxKV Core', icon: Database, level: 6 },
];

export const architectureLeaves = [
  { label: 'Storage', icon: HardDrive },
  { label: 'Monitoring', icon: Activity },
  { label: 'PubSub', icon: Radio },
  { label: 'Persistence', icon: Save },
  { label: 'TTL', icon: Clock },
];

// ─── Engineering Decisions ──────────────────────────────────────────────────
export const engineeringDecisions = [
  {
    title: 'Semantic Tool Routing',
    description: 'The LLM selects tools based on natural language understanding — not keyword matching or hardcoded rules.',
  },
  {
    title: 'Single Master Agent',
    description: 'One orchestrating agent with two-pass execution: first pass selects tools, second pass synthesizes results.',
  },
  {
    title: 'Server Sent Events',
    description: 'AI responses stream token-by-token via SSE for real-time, low-latency user experience.',
  },
  {
    title: 'Context-based React Architecture',
    description: 'Global state managed through React Context — no Redux overhead, clean component boundaries.',
  },
  {
    title: 'Modular Tool Design',
    description: 'Each tool is independently registered with name, description, parameter schema, and handler function.',
  },
  {
    title: 'Go Concurrency',
    description: 'Goroutines per connection, sync.RWMutex for the store, buffered channels for pub/sub, and worker pools with context cancellation.',
  },
  {
    title: 'Local AI Execution',
    description: 'All inference runs locally via Ollama with llama3.2:3b — no cloud APIs, full privacy and control.',
  },
  {
    title: 'Containerized Deployment',
    description: 'Three Docker services (RoxKV, Dashboard, Ollama) orchestrated with Docker Compose for one-command setup.',
  },
];

// ─── Technology Stack ───────────────────────────────────────────────────────
export const techStack = [
  { name: 'Go', color: '#00ADD8' },
  { name: 'React', color: '#61DAFB' },
  { name: 'TailwindCSS', color: '#06B6D4' },
  { name: 'SSE', color: '#A78BFA' },
  { name: 'Ollama', color: '#F9FAFB' },
  { name: 'Docker', color: '#2496ED' },
  { name: 'TCP', color: '#EF4444' },
  { name: 'HTTP', color: '#F59E0B' },
  { name: 'JSON', color: '#10B981' },
  { name: 'Goroutines', color: '#00ADD8' },
  { name: 'Mutex', color: '#8B5CF6' },
  { name: 'Channels', color: '#EC4899' },
  { name: 'Markdown', color: '#F9FAFB' },
  { name: 'Lucide', color: '#F97316' },
];

// ─── Demo ───────────────────────────────────────────────────────────────────
export const demo = {
  heading: 'RoxAI Is Under Active Development',
  description:
    'The live demo experience is still being refined. RoxAI is under continuous development, with ongoing work focused on making the orchestration flow more robust, improving reliability across tool execution paths, tightening edge-case handling, and polishing the end-to-end user experience before a full showcase is published.',
  details: [
    'The current focus is on strengthening the AI-to-tool execution pipeline so natural language requests behave more consistently under real usage.',
    'Work is continuing on stability improvements, better error recovery, cleaner streaming responses, and stronger operational visibility throughout the system.',
    'Rather than shipping a rushed demo section, this page now reflects that RoxAI is being actively improved to make the final experience more dependable, more accurate, and more production-ready.',
  ],
};

// ─── GitHub ─────────────────────────────────────────────────────────────────
export const github = {
  repoName: 'rahulkumarparida/roxkv',
  ownerName: 'rahulkumarparida',
  description:
    'A Redis-like in-memory key-value store built from scratch in Go, featuring 4 concurrent servers (TCP, AI TCP, SSE, AI HTTP), 45 AI tools, real-time dashboard, and containerized deployment.',
  stars: '★',
  commits: '100+',
  repoUrl: 'https://github.com/rahulkumarparida/roxkv',
  docsUrl: 'https://github.com/rahulkumarparida/roxkv#readme',
  profileUrl: 'https://github.com/rahulkumarparida/',
};

// ─── Run Locally ────────────────────────────────────────────────────────────
export const runLocally = {
  heading: 'Run Locally',
  description: 'Get RoxAI running on your machine in under a minute. The AI backend uses Ollama locally — no API keys required.',
  steps: [
    { label: 'Clone the repository', command: 'git clone https://github.com/rahulkumarparida/roxkv.git && cd roxkv' },
    { label: 'Start with Docker Compose', command: 'docker-compose up --build' },
    { label: 'Open the dashboard', command: 'open http://localhost:5173' },
    { label: 'Or connect via CLI', command: 'go run ./cli/cli.go' },
  ],
  note: 'Ollama runs locally inside Docker with the qwen3:1.7b model. No cloud APIs, full privacy.',
};

// ─── Gallery ────────────────────────────────────────────────────────────────
export const gallery = [
  { id: 'dashboard', label: 'Dashboard', description: 'Main dashboard with health, keys, and system metrics.' },
  { id: 'cli', label: 'CLI', description: 'Terminal-style interface for direct RoxKV command execution.' },
  { id: 'architecture', label: 'Architecture', description: 'High-level system view of the RoxAI orchestration flow.' },
  { id: 'monitoring', label: 'Monitoring', description: 'Live operational metrics for CPU, memory, and runtime activity.' },
];

// ─── Roadmap ────────────────────────────────────────────────────────────────
export const roadmap = [
  { title: 'Redis Protocol', description: 'Full RESP protocol compatibility for drop-in Redis client support.', status: 'planned' },
  { title: 'Distributed Storage', description: 'Shard data across multiple nodes with consistent hashing.', status: 'planned' },
  { title: 'Authentication', description: 'Token-based auth with user management and session handling.', status: 'planned' },
  { title: 'Role Based Access', description: 'Fine-grained RBAC for read, write, admin, and tool-level permissions.', status: 'planned' },
  { title: 'Historical Analytics', description: 'Time-series storage of metrics for trend analysis and anomaly detection.', status: 'planned' },
  { title: 'Smarter Planning', description: 'Multi-step query planning with tool chaining and dependency resolution.', status: 'planned' },
  { title: 'Cluster Support', description: 'Multi-node cluster with leader election and automatic failover.', status: 'planned' },
];

// ─── About ──────────────────────────────────────────────────────────────────
export const about = {
  heading: 'About the Project',
  description:
    'RoxKV was built completely from scratch as an educational systems programming project — a deep dive into building a Redis-like database server in Go. What started as a learning exercise in TCP servers, concurrent data structures, and protocol parsing evolved into a full-featured AI-powered operational platform with natural language database management, real-time monitoring, and containerized deployment.',
  repoUrl: 'https://github.com/rahulkumarparida/roxkv',
  profileUrl: 'https://github.com/rahulkumarparida/',
  author: 'Rahul Kumar Parida',
};

// ─── Footer ─────────────────────────────────────────────────────────────────
export const footer = {
  projectName: 'RoxAI',
  tagline: 'Talk to Your Database. Not Its Commands.',
  links: {
    product: [
      { label: 'GitHub', href: 'https://github.com/rahulkumarparida/roxkv' },
      { label: 'Documentation', href: 'https://github.com/rahulkumarparida/roxkv#readme' },
      { label: 'License', href: 'https://github.com/rahulkumarparida/roxkv/blob/main/LICENSE' },
    ],
    connect: [
      { label: 'GitHub', href: 'https://github.com/rahulkumarparida/' },
      { label: 'LinkedIn', href: 'https://www.linkedin.com/in/rahul-kumar-parida-b6219a292/' },
      { label: 'X', href: 'https://x.com/rahulkuparida' },
    ],
  },
  copyright: `© ${new Date().getFullYear()} RoxAI. All rights reserved.`,
  madeWith: 'Made with 🫶 Go + React + AI',
  madeBy: 'Made by rahulkumarparida 🖤'
};

// ─── Navigation ─────────────────────────────────────────────────────────────
export const navigation = [
  { label: 'Features', href: '/#capabilities' },
  { label: 'Architecture', href: '/#architecture' },
  { label: 'Tech Stack', href: '/#tech-stack' },
  { label: 'Demo', href: '/#demo' },
  { label: 'Run Locally', href: '/#run-locally' },
  { label: 'Resp Docs', href: '/docs' },
];
