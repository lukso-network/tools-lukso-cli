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
	Enter        keyBinding
	Quit         keyBinding
}

// A default mapping for the app - is overriden during init
var Mapping = appBinding{
	Left:         []string{"left"},
	Right:        []string{"right"},
	Up:           []string{"up"},
	Down:         []string{"down"},
	NextPage:     []string{"shift+right"},
	PreviousPage: []string{"shift+left"},
	Enter:        []string{"enter"},
	Quit:         []string{"esc", "ctrl+c"},
}

func init() {
	// override mappings
}

// Matches evaluates parsed keypress against active Mapping
func Matches(binding keyBinding, keypress string) bool {
	return slices.Contains(binding, keypress)
}
