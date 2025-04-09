package menu

type Item struct {
	Name    string
	Handler func() error
}
