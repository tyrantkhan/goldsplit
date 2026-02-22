<script lang="ts">
  import { Tooltip } from 'bits-ui';
  import { timerState } from '../stores/timer';
  import { settings } from '../stores/settings';
  import { deltas } from '../stores/splits';
  import { StartSplit, TogglePause, Reset, UndoSplit, SkipSplit, DiscardAttempt, GetDeltas } from '../../../wailsjs/go/main/App';

  async function fetchDeltas() {
    const d = await GetDeltas();
    if (d) deltas.set(d);
  }

  const primaryLabel = $derived(
    $timerState === 'idle' ? 'Start' :
    $timerState === 'running' ? 'Split' :
    $timerState === 'paused' ? 'Resume' :
    'Done'
  );

  const showPrimary = $derived($timerState !== 'finished');
  const showPause = $derived($timerState === 'running');
  const showReset = $derived($timerState === 'running' || $timerState === 'paused');
  const showUndo = $derived($timerState === 'running');
  const showSkip = $derived($timerState === 'running');
  const showFinished = $derived($timerState === 'finished');

  function displayKey(code: string): string {
    if (code.startsWith('Key')) return code.slice(3);
    if (code.startsWith('Digit')) return code.slice(5);
    if (code === 'ArrowLeft') return '\u2190';
    if (code === 'ArrowRight') return '\u2192';
    if (code === 'ArrowUp') return '\u2191';
    if (code === 'ArrowDown') return '\u2193';
    return code;
  }

  async function handlePrimary() {
    if ($timerState === 'idle' || $timerState === 'running') {
      await StartSplit();
      if (($timerState as string) !== 'finished') {
        fetchDeltas();
      }
    } else if ($timerState === 'paused') {
      TogglePause();
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) return;

    const code = e.code;
    const hk = $settings.hotkeys;

    if (code === hk.startSplit) {
      e.preventDefault();
      handlePrimary();
    } else if (code === hk.pause) {
      if ($timerState === 'running' || $timerState === 'paused') {
        TogglePause();
      }
    } else if (code === hk.reset) {
      if ($timerState !== 'idle') {
        Reset();
      }
    } else if (code === hk.undoSplit) {
      if ($timerState === 'running') {
        UndoSplit().then(() => fetchDeltas());
      }
    } else if (code === hk.skipSplit) {
      if ($timerState === 'running') {
        SkipSplit().then(() => fetchDeltas());
      }
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

<Tooltip.Provider delayDuration={300}>
  <div class="controls">
    {#if showPrimary}
      <Tooltip.Root>
        <Tooltip.Trigger>
          {#snippet child({ props })}
            <button {...props} class="btn primary" onclick={handlePrimary}>
              {primaryLabel}
            </button>
          {/snippet}
        </Tooltip.Trigger>
        <Tooltip.Content class="tooltip" sideOffset={6} side="top">
          {displayKey($settings.hotkeys.startSplit)}
        </Tooltip.Content>
      </Tooltip.Root>
    {/if}

    {#if showPause}
      <Tooltip.Root>
        <Tooltip.Trigger>
          {#snippet child({ props })}
            <button {...props} class="btn" onclick={() => TogglePause()}>
              Pause
            </button>
          {/snippet}
        </Tooltip.Trigger>
        <Tooltip.Content class="tooltip" sideOffset={6} side="top">
          {displayKey($settings.hotkeys.pause)}
        </Tooltip.Content>
      </Tooltip.Root>
    {/if}

    {#if showUndo}
      <Tooltip.Root>
        <Tooltip.Trigger>
          {#snippet child({ props })}
            <button {...props} class="btn small" onclick={() => UndoSplit()}>
              Undo
            </button>
          {/snippet}
        </Tooltip.Trigger>
        <Tooltip.Content class="tooltip" sideOffset={6} side="top">
          {displayKey($settings.hotkeys.undoSplit)}
        </Tooltip.Content>
      </Tooltip.Root>
    {/if}

    {#if showSkip}
      <Tooltip.Root>
        <Tooltip.Trigger>
          {#snippet child({ props })}
            <button {...props} class="btn small" onclick={() => SkipSplit()}>
              Skip
            </button>
          {/snippet}
        </Tooltip.Trigger>
        <Tooltip.Content class="tooltip" sideOffset={6} side="top">
          {displayKey($settings.hotkeys.skipSplit)}
        </Tooltip.Content>
      </Tooltip.Root>
    {/if}

    {#if showReset}
      <Tooltip.Root>
        <Tooltip.Trigger>
          {#snippet child({ props })}
            <button {...props} class="btn danger" onclick={() => Reset()}>
              Fail
            </button>
          {/snippet}
        </Tooltip.Trigger>
        <Tooltip.Content class="tooltip" sideOffset={6} side="top">
          {displayKey($settings.hotkeys.reset)}
        </Tooltip.Content>
      </Tooltip.Root>
    {/if}

    {#if showFinished}
      <button class="btn primary" onclick={() => Reset()}>
        Save
      </button>
      <button class="btn danger" onclick={() => DiscardAttempt()}>
        Discard
      </button>
    {/if}
  </div>
</Tooltip.Provider>

<style>
  .controls {
    display: flex;
    gap: 6px;
    padding: 8px 12px;
    flex-wrap: wrap;
  }

  .btn {
    flex: 1;
    min-width: 60px;
    padding: 8px 12px;
    border-radius: 6px;
    background: var(--bg-tertiary);
    color: var(--text-primary);
    font-size: 13px;
    font-weight: 500;
    transition: background 0.15s;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
  }

  .btn:hover {
    background: var(--bg-hover);
  }

  .btn.primary {
    background: var(--accent);
    color: white;
  }

  .btn.primary:hover {
    background: var(--accent-hover);
  }

  .btn.danger {
    color: var(--red);
  }

  .btn.danger:hover {
    background: rgba(255, 69, 58, 0.15);
  }

  .btn.small {
    flex: 0.5;
  }

  :global(.tooltip) {
    padding: 3px 8px;
    border-radius: 4px;
    background: var(--bg-primary);
    border: 1px solid var(--border);
    color: var(--text-secondary);
    font-size: 11px;
    font-weight: 400;
    white-space: nowrap;
    pointer-events: none;
    z-index: 50;
  }

</style>
