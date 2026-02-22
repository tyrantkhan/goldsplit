<script lang="ts">
  import type { Delta } from '../types';
  import { formatSplitTime } from '../utils/format';
  import DeltaDisplay from './DeltaDisplay.svelte';

  let {
    name,
    splitTimeMs = 0,
    comparisonSplitMs = 0,
    personalBestMs = 0,
    bestSegmentMs = 0,
    delta = null,
    isActive = false,
    isCompleted = false,
  }: {
    name: string;
    splitTimeMs?: number;
    comparisonSplitMs?: number;
    personalBestMs?: number;
    bestSegmentMs?: number;
    delta?: Delta | null;
    isActive?: boolean;
    isCompleted?: boolean;
  } = $props();

  const referenceSplitMs = $derived(comparisonSplitMs > 0 ? comparisonSplitMs : personalBestMs);
  const isFirstTime = $derived(isCompleted && bestSegmentMs === 0 && delta != null && !delta.skipped);
</script>

<div class="split-row" class:active={isActive} class:completed={isCompleted}>
  <span class="name">{name}</span>
  <span class="delta-col">
    {#if isFirstTime}
      <span class="new-segment">New</span>
    {:else if isCompleted}
      <DeltaDisplay {delta} />
    {/if}
  </span>
  <span class="time">
    {#if isCompleted}
      {formatSplitTime(splitTimeMs)}
    {:else if referenceSplitMs > 0}
      <span class="pb-time">{formatSplitTime(referenceSplitMs)}</span>
    {:else}
      <span class="no-time">-</span>
    {/if}
  </span>
</div>

<style>
  .split-row {
    display: flex;
    align-items: center;
    height: var(--split-row-height);
    padding: 0 12px;
    font-size: 13px;
    border-bottom: 1px solid var(--border);
  }

  .split-row.active {
    background: var(--bg-tertiary);
  }

  .split-row.completed {
    opacity: 0.85;
  }

  .name {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .delta-col {
    width: 75px;
    text-align: right;
    padding-right: 8px;
  }

  .time {
    width: 85px;
    text-align: right;
    font-family: var(--timer-font);
    font-size: 12px;
    font-variant-numeric: tabular-nums;
  }

  .pb-time {
    color: var(--text-muted);
  }

  .no-time {
    color: var(--text-muted);
  }

  .new-segment {
    font-family: var(--timer-font);
    font-size: 12px;
    color: var(--best-time);
  }
</style>
