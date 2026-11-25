package api

import (
	"encoding/json"
	"maps"
	"math/rand/v2"
	"net/http"
	"slices"
	"strings"
	"sync"
	"unicode"

	constellation "constellation/lib"
	"constellation/lib/middleware"
	bible_lib "constellation/src/lib"
)

var GuessTheMistakePassage = constellation.Endpoint{
	Path:   "/api/guess-the-mistake",
	Method: constellation.GET,
	Handler: constellation.HandlerSetup{
		Api: guessTheMistakePassage,
	},
}

var GuessTheMistakeCheck = constellation.Endpoint{
	Path:   "/api/guess-the-mistake",
	Method: constellation.POST,
	Handler: constellation.HandlerSetup{
		Api: guessTheMistakeCheck,
	},
}

type difficulty = bible_lib.Difficulty

const (
	difficultyEasy   = bible_lib.Easy
	difficultyMedium = bible_lib.Medium
	difficultyHard   = bible_lib.Hard
)

func parseDifficulty(value string) (difficulty, bool) {
	switch difficulty(value) {
	case difficultyEasy, difficultyMedium, difficultyHard:
		return difficulty(value), true
	default:
		return "", false
	}
}

type wordFix struct {
	Given      string
	Correction string
}

type gameRound struct {
	ID         string
	Difficulty difficulty
	Reference  string
	Version    string
	Words      []string
	Answers    []int
	Fixes      []wordFix
}

type mistakeMatch struct {
	start   int
	end     int
	correct string
	given   string
	wrong   []string
}

func tokenCore(token string) (string, string) {
	core := strings.TrimRightFunc(token, unicode.IsPunct)
	return core, token[len(core):]
}

func findMistake(correct string, tokens []string) (int, int, bool) {
	needles := strings.Fields(correct)
	if len(needles) == 0 || len(needles) > len(tokens) {
		return -1, -1, false
	}

	for start := 0; start+len(needles) <= len(tokens); start++ {
		matched := true
		for offset, needle := range needles {
			core, _ := tokenCore(tokens[start+offset])
			if !strings.EqualFold(core, needle) {
				matched = false
				break
			}
		}
		if matched {
			return start, start + len(needles) - 1, true
		}
	}

	return -1, -1, false
}

func sortMistakeMatches(matches []mistakeMatch) {
	slices.SortFunc(matches, func(a, b mistakeMatch) int {
		if a.start < b.start {
			return -1
		}
		if a.start > b.start {
			return 1
		}

		aLength := a.end - a.start
		bLength := b.end - b.start
		if aLength > bLength {
			return -1
		}
		if aLength < bLength {
			return 1
		}

		return strings.Compare(a.correct, b.correct)
	})
}

func buildRound(item bible_lib.Passage, level difficulty) gameRound {
	tokens := strings.Fields(item.Texts[level])
	round := gameRound{
		ID:         item.ID,
		Difficulty: level,
		Reference:  item.Reference,
		Version:    bible_lib.Versions[level],
		Words:      make([]string, 0, len(tokens)),
	}
	mistakes := item.Mistakes[level]
	matches := make([]mistakeMatch, 0, len(mistakes))

	for _, correct := range slices.Sorted(maps.Keys(mistakes)) {
		wrong := mistakes[correct]
		wrongTokens := strings.Fields(wrong)
		if len(wrongTokens) == 0 {
			continue
		}

		start, end, ok := findMistake(correct, tokens)
		if !ok {
			continue
		}

		matches = append(matches, mistakeMatch{
			start:   start,
			end:     end,
			correct: correct,
			given:   wrong,
			wrong:   wrongTokens,
		})
	}

	sortMistakeMatches(matches)

	for i, matchIndex := 0, 0; i < len(tokens); {
		for matchIndex < len(matches) && matches[matchIndex].start < i {
			matchIndex++
		}

		if matchIndex < len(matches) && matches[matchIndex].start == i {
			match := matches[matchIndex]
			_, suffix := tokenCore(tokens[match.end])
			replacement := append([]string(nil), match.wrong...)
			replacement[len(replacement)-1] += suffix

			for _, word := range replacement {
				round.Answers = append(round.Answers, len(round.Words))
				round.Words = append(round.Words, word)
			}
			round.Fixes = append(round.Fixes, wordFix{Given: match.given, Correction: match.correct})

			i = match.end + 1
			matchIndex++
			continue
		}

		round.Words = append(round.Words, tokens[i])
		i++
	}

	return round
}

func buildRounds(level difficulty) []gameRound {
	rounds := make([]gameRound, 0, len(bible_lib.Passages))
	for _, item := range bible_lib.Passages {
		rounds = append(rounds, buildRound(item, level))
	}
	return rounds
}

var roundsByDifficulty = map[difficulty][]gameRound{
	difficultyEasy:   buildRounds(difficultyEasy),
	difficultyMedium: buildRounds(difficultyMedium),
	difficultyHard:   buildRounds(difficultyHard),
}

func roundKey(level difficulty, id string) string {
	return string(level) + "/" + id
}

var roundsByID = func() map[string]gameRound {
	index := make(map[string]gameRound)
	for level, set := range roundsByDifficulty {
		for _, round := range set {
			index[roundKey(level, round.ID)] = round
		}
	}
	return index
}()

type roundDeck struct {
	mu        sync.Mutex
	rounds    []gameRound
	remaining []int
	last      int
}

func newRoundDeck(rounds []gameRound) *roundDeck {
	deck := &roundDeck{rounds: rounds, last: -1}
	deck.refill()
	return deck
}

func (deck *roundDeck) refill() {
	deck.remaining = make([]int, 0, len(deck.rounds))
	for i := range deck.rounds {
		deck.remaining = append(deck.remaining, i)
	}
	rand.Shuffle(len(deck.remaining), func(i, j int) {
		deck.remaining[i], deck.remaining[j] = deck.remaining[j], deck.remaining[i]
	})

	top := len(deck.remaining) - 1
	if len(deck.remaining) > 1 && deck.remaining[top] == deck.last {
		deck.remaining[0], deck.remaining[top] = deck.remaining[top], deck.remaining[0]
	}
}

func (deck *roundDeck) next() gameRound {
	deck.mu.Lock()
	defer deck.mu.Unlock()

	if len(deck.remaining) == 0 {
		deck.refill()
	}

	top := len(deck.remaining) - 1
	position := deck.remaining[top]
	deck.remaining = deck.remaining[:top]
	deck.last = position
	return deck.rounds[position]
}

var decks = struct {
	mu           sync.Mutex
	byDifficulty map[difficulty]*roundDeck
}{byDifficulty: make(map[difficulty]*roundDeck)}

func nextRound(level difficulty) gameRound {
	decks.mu.Lock()
	deck, ok := decks.byDifficulty[level]
	if !ok {
		deck = newRoundDeck(roundsByDifficulty[level])
		decks.byDifficulty[level] = deck
	}
	decks.mu.Unlock()

	return deck.next()
}

func guessTheMistakePassage(r *constellation.Request, w constellation.Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	level, ok := parseDifficulty(r.URL.Query().Get("difficulty"))
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		middleware.SendApi(w, map[string]any{
			"message": "Choose a difficulty before starting the game.",
		})
		return
	}

	round := nextRound(level)
	middleware.SendApi(w, map[string]any{
		"id":         roundKey(round.Difficulty, round.ID),
		"difficulty": string(round.Difficulty),
		"reference":  round.Reference,
		"version":    round.Version,
		"words":      round.Words,
	})
}

type guessPayload struct {
	ID    string `json:"id"`
	Words []int  `json:"words"`
	Word  *int   `json:"word"`
}

func (payload guessPayload) selectedWords() []int {
	if len(payload.Words) > 0 {
		return payload.Words
	}
	if payload.Word != nil {
		return []int{*payload.Word}
	}
	return nil
}

func describeFixes(fixes []wordFix) string {
	parts := make([]string, 0, len(fixes))
	for _, fix := range fixes {
		parts = append(parts, "“"+fix.Given+"” should be “"+fix.Correction+"”")
	}
	return strings.Join(parts, " and ")
}

func guessTheMistakeCheck(r *constellation.Request, w constellation.Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var payload guessPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		middleware.SendApi(w, map[string]any{
			"correct": false,
			"message": "That answer could not be read. Try again.",
		})
		return
	}

	round, ok := roundsByID[payload.ID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		middleware.SendApi(w, map[string]any{
			"correct": false,
			"message": "That passage is no longer in the quest.",
		})
		return
	}

	selected := make(map[int]bool)
	for _, index := range payload.selectedWords() {
		if index < 0 || index >= len(round.Words) {
			w.WriteHeader(http.StatusBadRequest)
			middleware.SendApi(w, map[string]any{
				"correct": false,
				"message": "One of those words is not in the passage.",
			})
			return
		}
		selected[index] = true
	}

	correct := false
	message := ""
	description := describeFixes(round.Fixes)

	if len(round.Answers) == 0 {
		correct = len(selected) == 0
		if correct {
			message = "Correct — this passage has no mistake."
		} else {
			message = "This passage was already right. There was no word to swap."
		}
	} else {
		correct = true
		for _, answer := range round.Answers {
			if !selected[answer] {
				correct = false
				break
			}
		}

		if correct {
			if len(round.Answers) > 1 {
				message = "All mistakes found! " + description + "."
			} else {
				message = "Spot on! " + description + "."
			}
		} else if len(selected) == 0 {
			message = "There is a mistake here: " + description + "."
		} else {
			message = "Not quite. " + description + "."
		}
	}

	indices := append([]int(nil), round.Answers...)
	if indices == nil {
		indices = []int{}
	}

	fixes := make([]map[string]string, 0, len(round.Fixes))
	for _, fix := range round.Fixes {
		fixes = append(fixes, map[string]string{"given": fix.Given, "correction": fix.Correction})
	}

	firstIndex := -1
	if len(indices) > 0 {
		firstIndex = indices[0]
	}

	middleware.SendApi(w, map[string]any{
		"correct":        correct,
		"message":        message,
		"mistakeIndex":   firstIndex,
		"mistakeIndices": indices,
		"fixes":          fixes,
		"hasMistake":     len(indices) > 0,
	})
}
