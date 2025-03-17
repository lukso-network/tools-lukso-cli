package binding

import "slices"

type keyBinding []string

type appBinding struct {
	Left         keyBinding
	Right        keyBinding
	Up           keyBinding
	Down         keyBinding
	NextPage     keyBinding
	PreviousPage keyBinding
	RightField   keyBinding
	LeftField    keyBinding
	UpField      keyBinding
	DownField    keyBinding
	Enter        keyBinding
	Quit         keyBinding
}

// A default mapping for the app - is overriden during init
var Mapping = appBinding{
	Left:         keyBinding{"left"},
	Right:        keyBinding{"right"},
	Up:           keyBinding{"up"},
	Down:         keyBinding{"down"},
	NextPage:     keyBinding{"shift+right"},
	PreviousPage: keyBinding{"shift+left"},
	RightField:   keyBinding{"ctrl+right"},
	LeftField:    keyBinding{"ctrl+left"},
	UpField:      keyBinding{"ctrl+up"},
	DownField:    keyBinding{"ctrl+down"},
	Enter:        keyBinding{"enter"},
	Quit:         keyBinding{"esc", "ctrl+c"},
}

func init() {
	// override mappings
}

// Matches evaluates parsed keypress against active Mapping
func Matches(binding keyBinding, keypress string) bool {
	return slices.Contains(binding, keypress)
}
