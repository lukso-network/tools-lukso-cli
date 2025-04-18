package progress

type stubProgress struct{}

func NewStubProgress() Progress {
	return stubProgress{}
}

var _ Progress = stubProgress{}

func (p stubProgress) Move(count float64) {}
func (p stubProgress) Set(goal float64)   {}
func (p stubProgress) Render() string     { return "" }
func (p stubProgress) Show()              {}
func (p stubProgress) Hide()              {}
func (p stubProgress) Visible() bool      { return false }
