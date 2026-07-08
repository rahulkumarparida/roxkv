import { useEffect, useRef, useState } from 'react';
import { cn } from '../../utils/classNames';

/**
 * Smoothly animates between numeric values using requestAnimationFrame.
 * @param {Object} props
 * @param {number} props.value - Target value
 * @param {number} [props.duration=500] - Animation duration in ms
 * @param {number} [props.decimals=0] - Decimal places
 * @param {string} [props.className] - Additional classes
 * @param {string} [props.suffix] - Suffix text (e.g., '%', 'MB/s')
 */
export default function AnimatedNumber({ value, duration = 500, decimals = 0, className, suffix = '' }) {
  const [displayValue, setDisplayValue] = useState(value);
  const previousValue = useRef(value);
  const rafRef = useRef(null);

  useEffect(() => {
    const from = previousValue.current;
    const to = value;
    const startTime = performance.now();

    const animate = (currentTime) => {
      const elapsed = currentTime - startTime;
      const progress = Math.min(elapsed / duration, 1);
      // Ease out cubic
      const eased = 1 - Math.pow(1 - progress, 3);
      const current = from + (to - from) * eased;
      setDisplayValue(current);
      if (progress < 1) {
        rafRef.current = requestAnimationFrame(animate);
      } else {
        previousValue.current = to;
      }
    };

    rafRef.current = requestAnimationFrame(animate);
    return () => { if (rafRef.current) cancelAnimationFrame(rafRef.current); };
  }, [value, duration]);

  return (
    <span className={cn('font-mono tabular-nums', className)}>
      {displayValue.toFixed(decimals)}{suffix}
    </span>
  );
}
