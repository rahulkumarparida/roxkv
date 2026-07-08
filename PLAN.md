
### June 30: Finalize the minimal parser and metadata structure.
    - Update all the data struct to keep metadata
    - Keep the parser as is minimal chat flow
    - Implement such that to get snapshots of the memory and data in a interval store it in Hard disk
    - a simple parse that understands the commands

### July 1: Connect an LLM API and implement basic chat.
    - Create a simple chat interface for ollama
    - Connect to ollama API and use a base model
    - Prepare prompt for tool calling

### July 2: Add memory retrieval from RoxKV.
    - add the memory layer for ollama 
    - give snapshot acess and summarixzations

### July 3: Implement tool calling (GET, SET, LIST through the agent).
    - prepare prompts for tool calling
    - parser that parses the output from llm and execute the required tools



### July 4: Build a simple CLI or web demo.
    - Build a simple  CLI and web demo


### July 5: Test, fix bugs, and prepare your demo.
    - Test all the feature will all possible values


### July 8 : Build the Web DashBoard

Use this prompt and image of the purple dashboard 

RoxAI Dashboard UI
Objective
Build a modern AI Operations Dashboard for RoxKV + RoxAI.

This is ONLY the frontend.
Do NOT implement any backend logic.
Do NOT implement any real API calls.
Use realistic fake endpoints, robust mock data providers, and a strictly centralized API layer so that I can later hot-swap them with my Go backend.

The goal is to create a production-quality frontend architecture, not just a static demo page. Focus heavily on separation of concerns, clean component contracts, and scalable state management.

Tech Stack
Use:

ReactJS (Vite)

React Router (for potential future routing, keep it configured)

TailwindCSS (use arbitrary values for precise pixel/percentage matching if needed)

Framer Motion (for complex orchestrations and layout transitions)

React Markdown & remark-gfm (for robust markdown and code block rendering)

Lucide React Icons

Recharts (for performance-optimized streaming data visualization)

React Context API

Server Sent Events (EventSource)

Axios (strictly for the HTTP LLM endpoints)

Tailwind Merge & clsx (for clean, dynamic utility class construction)

Keep the project strictly modular.

Theme & Aesthetics
Dark futuristic terminal style.

Color Palette & Textures:

Primary Accent: Electric Purple (e.g., #a855f7 or similar).

Background: The entire application uses a fullscreen abstract, dark, futuristic image background with a heavy dark overlay (bg-black/80).

Dashboard Surface: Floating above the background with heavy glassmorphism cards. Utilize Tailwind's backdrop-blur-md or backdrop-blur-lg combined with subtle translucent backgrounds (e.g., bg-white/5).

Borders: Soft, 1px borders using translucent purple or gray (e.g., border-white/10 or border-purple-500/20).

Border Radius: Rounded corners (e.g., rounded-xl or rounded-2xl).

Interactions:

Purple glow on hover for interactive elements (e.g., hover:shadow-[0_0_15px_rgba(168,85,247,0.5)]).

Very smooth, eased animations for all state changes.

NO bright or jarring colors. Warning/Error states should use muted, terminal-appropriate amber or red.

Typography: Use a clean sans-serif for UI elements, and a strict monospace font (like JetBrains Mono or Fira Code) for metrics, logs, code blocks, and data readouts.

Everything should feel like an immersive, high-end AI control center. Use the provided image as a baseline inspiration, but refine the UI/UX micro-interactions.

Layout & Grid Architecture
Three-column responsive dashboard. Use CSS Grid for the main layout to ensure strict spatial constraints.

LEFT SIDEBAR (Metrics & Health)
Contains:
1. System Overview

Status (Animated pulsing dot indicator)

Hostname

Connected Users (Animated number counter)

Server Uptime (Live ticking format)

Master Agent Status

2. Database Health

Persistence Status

Storage Health (Color-coded)

Snapshot Count

Database Size (Auto-formatting bytes)

TTL Active / TTL Expired counters

Health Score (Progress ring or bar)

3. Snapshot Information

Latest Snapshot Timestamp

Next Snapshot Countdown

Total Snapshot Size

Action Button: "View Snapshots" (Includes hover glow)

CENTER PANEL (RoxAI Command Center)
This is the primary workspace.

Header:

ROXX AI MASTER AGENT title

Connection Status (Animated)

Chat Area:

Chat Messages (Distinct styles for User vs. Master Agent)

Full Markdown Support (Render code blocks, bold, lists correctly)

Tool execution animations (e.g., "Agent is querying database..." with a loading shimmer)

Simulated token-by-token typing animation for realistic LLM output

Message timestamps

Quick Action Chips (Horizontal scroll if needed):

Show TTL, List Keys, Storage Status, Runtime, Clients, Snapshots

Input Area:

Multiline Textarea (Auto-resizing up to a max height)

Enter = Send, Shift+Enter = New Line

Floating Send Button (Disabled while generating)

RIGHT SIDEBAR (Telemetry & Logs)
1. System Monitor (Recharts)

CPU Usage, RAM Usage, Disk Usage, Network In, Network Out

Each metric MUST have an Animated Line Graph.

Display the Current Value prominently.

Ensure graphs do not re-render the entire component tree on every tick (use memoization).

2. Realtime Activity Feed

Scrollable list, newest events prepended at the top.

Event types should have color-coded indicators:

User Connected (Green)

Client Timeout (Yellow/Amber)

SET/GET/DELETE key (Blue/Purple)

Snapshot Created (Purple)

AI Query (White)

Smooth height animation when a new item enters the list.

UI Behaviour & Animations
Data Intake: Cards gently pulse or highlight their border briefly when receiving updated SSE data.

Continuous Flow: Graphs animate continuously without visual stuttering.

Feed Dynamics: Activity feed uses Framer Motion AnimatePresence to smoothly slide down existing entries when prepending new ones.

Chat UX: Auto-scrolls to bottom on new message. Includes a typing indicator (... bouncing dots) and loading skeleton/shimmer for tool execution.

Micro-interactions: Glass cards slightly lift (translate-y-[-2px]) on hover. Buttons have satisfying click scaling (active:scale-95).

Folder Structure
Implement this highly scalable architecture:

Plaintext
src/
 ├── api/                # Strictly pure functions, no React hooks here
 │   ├── chat.js
 │   ├── sse.js
 │   └── endpoints.js
 ├── context/            # React Context providers
 │   ├── ChatContext.jsx
 │   ├── DashboardContext.jsx
 │   └── SSEContext.jsx
 ├── hooks/              # Custom hook wrappers for Contexts
 │   ├── useChat.js
 │   ├── useDashboard.js
 │   └── useSSE.js
 ├── pages/              # Route level components
 │   └── Dashboard.jsx
 ├── layouts/            # Persistent layout shells
 │   └── DashboardLayout.jsx
 ├── components/         # Feature-based component directories
 │   ├── Chat/
 │   │   ├── ChatBox.jsx
 │   │   ├── ChatInput.jsx
 │   │   ├── ChatMessage.jsx
 │   │   └── QuickActions.jsx
 │   ├── Sidebar/
 │   │   ├── SystemOverview.jsx
 │   │   ├── DatabaseHealth.jsx
 │   │   └── SnapshotInfo.jsx
 │   ├── Monitor/
 │   │   ├── MetricCard.jsx      # Generic reusable card for CPU/RAM etc.
 │   │   └── MonitorGraph.jsx    # Reusable Recharts component
 │   ├── Activity/
 │   │   └── ActivityFeed.jsx
 │   └── Common/         # Reusable UI primitives
 │       ├── GlassCard.jsx
 │       ├── StatusBadge.jsx
 │       ├── LoadingSpinner.jsx
 │       ├── AnimatedNumber.jsx
 │       └── TypingIndicator.jsx
 ├── services/           # Logic layer between API and Context
 │   ├── ChatService.js
 │   ├── SSEService.js
 │   └── MockDataGenerator.js # Robust engine for simulating the Go backend
 ├── utils/              # Pure helper functions
 │   ├── formatBytes.js
 │   ├── formatTime.js
 │   ├── constants.js
 │   └── classNames.js   # clsx/tailwind-merge helper
Central API Directory & Architecture
ALL communication must happen through a centralized API layer.
No component should directly call fetch(), axios.post(), or new EventSource().
Components should ONLY use Context (via custom hooks).

Mandatory Flow:
Component -> useChat() / useDashboard() -> Context -> Service Class -> API Function

This strict separation ensures that when the Go backend is ready, only the api/ directory needs to be modified. Ensure data models returned by the mock API mimic strict struct-like JSON responses.

Chat Endpoint (Mock HTTP)
Create a mock Axios implementation that simulates network latency and LLM generation time.

POST /api/chat
Payload: { "message": "...", "context_keys": [...] }
Response: { "response": "Markdown formatted string...", "status": "success", "latency_ms": 450 }

Include logic in ChatService to simulate a streamed response (revealing the text chunk by chunk) for realistic UX.

Server Sent Events (SSE) Architecture
The simulated backend continuously sends stringified JSON.
Example payload: data: "{\"cpu\":10,\"ram\":30}"

Strict Handling Rules:

Frontend MUST receive string -> JSON.parse() -> Validate -> Update Context -> Re-render UI.

Do not assume objects. Always implement try/catch around JSON.parse().

Create a centralized EventSource manager inside SSEContext. Every component subscribes to the Context; NO component instantiates EventSource directly.

Implement automatic reconnection logic (exponential backoff) in the SSEService if the connection drops.

Fake SSE URLs (Constants only):

JavaScript
const SSE_ENDPOINTS = {
  monitor: 'http://localhost:8081/sse/monitor',
  storage: 'http://localhost:8081/sse/storage',
  database: 'http://localhost:8081/sse/database',
  activity: 'http://localhost:8081/sse/activity'
}
Mock SSE Implementation:
Since there is no real backend, build a MockDataGenerator.js utility that mimics EventSource behavior. It should emit synthetic events at randomized intervals (e.g., using setInterval) to the SSEService to populate the dashboard with realistic, fluctuating data.

Visualizations & Feed
Graphs (Recharts):

Animated Line Charts.

Generate smoothly changing values in the mock generator (e.g., using Perlin noise or smoothed random walks, no erratic fake spikes).

Manage independent updates for CPU, RAM, Disk, Network In, Network Out.

Activity Feed:

Mock events should arrive every 2-5 seconds.

Animate them pushing down older items.

Retain a maximum of 50 items in state to prevent DOM bloat.

Deliverables Requirements
Generate:

The complete React/Vite code for the files listed in the folder structure.

Clean, documented code utilizing JSDoc comments to define the expected shapes of the mock data payloads.

Production-ready styling ensuring perfect responsive behavior.

Functional mock layer that makes the app look completely alive and operational upon npm run dev.

No backend code. No "TODO: implement here" placeholders in the frontend logic.

Deliver a dashboard that serves as a flawless, plug-and-play foundation for a high-performance backend.

Only toch the WebDashBoard/ directory dont do any changes in any other folder except this one 
Keep the design as same as the image provided