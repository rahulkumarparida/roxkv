/**
 * Formats a byte value into a human-readable string.
 * @param {number} bytes - The number of bytes
 * @param {number} [decimals=1] - Number of decimal places
 * @returns {string} Formatted string (e.g., "52.9 KB")
 */
export function formatBytes(bytes, decimals = 1) {
  if (bytes === 0) return '0 B';
  if (!bytes || isNaN(bytes)) return '\u2014';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(Math.abs(bytes)) / Math.log(k));
  const value = bytes / Math.pow(k, i);
  return `${value.toFixed(decimals)} ${sizes[i]}`;
}
