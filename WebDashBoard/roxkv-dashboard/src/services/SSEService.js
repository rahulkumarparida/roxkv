/**
 * @fileoverview Centralized SSE connection manager.
 * Manages connections and implements reconnection with exponential backoff.
 */

import { createSSEConnection } from '../api/sse.js';
import { SSE_ENDPOINTS } from '../api/endpoints.js';

class SSEService {
  constructor() {
    /** @type {Map<string, {status: string, cleanup: function|null, retryCount: number, retryTimeout: number|null}>} */
    this.connections = new Map();
    /** @type {Map<string, Set<function>>} */
    this.subscribers = new Map();
    /** @type {Set<function>} */
    this.statusSubscribers = new Set();
  }

  /**
   * Subscribe to a specific data channel.
   * @param {string} channel - 'monitor' | 'storage' | 'database' | 'activity'
   * @param {function(Object): void} callback
   * @returns {function(): void} Unsubscribe function
   */
  subscribe(channel, callback) {
    if (!this.subscribers.has(channel)) this.subscribers.set(channel, new Set());
    this.subscribers.get(channel).add(callback);
    return () => {
      const subs = this.subscribers.get(channel);
      if (subs) subs.delete(callback);
    };
  }

  /**
   * Subscribe to connection status changes.
   * @param {function({channel: string, status: string, error?: Error|null}): void} callback
   * @returns {function(): void}
   */
  subscribeStatus(callback) {
    this.statusSubscribers.add(callback);
    return () => this.statusSubscribers.delete(callback);
  }

  /**
   * Notify status subscribers.
   * @param {string} channel
   * @param {string} status
   * @param {Error|null} [error]
   */
  _notifyStatus(channel, status, error = null) {
    this.statusSubscribers.forEach((cb) => {
      try {
        cb({ channel, status, error });
      } catch (err) {
        console.error(`[SSEService] Status subscriber error on ${channel}:`, err);
      }
    });
  }

  /**
   * Notify all subscribers on a channel.
   * @param {string} channel
   * @param {Object} data
   */
  _notify(channel, data) {
    const subs = this.subscribers.get(channel);

    if (subs) {
      subs.forEach((cb) => {
        try { cb(data); } catch (err) { console.error(`[SSEService] Subscriber error on ${channel}:`, err); }
      });
    }
  }

  /**
   * Connect to a specific SSE channel.
   * @param {string} channel
   */
  connect(channel) {
    const url = SSE_ENDPOINTS[channel];
    if (!url) {
      console.warn(`[SSEService] Unknown channel: ${channel}`);
      return;
    }

    const existing = this.connections.get(channel);
    if (existing?.status === 'connected' || existing?.status === 'connecting') return;

    if (existing?.retryTimeout) clearTimeout(existing.retryTimeout);

    this.connections.set(channel, {
      status: 'connecting',
      cleanup: null,
      retryCount: existing?.retryCount ?? 0,
      retryTimeout: null,
    });
    this._notifyStatus(channel, 'connecting');

    const cleanup = createSSEConnection(
      url,
      (data) => this._notify(channel, data),
      (error) => this._handleError(channel, error),
      () => {
        const current = this.connections.get(channel);
        if (!current) return;
        current.status = 'connected';
        current.retryCount = 0;
        this._notifyStatus(channel, 'connected');
      }
    );

    const conn = this.connections.get(channel);

    if (conn) {
      conn.cleanup = cleanup;
    }
  }

  /**
   * Handle connection errors with exponential backoff.
   * 1s → 2s → 4s → 8s → max 30s
   * @param {string} channel
   */
  _handleError(channel, error = null) {
    const conn = this.connections.get(channel);
    if (!conn) return;
    conn.status = 'error';
    
    if (conn.cleanup) conn.cleanup();
    if (conn.retryTimeout) clearTimeout(conn.retryTimeout);
    this._notifyStatus(channel, 'error', error);
    const delay = Math.min(1000 * Math.pow(2, conn.retryCount), 30000);
    conn.retryCount++;
    console.log(`[SSEService] Reconnecting ${channel} in ${delay}ms (attempt ${conn.retryCount})`);
    conn.retryTimeout = setTimeout(() => this.connect(channel), delay);
  }

  /** Disconnect a specific channel. */
  disconnect(channel) {
    const conn = this.connections.get(channel);
    if (conn) {
      if (conn.cleanup) conn.cleanup();
      if (conn.retryTimeout) clearTimeout(conn.retryTimeout);
      this.connections.delete(channel);
      this._notifyStatus(channel, 'disconnected');
    }
  }

  /** Disconnect all channels. */
  disconnectAll() {
    this.connections.forEach((_, channel) => this.disconnect(channel));
  }

  /** Get the connection status for a channel. */
  getStatus(channel) {
    const conn = this.connections.get(channel);
    return conn ? conn.status : 'disconnected';
  }
}

// Singleton instance
const sseService = new SSEService();
export default sseService;
