package page

import (
	"github.com/lukso-network/tools-lukso-cli/ui/common/context"
)

// Page represents what a ui page needs in order to correctly take input and render back output.
type Page interface {
	Title() string
	Handle(ctx context.Context)
	View() string
}
