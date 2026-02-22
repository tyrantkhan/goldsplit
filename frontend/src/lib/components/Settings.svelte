<script lang="ts">
  import { Tabs, Select, Switch } from 'bits-ui';
  import { settings, saveSettings } from '../stores/settings';
  import { closeSettings } from '../stores/splits';
  import TopNav from './TopNav.svelte';
  import type { Settings, HotkeyBindings, ColorSettings } from '../types';

  type Tab = 'general' | 'hotkeys' | 'colors';
  let activeTab: Tab = $state('general');

  let capturingKey: keyof HotkeyBindings | null = $state(null);

  const hotkeyLabels: { key: keyof HotkeyBindings; label: string }[] = [
    { key: 'startSplit', label: 'Start / Split' },
    { key: 'pause', label: 'Pause' },
    { key: 'reset', label: 'Reset' },
    { key: 'undoSplit', label: 'Undo Split' },
    { key: 'skipSplit', label: 'Skip Split' },
  ];

  function displayKey(code: string): string {
    if (code.startsWith('Key')) return code.slice(3);
    if (code.startsWith('Digit')) return code.slice(5);
    if (code === 'ArrowLeft') return '\u2190';
    if (code === 'ArrowRight') return '\u2192';
    if (code === 'ArrowUp') return '\u2191';
    if (code === 'ArrowDown') return '\u2193';
    return code;
  }

  function handleBack() {
    closeSettings();
  }

  async function handleAlwaysOnTopChange(checked: boolean) {
    const updated: Settings = { ...$settings, alwaysOnTop: checked };
    await saveSettings(updated);
  }

  const comparisonOptions = [
    { value: 'personal_best', label: 'Personal Best' },
    { value: 'best_segments', label: 'Best Segments' },
    { value: 'average_segments', label: 'Average Segments' },
    { value: 'latest_run', label: 'Latest Run' },
  ];

  const comparisonLabel = $derived(comparisonOptions.find(o => o.value === $settings.comparison)?.label ?? 'Personal Best');

  async function handleComparisonChange(value: string) {
    const updated: Settings = { ...$settings, comparison: value };
    await saveSettings(updated);
  }

  function startCapture(key: keyof HotkeyBindings) {
    capturingKey = key;
  }

  async function handleKeydown(e: KeyboardEvent) {
    if (!capturingKey) return;

    e.preventDefault();
    e.stopPropagation();

    if (e.code === 'Escape') {
      capturingKey = null;
      return;
    }

    const updated: Settings = {
      ...$settings,
      hotkeys: { ...$settings.hotkeys, [capturingKey]: e.code },
    };
    capturingKey = null;
    await saveSettings(updated);
  }

  const colorLabels: { key: keyof ColorSettings; label: string }[] = [
    { key: 'aheadGaining', label: 'Ahead & Gaining' },
    { key: 'aheadLosing', label: 'Ahead & Losing' },
    { key: 'behindGaining', label: 'Behind & Gaining' },
    { key: 'behindLosing', label: 'Behind & Losing' },
    { key: 'bestSegment', label: 'Best Segment' },
  ];

  const defaultColors: ColorSettings = {
    aheadGaining: '#30d158',
    aheadLosing: '#7ec890',
    behindGaining: '#cc6b65',
    behindLosing: '#ff453a',
    bestSegment: '#ffd60a',
  };

  async function updateColor(key: keyof ColorSettings, value: string) {
    const updated: Settings = {
      ...$settings,
      colors: { ...$settings.colors, [key]: value },
    };
    await saveSettings(updated);
  }

  async function resetColors() {
    const updated: Settings = {
      ...$settings,
      colors: { ...defaultColors },
    };
    await saveSettings(updated);
  }
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="settings">
  <TopNav title="Settings" onBack={handleBack} />

  <Tabs.Root bind:value={activeTab}>
    <Tabs.List class="tabs">
      <Tabs.Trigger value="general" class="tab">General</Tabs.Trigger>
      <Tabs.Trigger value="hotkeys" class="tab">Hotkeys</Tabs.Trigger>
      <Tabs.Trigger value="colors" class="tab">Colors</Tabs.Trigger>
    </Tabs.List>

    <div class="content">
      <Tabs.Content value="general">
        <section class="section">
          <div class="row">
            <span class="label">Always on Top</span>
            <Switch.Root
              checked={$settings.alwaysOnTop}
              onCheckedChange={handleAlwaysOnTopChange}
              class="toggle"
            >
              <Switch.Thumb class="toggle-knob" />
            </Switch.Root>
          </div>
        </section>

        <section class="section">
          <h3>Comparison</h3>
          <div class="row">
            <span class="label">Compare against</span>
            <Select.Root type="single" value={$settings.comparison} onValueChange={handleComparisonChange}>
              <Select.Trigger class="dropdown-btn">
                {comparisonLabel}
                <span class="dropdown-arrow">&#x25BE;</span>
              </Select.Trigger>
              <Select.Content class="dropdown-menu">
                {#each comparisonOptions as opt}
                  <Select.Item value={opt.value} label={opt.label} class="dropdown-item">
                    {opt.label}
                  </Select.Item>
                {/each}
              </Select.Content>
            </Select.Root>
          </div>
        </section>
      </Tabs.Content>

      <Tabs.Content value="hotkeys">
        <section class="section">
          {#each hotkeyLabels as { key, label }}
            <div class="row">
              <span class="label">{label}</span>
              <button
                class="key-btn"
                class:capturing={capturingKey === key}
                onclick={() => startCapture(key)}
              >
                {#if capturingKey === key}
                  Press a key...
                {:else}
                  {displayKey($settings.hotkeys[key])}
                {/if}
              </button>
            </div>
          {/each}
        </section>
      </Tabs.Content>

      <Tabs.Content value="colors">
        <section class="section">
          <h3>Delta Colors</h3>
          {#each colorLabels as { key, label }}
            <div class="row">
              <span class="label">{label}</span>
              <input
                type="color"
                class="color-picker"
                value={$settings.colors[key]}
                onchange={(e) => updateColor(key, e.currentTarget.value)}
              />
            </div>
          {/each}
        </section>
        <section class="section">
          <button class="reset-btn" onclick={resetColors}>Reset to Defaults</button>
        </section>
      </Tabs.Content>
    </div>
  </Tabs.Root>
</div>

<style>
  .settings {
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  .settings :global(.tabs) {
    display: flex;
    padding: 0 12px;
    gap: 0;
    border-bottom: 1px solid var(--border);
  }

  .settings :global(.tab) {
    padding: 8px 14px;
    font-size: 13px;
    font-weight: 500;
    color: var(--text-secondary);
    border-bottom: 2px solid transparent;
    margin-bottom: -1px;
    transition: color 0.15s;
  }

  .settings :global(.tab:hover) {
    color: var(--text-primary);
  }

  .settings :global(.tab[data-state=active]) {
    color: var(--text-primary);
    border-bottom-color: var(--accent);
  }

  .content {
    flex: 1;
    overflow-y: auto;
    padding: 12px 16px;
  }

  .section {
    margin-bottom: 20px;
  }

  h3 {
    font-size: 12px;
    font-weight: 600;
    color: var(--text-secondary);
    text-transform: uppercase;
    letter-spacing: 0.5px;
    margin-bottom: 8px;
  }

  .row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 0;
  }

  .label {
    font-size: 13px;
  }

  .settings :global(.dropdown-btn) {
    min-width: 80px;
    padding: 4px 10px;
    border-radius: 4px;
    background: var(--bg-tertiary);
    color: var(--text-primary);
    font-size: 12px;
    font-weight: 500;
    text-align: center;
    transition: background 0.15s;
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .settings :global(.dropdown-btn:hover) {
    background: var(--bg-hover);
  }

  .settings :global(.dropdown-arrow) {
    font-size: 10px;
    color: var(--text-secondary);
  }

  :global(.dropdown-menu) {
    min-width: 160px;
    background: var(--bg-secondary);
    border: 1px solid var(--border);
    border-radius: 6px;
    padding: 4px 0;
    z-index: 10;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  }

  :global(.dropdown-item) {
    width: 100%;
    padding: 6px 12px;
    font-size: 12px;
    text-align: left;
    color: var(--text-primary);
    transition: background 0.1s;
  }

  :global(.dropdown-item:hover),
  :global(.dropdown-item[data-highlighted]) {
    background: var(--bg-tertiary);
  }

  :global(.dropdown-item[data-state=checked]) {
    color: var(--accent);
  }

  .settings :global(.toggle) {
    width: 40px;
    height: 22px;
    border-radius: 11px;
    background: var(--bg-tertiary);
    position: relative;
    transition: background 0.2s;
    padding: 0;
  }

  .settings :global(.toggle[data-state=checked]) {
    background: var(--accent);
  }

  .settings :global(.toggle-knob) {
    position: absolute;
    top: 2px;
    left: 2px;
    width: 18px;
    height: 18px;
    border-radius: 50%;
    background: white;
    transition: transform 0.2s;
  }

  .settings :global(.toggle[data-state=checked] .toggle-knob) {
    transform: translateX(18px);
  }

  .key-btn {
    min-width: 80px;
    padding: 4px 10px;
    border-radius: 4px;
    background: var(--bg-tertiary);
    color: var(--text-primary);
    font-size: 12px;
    font-weight: 500;
    text-align: center;
    transition: background 0.15s;
  }

  .key-btn:hover {
    background: var(--bg-hover);
  }

  .key-btn.capturing {
    background: var(--accent);
    color: white;
    font-size: 11px;
  }

  .color-picker {
    width: 36px;
    height: 28px;
    padding: 2px;
    border: 1px solid var(--border);
    border-radius: 4px;
    background: var(--bg-tertiary);
    cursor: pointer;
  }

  .color-picker::-webkit-color-swatch-wrapper {
    padding: 2px;
  }

  .color-picker::-webkit-color-swatch {
    border: none;
    border-radius: 2px;
  }

  .reset-btn {
    width: 100%;
    padding: 8px 12px;
    border-radius: 4px;
    background: var(--bg-tertiary);
    color: var(--text-secondary);
    font-size: 12px;
    font-weight: 500;
    transition: background 0.15s, color 0.15s;
  }

  .reset-btn:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
  }
</style>
