<script lang="ts">
  let { size = 16 }: { size?: number } = $props();
</script>

<svg width={size} height={size} viewBox="0 0 16 16" fill="none">
  <circle cx="8" cy="8" r="6.5" stroke="currentColor" stroke-width="1.2"/>
  <path d="M8 7v4M8 5v1" stroke="currentColor" stroke-width="1.2" stroke-linecap="round"/>
</svg>
