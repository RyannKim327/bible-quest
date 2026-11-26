package web_components

import (
	"strings"

	constellation "constellation/lib"
	"constellation/lib/elements"
	bible_lib "constellation/src/lib"
)

func characterGuideIntro(prefix string) []any {
	var templates strings.Builder
	templates.WriteString(`<div id="` + prefix + `-character-templates" hidden>`)
	for _, character := range bible_lib.Characters {
		templates.WriteString(`<template data-character-key="` + character.Key + `">`)
		templates.WriteString(CharacterSprite(character.Key, "sprite sprite-guide"))
		templates.WriteString(`</template>`)
	}
	templates.WriteString(`</div>`)

	return []any{
		elements.Div(constellation.WebComponent{
			Class: "guide-scene",
			Children: []any{
				elements.Div(constellation.WebComponent{
					Class:      "guide-sprite",
					Attributes: map[string]string{"id": prefix + "-guide-sprite"},
				}),
				elements.Div(constellation.WebComponent{
					Class:      "guide-bubble",
					Attributes: map[string]string{"id": prefix + "-guide-bubble"},
					Children: []any{
						elements.H3(constellation.WebComponent{
							Class:      "guide-name",
							Attributes: map[string]string{"id": prefix + "-guide-name"},
							Text:       "Your guide",
						}),
						elements.P(constellation.WebComponent{
							Class:      "guide-role",
							Attributes: map[string]string{"id": prefix + "-guide-role"},
						}),
						elements.Div(constellation.WebComponent{
							Class:      "guide-note",
							Attributes: map[string]string{"id": prefix + "-guide-note"},
							Text:       "Choose a difficulty to meet your guide.",
						}),
					},
				}),
			},
		}),
		templates.String(),
	}
}
