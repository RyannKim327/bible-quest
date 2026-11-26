package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	constellation "constellation/lib"
	bible_lib "constellation/src/lib"
)

func TestNextCharacterUsesDifficultyPools(t *testing.T) {
	for _, level := range difficulties {
		pool := bible_lib.CharacterPools[bible_lib.Difficulty(level)]
		if len(pool) == 0 {
			t.Fatalf("difficulty %s has no characters", level)
		}

		seen := make(map[string]bool)
		for i := 0; i < 12; i++ {
			character := nextCharacter(level)
			found := false
			for _, expected := range pool {
				if character.Key == expected.Key && character.Name == expected.Name && character.Role == expected.Role && character.Verse == expected.Verse {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("difficulty %s returned unknown character: %+v", level, character)
			}
			seen[character.Key] = true
		}

		if len(seen) < 2 {
			t.Errorf("difficulty %s reused one character for every round", level)
		}
	}
}

func TestNextCharacterAvoidsImmediateRepeat(t *testing.T) {
	for _, level := range difficulties {
		first := nextCharacter(level)
		second := nextCharacter(level)
		if first.Key == second.Key {
			t.Errorf("difficulty %s repeated %s back-to-back", level, first.Key)
		}
	}
}

func TestGuessTheMistakeRoundHasCharacter(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/guess-the-mistake?difficulty=hard", nil)
	recorder := httptest.NewRecorder()
	guessTheMistakePassage((*constellation.Request)(request), recorder)

	if recorder.Code != http.StatusOK {
		t.Fatalf("GET status = %d, want 200", recorder.Code)
	}

	var round struct {
		Character characterNote `json:"character"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &round); err != nil {
		t.Fatalf("decode round: %v", err)
	}
	if round.Character.Key == "" || round.Character.Name == "" {
		t.Fatalf("round has no character guide: %+v", round)
	}
}
