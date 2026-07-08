import { createContext, useState, useCallback } from 'react';

export const DashboardContext = createContext(null);

/**
 * DashboardProvider now ONLY aggregates slow-moving/static state
 * (like UI toggles). High-frequency SSE data is now consumed
 * directly by components using specialized hooks.
 */
export function DashboardProvider({ children }) {
  const [autoRefresh, setAutoRefresh] = useState(true);

  const toggleAutoRefresh = useCallback(() => {
    setAutoRefresh((prev) => !prev);
  }, []);

  return (
    <DashboardContext.Provider value={{ autoRefresh, toggleAutoRefresh }}>
      {children}
    </DashboardContext.Provider>
  );
}
