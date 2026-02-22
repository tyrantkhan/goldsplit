<script lang="ts">
  import { timerState, interpolatedElapsed, elapsedMs, currentSegment, splitTimesMs } from '../stores/timer';
  import { currentAttempts, deltas } from '../stores/splits';
  import { formatRunTime, formatSegDelta } from '../utils/format';

  // Sum of Best: sum of all bestSegmentMs (null if any segment lacks a best)
  const sumOfBest = $derived.by(() => {
    const segs = $currentAttempts?.segments;
    if (!segs) return null;
    let sum = 0;
    for (const seg of segs) {
      if (!seg.bestSegmentMs) return null;
      sum += seg.bestSegmentMs;
    }
    return sum;
  });

  // Personal Best: last segment's personalBestMs (cumulative)
  const personalBest = $derived.by(() => {
    const segs = $currentAttempts?.segments;
    if (!segs || segs.length === 0) return null;
    const pb = segs[segs.length - 1].personalBestMs;
    return pb > 0 ? pb : null;
  });

  // Previous segment delta: last entry in deltas array
  // Previous segment delta: segment-level delta (not cumulative)
  const prevSegDelta = $derived.by(() => {
    const d = $deltas;
    if (!d || d.length === 0) return null;
    const last = d[d.length - 1];
    if (!last || last.skipped) return null;

    // Segment delta = this cumulative delta minus previous cumulative delta
    let segDeltaMs = last.deltaMs;
    if (d.length >= 2) {
      const prev = d[d.length - 2];
      if (prev && !prev.skipped) {
        segDeltaMs = last.deltaMs - prev.deltaMs;
      }
    }

    return { ...last, segDeltaMs };
  });

  const prevDeltaText = $derived(
    prevSegDelta && prevSegDelta.segDeltaMs !== 0
      ? formatSegDelta(prevSegDelta.segDeltaMs)
      : null
  );

  const prevDeltaColor = $derived(
    prevSegDelta
      ? prevSegDelta.isBestEver
        ? 'var(--best-time)'
        : prevSegDelta.segDeltaMs < 0
          ? prevSegDelta.gainedTime ? 'var(--ahead-gaining)' : 'var(--ahead-losing)'
          : prevSegDelta.segDeltaMs > 0
            ? prevSegDelta.gainedTime ? 'var(--behind-gaining)' : 'var(--behind-losing)'
            : undefined
      : undefined
  );

  // Predicted Time: lastSplit + max(timeInCurrentSeg, currentSegBest) + remaining bests
  const predictedTime = $derived.by(() => {
    const segs = $currentAttempts?.segments;
    if (!segs || $timerState === 'idle' || $timerState === 'finished') return null;

    const segIdx = $currentSegment;
    const splits = $splitTimesMs;
    const elapsed = $timerState === 'running' ? $interpolatedElapsed : $elapsedMs;

    // Need best segment data for current and remaining segments
    if (!segs[segIdx]?.bestSegmentMs) return null;

    const lastSplitTime = splits.length > 0 ? splits[splits.length - 1] : 0;
    const timeInCurrentSeg = elapsed - lastSplitTime;
    const currentSegBest = segs[segIdx].bestSegmentMs;

    let predicted = lastSplitTime + Math.max(timeInCurrentSeg, currentSegBest);

    // Add remaining segment bests
    for (let i = segIdx + 1; i < segs.length; i++) {
      if (!segs[i].bestSegmentMs) return null;
      predicted += segs[i].bestSegmentMs;
    }

    return predicted;
  });

  const isIdle = $derived($timerState === 'idle');
  const isFinished = $derived($timerState === 'finished');

  const showPredicted = $derived(!isIdle && !isFinished && predictedTime !== null);
  const showSumOfBest = $derived(!isIdle && sumOfBest !== null);
  const showPersonalBest = $derived(!isIdle && personalBest !== null);
  const showPrevDelta = $derived(!isIdle && prevSegDelta !== null);

  const hasAnyStats = $derived(showPredicted || showSumOfBest || showPersonalBest || showPrevDelta);
</script>

{#if hasAnyStats}
  <div class="timer-stats">
    {#if showPredicted}
      <span class="label">Predicted</span>
      <span class="value">{formatRunTime(Math.floor(predictedTime!))}</span>
    {/if}
    {#if showSumOfBest}
      <span class="label">Sum of Best</span>
      <span class="value">{formatRunTime(sumOfBest!)}</span>
    {/if}
    {#if showPersonalBest}
      <span class="label">Personal Best</span>
      <span class="value">{formatRunTime(personalBest!)}</span>
    {/if}
    {#if showPrevDelta}
      <span class="label">Prev Segment</span>
      <span class="value" style:color={prevDeltaColor}>
        {prevDeltaText ?? '-'}
      </span>
    {/if}
  </div>
{/if}

<style>
  .timer-stats {
    display: grid;
    grid-template-columns: auto auto;
    gap: 2px 12px;
    justify-content: center;
    padding: 0 16px 4px;
    font-family: var(--timer-font);
    font-size: 12px;
    font-variant-numeric: tabular-nums;
  }

  .label {
    text-align: right;
    color: var(--text-muted);
  }

  .value {
    text-align: left;
    color: var(--text-secondary);
  }
</style>
