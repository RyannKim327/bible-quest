package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode"

	constellation "constellation/lib"
	"constellation/lib/middleware"
	bible_lib "constellation/src/lib"
)

var OneWordRound = constellation.Endpoint{
	Path:   "/api/one-word",
	Method: constellation.GET,
	Handler: constellation.HandlerSetup{
		Api: oneWordRoundHandler,
	},
}

var OneWordCheck = constellation.Endpoint{
	Path:   "/api/one-word",
	Method: constellation.POST,
	Handler: constellation.HandlerSetup{
		Api: oneWordCheckHandler,
	},
}

type oneWordRound struct {
	ID         string
	Difficulty difficulty
	Word       string
	Kind       string
	Answer     string
	Aliases    []string
	Reference  string
}

func buildOneWordRounds(level difficulty) []oneWordRound {
	rounds := make([]oneWordRound, 0)
	for _, clue := range bible_lib.OneWordClues {
		if clue.Difficulty != level {
			continue
		}

		aliases := append([]string(nil), clue.Aliases...)
		rounds = append(rounds, oneWordRound{
			ID:         clue.ID,
			Difficulty: level,
			Word:       clue.Word,
			Kind:       clue.Kind,
			Answer:     clue.Answer,
			Aliases:    aliases,
			Reference:  clue.Reference,
		})
	}
	return rounds
}

var oneWordByDifficulty = map[difficulty][]oneWordRound{
	difficultyEasy:   buildOneWordRounds(difficultyEasy),
	difficultyMedium: buildOneWordRounds(difficultyMedium),
	difficultyHard:   buildOneWordRounds(difficultyHard),
}

var oneWordByID = func() map[string]oneWordRound {
	index := make(map[string]oneWordRound)
	for level, set := range oneWordByDifficulty {
		for _, round := range set {
			index[roundKey(level, round.ID)] = round
		}
	}
	return index
}()

var oneWordDecks = newDifficultyDecks(oneWordByDifficulty)

func normalizeAnswer(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.TrimFunc(value, unicode.IsPunct)
	return strings.Join(strings.Fields(value), " ")
}

func oneWordRoundHandler(r *constellation.Request, w constellation.Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	level, ok := parseDifficulty(r.URL.Query().Get("difficulty"))
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		middleware.SendApi(w, map[string]any{
			"message": "Choose a difficulty before starting the game.",
		})
		return
	}

	round, ok := oneWordDecks.next(level)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		middleware.SendApi(w, map[string]any{
			"message": "There are no clues for that difficulty yet.",
		})
		return
	}

	middleware.SendApi(w, map[string]any{
		"id":         roundKey(round.Difficulty, round.ID),
		"difficulty": string(round.Difficulty),
		"word":       round.Word,
		"kind":       round.Kind,
		"character":  nextCharacter(round.Difficulty),
	})
}

type oneWordPayload struct {
	ID     string `json:"id"`
	Answer string `json:"answer"`
}

func oneWordCheckHandler(r *constellation.Request, w constellation.Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var payload oneWordPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		middleware.SendApi(w, map[string]any{
			"correct": false,
			"message": "That answer could not be read. Try again.",
		})
		return
	}

	round, ok := oneWordByID[payload.ID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		middleware.SendApi(w, map[string]any{
			"correct": false,
			"message": "That clue is no longer in the game.",
		})
		return
	}

	guess := normalizeAnswer(payload.Answer)
	if guess == "" {
		w.WriteHeader(http.StatusBadRequest)
		middleware.SendApi(w, map[string]any{
			"correct": false,
			"message": "Type a name before submitting.",
		})
		return
	}

	correct := guess == normalizeAnswer(round.Answer)
	if !correct {
		for _, alias := range round.Aliases {
			if guess == normalizeAnswer(alias) {
				correct = true
				break
			}
		}
	}

	message := "Not quite. The answer is " + round.Answer + " — " + round.Reference + "."
	if correct {
		message = "You got it! " + round.Answer + " — " + round.Reference + "."
	}

	middleware.SendApi(w, map[string]any{
		"correct":   correct,
		"message":   message,
		"answer":    round.Answer,
		"kind":      round.Kind,
		"reference": round.Reference,
	})
}
