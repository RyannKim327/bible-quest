package web_components

import (
	"fmt"

	constellation "constellation/lib"
	"constellation/lib/elements"
)

func Characters() string {
	cards := make([]any, 0, len(pixelCharacters))

	for _, character := range pixelCharacters {
		cards = append(cards, elements.Article(constellation.WebComponent{
			Class: "character-card",
			Attributes: map[string]string{
				"style": fmt.Sprintf("--accent:%s", character.Style.trim),
			},
			Children: []any{
				elements.Div(constellation.WebComponent{
					Class:    "character-portrait",
					Children: []any{CharacterSprite(character.Key, "sprite sprite-card")},
				}),
				elements.Div(constellation.WebComponent{
					Class: "character-info",
					Children: []any{
						elements.Span(constellation.WebComponent{Class: "character-role", Text: character.Role}),
						elements.H3(constellation.WebComponent{Class: "character-name", Text: character.Name}),
						elements.Span(constellation.WebComponent{Class: "character-verse", Text: character.Verse}),
					},
				}),
			},
		}))
	}

	return elements.Section(constellation.WebComponent{
		Class:      "section characters",
		Attributes: map[string]string{"id": "characters"},
		Children: []any{
			elements.Div(constellation.WebComponent{
				Class: "section-heading",
				Children: []any{
					elements.Span(constellation.WebComponent{Class: "eyebrow", Text: "YOUR PARTY"}),
					elements.H2(constellation.WebComponent{Class: "section-title", Text: "Meet your companions"}),
					elements.P(constellation.WebComponent{
						Class: "section-lead",
						Text:  "Eight friends from the Word join you on the journey. Each one brings a story worth remembering.",
					}),
				},
			}),
			elements.Div(constellation.WebComponent{
				Class:    "character-grid",
				Children: cards,
			}),
		},
	})
}
