<script lang="ts">
  import { tick } from 'svelte';
  import { currentSegment, splitTimesMs } from '../stores/timer';
  import { currentAttempts, deltas } from '../stores/splits';
  import SplitRow from './SplitRow.svelte';

  let container: HTMLDivElement | undefined = $state();

  const segments = $derived($currentAttempts?.segments || []);
  const currentDeltas = $derived($deltas || []);

  $effect(() => {
    const idx = $currentSegment;
    if (idx >= 0 && container) {
      tick().then(() => {
        const row = container!.children[idx] as HTMLElement | undefined;
        if (row) row.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
      });
    }
  });
</script>

<div class="splits-list" bind:this={container}>
  {#each segments as segment, i}
    <SplitRow
      name={segment.name}
      splitTimeMs={$splitTimesMs[i] || 0}
      personalBestMs={segment.personalBestMs}
      delta={i < currentDeltas.length ? currentDeltas[i] : null}
      isActive={i === $currentSegment}
      isCompleted={i < $currentSegment}
    />
  {/each}
</div>

<style>
  .splits-list {
    flex: 1;
    overflow-y: auto;
    border-top: 1px solid var(--border);
    border-bottom: 1px solid var(--border);
  }
</style>
