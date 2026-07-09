/**
 * @fileoverview Centralized endpoint configuration.
 * This is the ONLY place where URLs are defined.
 * 
 * ChatWebServer (Port 6972):
 * - Endpoint: POST /api/chat/
 * - Request:  {"query": "user's question"}
 * - Response: Plain text string
 */

export const API_BASE = 'http://localhost:6971'; // Main monitoring server
export const CHAT_ENDPOINT = `http://localhost:6972`; // ChatWebServer port
export const ENDPOINTS = {
  // POST /api/chat/ - Send query to master agent
  // Request:  {"query": "What is the database status?"}
  // Response: Plain text string from agent
  chat: `${CHAT_ENDPOINT}/api/chat/`,
};

export const SSE_ENDPOINTS = {
  monitor: `${API_BASE}/api/events/monitor`,
  storage: `${API_BASE}/api/events/storage`,
  database: `${API_BASE}/api/events/database`,
  activity: `${API_BASE}/api/events/activity`,
};
