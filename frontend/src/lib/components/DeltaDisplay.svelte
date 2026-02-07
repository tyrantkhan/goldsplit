<script lang="ts">
  import type { Delta } from '../types';
  import { formatDelta } from '../utils/format';

  let { delta = null }: { delta?: Delta | null } = $props();

  const displayText = $derived(
    delta && !delta.skipped && delta.deltaMs !== 0
      ? formatDelta(delta.deltaMs)
      : ''
  );

  const colorVar = $derived(
    delta
      ? delta.isBestEver
        ? '--best-segment'
        : delta.isAhead
          ? delta.gainedTime ? '--ahead-gaining' : '--ahead-losing'
          : delta.deltaMs > 0
            ? delta.gainedTime ? '--behind-gaining' : '--behind-losing'
            : ''
      : ''
  );
</script>

{#if displayText}
  <span class="delta" style:color={colorVar ? `var(${colorVar})` : undefined}>{displayText}</span>
{:else}
  <span class="delta empty">-</span>
{/if}

<style>
  .delta {
    font-family: var(--timer-font);
    font-size: 12px;
    font-variant-numeric: tabular-nums;
  }

  .delta.empty {
    color: var(--text-muted);
  }
</style>
