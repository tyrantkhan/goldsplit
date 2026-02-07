<script lang="ts">
  import { CreateTemplate } from '../../../wailsjs/go/main/App';
  import { setTemplate, viewMode } from '../stores/splits';
  import TopNav from './TopNav.svelte';

  let name = $state('');
  let segmentsText = $state('');

  async function handleCreate() {
    const names = segmentsText
      .split('\n')
      .map(s => s.trim())
      .filter(s => s.length > 0);

    if (!name.trim() || names.length === 0) return;

    const data = await CreateTemplate(name.trim(), names);
    if (data) {
      setTemplate(data);
    }
  }

  function handleCancel() {
    viewMode.set('templates');
  }
</script>

<div class="page">
  <TopNav title="New Game" onBack={handleCancel} />

  <div class="form-content">
    <label>
      <span>Game Name</span>
      <input type="text" bind:value={name} placeholder="Super Mario 64" />
    </label>

    <label>
      <span>Segments <span class="sub">(one per line)</span></span>
      <textarea bind:value={segmentsText} rows="8" placeholder={"Bob-omb Battlefield\nWhomp's Fortress\nJolly Roger Bay"}></textarea>
    </label>

    <div class="actions">
      <button class="btn cancel" onclick={handleCancel}>Cancel</button>
      <button class="btn create" onclick={handleCreate} disabled={!name.trim() || !segmentsText.trim()}>
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

  .sub {
    opacity: 0.6;
    font-weight: 400;
  }

  textarea {
    resize: vertical;
    min-height: 100px;
    font-family: var(--timer-font);
    font-size: 13px;
    line-height: 1.5;
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
