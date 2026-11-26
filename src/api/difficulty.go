package api

import bible_lib "constellation/src/lib"

type difficulty = bible_lib.Difficulty

const (
	difficultyEasy   = bible_lib.Easy
	difficultyMedium = bible_lib.Medium
	difficultyHard   = bible_lib.Hard
)

var difficulties = []difficulty{difficultyEasy, difficultyMedium, difficultyHard}

func parseDifficulty(value string) (difficulty, bool) {
	switch difficulty(value) {
	case difficultyEasy, difficultyMedium, difficultyHard:
		return difficulty(value), true
	default:
		return "", false
	}
}

func roundKey(level difficulty, id string) string {
	return string(level) + "/" + id
}

func difficultyLabel(level difficulty) string {
	switch level {
	case difficultyEasy:
		return "Easy"
	case difficultyMedium:
		return "Medium"
	case difficultyHard:
		return "Hard"
	default:
		return "Difficulty"
	}
}
