<script lang="ts">
  import { onMount } from 'svelte';
  import { ListTemplates, LoadTemplate, DeleteTemplate, ConfirmDialog } from '../../../wailsjs/go/main/App';
  import { setTemplate, viewMode, openSettings, openAbout } from '../stores/splits';
  import TopNav from './TopNav.svelte';
  import IconInfo from '../icons/IconInfo.svelte';
  import IconSettings from '../icons/IconSettings.svelte';
  import IconPlus from '../icons/IconPlus.svelte';
  import IconTrash from '../icons/IconTrash.svelte';
  import type { TemplateSummary, TemplateData } from '../types';

  let templates: TemplateSummary[] = $state([]);

  onMount(async () => {
    await refreshList();
  });

  async function refreshList() {
    templates = (await ListTemplates()) || [];
    templates.sort((a, b) => b.updatedAt - a.updatedAt);
  }

  async function handleLoad(id: string) {
    const data = await LoadTemplate(id);
    if (data) {
      setTemplate(data as TemplateData);
    }
  }

  async function handleDelete(e: Event, id: string) {
    e.stopPropagation();
    if (!await ConfirmDialog('Delete Game', 'Delete this game and all its attempts?')) return;
    await DeleteTemplate(id);
    await refreshList();
  }

  function handleNew() {
    viewMode.set('template_setup');
  }
</script>

<div class="page">
  <TopNav title="Games">
    <button class="icon-btn" onclick={openAbout} title="About">
      <IconInfo />
    </button>
    <button class="icon-btn" onclick={openSettings} title="Settings">
      <IconSettings />
    </button>
    <button class="icon-btn accent" onclick={handleNew} title="New game">
      <IconPlus />
    </button>
  </TopNav>

  <div class="content">
    {#if templates.length === 0}
      <div class="empty">
        <p>No games yet</p>
        <p class="sub">Create your first game to get started</p>
      </div>
    {:else}
      <div class="list">
        {#each templates as tmpl}
          <div class="list-item" role="button" tabindex="0" onclick={() => handleLoad(tmpl.id)} onkeydown={(e) => e.key === 'Enter' && handleLoad(tmpl.id)}>
            <div class="info">
              <span class="name">{tmpl.name}</span>
              <span class="meta-text">{tmpl.segmentCount} segments</span>
            </div>
            <div class="item-actions">
              <button class="delete-btn" onclick={(e) => handleDelete(e, tmpl.id)}><IconTrash /></button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .page {
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  .content {
    flex: 1;
    overflow-y: auto;
    padding: 8px 12px;
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

  .icon-btn.accent {
    color: var(--accent);
  }

  .empty {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 4px;
    color: var(--text-secondary);
    padding: 40px 0;
  }

  .empty .sub {
    font-size: 12px;
    color: var(--text-muted);
  }

  .list {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .list-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 12px;
    border-radius: 6px;
    background: var(--bg-secondary);
    transition: background 0.15s;
    text-align: left;
    width: 100%;
    cursor: pointer;
  }

  .list-item:hover {
    background: var(--bg-tertiary);
  }

  .info {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .name {
    font-weight: 500;
    font-size: 14px;
  }

  .meta-text {
    font-size: 12px;
    color: var(--text-secondary);
  }

  .item-actions {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .delete-btn {
    width: 22px;
    height: 22px;
    border-radius: 4px;
    color: var(--red);
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .delete-btn:hover {
    background: rgba(255, 69, 58, 0.2);
  }
</style>
