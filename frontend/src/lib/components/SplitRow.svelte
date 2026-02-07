<script lang="ts">
  import type { Delta } from '../types';
  import { formatSplitTime } from '../utils/format';
  import DeltaDisplay from './DeltaDisplay.svelte';

  let {
    name,
    splitTimeMs = 0,
    personalBestMs = 0,
    delta = null,
    isActive = false,
    isCompleted = false,
  }: {
    name: string;
    splitTimeMs?: number;
    personalBestMs?: number;
    delta?: Delta | null;
    isActive?: boolean;
    isCompleted?: boolean;
  } = $props();
</script>

<div class="split-row" class:active={isActive} class:completed={isCompleted}>
  <span class="name">{name}</span>
  <span class="delta-col">
    {#if isCompleted}
      <DeltaDisplay {delta} />
    {/if}
  </span>
  <span class="time">
    {#if isCompleted}
      {formatSplitTime(splitTimeMs)}
    {:else if personalBestMs > 0}
      <span class="pb-time">{formatSplitTime(personalBestMs)}</span>
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
</style>
