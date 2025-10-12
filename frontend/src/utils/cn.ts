/**
 * Utility function for class names
 * Combines class names using clsx logic
 */
export function cn(...classes: (string | undefined | null | false)[]): string {
  return classes.filter(Boolean).join(' ')
}