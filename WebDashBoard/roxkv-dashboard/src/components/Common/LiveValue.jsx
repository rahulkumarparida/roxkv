import { useState, useEffect, memo } from 'react';
import sseService from '../../services/SSEService';

/**
 * A tiny wrapper that subscribes to an SSE channel and renders a specific value.
 * Using this prevents the parent component from re-rendering on every tick.
 */
const LiveValue = memo(function LiveValue({ channel, selector, fallback = 0 }) {
  const [value, setValue] = useState(fallback);

  useEffect(() => {
    return sseService.subscribe(channel, (data) => {
      const nextValue = selector(data);
      setValue(nextValue === undefined || nextValue === null || nextValue === '' ? fallback : nextValue);
    });
  }, [channel, selector, fallback]);

  return <>{value}</>;
});

export default LiveValue;
