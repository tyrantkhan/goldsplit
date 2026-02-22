package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"goldsplit/internal/hotkey"
	"goldsplit/internal/persist"
	"goldsplit/internal/split"
	"goldsplit/internal/timer"
)

// App is the Wails-bound application struct.
type App struct {
	ctx      context.Context
	engine   *timer.Engine
	store    *persist.Store
	hk       *hotkey.Manager
	tmpl     *split.Template
	attempts *split.Attempts
	settings persist.Settings
	version  string
}

// NewApp creates a new App instance.
func NewApp(version string) *App {
	return &App{version: version}
}

// GetAppInfo returns application metadata for the frontend.
func (a *App) GetAppInfo() map[string]any {
	return map[string]any{
		"version": a.version,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	baseDir, err := persist.DefaultBaseDir()
	if err != nil {
		fmt.Printf("Warning: could not determine data directory: %v\n", err)

		return
	}

	store, err := persist.NewStore(baseDir)
	if err != nil {
		fmt.Printf("Warning: could not create store: %v\n", err)

		return
	}

	a.store = store

	settings, err := store.LoadSettings()
	if err != nil {
		fmt.Printf("Warning: could not load settings: %v\n", err)
	}

	a.settings = settings
	runtime.WindowSetAlwaysOnTop(ctx, a.settings.AlwaysOnTop)

	a.engine = timer.New(nil, a.onTick, a.onStateChange)

	a.hk = hotkey.NewManager(a.onHotkey)
	if err := a.hk.Start(); err != nil {
		fmt.Printf("Warning: could not start hotkey manager: %v\n", err)
	}
}

func (a *App) shutdown(ctx context.Context) {
	if a.hk != nil {
		a.hk.Stop()
	}

	if a.engine != nil {
		a.engine.Reset()
	}
}

func (a *App) onTick(data timer.TickData) {
	runtime.EventsEmit(a.ctx, "timer:tick", data)
}

func (a *App) onStateChange(state timer.State) {
	runtime.EventsEmit(a.ctx, "timer:state", state.String())
}

func (a *App) onHotkey(action hotkey.Action) {
	switch action {
	case hotkey.ActionStartSplit:
		a.StartSplit()
	case hotkey.ActionPause:
		a.TogglePause()
	case hotkey.ActionReset:
		a.Reset()
	case hotkey.ActionUndoSplit:
		a.UndoSplit()
	case hotkey.ActionSkipSplit:
		a.SkipSplit()
	}
}

// StartSplit is the smart start/split action.
// If idle, starts the timer. If running, splits.
func (a *App) StartSplit() {
	state := a.engine.CurrentState()

	switch state {
	case timer.Idle:
		a.engine.Start()
	case timer.Running:
		a.engine.Split()
		a.emitDeltas()
		a.checkRunCompletion()
	case timer.Paused, timer.Finished:
		// No-op for these states.
	}
}

// TogglePause pauses if running, resumes if paused.
func (a *App) TogglePause() {
	state := a.engine.CurrentState()

	switch state {
	case timer.Running:
		a.engine.Pause()
	case timer.Paused:
		a.engine.Resume()
	case timer.Idle, timer.Finished:
		// No-op for these states.
	}
}

// Reset resets the timer. If a run was in progress (not finished), save the incomplete attempt.
func (a *App) Reset() {
	state := a.engine.CurrentState()
	if state == timer.Idle {
		return
	}

	// Only save an incomplete attempt if the run was in progress.
	// Finished runs are already saved by checkRunCompletion.
	if state != timer.Finished {
		a.saveAttempt(false)
	}

	a.engine.Reset()
}

// DiscardAttempt removes the last completed attempt, recalculates PB, and resets.
func (a *App) DiscardAttempt() {
	if a.engine.CurrentState() != timer.Finished {
		return
	}

	if a.attempts != nil && len(a.attempts.History) > 0 {
		a.attempts.History = a.attempts.History[:len(a.attempts.History)-1]
		a.attempts.AttemptCount--
		split.RecalculatePersonalBest(a.attempts)

		if a.store != nil {
			if err := a.store.SaveAttempts(a.attempts); err != nil {
				fmt.Printf("Warning: could not save attempts: %v\n", err)
			}
		}

		runtime.EventsEmit(a.ctx, "attempts:updated", a.getAttemptsData())
	}

	a.engine.Reset()
}

// UndoSplit undoes the last split.
func (a *App) UndoSplit() {
	a.engine.UndoSplit()
	a.emitDeltas()
}

// SkipSplit skips the current segment.
func (a *App) SkipSplit() {
	a.engine.SkipSplit()
	a.emitDeltas()
	a.checkRunCompletion()
}

func (a *App) emitDeltas() {
	if a.attempts == nil {
		return
	}

	d := split.ComputeSplitDeltas(a.attempts, a.engine.SplitTimesMS(), a.settings.Comparison)
	runtime.EventsEmit(a.ctx, "deltas:updated", d)
}

func (a *App) checkRunCompletion() {
	if a.engine.CurrentState() == timer.Finished {
		splits := a.engine.SplitTimesMS()
		completed := len(splits) > 0 && splits[len(splits)-1] > 0
		a.saveAttempt(completed)
	}
}

func (a *App) saveAttempt(completed bool) {
	if a.attempts == nil || a.store == nil {
		return
	}

	splits := a.engine.SplitTimesMS()
	a.attempts.AddAttempt(splits, completed)

	if completed {
		split.UpdatePersonalBest(a.attempts, splits)
	}

	if err := a.store.SaveAttempts(a.attempts); err != nil {
		fmt.Printf("Warning: could not save attempts: %v\n", err)
	}

	runtime.EventsEmit(a.ctx, "attempts:updated", a.getAttemptsData())
}

// CreateTemplate creates a new template and returns its data.
func (a *App) CreateTemplate(name string, segmentNames []string) map[string]any {
	id := uuid.New().String()
	tmpl := split.NewTemplate(id, name, segmentNames)
	a.tmpl = tmpl

	if a.store != nil {
		if err := a.store.SaveTemplate(tmpl); err != nil {
			fmt.Printf("Warning: could not save template: %v\n", err)
		}
	}

	return a.getTemplateData()
}

// LoadTemplate loads a template by ID and sets it as the current template.
func (a *App) LoadTemplate(id string) map[string]any {
	if a.store == nil {
		return nil
	}

	tmpl, err := a.store.LoadTemplate(id)
	if err != nil {
		fmt.Printf("Warning: could not load template: %v\n", err)

		return nil
	}

	a.tmpl = tmpl

	return a.getTemplateData()
}

// ListTemplates returns all saved templates.
func (a *App) ListTemplates() []persist.TemplateSummary {
	if a.store == nil {
		return nil
	}

	summaries, err := a.store.ListTemplates()
	if err != nil {
		fmt.Printf("Warning: could not list templates: %v\n", err)

		return nil
	}

	return summaries
}

// DeleteTemplate deletes a template and its associated attempts.
func (a *App) DeleteTemplate(id string) bool {
	if a.store == nil {
		return false
	}

	if err := a.store.DeleteTemplate(id); err != nil {
		fmt.Printf("Warning: could not delete template: %v\n", err)

		return false
	}

	if a.tmpl != nil && a.tmpl.ID == id {
		a.tmpl = nil
		a.attempts = nil
	}

	return true
}

// UpdateTemplate updates a template's name and segment names.
func (a *App) UpdateTemplate(id, name string, segmentNames []string) map[string]any {
	if a.store == nil {
		return nil
	}

	tmpl, err := a.store.LoadTemplate(id)
	if err != nil {
		fmt.Printf("Warning: could not load template: %v\n", err)

		return nil
	}

	tmpl.Name = name
	tmpl.SegmentNames = segmentNames
	tmpl.UpdatedAt = time.Now()

	if err := a.store.SaveTemplate(tmpl); err != nil {
		fmt.Printf("Warning: could not save template: %v\n", err)

		return nil
	}

	if a.tmpl != nil && a.tmpl.ID == id {
		a.tmpl = tmpl
	}

	return map[string]any{
		"id":           tmpl.ID,
		"name":         tmpl.Name,
		"segmentNames": tmpl.SegmentNames,
	}
}

// CreateAttempts creates a new attempts entry for the current template.
func (a *App) CreateAttempts(templateID, name, categoryName string) map[string]any {
	if a.store == nil {
		return nil
	}

	tmpl, err := a.store.LoadTemplate(templateID)
	if err != nil {
		fmt.Printf("Warning: could not load template: %v\n", err)

		return nil
	}

	id := uuid.New().String()
	att := split.NewAttempts(id, templateID, name, categoryName, tmpl.SegmentNames)
	a.attempts = att
	a.engine.SetSegments(att.SegmentNames())

	if err := a.store.SaveAttempts(att); err != nil {
		fmt.Printf("Warning: could not save attempts: %v\n", err)
	}

	return a.getAttemptsData()
}

// LoadAttempts loads attempts by ID and sets them as active.
func (a *App) LoadAttempts(id string) map[string]any {
	if a.store == nil {
		return nil
	}

	att, err := a.store.LoadAttempts(id)
	if err != nil {
		fmt.Printf("Warning: could not load attempts: %v\n", err)

		return nil
	}

	a.attempts = att
	a.engine.SetSegments(att.SegmentNames())

	return a.getAttemptsData()
}

// ListAttemptsForTemplate returns all attempts for a given template.
func (a *App) ListAttemptsForTemplate(templateID string) []persist.AttemptsSummary {
	if a.store == nil {
		return nil
	}

	summaries, err := a.store.ListAttemptsForTemplate(templateID)
	if err != nil {
		fmt.Printf("Warning: could not list attempts: %v\n", err)

		return nil
	}

	return summaries
}

// DeleteAttempts deletes an attempts entry by ID.
func (a *App) DeleteAttempts(id string) bool {
	if a.store == nil {
		return false
	}

	if err := a.store.DeleteAttempts(id); err != nil {
		fmt.Printf("Warning: could not delete attempts: %v\n", err)

		return false
	}

	if a.attempts != nil && a.attempts.ID == id {
		a.attempts = nil
	}

	return true
}

// UpdateCategoryName renames a category on an attempts entry.
func (a *App) UpdateCategoryName(attemptsID, newName string) map[string]any {
	if a.store == nil {
		return nil
	}

	att, err := a.store.LoadAttempts(attemptsID)
	if err != nil {
		fmt.Printf("Warning: could not load attempts: %v\n", err)

		return nil
	}

	att.CategoryName = newName
	att.UpdatedAt = time.Now()

	if err := a.store.SaveAttempts(att); err != nil {
		fmt.Printf("Warning: could not save attempts: %v\n", err)

		return nil
	}

	if a.attempts != nil && a.attempts.ID == attemptsID {
		a.attempts = att
	}

	return a.buildAttemptsData(att)
}

// DeleteSingleAttempt removes a single attempt from an attempts entry.
func (a *App) DeleteSingleAttempt(attemptsID string, attemptID int) map[string]any {
	if a.store == nil {
		return nil
	}

	att, err := a.store.LoadAttempts(attemptsID)
	if err != nil {
		fmt.Printf("Warning: could not load attempts: %v\n", err)

		return nil
	}

	if !att.DeleteAttempt(attemptID) {
		return nil
	}

	if err := a.store.SaveAttempts(att); err != nil {
		fmt.Printf("Warning: could not save attempts: %v\n", err)

		return nil
	}

	if a.attempts != nil && a.attempts.ID == attemptsID {
		a.attempts = att
	}

	return a.buildAttemptsData(att)
}

// EditAttemptSplits edits the split times of a single attempt.
func (a *App) EditAttemptSplits(attemptsID string, attemptID int, newSplits []int64) map[string]any {
	if a.store == nil {
		return nil
	}

	att, err := a.store.LoadAttempts(attemptsID)
	if err != nil {
		fmt.Printf("Warning: could not load attempts: %v\n", err)

		return nil
	}

	if !att.EditAttemptSplits(attemptID, newSplits) {
		return nil
	}

	if err := a.store.SaveAttempts(att); err != nil {
		fmt.Printf("Warning: could not save attempts: %v\n", err)

		return nil
	}

	if a.attempts != nil && a.attempts.ID == attemptsID {
		a.attempts = att
	}

	return a.buildAttemptsData(att)
}

// GetAttemptHistory returns the full attempt history for an attempts entry.
func (a *App) GetAttemptHistory(attemptsID string) []split.Attempt {
	if a.store == nil {
		return nil
	}

	att, err := a.store.LoadAttempts(attemptsID)
	if err != nil {
		fmt.Printf("Warning: could not load attempts: %v\n", err)

		return nil
	}

	return att.History
}

// GetCurrentTemplate returns the currently selected template data.
func (a *App) GetCurrentTemplate() map[string]any {
	return a.getTemplateData()
}

// GetCurrentAttempts returns the currently active attempts data.
func (a *App) GetCurrentAttempts() map[string]any {
	return a.getAttemptsData()
}

// GetDeltas returns the current deltas compared to PB.
func (a *App) GetDeltas() []split.Delta {
	if a.attempts == nil {
		return nil
	}

	return split.ComputeSplitDeltas(a.attempts, a.engine.SplitTimesMS(), a.settings.Comparison)
}

// GetSettings returns the current application settings.
func (a *App) GetSettings() persist.Settings {
	return a.settings
}

// UpdateSettings saves new settings, applies side effects, and emits an event.
func (a *App) UpdateSettings(settings persist.Settings) bool {
	if a.store != nil {
		if err := a.store.SaveSettings(settings); err != nil {
			fmt.Printf("Warning: could not save settings: %v\n", err)

			return false
		}
	}

	a.settings = settings
	runtime.WindowSetAlwaysOnTop(a.ctx, settings.AlwaysOnTop)
	runtime.EventsEmit(a.ctx, "settings:updated", settings)

	return true
}

func (a *App) getTemplateData() map[string]any {
	if a.tmpl == nil {
		return nil
	}

	return map[string]any{
		"id":           a.tmpl.ID,
		"name":         a.tmpl.Name,
		"segmentNames": a.tmpl.SegmentNames,
	}
}

func (a *App) getAttemptsData() map[string]any {
	if a.attempts == nil {
		return nil
	}

	return a.buildAttemptsData(a.attempts)
}

func (a *App) buildAttemptsData(att *split.Attempts) map[string]any {
	segments := make([]map[string]any, len(att.Segments))
	for i, s := range att.Segments {
		segments[i] = map[string]any{
			"name":           s.Name,
			"personalBestMs": s.PersonalBestMS,
			"bestSegmentMs":  s.BestSegmentMS,
		}
	}

	return map[string]any{
		"id":           att.ID,
		"templateId":   att.TemplateID,
		"name":         att.Name,
		"categoryName": att.CategoryName,
		"segments":     segments,
		"attemptCount": att.AttemptCount,
	}
}
