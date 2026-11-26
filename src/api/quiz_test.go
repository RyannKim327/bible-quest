package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	constellation "constellation/lib"
	bible_lib "constellation/src/lib"
)

func TestQuizData(t *testing.T) {
	counts := make(map[difficulty]int)
	seen := make(map[string]bool)

	for _, question := range bible_lib.QuizQuestions {
		counts[question.Difficulty]++
		if _, ok := parseDifficulty(string(question.Difficulty)); !ok {
			t.Errorf("%s: unknown difficulty %q", question.ID, question.Difficulty)
		}
		if len(question.Choices) < 2 {
			t.Errorf("%s: fewer than two choices", question.ID)
		}
		if question.Answer < 0 || question.Answer >= len(question.Choices) {
			t.Errorf("%s: answer index %d is out of range", question.ID, question.Answer)
		}
		if question.Prompt == "" || question.Reference == "" {
			t.Errorf("%s: missing prompt or reference", question.ID)
		}

		key := roundKey(question.Difficulty, question.ID)
		if seen[key] {
			t.Errorf("%s: duplicate question key %q", question.ID, key)
		}
		seen[key] = true
	}

	for _, level := range difficulties {
		if counts[level] == 0 {
			t.Errorf("difficulty %s has no questions", level)
		}
	}
}

type quizRoundResponse struct {
	ID        string        `json:"id"`
	Question  string        `json:"question"`
	Choices   []string      `json:"choices"`
	Character characterNote `json:"character"`
}

func fetchQuizRound(t *testing.T) (quizRoundResponse, string) {
	t.Helper()

	request := httptest.NewRequest(http.MethodGet, "/api/quiz?difficulty=easy", nil)
	recorder := httptest.NewRecorder()
	quizRoundHandler((*constellation.Request)(request), recorder)

	if recorder.Code != http.StatusOK {
		t.Fatalf("GET status = %d, want 200", recorder.Code)
	}

	var round quizRoundResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &round); err != nil {
		t.Fatalf("decode round: %v", err)
	}
	if round.ID == "" || round.Question == "" || len(round.Choices) == 0 {
		t.Fatalf("incomplete round: %+v", round)
	}
	if round.Character.Key == "" || round.Character.Name == "" {
		t.Fatalf("round has no character guide: %+v", round)
	}

	var fields map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &fields); err != nil {
		t.Fatalf("decode fields: %v", err)
	}
	if _, ok := fields["answer"]; ok {
		t.Errorf("GET response leaked the answer: %s", recorder.Body.String())
	}

	return round, recorder.Body.String()
}

func TestQuizRoundAndCheck(t *testing.T) {
	round, _ := fetchQuizRound(t)
	stored, ok := quizByID[round.ID]
	if !ok {
		t.Fatalf("round %q is not indexed", round.ID)
	}

	got := append([]string(nil), round.Choices...)
	want := append([]string(nil), stored.Choices...)
	slices.Sort(got)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Errorf("choices = %v, want a permutation of %v", round.Choices, stored.Choices)
	}

	correct := stored.Choices[stored.Answer]
	checkQuizAnswer(t, round.ID, correct, true)
	checkQuizAnswer(t, round.ID, strings.ToUpper(correct), true)

	wrong := ""
	for _, choice := range stored.Choices {
		if choice != correct {
			wrong = choice
			break
		}
	}
	if wrong == "" {
		t.Fatal("question has no wrong choice")
	}
	checkQuizAnswer(t, round.ID, wrong, false)
}

func checkQuizAnswer(t *testing.T, id, answer string, wantCorrect bool) {
	t.Helper()

	body, err := json.Marshal(map[string]string{"id": id, "answer": answer})
	if err != nil {
		t.Fatalf("marshal answer: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/quiz", strings.NewReader(string(body)))
	recorder := httptest.NewRecorder()
	quizCheckHandler((*constellation.Request)(request), recorder)

	if recorder.Code != http.StatusOK {
		t.Fatalf("POST status = %d, want 200: %s", recorder.Code, recorder.Body.String())
	}

	var result struct {
		Correct    bool   `json:"correct"`
		Message    string `json:"message"`
		AnswerText string `json:"answerText"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if result.Correct != wantCorrect {
		t.Errorf("correct = %v, want %v", result.Correct, wantCorrect)
	}
	if result.Message == "" || result.AnswerText == "" {
		t.Errorf("incomplete result: %+v", result)
	}
}

func TestQuizChoicesAreShuffled(t *testing.T) {
	shuffled := false

	for i := 0; i < 40 && !shuffled; i++ {
		round, _ := fetchQuizRound(t)
		stored, ok := quizByID[round.ID]
		if !ok {
			t.Fatalf("round %q is not indexed", round.ID)
		}
		shuffled = !slices.Equal(round.Choices, stored.Choices)
	}

	if !shuffled {
		t.Error("choices were never shuffled across 40 rounds")
	}
}

func TestQuizRejectsInvalidInput(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/quiz?difficulty=impossible", nil)
	recorder := httptest.NewRecorder()
	quizRoundHandler((*constellation.Request)(request), recorder)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("invalid difficulty status = %d, want 400", recorder.Code)
	}

	request = httptest.NewRequest(http.MethodPost, "/api/quiz", strings.NewReader(`{"id":"easy/ark-builder"}`))
	recorder = httptest.NewRecorder()
	quizCheckHandler((*constellation.Request)(request), recorder)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("missing answer status = %d, want 400", recorder.Code)
	}

	request = httptest.NewRequest(http.MethodPost, "/api/quiz", strings.NewReader(`{"id":"missing/question","answer":"Noah"}`))
	recorder = httptest.NewRecorder()
	quizCheckHandler((*constellation.Request)(request), recorder)
	if recorder.Code != http.StatusNotFound {
		t.Errorf("unknown question status = %d, want 404", recorder.Code)
	}
}

func TestQuizDeckAvoidsImmediateRepeat(t *testing.T) {
	deck := newShuffledDeck([]string{"a", "b", "c"})
	previous := ""

	for i := 0; i < 24; i++ {
		current, ok := deck.next()
		if !ok {
			t.Fatal("deck returned no item")
		}
		if current == previous {
			t.Fatalf("deck repeated %q back-to-back", current)
		}
		previous = current
	}
}
