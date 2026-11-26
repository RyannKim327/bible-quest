# Bible Quest

Bible Quest is a small Go web app for a pixel-style Bible adventure site.
It includes a landing page and three games with difficulty levels:

- **Bible Quiz** — answer a multiple-choice question.
- **One Word** — identify the Bible character or book behind a single word.
- **Guess the Mistake** — find the changed words in a Bible passage.

It is built on [Constellation](https://github.com/RyannKim327/constellation), a Go framework for web applications.
The framework code is included under `lib/`, and the app has no third-party Go module dependencies.
Pages are rendered on the server with Constellation's HTML component helper, styles are embedded into the binary at build time, and the browser game uses vanilla JavaScript.
Tailwind CSS is loaded from a CDN at runtime by the page template.

## Features

- Server-rendered landing page with links to all three games
- Three difficulty levels shared by every game: easy, medium, and hard
- A randomized Bible character guide in every game whose side notes point back to the verse
- Bible Quiz with four randomly ordered choices, answer review, and the Bible reference after each question
- One Word clues that ask for a character or a book, with case-insensitive answers and aliases
- Guess the Mistake with translations that change by difficulty:
  - Easy — NIV
  - Medium — NLT
  - Hard — KJV
- Shuffled decks that avoid repeating a round back-to-back
- Multi-word mistake detection and phrase corrections
- Immediate feedback and scoring in the browser
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
| `GET` | `/quiz` | Bible Quiz page |
| `GET` | `/one-word` | One Word page |
| `GET` | `/guess-the-mistake` | Guess the Mistake page |
| `GET` | `/api/quiz?difficulty=easy` | Get a quiz question |
| `POST` | `/api/quiz` | Check a quiz answer |
| `GET` | `/api/one-word?difficulty=easy` | Get a One Word clue |
| `POST` | `/api/one-word` | Check a One Word answer |
| `GET` | `/api/guess-the-mistake?difficulty=easy` | Get a passage round |
| `POST` | `/api/guess-the-mistake` | Check a submitted answer |

Every GET endpoint accepts `difficulty` values `easy`, `medium`, or `hard` and returns an opaque round `id`. POST endpoints use that ID so the difficulty cannot be forged by the client.

Each GET response also includes a `character` with a `key`, `name`, `role`, and `verse`. The games use it to show a randomized guide on the side and to deliver the answer note as if the character were speaking.

Quiz choices are shuffled on the server for every question. POST body:

```json
{
  "id": "easy/ark-builder",
  "answer": "Noah"
}
```

`answer` is the text of the selected choice, compared case-insensitively. The response reports whether the answer was correct, the right answer text, and the Bible reference.

One Word POST body:

```json
{
  "id": "easy/ark-noah",
  "answer": "Noah"
}
```

Answers are case-insensitive. The response reports correctness, the expected answer, its kind, and the Bible reference.

Guess the Mistake POST body:

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
src/api/                    Game APIs, round building, difficulty helpers, and decks
src/lib/passage.go          Passage texts, translations, and deliberate mistakes
src/lib/characters.go       Bible character guide pools by difficulty
src/lib/quiz.go             Quiz questions and choices
src/lib/one_word.go         One Word clues and answers
src/web/                    Page endpoints
src/web/components/         Server-rendered UI components and browser game scripts
src/web/components/character.go  Shared character guide markup and sprite templates
static/                     Embedded CSS and logo
```

## License

This project is licensed under the Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License (CC BY-NC-SA 4.0).

You may share and adapt the material for noncommercial purposes, with attribution, and under the same license. See [LICENSE.md](LICENSE.md) for the full legal text.
