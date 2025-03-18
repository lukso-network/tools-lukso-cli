package page

import "github.com/lukso-network/tools-lukso-cli/ui/common/context"

type InstallPage struct{}

var _ Page = &InstallPage{}

func (p *InstallPage) Title() string {
	return "Install"
}

func (p *InstallPage) Handle(ctx context.Context) {
}

func (p *InstallPage) View() string {
	return "Hello world! From the Install page"
}
