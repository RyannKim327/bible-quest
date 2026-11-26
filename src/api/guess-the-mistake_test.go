package api

import (
	"slices"
	"testing"

	bible_lib "constellation/src/lib"
)

func setTestPassages(t *testing.T, passages []bible_lib.Passage) {
	t.Helper()
	original := bible_lib.Passages
	bible_lib.Passages = passages
	t.Cleanup(func() {
		bible_lib.Passages = original
	})
}

func TestPassageTextsArePresent(t *testing.T) {
	wantVersions := map[bible_lib.Difficulty]string{
		bible_lib.Easy:   "NIV",
		bible_lib.Medium: "NLT",
		bible_lib.Hard:   "KJV",
	}
	for level, want := range wantVersions {
		if got := bible_lib.Versions[level]; got != want {
			t.Errorf("Versions[%s] = %q, want %q", level, got, want)
		}
	}

	for _, passage := range bible_lib.Passages {
		for _, level := range []bible_lib.Difficulty{bible_lib.Easy, bible_lib.Medium, bible_lib.Hard} {
			text, ok := passage.Texts[level]
			if !ok {
				t.Errorf("%s: missing %s text", passage.ID, level)
				continue
			}
			if text == "" {
				t.Errorf("%s: %s text is empty, want passage text", passage.ID, level)
			}
		}
	}
}

func TestBuildRoundsPhraseMistakes(t *testing.T) {
	setTestPassages(t, []bible_lib.Passage{
		{
			ID:        "phrase-test",
			Reference: "Phrase Test",
			Texts: bible_lib.PassageTexts{
				bible_lib.Easy:   "lean not on your own understanding;",
				bible_lib.Medium: "lean not on your own heart;",
				bible_lib.Hard:   "lean not upon your own understanding;",
			},
			Mistakes: bible_lib.PassageMistakes{
				bible_lib.Easy:   {"not on": "unto"},
				bible_lib.Medium: {"own heart": "own"},
				bible_lib.Hard:   {"lean not": "lean"},
			},
		},
	})

	tests := []struct {
		level   difficulty
		version string
		words   []string
		answers []int
		fixes   []wordFix
	}{
		{
			level:   difficultyEasy,
			version: "NIV",
			words:   []string{"lean", "unto", "your", "own", "understanding;"},
			answers: []int{1},
			fixes:   []wordFix{{Given: "unto", Correction: "not on"}},
		},
		{
			level:   difficultyMedium,
			version: "NLT",
			words:   []string{"lean", "not", "on", "your", "own;"},
			answers: []int{4},
			fixes:   []wordFix{{Given: "own", Correction: "own heart"}},
		},
		{
			level:   difficultyHard,
			version: "KJV",
			words:   []string{"lean", "upon", "your", "own", "understanding;"},
			answers: []int{0},
			fixes:   []wordFix{{Given: "lean", Correction: "lean not"}},
		},
	}

	for _, test := range tests {
		t.Run(string(test.level), func(t *testing.T) {
			rounds := buildRounds(test.level)
			if len(rounds) != 1 {
				t.Fatalf("got %d rounds, want 1", len(rounds))
			}

			round := rounds[0]
			if round.Version != test.version {
				t.Errorf("version = %q, want %q", round.Version, test.version)
			}
			if !slices.Equal(round.Words, test.words) {
				t.Errorf("words = %q, want %q", round.Words, test.words)
			}
			if !slices.Equal(round.Answers, test.answers) {
				t.Errorf("answers = %v, want %v", round.Answers, test.answers)
			}
			if !slices.Equal(round.Fixes, test.fixes) {
				t.Errorf("fixes = %v, want %v", round.Fixes, test.fixes)
			}
		})
	}
}

func TestBuildRoundsShiftsLaterAnswers(t *testing.T) {
	setTestPassages(t, []bible_lib.Passage{
		{
			ID:        "shift-test",
			Reference: "Shift Test",
			Texts: bible_lib.PassageTexts{
				bible_lib.Easy: "Trust not to your heart and love him.",
			},
			Mistakes: bible_lib.PassageMistakes{
				bible_lib.Easy: {
					"not to": "unto",
					"love":   "hate",
				},
			},
		},
	})

	rounds := buildRounds(difficultyEasy)
	if len(rounds) != 1 {
		t.Fatalf("got %d rounds, want 1", len(rounds))
	}

	round := rounds[0]
	if round.Version != "NIV" {
		t.Errorf("version = %q, want NIV", round.Version)
	}

	wantWords := []string{"Trust", "unto", "your", "heart", "and", "hate", "him."}
	if !slices.Equal(round.Words, wantWords) {
		t.Errorf("words = %q, want %q", round.Words, wantWords)
	}

	wantAnswers := []int{1, 5}
	if !slices.Equal(round.Answers, wantAnswers) {
		t.Errorf("answers = %v, want %v", round.Answers, wantAnswers)
	}

	wantFixes := []wordFix{
		{Given: "unto", Correction: "not to"},
		{Given: "hate", Correction: "love"},
	}
	if !slices.Equal(round.Fixes, wantFixes) {
		t.Errorf("fixes = %v, want %v", round.Fixes, wantFixes)
	}
}

func TestBuildRoundsPhraseWithMultipleWrongWords(t *testing.T) {
	setTestPassages(t, []bible_lib.Passage{
		{
			ID:        "long-wrong-test",
			Reference: "Long Wrong Test",
			Texts: bible_lib.PassageTexts{
				bible_lib.Hard: "lean not on your own",
			},
			Mistakes: bible_lib.PassageMistakes{
				bible_lib.Hard: {"not on": "not unto"},
			},
		},
	})

	rounds := buildRounds(difficultyHard)
	if len(rounds) != 1 {
		t.Fatalf("got %d rounds, want 1", len(rounds))
	}

	round := rounds[0]
	if round.Version != "KJV" {
		t.Errorf("version = %q, want KJV", round.Version)
	}

	wantWords := []string{"lean", "not", "unto", "your", "own"}
	if !slices.Equal(round.Words, wantWords) {
		t.Errorf("words = %q, want %q", round.Words, wantWords)
	}

	wantAnswers := []int{1, 2}
	if !slices.Equal(round.Answers, wantAnswers) {
		t.Errorf("answers = %v, want %v", round.Answers, wantAnswers)
	}

	wantFixes := []wordFix{{Given: "not unto", Correction: "not on"}}
	if !slices.Equal(round.Fixes, wantFixes) {
		t.Errorf("fixes = %v, want %v", round.Fixes, wantFixes)
	}
}

func TestBuildRoundUsesDifficultyText(t *testing.T) {
	passage := bible_lib.Passage{
		ID:        "proverbs-3-5",
		Reference: "Proverbs 3:5",
		Texts: bible_lib.PassageTexts{
			bible_lib.Hard: "Trust in the Lord with all your heart and lean not on your own understanding;",
		},
		Mistakes: bible_lib.PassageMistakes{
			bible_lib.Hard: {"not on": "unto"},
		},
	}

	round := buildRound(passage, difficultyHard)
	wantWords := []string{
		"Trust", "in", "the", "Lord", "with", "all", "your", "heart",
		"and", "lean", "unto", "your", "own", "understanding;",
	}
	if !slices.Equal(round.Words, wantWords) {
		t.Errorf("words = %q, want %q", round.Words, wantWords)
	}

	wantAnswers := []int{10}
	if !slices.Equal(round.Answers, wantAnswers) {
		t.Errorf("answers = %v, want %v", round.Answers, wantAnswers)
	}

	wantFixes := []wordFix{{Given: "unto", Correction: "not on"}}
	if !slices.Equal(round.Fixes, wantFixes) {
		t.Errorf("fixes = %v, want %v", round.Fixes, wantFixes)
	}
}
