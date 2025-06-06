package display

import "github.com/lukso-network/tools-lukso-cli/common/display/input"

type Display interface {
	Listen()
	Render() string
	AddInputHandler(h input.Handler)
	InputHandler() input.Handler
}
