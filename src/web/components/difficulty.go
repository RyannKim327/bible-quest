package web_components

import (
	constellation "constellation/lib"
	"constellation/lib/elements"
	bible_lib "constellation/src/lib"
)

func gameDifficultyButtons(details map[bible_lib.Difficulty]string) []any {
	options := make([]any, 0, len(details))
	for _, option := range []struct {
		level bible_lib.Difficulty
		label string
	}{
		{bible_lib.Easy, "Easy"},
		{bible_lib.Medium, "Medium"},
		{bible_lib.Hard, "Hard"},
	} {
		options = append(options, elements.Button(constellation.WebComponent{
			Class: "difficulty-option",
			Attributes: map[string]string{
				"type":            "button",
				"data-difficulty": string(option.level),
			},
			Children: []any{
				elements.Strong(constellation.WebComponent{Text: option.label}),
				elements.Span(constellation.WebComponent{
					Class: "difficulty-detail",
					Text:  details[option.level],
				}),
			},
		}))
	}
	return options
}
