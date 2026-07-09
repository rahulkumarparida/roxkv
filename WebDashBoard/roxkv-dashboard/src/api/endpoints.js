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

export const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:6971';
export const CHAT_ENDPOINT = import.meta.env.VITE_CHAT_API_BASE_URL || 'http://localhost:6972';
export const ENDPOINTS = {
  chat: '/api/chat/',
};

export const SSE_ENDPOINTS = {
  monitor: joinUrl(API_BASE, '/api/events/monitor'),
  storage: joinUrl(API_BASE, '/api/events/storage'),
  database: joinUrl(API_BASE, '/api/events/database'),
  activity: joinUrl(API_BASE, '/api/events/activity'),
};
