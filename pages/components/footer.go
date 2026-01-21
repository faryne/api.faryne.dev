package components

import goapp "github.com/maxence-charriere/go-app/v9/pkg/app"

type FooterComponent struct {
	goapp.Compo
}

func NewFooterComponent() *FooterComponent {
	return &FooterComponent{}
}

func (c *FooterComponent) Render() goapp.UI {
	return goapp.Head().Body(
		goapp.P().Text("Copyright Hello World"),
	)
}
