import { writable } from 'svelte/store';
import { EventsOn } from '../../../wailsjs/runtime/runtime';
import type { TemplateData, AttemptsData, Delta, ViewMode } from '../types';

export const currentTemplate = writable<TemplateData | null>(null);
export const currentAttempts = writable<AttemptsData | null>(null);
export const deltas = writable<Delta[]>([]);
export const viewMode = writable<ViewMode>('templates');

let _previousView: ViewMode = 'templates';

export function openSettings() {
  let current: ViewMode = 'templates';
  viewMode.subscribe(v => current = v)();
  _previousView = current;
  viewMode.set('settings');
}

export function closeSettings() {
  viewMode.set(_previousView);
}

export function openAbout() {
  let current: ViewMode = 'templates';
  viewMode.subscribe(v => current = v)();
  _previousView = current;
  viewMode.set('about');
}

export function closeAbout() {
  viewMode.set(_previousView);
}

export function initSplitEvents() {
  EventsOn('attempts:updated', (data: AttemptsData) => {
    if (data) {
      currentAttempts.set(data);
    }
  });
}

export function setTemplate(data: TemplateData | null) {
  currentTemplate.set(data);
  if (data) {
    viewMode.set('template_detail');
  }
}

export function setAttempts(data: AttemptsData | null) {
  currentAttempts.set(data);
  if (data) {
    viewMode.set('timer');
  }
}

export function backToTemplates() {
  currentTemplate.set(null);
  currentAttempts.set(null);
  deltas.set([]);
  viewMode.set('templates');
}

export function backToTemplateDetail() {
  currentAttempts.set(null);
  deltas.set([]);
  viewMode.set('template_detail');
}
