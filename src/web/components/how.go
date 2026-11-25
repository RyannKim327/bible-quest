package web_components

import (
	constellation "constellation/lib"
	"constellation/lib/elements"
)

type step struct {
	number string
	title  string
	text   string
}

var steps = []step{
	{"01", "Pick a quest", "Choose a chapter of the Bible and set your difficulty."},
	{"02", "Answer the call", "Read each question and type the answer, tile by tile."},
	{"03", "Grow your party", "Collect companions as your score climbs."},
}

func HowItWorks() string {
	items := make([]any, 0, len(steps))

	for _, item := range steps {
		items = append(items, elements.Li(constellation.WebComponent{
			Class: "step",
			Children: []any{
				elements.Span(constellation.WebComponent{Class: "step-number", Text: item.number}),
				elements.H3(constellation.WebComponent{Class: "step-title", Text: item.title}),
				elements.P(constellation.WebComponent{Class: "step-text", Text: item.text}),
			},
		}))
	}

	return elements.Section(constellation.WebComponent{
		Class:      "section how",
		Attributes: map[string]string{"id": "how-it-works"},
		Children: []any{
			elements.Div(constellation.WebComponent{
				Class: "how-panel",
				Children: []any{
					elements.Div(constellation.WebComponent{
						Class: "how-intro",
						Children: []any{
							elements.Span(constellation.WebComponent{Class: "eyebrow", Text: "THE ADVENTURE LOOP"}),
							elements.H2(constellation.WebComponent{Class: "section-title", Text: "Simple to start, hard to master"}),
							elements.P(constellation.WebComponent{
								Class: "section-lead",
								Text:  "A familiar quiz with a little more story. Every right answer moves you forward.",
							}),
							elements.A(constellation.WebComponent{
								Class:      "pixel-button pixel-button-primary",
								Attributes: map[string]string{"href": "#quests"},
								Text:       "Take the first quest",
							}),
						},
					}),
					elements.Ol(constellation.WebComponent{
						Class:    "step-list",
						Children: items,
					}),
				},
			}),
		},
	})
}
