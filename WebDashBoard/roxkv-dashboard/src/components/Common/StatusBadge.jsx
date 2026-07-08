import { cn } from '../../utils/classNames';

const STATUS_STYLES = {
  active: { text: 'text-green-400', dot: 'bg-green-400', label: 'Active' },
  valid: { text: 'text-green-400', dot: 'bg-green-400', label: 'Valid' },
  optimal: { text: 'text-green-400', dot: 'bg-green-400', label: 'Optimal' },
  online: { text: 'text-green-400', dot: 'bg-green-400', label: 'Online' },
  connected: { text: 'text-green-400', dot: 'bg-green-400', label: 'Connected' },
  warning: { text: 'text-amber-400', dot: 'bg-amber-400', label: 'Warning' },
  degraded: { text: 'text-amber-400', dot: 'bg-amber-400', label: 'Degraded' },
  error: { text: 'text-red-400', dot: 'bg-red-400', label: 'Error' },
  offline: { text: 'text-red-400', dot: 'bg-red-400', label: 'Offline' },
  invalid: { text: 'text-red-400', dot: 'bg-red-400', label: 'Invalid' },
  live: { text: 'text-green-400', dot: 'bg-green-400', label: 'LIVE' },
};

/**
 * Status badge with optional animated pulsing dot indicator.
 * @param {Object} props
 * @param {string} props.status - Status type key
 * @param {string} [props.label] - Override default label
 * @param {boolean} [props.dot=true] - Show animated dot
 * @param {string} [props.className] - Additional classes
 */
export default function StatusBadge({ status, label, dot = true, className }) {
  const style = STATUS_STYLES[status?.toLowerCase()] || STATUS_STYLES.active;
  const displayLabel = label || style.label;
  return (
    <span className={cn('inline-flex items-center gap-1.5', className)}>
      {dot && (
        <span className="relative flex h-2 w-2">
          <span className={cn('absolute inline-flex h-full w-full rounded-full opacity-75 animate-ping', style.dot)} />
          <span className={cn('relative inline-flex h-2 w-2 rounded-full', style.dot)} />
        </span>
      )}
      <span className={cn('text-xs font-medium', style.text)}>{displayLabel}</span>
    </span>
  );
}
