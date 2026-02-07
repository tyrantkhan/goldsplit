export function formatTime(ms: number): string {
  if (ms < 0) ms = 0;

  const totalSeconds = Math.floor(ms / 1000);
  const milliseconds = ms % 1000;
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;

  const msStr = String(milliseconds).padStart(3, '0');
  const secStr = String(seconds).padStart(2, '0');

  if (hours > 0) {
    const minStr = String(minutes).padStart(2, '0');
    return `${hours}:${minStr}:${secStr}.${msStr}`;
  }

  if (minutes > 0) {
    return `${minutes}:${secStr}.${msStr}`;
  }

  return `${seconds}.${msStr}`;
}

export function formatDelta(ms: number): string {
  const sign = ms < 0 ? '-' : '+';
  const abs = Math.abs(ms);

  const totalSeconds = Math.floor(abs / 1000);
  const milliseconds = Math.floor((abs % 1000) / 10); // Show centiseconds
  const minutes = Math.floor(totalSeconds / 60);
  const seconds = totalSeconds % 60;

  const csStr = String(milliseconds).padStart(2, '0');

  if (minutes > 0) {
    const secStr = String(seconds).padStart(2, '0');
    return `${sign}${minutes}:${secStr}.${csStr}`;
  }

  return `${sign}${seconds}.${csStr}`;
}

export function formatSplitTime(ms: number): string {
  if (ms === 0) return '-';
  return formatTime(ms);
}

/** Always shows at least M:SS.mmm — suitable for total run times. */
export function formatRunTime(ms: number): string {
  if (ms <= 0) return '-';

  const totalSeconds = Math.floor(ms / 1000);
  const milliseconds = ms % 1000;
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;

  const msStr = String(milliseconds).padStart(3, '0');
  const secStr = String(seconds).padStart(2, '0');

  if (hours > 0) {
    const minStr = String(minutes).padStart(2, '0');
    return `${hours}:${minStr}:${secStr}.${msStr}`;
  }

  const minStr = String(minutes);
  return `${minStr}:${secStr}.${msStr}`;
}

/**
 * Parses a human-readable time string back to milliseconds.
 * Supports: H:MM:SS.mmm, M:SS.mmm, SS.mmm, S.mmm
 * Returns null for invalid input, 0 for empty or "-".
 */
export function parseTime(input: string): number | null {
  const s = input.trim();
  if (s === '' || s === '-') return 0;

  // Split on colons to determine format
  const parts = s.split(':');
  if (parts.length > 3) return null;

  let hours = 0;
  let minutes = 0;
  let secPart: string;

  if (parts.length === 3) {
    // H:MM:SS.mmm
    hours = parseInt(parts[0], 10);
    minutes = parseInt(parts[1], 10);
    secPart = parts[2];
  } else if (parts.length === 2) {
    // M:SS.mmm
    minutes = parseInt(parts[0], 10);
    secPart = parts[1];
  } else {
    // SS.mmm or S.mmm
    secPart = parts[0];
  }

  if (isNaN(hours) || isNaN(minutes) || hours < 0 || minutes < 0) return null;

  // Parse seconds.milliseconds
  const dotParts = secPart.split('.');
  if (dotParts.length > 2) return null;

  const seconds = parseInt(dotParts[0], 10);
  if (isNaN(seconds) || seconds < 0) return null;

  let ms = 0;
  if (dotParts.length === 2) {
    const frac = dotParts[1];
    if (frac.length === 0 || frac.length > 3) return null;
    // Pad to 3 digits: "4" -> "400", "45" -> "450", "456" -> "456"
    ms = parseInt(frac.padEnd(3, '0'), 10);
    if (isNaN(ms)) return null;
  }

  const totalMs = ((hours * 3600 + minutes * 60 + seconds) * 1000) + ms;
  return totalMs;
}
