/**
 * @fileoverview Centralized endpoint configuration.
 * This is the ONLY place where URLs are defined.
 * 
 * ChatWebServer (Port 6972):
 * - Endpoint: POST /api/chat/
 * - Request:  {"query": "user's question"}
 * - Response: Plain text string
 */

function joinUrl(base, path) {
  const normalizedBase = (base || '').replace(/\/+$/, '');
  return normalizedBase ? `${normalizedBase}${path}` : path;
}

export const API_BASE = import.meta.env.VITE_API_BASE || ''; // Main monitoring server
export const CHAT_ENDPOINT = import.meta.env.VITE_CHAT_BASE || 'localhost:6972'; // ChatWebServer port
export const ENDPOINTS = {
  // POST /api/chat/ - Send query to master agent
  // Request:  {"query": "What is the database status?"}
  // Response: Plain text string from agent
  chat: joinUrl(CHAT_ENDPOINT, '/api/chat/'),
};

export const SSE_ENDPOINTS = {
  monitor: joinUrl(API_BASE, '/api/events/monitor'),
  storage: joinUrl(API_BASE, '/api/events/storage'),
  database: joinUrl(API_BASE, '/api/events/database'),
  activity: joinUrl(API_BASE, '/api/events/activity'),
};
