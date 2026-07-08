# RoxAI Dashboard — SSE Data Flow Documentation

This document explains the architecture of how the frontend React application consumes, processes, and displays real-time Server-Sent Events (SSE) from the RoxKV Go backend.

---

## 1. Endpoint Configuration (`src/api/endpoints.js`)

All communication points are centralized in `endpoints.js`. This ensures that changing the backend host only requires modifying a single variable.

```javascript
export const API_BASE = 'http://localhost:8081';

export const SSE_ENDPOINTS = {
  monitor: `${API_BASE}/api/events/monitor`,
  storage: `${API_BASE}/api/events/storage`,
  database: `${API_BASE}/api/events/database`,
  activity: `${API_BASE}/api/events/activity`,
};
```

---

## 2. Connection Management (`src/services/SSEService.js`)

The `SSEService` acts as a robust network layer that manages the lifecycle of `EventSource` connections independently of the React component tree.

### Core Responsibilities
- **Connection Creation**: Connects to the defined endpoints using standard browser `EventSource` APIs.
- **Robust Reconnection**: If the server drops the connection or goes offline, `SSEService` automatically attempts to reconnect using an **Exponential Backoff Strategy** (delaying retries from 1s, 2s, 4s, 8s, up to 30s) to prevent browser network thrashing.
- **Pub/Sub Notification System**: Instead of passing raw data directly into the DOM, it acts as an Event Bus. Components (or Contexts) "subscribe" to a channel (e.g., `monitor`) and receive a callback whenever a new JSON payload arrives.

---

## 3. Data Ingestion & State Processing (`src/context/DashboardContext.jsx`)

The `DashboardContext` is the "brain" of the frontend. It mounts once at the root level and subscribes to `SSEService`. It listens to the incoming streams and transforms the raw backend data into structured React state.

### How Each Endpoint is Processed:

#### A. `/api/events/monitor`
Provides high-frequency system telemetry (CPU, Goroutines, Memory).
- **Processing**: The Context parses the incoming `allocatedMemMB`, `systemMemMB`, `heapAllocMB`, `goroutines`, and `gcCycles`.
- **Graph History**: To render the sparkline graphs, it appends the new values to a `history` array inside the state. Once the array length exceeds `MAX_GRAPH_DATA_POINTS` (e.g., 30 data points), it slices the oldest data point off to create a smooth rolling window effect.
- **Derived Fields**: It also pulls out `hostname` (`data.user`), `connectedUsers` (`data.connectedclient`), and `uptimeSeconds` (`data.serverUptime`) to feed the System Overview.

#### B. `/api/events/storage`
Provides snapshot metrics.
- **Processing**: Updates the `snapshotCountdown` by converting the incoming `data.nextSnap` (which is in nanoseconds) to seconds (divided by `1,000,000,000`). It also extracts `totalSnapshots` and `latestSnapshotTimestamp` to update the Database Health context.

#### C. `/api/events/database`
Provides data size and TTL metrics.
- **Processing**: Directly extracts the `databaseSizeBytes`, and gracefully accesses `data.ttlMetrics.activeTTLKeys` and `data.ttlMetrics.totalExpiredKeys` using optional chaining to prevent crashes if the object is malformed.

#### D. `/api/events/activity`
Provides the raw text logs (e.g., `"2026-07-08 12:34:56 INFO key set"`).
- **Processing**: 
  - Validates that the event array doesn't just contain `"Not Data Found"`.
  - Splits the string by space to isolate the Date, Time, Flag (`INFO`, `SUCCESS`, `ERROR`), and the rest of the message description.
  - Transforms these strings into structured JS objects with `id`, `type`, `description`, and `timestamp`.
  - Limits the stored state to `MAX_ACTIVITY_ITEMS` to prevent memory leaks from a growing DOM.

---

## 4. UI Rendering (`src/hooks/useDashboard.js` & Components)

Because all of this data updates rapidly (often multiple times a second), the application uses **custom hooks with granular selectors** to prevent unnecessary re-renders.

```javascript
// Example of a granular hook
export function useMetrics() {
  const { metrics } = useDashboard();
  return metrics;
}
```

### Component Consumption
- **`SystemMonitor` / `MetricCard`**: Subscribes only to `useMetrics()`. Passes the rolling `history` arrays into memoized Recharts components. Because the graphs are memoized (`React.memo`), they only repaint the SVG when the specific array reference updates, maintaining high 60fps performance.
- **`ActivityFeed`**: Consumes `useActivityFeed()`. It maps over the processed event objects. A simple conditional block applies Tailwind color classes based entirely on the parsed type:
  - `SUCCESS` → Green (`text-green-400`)
  - `INFO` → Blue (`text-blue-400`)
  - `ERROR` → Red (`text-red-400`)
  - It wraps items in Framer Motion's `<AnimatePresence>` so that new events slide in gracefully from the top while old ones fade out.

## Summary Flowchart

\`\`\`mermaid
graph TD
    A[Go Server] -->|SSE Stream| B(SSEService.js)
    B -->|Pub/Sub Event Bus| C(DashboardContext.jsx)
    
    C -->|Extracts History & Metrics| D{useMetrics}
    C -->|Extracts TTL & Snapshots| E{useDatabaseHealth}
    C -->|Parses Logs into Flags| F{useActivityFeed}
    
    D --> G[Monitor Graphs]
    E --> H[Sidebar Cards]
    F --> I[Activity Log UI]
\`\`\`
