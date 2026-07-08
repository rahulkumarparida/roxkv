/**
 * @fileoverview Centralized endpoint configuration.
 * This is the ONLY place where URLs are defined.
 * When swapping to the real Go backend, only this file needs updating.
 */

export const API_BASE = 'http://localhost:6971'; // Reverted to correct port
export const CHAT_ENDPOINT = `http://localhost:6972`;
export const ENDPOINTS = {
  chat: `${CHAT_ENDPOINT}/api/chat`,
};

export const SSE_ENDPOINTS = {
  monitor: `${API_BASE}/api/events/monitor`,
  storage: `${API_BASE}/api/events/storage`,
  database: `${API_BASE}/api/events/database`,
  activity: `${API_BASE}/api/events/activity`,
};
