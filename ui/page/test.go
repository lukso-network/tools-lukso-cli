package page

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/lukso-network/tools-lukso-cli/ui/common/context"
)

type TestPage struct {
	h int
	w int
}

var _ Page = &TestPage{}

func (p *TestPage) Title() string {
	return "Test"
}

func (p *TestPage) Handle(ctx context.Context) {
	switch msg := ctx.Msg.(type) {
	case tea.WindowSizeMsg:
		p.h = msg.Height
		p.w = msg.Width
	}
}

func (p *TestPage) View() string {
	return fmt.Sprintf("height: %d, width: %d", p.h, p.w)
}
