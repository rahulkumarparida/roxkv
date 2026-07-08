/**
 * @fileoverview Central constants for the RoxAI Dashboard.
 * All endpoint URLs, event types, and configuration values live here.
 */

/** Base URL for the Go backend API */
export const API_BASE = 'http://localhost:6971';

/** Chat endpoint for LLM communication */
export const CHAT_ENDPOINT = `${API_BASE}/api/chat`;

/** Server Sent Events endpoint configuration */
export const SSE_ENDPOINTS = {
  monitor: `${API_BASE}/sse/monitor`,
  storage: `${API_BASE}/sse/storage`,
  database: `${API_BASE}/sse/database`,
  activity: `${API_BASE}/sse/activity`,
};

/** Activity event type definitions with display properties */
export const EVENT_TYPES = {
  USER_CONNECTED: { label: 'User connected', color: 'text-green-400', dotColor: 'bg-green-400' },
  CLIENT_TIMEOUT: { label: 'Client inactive timeout', color: 'text-amber-400', dotColor: 'bg-amber-400', badge: 'WARNING', badgeColor: 'text-amber-400' },
  SET_KEY: { label: 'SET', color: 'text-purple-400', dotColor: 'bg-purple-400' },
  GET_KEY: { label: 'GET', color: 'text-blue-400', dotColor: 'bg-blue-400' },
  DELETE_KEY: { label: 'DELETE', color: 'text-red-400', dotColor: 'bg-red-400' },
  SNAPSHOT_CREATED: { label: 'SNAPSHOT created', color: 'text-purple-300', dotColor: 'bg-purple-300' },
  AI_QUERY: { label: 'AI Query', color: 'text-white', dotColor: 'bg-white' },
};

/** Maximum number of activity feed items to retain in state */
export const MAX_ACTIVITY_ITEMS = 50;

/** Maximum number of data points per metric graph */
export const MAX_GRAPH_DATA_POINTS = 30;

/** Mock data update interval ranges (ms) */
export const MOCK_INTERVALS = {
  metrics: 2000,
  activity: { min: 2000, max: 5000 },
  systemOverview: 1000,
  database: 5000,
};

/** Quick action commands for the chat */
export const QUICK_ACTIONS = [
  { label: 'Show TTL stats', message: 'Show me the TTL stats?' },
  { label: 'List all keys', message: 'List all keys in the database' },
  { label: 'Database size', message: 'What is the current database size?' },
  { label: 'Snapshot details', message: 'Show me the snapshot details' },
  { label: 'Active connections', message: 'How many active connections are there?' },
];

/** Application metadata */
export const APP_META = {
  version: 'v2.1.0',
  name: 'ROXX KV',
  tagline: 'In-Memory Database with AI Agent',
  connectionHost: 'localhost:6971',
};
