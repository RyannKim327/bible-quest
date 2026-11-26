package main

import (
	"embed"

	constellation "constellation/lib"
	"constellation/lib/status"
	"constellation/src/api"
	"constellation/src/web"
)

// StaticFS contains the assets under static, embedded at build time.
//
//go:embed all:static
var StaticFS embed.FS

func main() {
	constellation.SetErrorPages(constellation.ErrorPages{
		NotFound:  status.Status404,
		Forbidden: status.Status403,
	})

	constellation.AddEndpoint(web.Index)
	constellation.AddEndpoint(web.GuessTheMistake)
	constellation.AddEndpoint(web.Quiz)
	constellation.AddEndpoint(web.OneWord)
	constellation.AddEndpoint(api.GuessTheMistakePassage)
	constellation.AddEndpoint(api.GuessTheMistakeCheck)
	constellation.AddEndpoint(api.QuizRound)
	constellation.AddEndpoint(api.QuizCheck)
	constellation.AddEndpoint(api.OneWordRound)
	constellation.AddEndpoint(api.OneWordCheck)

	css, err := StaticFS.ReadFile("static/styles.css")
	if err != nil {
		panic(err)
	}

	constellation.StartServer(css)
}
