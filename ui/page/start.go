package page

import "github.com/lukso-network/tools-lukso-cli/ui/common/context"

type StartPage struct{}

var _ Page = &StartPage{}

func (p *StartPage) Title() string {
	return "Start/Stop"
}

func (p *StartPage) Handle(ctx context.Context) {
}

func (p *StartPage) View() string {
	return "Hello world! From the start page"
}
