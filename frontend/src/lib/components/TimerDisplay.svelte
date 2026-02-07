<script lang="ts">
  import { interpolatedElapsed, timerState, elapsedMs } from '../stores/timer';
  import { formatTime } from '../utils/format';

  const displayMs = $derived($timerState === 'running' ? $interpolatedElapsed : $elapsedMs);
  const displayTime = $derived(formatTime(Math.floor(displayMs)));
  const stateClass = $derived($timerState);
</script>

<div class="timer-display {stateClass}">
  <span class="time">{displayTime}</span>
</div>

<style>
  .timer-display {
    padding: 16px 0 8px;
    text-align: center;
  }

  .time {
    font-family: var(--timer-font);
    font-size: var(--timer-size);
    font-weight: 700;
    font-variant-numeric: tabular-nums;
    letter-spacing: -1px;
    color: var(--text-primary);
  }

  .timer-display.running .time {
    color: var(--text-primary);
  }

  .timer-display.paused .time {
    color: var(--text-secondary);
  }

  .timer-display.finished .time {
    color: var(--green);
  }
</style>
