<script lang="ts">
  import { onMount } from 'svelte';
  import { GetAttemptHistory, DeleteSingleAttempt, EditAttemptSplits, UpdateCategoryName } from '../../../wailsjs/go/main/App';
  import { viewMode } from '../stores/splits';
  import { formatSplitTime, formatRunTime, parseTime } from '../utils/format';
  import TopNav from './TopNav.svelte';
  import type { AttemptEntry } from '../types';

  const props: {
    attemptsId: string;
    categoryName: string;
    segmentNames: string[];
  } = $props();

  let { attemptsId, segmentNames } = $derived(props);
  let displayCategoryName = $state(props.categoryName);
  let history: AttemptEntry[] = $state([]);
  let editingName = $state(false);
  let newName = $state('');
  let editingAttemptId: number | null = $state(null);
  let editInputs: string[] = $state([]);
  let editError = $state('');

  const idWidth = $derived(String(history.length).length);
  function padId(id: number): string {
    return String(id).padStart(idWidth, '0');
  }

  onMount(async () => {
    await refreshHistory();
  });

  async function refreshHistory() {
    history = (await GetAttemptHistory(attemptsId)) || [];
  }

  function startEditName() {
    newName = displayCategoryName;
    editingName = true;
  }

  async function saveName() {
    if (!newName.trim()) return;
    const data = await UpdateCategoryName(attemptsId, newName.trim());
    if (data) {
      displayCategoryName = data.categoryName as string;
    }
    editingName = false;
  }

  function cancelEditName() {
    editingName = false;
  }

  async function handleDeleteAttempt(attemptId: number) {
    await DeleteSingleAttempt(attemptsId, attemptId);
    await refreshHistory();
  }

  function segmentTime(splits: number[], index: number): number {
    if (splits[index] === 0) return 0;
    if (index === 0) return splits[0];
    const prev = splits[index - 1];
    if (prev === 0) return 0;
    return splits[index] - prev;
  }

  function startEditSplits(attempt: AttemptEntry) {
    editingAttemptId = attempt.id;
    editInputs = attempt.splitTimesMs.map(ms => formatSplitTime(ms));
    editError = '';
  }

  const derivedSegTimes = $derived(
    editInputs.map((_, i) => {
      const parsed = parseTime(editInputs[i]);
      if (parsed === null || parsed === 0) return '-';
      if (i === 0) return formatSplitTime(parsed);
      const prev = parseTime(editInputs[i - 1]);
      if (prev === null || prev === 0) return '-';
      const diff = parsed - prev;
      if (diff <= 0) return '-';
      return formatSplitTime(diff);
    })
  );

  async function saveEditSplits() {
    if (editingAttemptId === null) return;
    editError = '';

    const splits: number[] = [];
    for (let i = 0; i < editInputs.length; i++) {
      const name = i < segmentNames.length ? segmentNames[i] : `Segment ${i + 1}`;
      const val = parseTime(editInputs[i]);
      if (val === null) {
        editError = `Invalid time for "${name}"`;
        return;
      }
      splits.push(val);
    }

    // Validate monotonically increasing for non-zero values
    let lastNonZero = 0;
    for (let i = 0; i < splits.length; i++) {
      if (splits[i] === 0) continue;
      if (splits[i] <= lastNonZero) {
        editError = `"${segmentNames[i]}" must be greater than the previous split`;
        return;
      }
      lastNonZero = splits[i];
    }

    const ok = await EditAttemptSplits(attemptsId, editingAttemptId, splits);
    if (!ok) {
      editError = 'Failed to save (backend validation)';
      return;
    }
    editingAttemptId = null;
    editInputs = [];
    editError = '';
    await refreshHistory();
  }

  function cancelEditSplits() {
    editingAttemptId = null;
    editInputs = [];
    editError = '';
  }

  function handleBack() {
    viewMode.set('template_detail');
  }
</script>

<div class="editor">
  <TopNav title={displayCategoryName} onBack={handleBack}>
    {#if !editingName}
      <button class="edit-btn" onclick={startEditName}>Rename</button>
    {/if}
  </TopNav>

  {#if editingName}
    <div class="rename-bar">
      <input class="edit-name" type="text" bind:value={newName} />
      <button class="small-btn" onclick={saveName}>Save</button>
      <button class="small-btn" onclick={cancelEditName}>Cancel</button>
    </div>
  {/if}

  <div class="content">
    {#if history.length === 0}
      <div class="empty">No attempts yet</div>
    {:else}
      <div class="attempt-list">
        {#each history as attempt}
          <div class="attempt-item">
            {#if editingAttemptId === attempt.id}
              <div class="edit-sticky-header">
                <div class="attempt-header">
                  <span class="attempt-id">#{padId(attempt.id)}</span>
                  <span class="attempt-status" class:completed={attempt.completed}>
                    {attempt.completed ? 'Completed' : 'Incomplete'}
                  </span>
                  <div class="attempt-actions">
                    <button class="small-btn" onclick={saveEditSplits}>Save</button>
                    <button class="small-btn" onclick={cancelEditSplits}>Cancel</button>
                  </div>
                </div>
                <div class="edit-grid col-labels">
                  <span>Segment</span>
                  <span>Split Time</span>
                  <span>Seg. Time</span>
                </div>
              </div>
              <div class="edit-rows">
                {#each editInputs as _, i}
                  <div class="edit-grid edit-row">
                    <span class="seg-name">{i < segmentNames.length ? segmentNames[i] : `Segment ${i + 1}`}</span>
                    <span><input class="split-input" type="text" bind:value={editInputs[i]} /></span>
                    <span class="seg-time">{derivedSegTimes[i]}</span>
                  </div>
                {/each}
              </div>
              {#if editError}
                <div class="edit-error">{editError}</div>
              {/if}
            {:else}
              <div class="attempt-header">
                <span class="attempt-id">#{padId(attempt.id)}</span>
                <span class="attempt-status" class:completed={attempt.completed}>
                  {attempt.completed ? 'Completed' : 'Incomplete'}
                </span>
                {#if attempt.completed}
                  <span class="attempt-total">&middot; {formatRunTime(attempt.splitTimesMs[attempt.splitTimesMs.length - 1] || 0)}</span>
                {/if}
                <div class="attempt-actions">
                  <button class="small-btn" onclick={() => startEditSplits(attempt)}>Edit</button>
                  <button class="small-btn danger" onclick={() => handleDeleteAttempt(attempt.id)}>Delete</button>
                </div>
              </div>
            {/if}
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .editor {
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  .edit-btn {
    padding: 4px 10px;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 500;
    color: var(--text-secondary);
    background: var(--bg-tertiary);
  }

  .edit-btn:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
  }

  .rename-bar {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    border-bottom: 1px solid var(--border);
  }

  .edit-name {
    flex: 1;
    font-size: 14px;
    padding: 4px 8px;
    border-radius: 4px;
    background: var(--bg-tertiary);
  }

  .content {
    flex: 1;
    overflow-y: auto;
    padding: 0 16px;
  }

  .empty {
    padding: 24px 0;
    text-align: center;
    color: var(--text-muted);
  }

  .attempt-list {
    display: flex;
    flex-direction: column;
    gap: 6px;
    padding: 8px 0;
  }

  .attempt-item {
    padding: 8px 10px;
    border-radius: 6px;
    background: var(--bg-secondary);
  }

  .attempt-item:has(.edit-sticky-header) {
    padding-top: 0;
  }

  .attempt-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 4px;
  }

  .attempt-id {
    font-weight: 600;
    font-size: 12px;
    color: var(--text-secondary);
  }

  .attempt-status {
    font-size: 11px;
    color: var(--text-muted);
  }

  .attempt-status.completed {
    color: var(--green, #30d158);
  }

  .attempt-total {
    font-family: var(--timer-font);
    color: var(--text-secondary);
  }

  .attempt-actions {
    margin-left: auto;
    display: flex;
    gap: 4px;
  }

  .small-btn {
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 11px;
    font-weight: 500;
    color: var(--text-secondary);
    background: var(--bg-tertiary);
  }

  .small-btn:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
  }

  .small-btn.danger {
    color: var(--red);
  }

  .small-btn.danger:hover {
    background: rgba(255, 69, 58, 0.15);
  }

  .edit-sticky-header {
    position: sticky;
    top: 0;
    z-index: 1;
    background: var(--bg-secondary);
  }

  .edit-sticky-header .attempt-header {
    margin-bottom: 0;
    padding: 8px 4px 4px;
  }

  .edit-grid {
    display: grid;
    grid-template-columns: 4fr 3fr 3fr;
    font-size: 12px;
  }

  .edit-grid > span {
    padding: 3px 6px;
  }

  .col-labels {
    border-bottom: 1px solid var(--border);
  }

  .col-labels > span {
    font-size: 10px;
    font-weight: 600;
    color: var(--text-muted);
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .seg-name {
    color: var(--text-secondary);
  }

  .seg-time {
    font-family: var(--timer-font);
    color: var(--text-secondary);
  }

  .split-input {
    width: 100%;
    font-size: 12px;
    font-family: var(--timer-font);
    padding: 2px 6px;
    border-radius: 4px;
    background: var(--bg-tertiary);
    box-sizing: border-box;
  }

  .edit-error {
    font-size: 11px;
    color: var(--red, #ff453a);
    padding: 2px 6px;
  }
</style>
