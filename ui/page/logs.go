package page

import "github.com/lukso-network/tools-lukso-cli/ui/common/context"

type LogsPage struct{}

var _ Page = &LogsPage{}

func (p *LogsPage) Title() string {
	return "Logs"
}

func (p *LogsPage) Handle(ctx context.Context) {
}

func (p *LogsPage) View() string {
	return "Hello world! From the Logs page"
}
