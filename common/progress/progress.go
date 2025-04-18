// package progress provides a simple way for steps to sends its progress to a parent function.
package progress

import (
	bar "github.com/charmbracelet/bubbles/progress"
)

type Progress interface {
	// Visible returns whether progress should be visible or not.
	Visible() bool
	// Show toggles the progress bar on.
	Show()
	// Hode toggles the progress bar off.
	Hide()

	// Move moves to progress by count steps.
	Move(count float64)
	// Set resets the progress and sets a new goal.
	Set(goal float64)

	// Render calls the bubbletea.View() method on the underlying tea Progress.
	Render() string
}

type progress struct {
	done    float64
	goal    float64
	bar     bar.Model
	visible bool
}

// NewProgress returns a new progress with a goal of 0. Goal needs to be initialized before every progressable action.
func NewProgress() Progress {
	return &progress{
		done:    0,
		goal:    1,
		bar:     bar.New(bar.WithDefaultGradient()),
		visible: false,
	}
}

func (p *progress) Move(count float64) {
	p.done += count
}

func (p *progress) Set(goal float64) {
	p.done = 0
	p.goal = goal
}

func (p *progress) Render() string {
	return p.bar.ViewAs(p.done / p.goal)
}

func (p *progress) Show() {
	p.visible = true
}

func (p *progress) Hide() {
	p.visible = false
}

func (p *progress) Visible() bool {
	return p.visible
}
