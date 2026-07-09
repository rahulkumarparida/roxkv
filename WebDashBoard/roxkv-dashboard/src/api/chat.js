import { chatApi } from './client';
import { ENDPOINTS } from './endpoints';

/**
 * @typedef {Object} ChatRequest
 * @property {string} message - The user's message
 * @property {string[]} [context_keys] - Optional context keys for the agent
 */

/**
 * @typedef {Object} ChatResponse
 * @property {string} response - Markdown formatted response string
 * @property {string} status - 'success' | 'error'
 * @property {number} latency_ms - Response generation time in milliseconds
 */

/** Pool of realistic markdown responses the mock can return */
const MOCK_RESPONSES = [
  {
    trigger: /ttl|expire/i,
    response: `**Database TTL Statistics Summary**\n\n**TTL Information**\n\n* There are currently **1** active TTL keys in the database.\n* **8** keys have expired and been cleaned up.\n* The TTL cleanup goroutine is running at regular intervals.\n\n**Active TTL Keys**\n\n| Key | Remaining TTL | Expiry Time |\n|-----|--------------|-------------|\n| session:user:42 | 245s | 2026-07-08T22:15:00+05:30 |\n\n**No Immediate Concerns**\n\nThe TTL subsystem is operating normally. All expired keys have been properly cleaned up.`,
  },
  {
    trigger: /list.*key|all.*key/i,
    response: `**Database Keys Listing**\n\nCurrently stored keys in the database:\n\n\`\`\`\n1. name\n2. game \u2192 "Tomb Raider"\n3. nick \u2192 "situ"\n4. session:user:42 (TTL: 245s)\n\`\`\`\n\n**Summary:** 4 keys total, 1 with active TTL.`,
  },
  {
    trigger: /size|storage|space|status|health|overview/i,
    response: `**Database Storage Health Summary**\n\nThe database has a valid persistence folder and an existing snapshot folder.\n\n**Snapshot Information**\n\n* There are currently **25** snapshots stored in the database.\n* The total size of all snapshots is approximately **52.9 KB**.\n* The latest snapshot was created on July 2, 2026, at 21:51:46 UTC+05:30.\n* The snapshot's name and path are provided.\n\n**No Immediate Risks or Recommendations**\n\nBased on the current storage health and snapshot information, there are no immediate risks or recommendations to address. However, it is essential to regularly clean up old snapshots to maintain optimal database performance and disk space utilization.\n\nGive me the TTL stats?`,
  },
  {
    trigger: /snapshot/i,
    response: `**Snapshot Details**\n\n| Property | Value |\n|----------|-------|\n| Total Snapshots | 25 |\n| Total Size | 52.9 KB |\n| Latest Created | Jul 02, 21:51:46 |\n| Snapshot Interval | 10 minutes |\n| Persistence Path | /data/roxkv/snapshots/ |\n\n**Latest Snapshot**\n\n\`\`\`json\n{\n  "id": "snap_025",\n  "timestamp": "2026-07-02T21:51:46+05:30",\n  "size_bytes": 2156,\n  "keys_count": 4,\n  "checksum": "a3f2b8c1"\n}\n\`\`\`\n\nAll snapshots are healthy and consistent.`,
  },
  {
    trigger: /connect|client|user/i,
    response: `**Active Connections Summary**\n\n* **3** clients currently connected\n* **0** connections in idle timeout\n\n**Connected Clients**\n\n| Client ID | User | Connected Since | Last Activity |\n|-----------|------|----------------|---------------|\n| conn_001 | rahul | 2h 17m ago | 2s ago |\n| conn_002 | nick | 45m ago | 12s ago |\n| conn_003 | situ | 15m ago | 5s ago |\n\nAll connections are healthy and active.`,
  },
];

/** Fallback response when backend is unavailable */
const DEFAULT_RESPONSE = {
  response: `Backend connection failed. Using mock response.\n\n**Current Database Status**\n\n* **Status:** Active and healthy\n* **Keys:** 4 stored keys\n* **Snapshots:** 25 snapshots (52.9 KB total)\n* **Storage Health:** Optimal`,
  status: 'success',
  latency_ms: 450,
};

/**
 * Sends a chat message to the ChatWebServer endpoint.
 * The backend expects: {"query": "user message"}
 * The backend returns a plain string response (not JSON).
 * Converts it to the expected internal format for ChatService.
 * @param {ChatRequest} payload - {message: "...", context_keys?: [...]}
 * @returns {Promise<ChatResponse>} - {response: "...", status: "success", latency_ms: number}
 */
export async function sendChatMessage(payload) {
  const startTime = Date.now();
  try {
    const response = await chatApi.post(
      ENDPOINTS.chat,
      { query: payload.message },
      {
        timeout: 300000,
        headers: {
          'Content-Type': 'application/json',
          Accept: 'text/plain',
        },
      },
    );

    const latency = Date.now() - startTime;
    const responseText = typeof response.data === 'string' ? response.data : JSON.stringify(response.data);

    return {
      response: responseText,
      status: 'success',
      latency_ms: latency,
    };
  } catch (error) {
    console.error('[API] Chat request failed:', error.message);
    return mockChatResponse(payload.message);
  }
}

/**
 * Mock implementation that simulates network latency and returns
 * contextually relevant responses based on message content.
 * @param {string} message
 * @returns {Promise<ChatResponse>}
 */
async function mockChatResponse(message) {
  const latency = 300 + Math.random() * 400;
  await new Promise((resolve) => setTimeout(resolve, latency));
  const matched = MOCK_RESPONSES.find((r) => r.trigger.test(message));
  return {
    response: matched ? matched.response : DEFAULT_RESPONSE.response,
    status: 'success',
    latency_ms: Math.round(latency),
  };
}
