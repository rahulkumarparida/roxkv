import { createContext, useEffect, useState, useCallback, useRef } from 'react';
import sseService from '../services/SSEService';

export const SSEContext = createContext(null);
const CHANNELS = ['monitor', 'storage', 'database', 'activity'];

function getOverallStatus(statusByChannel) {
  const statuses = Object.values(statusByChannel);

  if (statuses.length === 0) return 'connecting';
  if (statuses.some((status) => status === 'connected')) return 'connected';
  if (statuses.some((status) => status === 'connecting')) return 'connecting';
  if (statuses.some((status) => status === 'error')) return 'error';
  return 'disconnected';
}

/**
 * SSEProvider manages SSE lifecycle and provides connection status.
 */
export function SSEProvider({ children }) {
  const [statusByChannel, setStatusByChannel] = useState(() =>
    Object.fromEntries(CHANNELS.map((channel) => [channel, 'connecting']))
  );
  const [error, setError] = useState(null);
  const initialized = useRef(false);

  useEffect(() => {
    if (initialized.current) return;
    initialized.current = true;
    console.log("SSEProvider mounted");

    const unsubscribeStatus = sseService.subscribeStatus(({ channel, status, error: nextError }) => {
      setStatusByChannel((prev) => ({ ...prev, [channel]: status }));
      if (nextError) {
        setError(nextError.message);
      } else if (status === 'connected') {
        setError(null);
      }
    });

    CHANNELS.forEach((channel) => sseService.connect(channel));

    return () => {
      console.log("SSEProvider unmounted");
      unsubscribeStatus();
      sseService.disconnectAll();
    };
  }, []);

  const reconnect = useCallback(() => {
    setError(null);
    setStatusByChannel(Object.fromEntries(CHANNELS.map((channel) => [channel, 'connecting'])));
    sseService.disconnectAll();
    CHANNELS.forEach((channel) => sseService.connect(channel));
  }, []);

  const connectionStatus = getOverallStatus(statusByChannel);

  return (
    <SSEContext.Provider value={{ connectionStatus, error, reconnect, sseService, statusByChannel }}>
      {children}
    </SSEContext.Provider>
  );
}
