// package progress provides a simple way for steps to sends its progress to a parent function.
package progress

import bar "github.com/charmbracelet/bubbles/progress"

type Progress interface {
	// Move moves to progress by count steps.
	Move(count float64)
	// Set resets the progress and sets a new goal.
	Set(goal float64)

	// Render calls the bubbletea.View() method
	Render() string
}

type progress struct {
	done float64
	goal float64
	bar  bar.Model
}

// NewProgress returns a new progress with a goal of 0. Goal needs to be initialized before every progressable action.
func NewProgress() Progress {
	return &progress{
		done: 0,
		goal: 1,
		bar:  bar.New(bar.WithDefaultGradient()),
	}
}

func (p *progress) Move(count float64) {
	p.done += count
	p.bar.SetPercent(p.done / p.goal)
}

func (p *progress) Set(goal float64) {
	p.done = 0
	p.goal = goal
	p.bar.SetPercent(0)
}

func (p *progress) Render() string {
	return p.bar.View()
}
