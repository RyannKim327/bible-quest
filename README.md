# Bible Quest

Bible Quest is a small Go web app for a pixel-style Bible adventure site.
It includes a landing page and a **Guess the Mistake** game where players choose a difficulty, read a Bible passage, and select the words that have been changed.

It is built on [Constellation](https://github.com/RyannKim327/constellation), a Go framework for web applications.
The framework code is included under `lib/`, and the app has no third-party Go module dependencies.
Pages are rendered on the server with Constellation's HTML component helper, styles are embedded into the binary at build time, and the browser game uses vanilla JavaScript.
Tailwind CSS is loaded from a CDN at runtime by the page template.

## Features

- Server-rendered landing page and game page
- Three difficulty levels with different Bible translations:
  - Easy — NIV
  - Medium — NLT
  - Hard — KJV
- Shuffled passage deck that avoids repeating a round back-to-back
- Multi-word mistake detection and phrase corrections
- Immediate feedback, scoring, and answer review in the browser
- Custom 404 and 403 pages
- Built-in development server with hot reload

## Requirements

- Go 1.27.0 or newer
- A modern browser with JavaScript enabled

No third-party Go modules are required.

## Getting started

Run the app:

```bash
go run .
```

Then open <http://localhost:8000>.

To use a different port:

```bash
PORT=3000 go run .
```

Build a binary:

```bash
go build -o bible-quest .
./bible-quest
```

## Development

The repository includes a small hot-reload wrapper:

```bash
go run ./lib/dev
```

It watches `.go`, `.css`, and `.md` files, rebuilds the app, and reloads connected browser tabs. It serves the reloaded app at `http://localhost:$PORT` (default `8000`) while running the actual server on an internal free port.

## Tests and checks

```bash
go test ./...
go vet ./...
```

## Routes

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/` | Landing page |
| `GET` | `/guess-the-mistake` | Game page |
| `GET` | `/api/guess-the-mistake?difficulty=easy` | Get a passage round |
| `POST` | `/api/guess-the-mistake` | Check a submitted answer |

The GET endpoint accepts `difficulty` values `easy`, `medium`, or `hard`, and returns the round ID, reference, translation, and word list.

The POST endpoint accepts JSON like:

```json
{
  "id": "easy/john-3-16",
  "words": [4]
}
```

`words` contains zero-based indexes into the `words` array from the GET response. The response reports whether the answer was correct, the mistake indexes, and the expected corrections. A legacy `word` integer field is also accepted for single-word guesses.

## Project structure

```text
main.go                     Registers routes and embeds static assets
lib/                        Constellation web framework and dev hot-reload server
src/api/                    Guess the Mistake API and round/scoring logic
src/lib/passage.go          Passage texts, translations, and deliberate mistakes
src/web/                    Page endpoints
src/web/components/         Server-rendered UI components and browser game script
static/                     Embedded CSS and logo
```

## License

This project is licensed under the Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License (CC BY-NC-SA 4.0).

You may share and adapt the material for noncommercial purposes, with attribution, and under the same license. See [LICENSE.md](LICENSE.md) for the full legal text.
