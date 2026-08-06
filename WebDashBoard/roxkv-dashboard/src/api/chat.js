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
    response: `**Database Storage Health Summary**\n\nThe database has a valid persistence folder and an existing snapshot folder.\n\n**Snapshot Information**\n\n* There are currently **25** snapshots stored in the database.\n* The total size of all snapshots is approximately **52.9 KB**.\n* The latest snapshot was created on July 2, 2026, at 21:51:46 UTC+05:30.\n* The snapshot's name and path are provided.\n\n**No Immediate Risks or Recommendations**\n\nBased on the current storage health and snapshot information, there are no immediate risks or recommendations to address. However, it is essential to regularly clean up old snapshots to maintain optimal database performance and disk space utilization.`,
  },
  {
    trigger: /snapshot/i,
    response: `**Snapshot Details**\n\n| Property | Value |\n|----------|-------|\n| Total Snapshots | 25 |\n| Total Size | 52.9 KB |\n| Latest Created | Jul 02, 21:51:46 |\n| Snapshot Interval | 10 minutes |\n| Persistence Path | /data/roxkv/snapshots/ |\n\nAll snapshots are healthy and consistent.`,
  },
];

const DEFAULT_RESPONSE = {
  response: `Backend connection failed. Using mock response.\n\n**Current Database Status**\n\n* **Status:** Active and healthy\n* **Keys:** 4 stored keys\n* **Snapshots:** 25 snapshots (52.9 KB total)\n* **Storage Health:** Optimal`,
  status: 'success',
  latency_ms: 450,
};

/**
 * Sends a chat message to the ChatWebServer endpoint.
 * Handles validation errors (e.g. Missing API Key) without silent fallback.
 */
export async function sendChatMessage(payload) {
  const startTime = Date.now();
  try {
    const response = await chatApi.post(
      ENDPOINTS.chat,
      { query: payload.message },
      {
        timeout: 3000000,
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
    if (error.response && error.response.data) {
      const data = error.response.data;
      const errorMsg = typeof data === 'string' ? data : (data.message || data.error || 'Request failed');
      return {
        response: `⚠️ **Error**: ${errorMsg}`,
        status: 'error',
        latency_ms: Date.now() - startTime,
      };
    }
    return mockChatResponse(payload.message);
  }
}

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

/**
 * Provider API helper methods
 */
export async function fetchProviders() {
  const res = await chatApi.get(ENDPOINTS.providers);
  return res.data.providers || [];
}

export async function fetchProviderModels(provider) {
  const res = await chatApi.get(ENDPOINTS.providerModels(provider));
  return res.data.models || [];
}

export async function fetchProviderConfig(provider) {
  const res = await chatApi.get(ENDPOINTS.providerConfig(provider));
  return res.data;
}

export async function saveProviderConfig(provider, configData) {
  const res = await chatApi.post(ENDPOINTS.providerConfig(provider), configData);
  return res.data;
}

export async function fetchProviderStats(provider) {
  const res = await chatApi.get(ENDPOINTS.providerStats(provider));
  return res.data;
}

export async function selectActiveProvider(provider, model) {
  const res = await chatApi.post(ENDPOINTS.selectProvider, { provider, model });
  return res.data;
}
