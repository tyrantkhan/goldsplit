<script lang="ts">
  import { CreateAttempts } from '../../../wailsjs/go/main/App';
  import { currentTemplate, setAttempts, viewMode } from '../stores/splits';
  import TopNav from './TopNav.svelte';
  import type { AttemptsData } from '../types';

  let name = $state('');
  let categoryName = $state('');

  async function handleCreate() {
    if (!name.trim() || !categoryName.trim() || !$currentTemplate) return;

    const data = await CreateAttempts($currentTemplate.id, name.trim(), categoryName.trim());
    if (data) {
      setAttempts(data as AttemptsData);
    }
  }

  function handleCancel() {
    viewMode.set('template_detail');
  }
</script>

<div class="page">
  <TopNav title="New Category" subtitle="For {$currentTemplate?.name || ''}" onBack={handleCancel} />

  <div class="form-content">
    <label>
      <span>Name</span>
      <input type="text" bind:value={name} placeholder="16 Star" />
    </label>

    <label>
      <span>Category</span>
      <input type="text" bind:value={categoryName} placeholder="Any%" />
    </label>

    <div class="actions">
      <button class="btn cancel" onclick={handleCancel}>Cancel</button>
      <button class="btn create" onclick={handleCreate} disabled={!name.trim() || !categoryName.trim()}>
        Create
      </button>
    </div>
  </div>
</div>

<style>
  .page {
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  .form-content {
    flex: 1;
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 14px;
  }

  label {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  label span {
    font-size: 12px;
    color: var(--text-secondary);
    font-weight: 500;
  }

  .actions {
    display: flex;
    gap: 8px;
    margin-top: auto;
  }

  .btn {
    flex: 1;
    padding: 10px;
    border-radius: 6px;
    font-size: 14px;
    font-weight: 500;
  }

  .cancel {
    background: var(--bg-tertiary);
  }

  .cancel:hover {
    background: var(--bg-hover);
  }

  .create {
    background: var(--accent);
    color: white;
  }

  .create:hover:not(:disabled) {
    background: var(--accent-hover);
  }

  .create:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }
</style>
