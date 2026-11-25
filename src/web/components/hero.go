package web_components

import (
	constellation "constellation/lib"
	"constellation/lib/elements"
)

func Hero() string {
	return elements.Section(constellation.WebComponent{
		Class:      "hero",
		Attributes: map[string]string{"id": "top"},
		Children: []any{
			elements.Div(constellation.WebComponent{
				Class: "hero-copy",
				Children: []any{
					elements.Span(constellation.WebComponent{
						Class: "eyebrow",
						Text:  "A PIXEL BIBLE QUIZ",
					}),
					elements.H1(constellation.WebComponent{
						Class: "hero-title",
						Children: []any{
							elements.Span(constellation.WebComponent{Text: "Quest through the "}),
							elements.Span(constellation.WebComponent{Class: "title-highlight", Text: "Word."}),
						},
					}),
					elements.P(constellation.WebComponent{
						Class: "hero-lead",
						Text:  "Gather your companions, travel through Scripture, and prove what you know one question at a time.",
					}),
					elements.Div(constellation.WebComponent{
						Class: "hero-actions",
						Children: []any{
							elements.A(constellation.WebComponent{
								Class:      "pixel-button pixel-button-primary",
								Attributes: map[string]string{"href": "#quests"},
								Text:       "Start a quest",
							}),
							elements.A(constellation.WebComponent{
								Class:      "pixel-button pixel-button-secondary",
								Attributes: map[string]string{"href": "#characters"},
								Text:       "Meet the crew",
							}),
						},
					}),
					elements.Div(constellation.WebComponent{
						Class: "hero-stats",
						Children: []any{
							stat("120+", "Questions"),
							stat("08", "Biblical worlds"),
							stat("16", "Pixel friends"),
						},
					}),
				},
			}),
			elements.Div(constellation.WebComponent{
				Class: "hero-art",
				Children: []any{
					elements.Div(constellation.WebComponent{
						Class: "scene",
						Children: []any{
							elements.Div(constellation.WebComponent{Class: "scene-sun"}),
							elements.Div(constellation.WebComponent{Class: "scene-cloud cloud-one"}),
							elements.Div(constellation.WebComponent{Class: "scene-cloud cloud-two"}),
							elements.Div(constellation.WebComponent{
								Class: "scene-party",
								Children: []any{
									elements.Div(constellation.WebComponent{Class: "party-member member-left", Children: []any{CharacterSprite("peter", "sprite sprite-sm")}}),
									elements.Div(constellation.WebComponent{Class: "party-member member-center", Children: []any{CharacterSprite("jesus", "sprite sprite-lg")}}),
									elements.Div(constellation.WebComponent{Class: "party-member member-right", Children: []any{CharacterSprite("mary", "sprite sprite-sm")}}),
								},
							}),
							elements.Div(constellation.WebComponent{Class: "scene-ground"}),
							elements.Div(constellation.WebComponent{Class: "scene-tag", Text: "CHAPTER 01 · THE BEGINNING"}),
						},
					}),
				},
			}),
		},
	})
}

func stat(value, label string) string {
	return elements.Div(constellation.WebComponent{
		Class: "stat",
		Children: []any{
			elements.Span(constellation.WebComponent{Class: "stat-value", Text: value}),
			elements.Span(constellation.WebComponent{Class: "stat-label", Text: label}),
		},
	})
}
