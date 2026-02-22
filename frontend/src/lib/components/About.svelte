<script lang="ts">
  import { onMount } from 'svelte';
  import { GetAppInfo } from '../../../wailsjs/go/main/App';
  import { BrowserOpenURL } from '../../../wailsjs/runtime/runtime';
  import { closeAbout } from '../stores/splits';
  import TopNav from './TopNav.svelte';

  let version = $state('');

  onMount(async () => {
    const info = await GetAppInfo();
    if (info) {
      version = info.version;
    }
  });

  function handleBack() {
    closeAbout();
  }

  function openLink(url: string) {
    BrowserOpenURL(url);
  }
</script>

<div class="about">
  <TopNav title="About" onBack={handleBack} />

  <div class="content">
    <section class="hero">
      <h1>Goldsplit</h1>
      <span class="version">{version}</span>
    </section>

    <section class="section">
      <div class="row">
        <span class="label">Author</span>
        <span class="value">Haris Khan</span>
      </div>
    </section>

    <section class="section">
      <h3>Links</h3>
      <button class="link-btn" onclick={() => openLink('https://github.com/tyrantkhan/goldsplit')}>
        GitHub Repository
      </button>
      <button class="link-btn" onclick={() => openLink('https://github.com/tyrantkhan/goldsplit/blob/main/LICENSE')}>
        License
      </button>
      <button class="link-btn" onclick={() => openLink('https://github.com/tyrantkhan/goldsplit/issues')}>
        Report an Issue
      </button>
    </section>
  </div>
</div>

<style>
  .about {
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  .content {
    flex: 1;
    overflow-y: auto;
    padding: 12px 16px;
  }

  .hero {
    text-align: center;
    padding: 24px 0;
  }

  h1 {
    font-size: 20px;
    font-weight: 600;
    margin-bottom: 4px;
  }

  .version {
    font-size: 13px;
    color: var(--text-secondary);
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

  .value {
    font-size: 13px;
    color: var(--text-secondary);
  }

  .link-btn {
    display: block;
    width: 100%;
    padding: 10px 12px;
    border-radius: 6px;
    background: var(--bg-secondary);
    color: var(--text-primary);
    font-size: 13px;
    text-align: left;
    transition: background 0.15s;
    margin-bottom: 4px;
  }

  .link-btn:hover {
    background: var(--bg-tertiary);
  }
</style>
