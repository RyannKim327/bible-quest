package web_components

import (
	constellation "constellation/lib"
	"constellation/lib/elements"
)

func Footer() string {
	return elements.Footer(constellation.WebComponent{
		Class: "site-footer",
		Children: []any{
			elements.Span(constellation.WebComponent{
				Class:    "footer-brand",
				Children: []any{PixelCross("pixel-cross pixel-cross-small"), elements.Span(constellation.WebComponent{Text: "Bible Quest"})},
			}),
			elements.Span(constellation.WebComponent{
				Class: "footer-note",
				Text:  "A little pixel adventure through the big story.",
			}),
		},
	})
}
