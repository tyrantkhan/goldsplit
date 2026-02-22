export interface TickData {
  elapsedMs: number;
  state: TimerState;
  currentSegment: number;
  splitTimesMs: number[];
  segmentTimesMs: number[];
  splitNames: string[];
}

export type TimerState = 'idle' | 'running' | 'paused' | 'finished';

export interface Segment {
  name: string;
  personalBestMs: number;
  bestSegmentMs: number;
  comparisonSplitMs: number;
}

export interface TemplateData {
  id: string;
  name: string;
  segmentNames: string[];
}

export interface TemplateSummary {
  id: string;
  name: string;
  segmentCount: number;
  updatedAt: number;
}

export interface AttemptsData {
  id: string;
  templateId: string;
  name: string;
  categoryName: string;
  segments: Segment[];
  attemptCount: number;
}

export interface AttemptsSummary {
  id: string;
  templateId: string;
  name: string;
  categoryName: string;
  attemptCount: number;
  updatedAt: number;
}

export interface AttemptEntry {
  id: number;
  startedAt: string;
  splitTimesMs: number[];
  completed: boolean;
}

export interface Delta {
  segmentIndex: number;
  deltaMs: number;
  isBestEver: boolean;
  isAhead: boolean;
  gainedTime: boolean;
  skipped: boolean;
}

export type ViewMode = 'templates' | 'template_detail' | 'template_setup' | 'attempts_setup' | 'timer' | 'settings' | 'attempt_editor' | 'about';

export interface HotkeyBindings {
  startSplit: string;
  pause: string;
  reset: string;
  undoSplit: string;
  skipSplit: string;
}

export interface ColorSettings {
  aheadGaining: string;
  aheadLosing: string;
  behindGaining: string;
  behindLosing: string;
  bestTime: string;
}

export interface Settings {
  alwaysOnTop: boolean;
  hotkeys: HotkeyBindings;
  comparison: string;
  colors: ColorSettings;
}
