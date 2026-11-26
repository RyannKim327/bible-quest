package web_components

import (
	constellation "constellation/lib"
	"constellation/lib/elements"
)

type game struct {
	path   string
	number string
	name   string
	detail string
	count  string
	tone   string
}

var games = []game{
	{"/quiz", "01", "Bible Quiz", "Pick from four choices and prove what you know.", "18 questions", "gold"},
	{"/one-word", "02", "One Word", "One small clue points to a person or a book.", "18 clues", "green"},
	{"/guess-the-mistake", "03", "Guess the Mistake", "Find the words that were swapped in a passage.", "3 translations", "blue"},
}

func Games() string {
	cards := make([]any, 0, len(games))

	for _, item := range games {
		cards = append(cards, elements.A(constellation.WebComponent{
			Class:      "game-card tone-" + item.tone,
			Attributes: map[string]string{"href": item.path},
			Children: []any{
				elements.Div(constellation.WebComponent{
					Class: "game-card-top",
					Children: []any{
						elements.Span(constellation.WebComponent{Class: "game-number", Text: item.number}),
						elements.Span(constellation.WebComponent{Class: "game-play", Text: "PLAY →"}),
					},
				}),
				elements.H3(constellation.WebComponent{Class: "game-name", Text: item.name}),
				elements.P(constellation.WebComponent{Class: "game-detail", Text: item.detail}),
				elements.Div(constellation.WebComponent{
					Class: "game-meta",
					Children: []any{
						elements.Span(constellation.WebComponent{Text: item.count}),
						elements.Span(constellation.WebComponent{Class: "game-level", Text: "3 levels"}),
					},
				}),
			},
		}))
	}

	return elements.Section(constellation.WebComponent{
		Class:      "section games",
		Attributes: map[string]string{"id": "games"},
		Children: []any{
			elements.Div(constellation.WebComponent{
				Class: "section-heading",
				Children: []any{
					elements.Span(constellation.WebComponent{Class: "eyebrow", Text: "ENTER THE ARENA"}),
					elements.H2(constellation.WebComponent{Class: "section-title", Text: "Three ways to play"}),
					elements.P(constellation.WebComponent{
						Class: "section-lead",
						Text:  "Every game has easy, medium, and hard levels. Choose your path and see how much Scripture you can recall.",
					}),
				},
			}),
			elements.Div(constellation.WebComponent{
				Class:    "game-grid",
				Children: cards,
			}),
		},
	})
}
