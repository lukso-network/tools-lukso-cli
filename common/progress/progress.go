// package progress provides a simple way for steps to sends its progress to a parent function.
package progress

import (
	"time"

	bar "github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lukso-network/tools-lukso-cli/ui/style/color"
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
	done     float64
	goal     float64
	bar      bar.Model
	spin     spinner.Model
	ch       chan tea.Msg
	tickerCh chan bool

	visible bool
}

// NewProgress returns a new progress with a goal of 0. Goal needs to be initialized before every progressable action.
func NewProgress(ch chan tea.Msg) Progress {
	p := &progress{
		done:     0,
		goal:     1,
		bar:      bar.New(bar.WithSolidFill(string(color.LuksoPink))),
		spin:     spinner.New(spinner.WithSpinner(spinner.Line)),
		visible:  false,
		ch:       ch,
		tickerCh: make(chan bool),
	}

	return p
}

func (p *progress) Move(count float64) {
	p.done += count
}

func (p *progress) Set(goal float64) {
	p.done = 0
	p.goal = goal
}

func (p *progress) Render() string {
	return lipgloss.JoinHorizontal(lipgloss.Center, p.spin.View()+" ", p.bar.ViewAs(p.done/p.goal))
}

func (p *progress) Show() {
	go p.runTicker()
	p.visible = true
}

func (p *progress) Hide() {
	p.visible = false
	p.stopTicker()
}

func (p *progress) Visible() bool {
	return p.visible
}

func (p *progress) runTicker() {
	t := time.NewTicker(time.Millisecond * 100)

loop:
	for {
		select {
		case <-t.C:
			p.spin, _ = p.spin.Update(spinner.TickMsg{})
			p.ch <- tea.Cmd(p.spin.Tick)
		case <-p.tickerCh:
			break loop
		}
	}
}

func (p *progress) stopTicker() {
	p.tickerCh <- true
}
