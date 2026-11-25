package web_components

import (
	constellation "constellation/lib"
	"constellation/lib/elements"
)

type quest struct {
	number   string
	name     string
	category string
	count    string
	level    string
	tone     string
}

var quests = []quest{
	{"01", "In the Beginning", "Genesis", "8 questions", "Easy", "gold"},
	{"02", "Shepherds & Crowns", "History", "12 questions", "Medium", "green"},
	{"03", "The Good News", "Gospels", "16 questions", "Challenge", "blue"},
}

func Quests() string {
	cards := make([]any, 0, len(quests))

	for _, item := range quests {
		cards = append(cards, elements.Article(constellation.WebComponent{
			Class: "quest-card tone-" + item.tone,
			Children: []any{
				elements.Div(constellation.WebComponent{
					Class: "quest-top",
					Children: []any{
						elements.Span(constellation.WebComponent{Class: "quest-number", Text: item.number}),
						elements.Span(constellation.WebComponent{Class: "quest-category", Text: item.category}),
					},
				}),
				elements.H3(constellation.WebComponent{Class: "quest-name", Text: item.name}),
				elements.Div(constellation.WebComponent{
					Class: "quest-meta",
					Children: []any{
						elements.Span(constellation.WebComponent{Text: item.count}),
						elements.Span(constellation.WebComponent{Class: "quest-level", Text: item.level}),
					},
				}),
			},
		}))
	}

	return elements.Section(constellation.WebComponent{
		Class:      "section quests",
		Attributes: map[string]string{"id": "quests"},
		Children: []any{
			elements.Div(constellation.WebComponent{
				Class: "section-heading",
				Children: []any{
					elements.Span(constellation.WebComponent{Class: "eyebrow", Text: "CHOOSE YOUR PATH"}),
					elements.H2(constellation.WebComponent{Class: "section-title", Text: "Quests for every chapter"}),
					elements.P(constellation.WebComponent{
						Class: "section-lead",
						Text:  "Pick a world, answer the questions, and see how much of the story you remember.",
					}),
				},
			}),
			elements.Div(constellation.WebComponent{
				Class:    "quest-grid",
				Children: cards,
			}),
		},
	})
}
