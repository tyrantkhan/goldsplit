<script lang="ts">
  import { onMount } from 'svelte';
  import { initTimerEvents, timerState } from './lib/stores/timer';
  import { initSplitEvents, viewMode, currentTemplate, currentAttempts, backToTemplateDetail, openSettings } from './lib/stores/splits';
  import { initSettings } from './lib/stores/settings';
  import { deltas } from './lib/stores/splits';
  import { EventsOn } from '../wailsjs/runtime/runtime';
  import TopNav from './lib/components/TopNav.svelte';
  import IconEdit from './lib/icons/IconEdit.svelte';
  import IconSettings from './lib/icons/IconSettings.svelte';
  import TimerDisplay from './lib/components/TimerDisplay.svelte';
  import SplitsList from './lib/components/SplitsList.svelte';
  import Controls from './lib/components/Controls.svelte';
  import TemplateList from './lib/components/TemplateList.svelte';
  import TemplateSetup from './lib/components/TemplateSetup.svelte';
  import TemplateDetail from './lib/components/TemplateDetail.svelte';
  import AttemptsSetup from './lib/components/AttemptsSetup.svelte';
  import AttemptEditor from './lib/components/AttemptEditor.svelte';
  import Settings from './lib/components/Settings.svelte';

  onMount(() => {
    initTimerEvents();
    initSplitEvents();
    initSettings();

    EventsOn('deltas:updated', (d) => {
      if (d) deltas.set(d);
    });

    EventsOn('timer:state', (state) => {
      if (state === 'idle') deltas.set([]);
    });
  });

  function handleBack() {
    backToTemplateDetail();
  }

  // openSettings imported from store

  function openAttemptEditor() {
    viewMode.set('attempt_editor');
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.code === 'Escape' && $viewMode === 'timer' && $timerState === 'idle') {
      handleBack();
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="app">
  {#if $viewMode === 'timer' && $currentAttempts}
    <TopNav
      title={$currentAttempts.name}
      subtitle="{$currentAttempts.categoryName} · #{$currentAttempts.attemptCount}"
      onBack={$timerState === 'idle' ? handleBack : null}
    >
      {#if $timerState === 'idle'}
        <button class="icon-btn" onclick={openAttemptEditor} title="Edit attempts">
          <IconEdit />
        </button>
        <button class="icon-btn" onclick={openSettings} title="Settings">
          <IconSettings />
        </button>
      {/if}
    </TopNav>

    <SplitsList />
    <TimerDisplay />
    <Controls />
  {:else if $viewMode === 'template_setup'}
    <TemplateSetup />
  {:else if $viewMode === 'template_detail'}
    <TemplateDetail />
  {:else if $viewMode === 'attempts_setup'}
    <AttemptsSetup />
  {:else if $viewMode === 'attempt_editor' && $currentAttempts}
    <AttemptEditor attemptsId={$currentAttempts.id} categoryName={$currentAttempts.categoryName} segmentNames={$currentAttempts.segments.map(s => s.name)} />
  {:else if $viewMode === 'settings'}
    <Settings />
  {:else}
    <TemplateList />
  {/if}
</div>

<style>
  .app {
    height: 100vh;
    display: flex;
    flex-direction: column;
  }

  .icon-btn {
    width: 28px;
    height: 28px;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--text-secondary);
  }

  .icon-btn:hover {
    background: var(--bg-tertiary);
    color: var(--text-primary);
  }
</style>
