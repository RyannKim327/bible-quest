package web_components

import (
	constellation "constellation/lib"
	"constellation/lib/elements"
)

const heroBirdScript = `(function () {
  var scene = document.querySelector('.hero-art .scene');
  var template = document.getElementById('pixel-bird-template');
  if (!scene || !template) return;

  var birds = 5 + Math.floor(Math.random() * 3);
  for (var i = 0; i < birds; i++) {
    var bird = document.createElement('span');
    var direction = Math.random() > 0.5 ? 1 : -1;
    var duration = 7 + Math.random() * 6;

    bird.className = 'pixel-bird';
    bird.style.setProperty('--bird-top', (6 + Math.random() * 36).toFixed(1) + '%');
    bird.style.setProperty('--bird-from', direction > 0 ? '-12%' : '112%');
    bird.style.setProperty('--bird-to', direction > 0 ? '112%' : '-12%');
    bird.style.setProperty('--bird-face', direction > 0 ? '1' : '-1');
    bird.style.setProperty('--bird-scale', (0.7 + Math.random() * 0.6).toFixed(2));
    bird.style.setProperty('--bird-duration', duration.toFixed(2) + 's');
    bird.style.setProperty('--bird-delay', (-Math.random() * duration).toFixed(2) + 's');
    bird.style.setProperty('--bird-bob', (Math.random() * 26 - 13).toFixed(1) + 'px');
    bird.appendChild(template.content.cloneNode(true));
    scene.appendChild(bird);
  }
})();`

func Hero() string {
	return Parent([]any{
		elements.Section(constellation.WebComponent{
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
								"<template id=\"pixel-bird-template\">" + PixelBird("sprite bird-sprite") + "</template>",
							},
						}),
					},
				}),
			},
		}),
		elements.Script(constellation.WebComponent{Children: []any{heroBirdScript}}),
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
