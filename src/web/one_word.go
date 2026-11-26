package web

import (
	constellation "constellation/lib"
	web_components "constellation/src/web/components"
)

var OneWord = constellation.Endpoint{
	Path:  "/one-word",
	Title: "One Word · Bible Quest",
	Handler: constellation.HandlerSetup{
		Web: oneWordPage,
	},
}

func oneWordPage(_ *constellation.Request) string {
	return web_components.OneWord()
}
