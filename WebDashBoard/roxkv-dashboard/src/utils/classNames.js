import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

/**
 * Merges class names using clsx for conditional logic
 * and tailwind-merge for deduplication of Tailwind classes.
 * @param {...(string|Object|Array)} inputs - Class name inputs
 * @returns {string} Merged class string
 */
export function cn(...inputs) {
  return twMerge(clsx(inputs));
}
