import { writable } from 'svelte/store';
import { EventsOn } from '../../../wailsjs/runtime/runtime';
import { GetSettings, UpdateSettings } from '../../../wailsjs/go/main/App';
import type { Settings, ColorSettings } from '../types';

const defaultSettings: Settings = {
  alwaysOnTop: false,
  hotkeys: {
    startSplit: 'Space',
    pause: 'KeyP',
    reset: 'KeyR',
    undoSplit: 'Backspace',
    skipSplit: 'KeyS',
  },
  comparison: 'personal_best',
  colors: {
    aheadGaining: '#30d158',
    aheadLosing: '#7ec890',
    behindGaining: '#cc6b65',
    behindLosing: '#ff453a',
    bestSegment: '#ffd60a',
  },
};

export const settings = writable<Settings>(defaultSettings);

function applyColorVars(colors: ColorSettings) {
  const s = document.documentElement.style;
  s.setProperty('--ahead-gaining', colors.aheadGaining);
  s.setProperty('--ahead-losing', colors.aheadLosing);
  s.setProperty('--behind-gaining', colors.behindGaining);
  s.setProperty('--behind-losing', colors.behindLosing);
  s.setProperty('--best-segment', colors.bestSegment);
}

export async function initSettings() {
  const s = await GetSettings();
  if (s) {
    settings.set(s);
    applyColorVars(s.colors);
  }

  EventsOn('settings:updated', (s: Settings) => {
    if (s) {
      settings.set(s);
      applyColorVars(s.colors);
    }
  });
}

export async function saveSettings(updated: Settings): Promise<boolean> {
  const ok = await UpdateSettings(updated as any);
  if (ok) {
    settings.set(updated);
  }
  return ok;
}
