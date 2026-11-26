package api

import (
	"encoding/json"
	"math/rand/v2"
	"net/http"

	constellation "constellation/lib"
	"constellation/lib/middleware"
	bible_lib "constellation/src/lib"
)

var QuizRound = constellation.Endpoint{
	Path:   "/api/quiz",
	Method: constellation.GET,
	Handler: constellation.HandlerSetup{
		Api: quizRoundHandler,
	},
}

var QuizCheck = constellation.Endpoint{
	Path:   "/api/quiz",
	Method: constellation.POST,
	Handler: constellation.HandlerSetup{
		Api: quizCheckHandler,
	},
}

type quizRound struct {
	ID         string
	Difficulty difficulty
	Prompt     string
	Choices    []string
	Answer     int
	Reference  string
}

func buildQuizRounds(level difficulty) []quizRound {
	rounds := make([]quizRound, 0)
	for _, question := range bible_lib.QuizQuestions {
		if question.Difficulty != level {
			continue
		}

		choices := append([]string(nil), question.Choices...)
		rounds = append(rounds, quizRound{
			ID:         question.ID,
			Difficulty: level,
			Prompt:     question.Prompt,
			Choices:    choices,
			Answer:     question.Answer,
			Reference:  question.Reference,
		})
	}
	return rounds
}

var quizByDifficulty = map[difficulty][]quizRound{
	difficultyEasy:   buildQuizRounds(difficultyEasy),
	difficultyMedium: buildQuizRounds(difficultyMedium),
	difficultyHard:   buildQuizRounds(difficultyHard),
}

var quizByID = func() map[string]quizRound {
	index := make(map[string]quizRound)
	for level, set := range quizByDifficulty {
		for _, round := range set {
			index[roundKey(level, round.ID)] = round
		}
	}
	return index
}()

var quizDecks = newDifficultyDecks(quizByDifficulty)

func quizRoundHandler(r *constellation.Request, w constellation.Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	level, ok := parseDifficulty(r.URL.Query().Get("difficulty"))
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		middleware.SendApi(w, map[string]any{
			"message": "Choose a difficulty before starting the quiz.",
		})
		return
	}

	round, ok := quizDecks.next(level)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		middleware.SendApi(w, map[string]any{
			"message": "There are no questions for that difficulty yet.",
		})
		return
	}

	choices := append([]string(nil), round.Choices...)
	rand.Shuffle(len(choices), func(i, j int) {
		choices[i], choices[j] = choices[j], choices[i]
	})

	middleware.SendApi(w, map[string]any{
		"id":         roundKey(round.Difficulty, round.ID),
		"difficulty": string(round.Difficulty),
		"question":   round.Prompt,
		"choices":    choices,
		"reference":  round.Reference,
		"character":  nextCharacter(round.Difficulty),
	})
}

type quizPayload struct {
	ID     string `json:"id"`
	Answer string `json:"answer"`
}

func quizCheckHandler(r *constellation.Request, w constellation.Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var payload quizPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		middleware.SendApi(w, map[string]any{
			"correct": false,
			"message": "That answer could not be read. Try again.",
		})
		return
	}

	round, ok := quizByID[payload.ID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		middleware.SendApi(w, map[string]any{
			"correct": false,
			"message": "That question is no longer in the quiz.",
		})
		return
	}

	answer := normalizeAnswer(payload.Answer)
	if answer == "" {
		w.WriteHeader(http.StatusBadRequest)
		middleware.SendApi(w, map[string]any{
			"correct": false,
			"message": "Choose an answer before submitting.",
		})
		return
	}

	answerText := round.Choices[round.Answer]
	correct := answer == normalizeAnswer(answerText)
	message := "Not quite. The answer is " + answerText + " — " + round.Reference + "."
	if correct {
		message = "Correct! " + answerText + " — " + round.Reference + "."
	}

	middleware.SendApi(w, map[string]any{
		"correct":    correct,
		"message":    message,
		"answerText": answerText,
		"reference":  round.Reference,
	})
}
