package components

import goapp "github.com/maxence-charriere/go-app/v9/pkg/app"

type HeaderComponent struct {
	goapp.Compo
}

func NewHeaderComponent() *HeaderComponent {
	return &HeaderComponent{}
}

func (c *HeaderComponent) Render() goapp.UI {
	return goapp.Head().Body(
		goapp.P().Text("Hello World"),
	)
}
