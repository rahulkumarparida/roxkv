# Agent Router Refactor

## 1. Overview

RoxKV uses a local LLM (Llama 3.2:3b via Ollama) as an orchestration layer. The Master Agent receives a user query, presents available tools to the LLM, and the LLM decides which tools to invoke and in what order.

Previously, **every request** exposed the complete list of ~46 tools to the LLM, regardless of what the user actually asked. For a 3B-parameter model running on consumer hardware with a 2048-token context window, this created significant problems:

- **Prompt bloat**: Tool descriptions consumed a large portion of the context window.
- **Poor tool selection**: The LLM struggled to pick the right tool from dozens of candidates.
- **Hallucinations**: Confused by irrelevant tool descriptions, the model would fabricate tool names or arguments.
- **High latency**: Reading and reasoning over 46 tool descriptions slowed every inference call.

This refactor introduces a **metadata-driven tool routing system** that scores every tool against the user query and exposes only the relevant subset. The LLM now sees fewer, more relevant tools and makes better decisions faster.

---

## 2. Previous Architecture

```
User Query
    │
    ▼
Master Agent
    │
    ▼
Expose ALL tools (~46)
    │
    ▼
LLM reads all tool descriptions
    │
    ▼
LLM selects tools (error-prone)
    │
    ▼
Execute tools
    │
    ▼
Synthesize response
```

### Problems

| Issue | Impact |
|---|---|
| Prompt size | ~46 tool descriptions consumed most of the 2048-token context |
| Tool overload | The LLM had to reason over dozens of irrelevant tools per query |
| Hallucination risk | Confusion between similarly-named tools led to invented arguments |
| Latency | Every request paid the inference cost of processing all tool descriptions |
| Maintainability | Adding a tool required appending to a manually maintained 46-line slice |

---

## 3. Updated Architecture

```
User Query
    │
    ▼
Metadata Router
    │
    ▼
Score every tool against the query
    │
    ▼
Select top-N relevant tools
    │
    ▼
Expose filtered tool list (~5-10)
    │
    ▼
LLM reads relevant descriptions only
    │
    ▼
LLM selects tools (accurate)
    │
    ▼
Execute tools
    │
    ▼
Synthesize response
```

### Key improvements

- **Routing**: A lightweight, rule-based router scores each tool using metadata.
- **Metadata**: Every tool carries keywords, synonyms, example phrases, and priority weights.
- **Scoring**: Relevance is computed per-tool from keyword matches, synonym matches, example overlaps, and category bonuses — all weighted by the tool's priority.
- **Reduced context**: The LLM sees only the top-scoring tools (default: 10), dramatically shrinking the prompt.
- **Better tool selection**: Fewer, more relevant candidates mean the LLM makes better choices.
- **Fallback**: If no tool matches (e.g., a greeting), all tools are exposed to preserve backward compatibility.

---

## 4. Metadata System

### ToolMetadata struct

```go
type ToolMetadata struct {
    Name        string       // Unique identifier matching api.ToolFunction.Name
    Tool        api.Tool     // The actual api.Tool value (unchanged)
    Category    string       // Logical grouping: "KV", "Monitoring", "PubSub", "Storage", "Persistence"
    Description string       // Routing-only description (NOT sent to LLM)
    Keywords    []string     // Primary routing keywords
    Synonyms    []string     // Natural language alternatives
    Examples    []string     // Representative user request phrases
    Priority    int          // Numeric routing weight (higher = stronger signal)
    IsMutation  bool         // Whether the tool changes application state
}
```

### Field details

| Field | Purpose | Routing role |
|---|---|---|
| `Name` | Unique ID, matches `api.ToolFunction.Name` | Identification only |
| `Tool` | The `api.Tool` value passed to the LLM | Returned when selected |
| `Category` | Logical grouping (e.g., "KV", "PubSub") | Category bonus scoring |
| `Description` | Human-readable explanation for maintainability | Not used in scoring |
| `Keywords` | Primary trigger words (e.g., "delete", "key") | Exact-match scoring (+3 × priority/10) |
| `Synonyms` | Alternative phrasing (e.g., "remove", "erase") | Synonym-match scoring (+2 × priority/10) |
| `Examples` | Representative queries (e.g., "show all keys") | Example-overlap scoring (+2 × priority/10) |
| `Priority` | Confidence weight (1–10 scale) | Multiplier for all match scores |
| `IsMutation` | Whether the tool writes data | Available for future mutation-aware routing |

---

## 5. Routing Algorithm

### Scoring process

For each registered tool, the router computes a relevance score:

```
For each query word (after stop-word removal):
  For each keyword:
    if exact match:   score += 3.0 × (priority / 10)
    if partial match:  score += 1.0

  For each synonym:
    if exact match:   score += 2.0 × (priority / 10)
    if partial match:  score += 1.0

  For each example phrase:
    count overlapping words with query
    score += overlap × 2.0 × (priority / 10)

  If query contains the category name:
    score += 1.0
```

### Selection

1. Tools with `score > 0` are collected.
2. Sorted descending by score.
3. Top N (default: 10) are returned.
4. If **no** tool scores above zero, **all** tools are returned (fallback).

### Why the LLM performs better

- **Fewer choices**: Instead of 46 tools, the LLM sees 5–10 highly relevant ones.
- **Reduced prompt**: Tool descriptions consume fewer tokens, leaving more room for reasoning.
- **Clearer intent**: When all presented tools are relevant, the LLM's selection accuracy improves.
- **Faster inference**: Less input text means faster token processing.

---

## 6. Key Code Changes

### New files

#### `agents/registry/metadata.go`
- **Responsibility**: Defines the `ToolMetadata` struct and a thread-safe global registry.
- **Exports**: `Register()`, `AllMetadata()`, `AllTools()`.
- **Why**: Provides the central data structure and storage for tool metadata, decoupled from any specific tool implementation.

#### `agents/router/router.go`
- **Responsibility**: Implements the `RouteTools()` scoring engine.
- **Exports**: `RouteTools()`, `DefaultMaxTools`.
- **Why**: Encapsulates the metadata-driven routing logic. Adding new tools requires zero changes to the router — it reads metadata dynamically.

#### `agents/Master/tools_registry.go`
- **Responsibility**: Registers all 46 tools with rich metadata via `init()`.
- **Why**: Lives in the `master` package (alongside the tool-definition functions) to avoid import cycles. The `init()` function runs at startup, populating the global registry before any query is processed.

### Modified files

#### `agents/Master/Masteragent.go`
- **Previous responsibility**: Manually assembled a 46-element `[]api.Tool` slice and passed it to `session.Run()`.
- **Updated responsibility**: Calls `router.RouteTools(query, registry.AllMetadata(), router.DefaultMaxTools)` to get only the relevant tools.
- **What changed**: The 46-line tool slice was replaced with a single line. Two imports were added (`registry`, `router`). The entire tool-execution `switch` block and all helper functions remain untouched.

### Unchanged files

| File | Reason unchanged |
|---|---|
| `agents/Master/constants.go` | Tool-definition functions are called by `tools_registry.go` — no modifications needed |
| `agents/session.go` | Session management is orthogonal to tool selection |
| `agents/results.go` | Result formatting is orthogonal to tool selection |
| `agents/constants.go` | Model constants are unchanged |
| `agents/KvAgent/*` | Sub-agent has its own small tool list (6 tools), no routing needed |
| `agents/MonitorAgent/*` | Sub-agent has its own small tool list (8 tools), no routing needed |
| `agents/PubSubAgent/*` | Sub-agent has its own small tool list (13 tools), no routing needed |
| `agents/StorageAgent/*` | Sub-agent has its own small tool list (18 tools), no routing needed |

---

## 7. Benefits

| Benefit | Before | After |
|---|---|---|
| **Prompt size** | ~46 tool descriptions per request | 5–10 relevant tool descriptions |
| **Tool selection accuracy** | LLM confused by irrelevant options | LLM sees only candidates that match the query |
| **Hallucination rate** | Higher — model confused by similar tools | Lower — fewer, clearer choices |
| **Inference latency** | Slower — processing all tool descriptions | Faster — fewer tokens to process |
| **Maintainability** | Manual 46-line slice, easy to break | `Register()` call with metadata, self-documenting |
| **Extensibility** | Adding a tool = edit the slice + edit the switch | Adding a tool = add one `Register()` call |
| **Scalability** | Degrades as tools grow | Router scales linearly; only top-N are exposed |

### Adding a new tool

Before:
1. Write the tool-definition function in `constants.go`
2. Add the function call to the 46-line slice in `Masteragent.go`
3. Add a `case` to the execution `switch` in `Masteragent.go`

After:
1. Write the tool-definition function in `constants.go`
2. Add a `registry.Register()` call in `tools_registry.go` with metadata
3. Add a `case` to the execution `switch` in `Masteragent.go`

The routing logic **never needs modification** — it reads metadata dynamically.

---

## 8. Future Improvements

### Embedding-based routing
Replace keyword matching with vector embeddings. Encode tool descriptions and user queries into the same vector space, then select tools by cosine similarity. This would handle novel phrasings that keyword matching misses.

### Semantic similarity scoring
Use a lightweight sentence-similarity model (e.g., all-MiniLM-L6-v2) to score query-to-tool similarity. Hybrid approach: combine embedding scores with keyword scores for robustness.

### Fuzzy matching
Add Levenshtein distance or n-gram similarity to handle typos and morphological variations (e.g., "delet" → "delete", "subscribing" → "subscriber").

### Hybrid router
Combine multiple scoring strategies (keyword, embedding, fuzzy, historical frequency) with configurable weights. Allow tuning per-deployment.

### Learning-based routing
Track which tools the LLM actually invokes for various queries. Over time, build a frequency model that biases routing toward historically successful tool selections.

### Dynamic routing thresholds
Instead of a fixed top-N cutoff, dynamically adjust the threshold based on the score distribution. If scores are tightly clustered, include more tools; if one tool dominates, include fewer.

### Mutation-aware routing
Use the `IsMutation` field to implement safety checks: require confirmation for mutation tools, or prefer read-only tools when the query is ambiguous.

### Context-aware routing
Consider conversation history when scoring tools. If the previous turn used a KV tool, bias toward KV tools for follow-up queries.
