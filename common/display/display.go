package display

type Display interface {
	Listen()
	Render() string
}
