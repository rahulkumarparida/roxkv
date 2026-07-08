import { useState, useEffect, useContext } from 'react';
import { DashboardContext } from '../context/DashboardContext';
import sseService from '../services/SSEService';
import { MAX_ACTIVITY_ITEMS, MAX_GRAPH_DATA_POINTS } from '../utils/constants';

// --- Global UI State ---
export function useDashboard() {
  const context = useContext(DashboardContext);
  if (!context) throw new Error('useDashboard must be used within a DashboardProvider');
  return context;
}

// --- Fine-Grained Reactive Hooks ---

/** Hook for Database Health */
export function useDatabaseHealth() {
  const [health, setHealth] = useState({
    persistence: 'Valid', storageHealth: 'Optimal', snapshotCount: 0,
    databaseSizeBytes: 0, permanentKeys: 0, pubsubTopicsCount: 0, healthScore: 100,
    latestSnapshotTimestamp: null, totalSnapshotSizeBytes: 0, totalSnapshots: 0,
  });

  useEffect(() => {
    const unsubStorage = sseService.subscribe('storage', (data) => {
      setHealth(prev => ({
        ...prev,
        totalSnapshots: data.totalSnapshots || 0,
        snapshotCount: data.totalSnapshots || 0,
        totalSnapshotSizeBytes: data.totalSnapShotSize || 0,
        latestSnapshotTimestamp: data.latestSnapshot || data.lastCreated,
      }));
    });

    const unsubDatabase = sseService.subscribe('database', (data) => {
      const topicList = Array.isArray(data.pubsubTopics) ? data.pubsubTopics : [];
      const totalTopics = topicList.reduce((count, group) => {
        if (typeof group?.total === 'number') return count + group.total;
        if (Array.isArray(group?.topics)) return count + group.topics.length;
        return count;
      }, 0);

      setHealth(prev => ({
        ...prev,
        databaseSizeBytes: data.savedDataSize || 0,
        permanentKeys: data.ttlMetrics?.permanentKeys || 0,
        pubsubTopicsCount: totalTopics,
      }));
    });

    return () => { unsubStorage(); unsubDatabase(); };
  }, []);

  return health;
}

/** Hook for a specific metric's graph history and current value */
export function useMetric(key) {
  const [state, setState] = useState({ current: 0, history: [] });

  useEffect(() => {
    return sseService.subscribe('monitor', (data) => {
      const now = Date.now();
      let val = 0;

      if (typeof key === 'function') {
        val = key(data);
      } else {
        val = data?.[key];
      }

      if (typeof val === 'string') {
        const parsed = parseFloat(val);
        val = Number.isNaN(parsed) ? 0 : parsed;
      } else if (typeof val !== 'number') {
        val = 0;
      }
      
      setState(prev => {
        const newHistory = [...prev.history, { time: now, value: val }];
        if (newHistory.length > MAX_GRAPH_DATA_POINTS) {
          newHistory.shift();
        }
        return { current: val, history: newHistory };
      });
    });
  }, [key]);

  return state;
}

/** Hook for Activity Feed */
export function useActivityFeed() {
  const [activityFeed, setActivityFeed] = useState([]);

  useEffect(() => {
    return sseService.subscribe('activity', (events) => {
       if (Array.isArray(events)) {
         if (events.length > 0 && events[0] === "Not Data Found") return;

         const parsedEvents = events
         .filter((ev) => typeof ev === 'string' && ev.trim() !== '')
         .map((ev, i) => {
            const parts = ev.split(' ');
            if (parts.length >= 4) {
               return {
                 id: `${parts[0]}-${parts[1]}-${i}`,
                 type: parts[2].toUpperCase(),
                 description: parts.slice(3).join(' '),
                 timestamp: `${parts[0]}, ${parts[1]}`,
               };
            }
            return {
               id: `unknown-${i}-${Date.now()}`,
               type: 'INFO',
               description: ev,
               timestamp: new Date().toLocaleTimeString(),
            };
         });
         
         setActivityFeed(parsedEvents.slice(0, MAX_ACTIVITY_ITEMS));
       }
    });
  }, []);

  return activityFeed;
}
