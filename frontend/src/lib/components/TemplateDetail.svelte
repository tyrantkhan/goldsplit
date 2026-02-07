<script lang="ts">
  import { onMount } from 'svelte';
  import { ListAttemptsForTemplate, LoadAttempts, DeleteAttempts, UpdateTemplate } from '../../../wailsjs/go/main/App';
  import { currentTemplate, setAttempts, backToTemplates, viewMode, openSettings } from '../stores/splits';
  import TopNav from './TopNav.svelte';
  import IconSettings from '../icons/IconSettings.svelte';
  import IconEdit from '../icons/IconEdit.svelte';
  import IconPlus from '../icons/IconPlus.svelte';
  import type { AttemptsSummary } from '../types';

  let categories: AttemptsSummary[] = $state([]);
  let editing = $state(false);
  let editName = $state('');
  let editSegments = $state('');

  const templateId = $derived($currentTemplate?.id || '');

  onMount(async () => {
    await refreshList();
  });

  async function refreshList() {
    if (!templateId) return;
    categories = (await ListAttemptsForTemplate(templateId)) || [];
    categories.sort((a, b) => b.updatedAt - a.updatedAt);
  }

  async function handleLoad(id: string) {
    const data = await LoadAttempts(id);
    if (data) {
      setAttempts(data);
    }
  }

  async function handleDelete(e: Event, id: string) {
    e.stopPropagation();
    await DeleteAttempts(id);
    await refreshList();
  }

  function handleNewCategory() {
    viewMode.set('attempts_setup');
  }

  function startEdit() {
    if (!$currentTemplate) return;
    editName = $currentTemplate.name;
    editSegments = $currentTemplate.segmentNames.join('\n');
    editing = true;
  }

  async function saveEdit() {
    if (!$currentTemplate) return;
    const names = editSegments
      .split('\n')
      .map(s => s.trim())
      .filter(s => s.length > 0);

    if (!editName.trim() || names.length === 0) return;

    const data = await UpdateTemplate($currentTemplate.id, editName.trim(), names);
    if (data) {
      currentTemplate.set(data);
    }
    editing = false;
  }

  function cancelEdit() {
    editing = false;
  }
</script>

<div class="detail">
  {#if editing}
    <TopNav title="Edit Game" onBack={cancelEdit} />

    <div class="edit-section">
      <label>
        <span>Game Name</span>
        <input class="edit-input" type="text" bind:value={editName} />
      </label>
      <label>
        <span>Segments <span class="sub">(one per line)</span></span>
        <textarea bind:value={editSegments} rows="6"></textarea>
      </label>
      <div class="edit-actions">
        <button class="btn cancel" onclick={cancelEdit}>Cancel</button>
        <button class="btn save" onclick={saveEdit} disabled={!editName.trim() || !editSegments.trim()}>Save</button>
      </div>
    </div>
  {:else}
    <TopNav title={$currentTemplate?.name || ''} subtitle="{$currentTemplate?.segmentNames.length || 0} segments" onBack={backToTemplates}>
      <button class="icon-btn" onclick={openSettings} title="Settings">
        <IconSettings />
      </button>
      <button class="icon-btn" onclick={startEdit} title="Edit game">
        <IconEdit />
      </button>
      <button class="icon-btn accent" onclick={handleNewCategory} title="New attempt">
        <IconPlus />
      </button>
    </TopNav>

    <div class="section-divider"><span>Attempts</span></div>

    {#if categories.length === 0}
      <div class="empty">
        <p>No attempts yet</p>
        <p class="sub">Create an attempt to start tracking</p>
      </div>
    {:else}
      <div class="list">
        {#each categories as cat}
          <div class="list-item" role="button" tabindex="0" onclick={() => handleLoad(cat.id)} onkeydown={(e) => e.key === 'Enter' && handleLoad(cat.id)}>
            <div class="info">
              <span class="name">{cat.name}</span>
              <span class="meta-text">{cat.categoryName} &middot; {cat.attemptCount} attempts</span>
            </div>
            <div class="actions">
              <button class="delete-btn" onclick={(e) => handleDelete(e, cat.id)}>x</button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  {/if}
</div>

<style>
  .detail {
    height: 100%;
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

  .icon-btn.accent {
    color: var(--accent);
  }

  .icon-btn.accent:hover {
    background: rgba(var(--accent-rgb, 50, 130, 240), 0.15);
  }

  .edit-section {
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 12px;
    flex: 1;
  }

  .edit-section label {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .edit-section label span {
    font-size: 12px;
    color: var(--text-secondary);
    font-weight: 500;
  }

  .sub {
    opacity: 0.6;
    font-weight: 400;
  }

  .edit-input {
    font-size: 14px;
    padding: 6px 8px;
    border-radius: 4px;
    background: var(--bg-tertiary);
  }

  .edit-section textarea {
    resize: vertical;
    min-height: 80px;
    font-family: var(--timer-font);
    font-size: 13px;
    line-height: 1.5;
  }

  .edit-actions {
    display: flex;
    gap: 8px;
    margin-top: auto;
  }

  .btn {
    flex: 1;
    padding: 8px;
    border-radius: 6px;
    font-size: 13px;
    font-weight: 500;
  }

  .cancel {
    background: var(--bg-tertiary);
  }

  .cancel:hover {
    background: var(--bg-hover);
  }

  .save {
    background: var(--accent);
    color: white;
  }

  .save:hover:not(:disabled) {
    background: var(--accent-hover);
  }

  .save:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  .section-divider {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px 16px 4px;
  }

  .section-divider::before,
  .section-divider::after {
    content: '';
    flex: 1;
    height: 1px;
    background: var(--border);
  }

  .section-divider span {
    font-size: 11px;
    font-weight: 600;
    color: var(--text-muted);
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .empty {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 4px;
    color: var(--text-secondary);
  }

  .empty .sub {
    font-size: 12px;
    color: var(--text-muted);
  }

  .list {
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding: 8px 12px;
    overflow-y: auto;
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

  .actions {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .delete-btn {
    width: 22px;
    height: 22px;
    border-radius: 4px;
    font-size: 12px;
    color: var(--text-muted);
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .delete-btn:hover {
    background: rgba(255, 69, 58, 0.2);
    color: var(--red);
  }
</style>
