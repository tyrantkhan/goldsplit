import { writable, derived, get } from 'svelte/store';
import { EventsOn } from '../../../wailsjs/runtime/runtime';
import type { TickData, TimerState } from '../types';

export const timerState = writable<TimerState>('idle');
export const elapsedMs = writable<number>(0);
export const currentSegment = writable<number>(0);
export const splitTimesMs = writable<number[]>([]);
export const segmentTimesMs = writable<number[]>([]);
export const splitNames = writable<string[]>([]);

// For rAF interpolation
let lastTickTime = 0;
let lastTickElapsed = 0;
let rafId: number | null = null;

// Interpolated elapsed time for smooth display
export const interpolatedElapsed = writable<number>(0);

function startInterpolation() {
  function tick() {
    const state = get(timerState);
    if (state === 'running') {
      const now = performance.now();
      const interpolated = lastTickElapsed + (now - lastTickTime);
      interpolatedElapsed.set(interpolated);
      rafId = requestAnimationFrame(tick);
    }
  }
  if (rafId !== null) cancelAnimationFrame(rafId);
  rafId = requestAnimationFrame(tick);
}

function stopInterpolation() {
  if (rafId !== null) {
    cancelAnimationFrame(rafId);
    rafId = null;
  }
}

export function initTimerEvents() {
  EventsOn('timer:tick', (data: TickData) => {
    lastTickTime = performance.now();
    lastTickElapsed = data.elapsedMs;

    elapsedMs.set(data.elapsedMs);
    currentSegment.set(data.currentSegment);
    splitTimesMs.set(data.splitTimesMs || []);
    segmentTimesMs.set(data.segmentTimesMs || []);
    splitNames.set(data.splitNames || []);

    if (data.state === 'running') {
      interpolatedElapsed.set(data.elapsedMs);
    }
  });

  EventsOn('timer:state', (state: TimerState) => {
    timerState.set(state);

    if (state === 'running') {
      startInterpolation();
    } else {
      stopInterpolation();
      // Set exact value when not running
      interpolatedElapsed.set(get(elapsedMs));
    }
  });
}

export const isIdle = derived(timerState, $s => $s === 'idle');
export const isRunning = derived(timerState, $s => $s === 'running');
export const isPaused = derived(timerState, $s => $s === 'paused');
export const isFinished = derived(timerState, $s => $s === 'finished');
