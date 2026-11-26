package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	constellation "constellation/lib"
	bible_lib "constellation/src/lib"
)

func TestOneWordData(t *testing.T) {
	counts := make(map[difficulty]int)
	seen := make(map[string]bool)

	for _, clue := range bible_lib.OneWordClues {
		counts[clue.Difficulty]++
		if _, ok := parseDifficulty(string(clue.Difficulty)); !ok {
			t.Errorf("%s: unknown difficulty %q", clue.ID, clue.Difficulty)
		}
		if len(strings.Fields(clue.Word)) != 1 {
			t.Errorf("%s: clue %q is not one word", clue.ID, clue.Word)
		}
		if clue.Kind != bible_lib.OneWordCharacter && clue.Kind != bible_lib.OneWordBook {
			t.Errorf("%s: unknown kind %q", clue.ID, clue.Kind)
		}
		if clue.Answer == "" || clue.Reference == "" {
			t.Errorf("%s: missing answer or reference", clue.ID)
		}

		key := roundKey(clue.Difficulty, clue.ID)
		if seen[key] {
			t.Errorf("%s: duplicate clue key %q", clue.ID, key)
		}
		seen[key] = true
	}

	for _, level := range difficulties {
		if counts[level] == 0 {
			t.Errorf("difficulty %s has no clues", level)
		}
	}
}

func TestNormalizeAnswer(t *testing.T) {
	cases := map[string]string{
		"  JOHN THE BAPTIST. ": "john the baptist",
		"Psalms!":              "psalms",
		"  Noah  ":             "noah",
	}
	for input, want := range cases {
		if got := normalizeAnswer(input); got != want {
			t.Errorf("normalizeAnswer(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestOneWordRoundAndCheck(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/one-word?difficulty=medium", nil)
	recorder := httptest.NewRecorder()
	oneWordRoundHandler((*constellation.Request)(request), recorder)

	if recorder.Code != http.StatusOK {
		t.Fatalf("GET status = %d, want 200", recorder.Code)
	}

	var round struct {
		ID        string        `json:"id"`
		Word      string        `json:"word"`
		Kind      string        `json:"kind"`
		Character characterNote `json:"character"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &round); err != nil {
		t.Fatalf("decode round: %v", err)
	}
	if round.ID == "" || round.Word == "" || round.Kind == "" {
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

	stored, ok := oneWordByID[round.ID]
	if !ok {
		t.Fatalf("round %q is not indexed", round.ID)
	}

	result := checkOneWordAnswer(t, round.ID, strings.ToLower(stored.Answer), true)
	if result["answer"] != stored.Answer {
		t.Errorf("revealed answer = %v, want %q", result["answer"], stored.Answer)
	}

	checkOneWordAnswer(t, round.ID, "definitely not right", false)
}

func checkOneWordAnswer(t *testing.T, id, answer string, wantCorrect bool) map[string]any {
	t.Helper()

	body, err := json.Marshal(map[string]string{"id": id, "answer": answer})
	if err != nil {
		t.Fatalf("marshal answer: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/one-word", strings.NewReader(string(body)))
	recorder := httptest.NewRecorder()
	oneWordCheckHandler((*constellation.Request)(request), recorder)

	if recorder.Code != http.StatusOK {
		t.Fatalf("POST status = %d, want 200: %s", recorder.Code, recorder.Body.String())
	}

	var result map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if result["correct"] != wantCorrect {
		t.Errorf("correct = %v, want %v", result["correct"], wantCorrect)
	}
	return result
}

func TestOneWordAcceptsAlias(t *testing.T) {
	id := roundKey(difficultyHard, "camel-john")
	checkOneWordAnswer(t, id, "John", true)
}

func TestOneWordRejectsInvalidInput(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/one-word?difficulty=impossible", nil)
	recorder := httptest.NewRecorder()
	oneWordRoundHandler((*constellation.Request)(request), recorder)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("invalid difficulty status = %d, want 400", recorder.Code)
	}

	request = httptest.NewRequest(http.MethodPost, "/api/one-word", strings.NewReader(`{"id":"easy/ark-noah","answer":"   "}`))
	recorder = httptest.NewRecorder()
	oneWordCheckHandler((*constellation.Request)(request), recorder)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("blank answer status = %d, want 400", recorder.Code)
	}

	request = httptest.NewRequest(http.MethodPost, "/api/one-word", strings.NewReader(`{"id":"missing/clue","answer":"Noah"}`))
	recorder = httptest.NewRecorder()
	oneWordCheckHandler((*constellation.Request)(request), recorder)
	if recorder.Code != http.StatusNotFound {
		t.Errorf("unknown clue status = %d, want 404", recorder.Code)
	}
}
