<script lang="ts">
  import { interpolatedElapsed, timerState, elapsedMs, splitTimesMs } from '../stores/timer';
  import { deltas, currentAttempts } from '../stores/splits';
  import { formatTime } from '../utils/format';

  const displayMs = $derived($timerState === 'running' ? $interpolatedElapsed : $elapsedMs);
  const displayTime = $derived(formatTime(Math.floor(displayMs)));
  const stateClass = $derived($timerState);

  // On PB pace: last completed split's cumulative <= PB cumulative at that split.
  const onPBPace = $derived.by(() => {
    const splits = $splitTimesMs;
    const segs = $currentAttempts?.segments;
    if (!splits || splits.length === 0 || !segs) return false;

    for (let i = splits.length - 1; i >= 0; i--) {
      if (splits[i] === 0) continue;
      const pb = segs[i]?.personalBestMs;
      if (!pb || pb === 0) continue;
      return splits[i] <= pb;
    }

    return false;
  });

  const timerColor = $derived.by(() => {
    const state = $timerState;
    if (state === 'idle') return undefined;
    if (state === 'paused') return 'var(--text-secondary)';

    // Gold when on PB pace (neutral-or-better from here means a new PB).
    if (onPBPace) return 'var(--best-time)';

    const d = $deltas;
    if (!d || d.length === 0) return undefined;

    const last = d[d.length - 1];
    if (!last || last.skipped || last.deltaMs === 0) return undefined;

    if (last.isAhead) {
      return last.gainedTime ? 'var(--ahead-gaining)' : 'var(--ahead-losing)';
    }
    return last.gainedTime ? 'var(--behind-gaining)' : 'var(--behind-losing)';
  });
</script>

<div class="timer-display {stateClass}">
  <span class="time" style:color={timerColor}>{displayTime}</span>
</div>

<style>
  .timer-display {
    padding: 16px 0 8px;
    text-align: center;
  }

  .time {
    font-family: var(--timer-font);
    font-size: var(--timer-size);
    font-weight: 700;
    font-variant-numeric: tabular-nums;
    letter-spacing: -1px;
    color: var(--text-primary);
  }
</style>
