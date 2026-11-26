package web

import (
	constellation "constellation/lib"
	web_components "constellation/src/web/components"
)

var Quiz = constellation.Endpoint{
	Path:  "/quiz",
	Title: "Bible Quiz · Bible Quest",
	Handler: constellation.HandlerSetup{
		Web: quizPage,
	},
}

func quizPage(_ *constellation.Request) string {
	return web_components.Quiz()
}
