package web

import (
	constellation "constellation/lib"
	web_components "constellation/src/web/components"
)

var GuessTheMistake = constellation.Endpoint{
	Path:  "/guess-the-mistake",
	Title: "Guess the Mistake · Bible Quest",
	Handler: constellation.HandlerSetup{
		Web: guessTheMistakePage,
	},
}

func guessTheMistakePage(_ *constellation.Request) string {
	return web_components.GuessTheMistake()
}
