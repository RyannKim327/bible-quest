package web_components

import (
	constellation "constellation/lib"
	"constellation/lib/elements"
)

func Header() string {
	return elements.Header(constellation.WebComponent{
		Class: "site-header",
		Children: []any{
			elements.A(constellation.WebComponent{
				Class:      "brand",
				Attributes: map[string]string{"href": "#top"},
				Children: []any{
					elements.Span(constellation.WebComponent{
						Class:    "brand-mark",
						Children: []any{PixelCross("pixel-cross")},
					}),
					elements.Span(constellation.WebComponent{
						Class: "brand-name",
						Text:  "Bible Quest",
					}),
				},
			}),
			elements.Nav(constellation.WebComponent{
				Class: "site-nav",
				Children: []any{
					elements.A(constellation.WebComponent{Class: "nav-link", Attributes: map[string]string{"href": "#quests"}, Text: "Quests"}),
					elements.A(constellation.WebComponent{Class: "nav-link", Attributes: map[string]string{"href": "#characters"}, Text: "Characters"}),
					elements.A(constellation.WebComponent{Class: "nav-link", Attributes: map[string]string{"href": "/guess-the-mistake"}, Text: "Guess the mistake"}),
					elements.A(constellation.WebComponent{Class: "nav-link", Attributes: map[string]string{"href": "#how-it-works"}, Text: "How it works"}),
				},
			}),
			elements.A(constellation.WebComponent{
				Class:      "pixel-button header-cta",
				Attributes: map[string]string{"href": "/guess-the-mistake"},
				Text:       "Play now",
			}),
		},
	})
}
