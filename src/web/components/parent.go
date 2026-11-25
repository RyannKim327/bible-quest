package web_components

import (
	constellation "constellation/lib"
	"constellation/lib/elements"
)

func Parent(children []any) string {
	return elements.Div(constellation.WebComponent{
		Class:    "page-shell flex min-h-dvh w-full flex-col",
		Children: children,
	})
}
