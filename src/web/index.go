package web

import (
	constellation "constellation/lib"
	web_components "constellation/src/web/components"
)

var Index = constellation.Endpoint{
	Path:  "/",
	Title: "Bible Quest · A Pixel Bible Adventure",
	Handler: constellation.HandlerSetup{
		Web: index,
	},
}

func index(_ *constellation.Request) string {
	return web_components.Parent([]any{
		web_components.Header(),
		web_components.Hero(),
		web_components.Quests(),
		web_components.Characters(),
		web_components.HowItWorks(),
		web_components.Footer(),
	})
}
