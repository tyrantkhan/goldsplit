<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    title,
    subtitle = '',
    onBack = null,
    children,
  }: {
    title: string;
    subtitle?: string;
    onBack?: (() => void) | null;
    children?: Snippet;
  } = $props();
</script>

<header class="topnav">
  {#if onBack}
    <button class="back-btn" onclick={onBack}>&larr;</button>
  {/if}
  <div class="title-area">
    <span class="title">{title}</span>
    {#if subtitle}
      <span class="subtitle">{subtitle}</span>
    {/if}
  </div>
  <div class="actions">
    {@render children?.()}
  </div>
</header>

<style>
  .topnav {
    display: flex;
    align-items: center;
    padding: 8px 12px;
    gap: 8px;
    min-height: 44px;
    border-bottom: 1px solid var(--border);
  }

  .back-btn {
    width: 28px;
    height: 28px;
    border-radius: 4px;
    font-size: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--text-secondary);
    flex-shrink: 0;
  }

  .back-btn:hover {
    background: var(--bg-tertiary);
    color: var(--text-primary);
  }

  .title-area {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 1px;
    min-width: 0;
  }

  .title {
    font-size: 14px;
    font-weight: 600;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .subtitle {
    font-size: 11px;
    color: var(--text-secondary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .actions {
    display: flex;
    align-items: center;
    gap: 4px;
    flex-shrink: 0;
  }
</style>
