package api

import (
	"math/rand/v2"
	"sync"

	bible_lib "constellation/src/lib"
)

type characterNote struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Verse string `json:"verse"`
}

var characterState = struct {
	mu   sync.Mutex
	last map[difficulty]string
}{last: make(map[difficulty]string)}

func nextCharacter(level difficulty) characterNote {
	characterState.mu.Lock()
	defer characterState.mu.Unlock()

	pool := bible_lib.CharacterPools[bible_lib.Difficulty(level)]
	if len(pool) == 0 {
		return characterNote{}
	}

	index := rand.IntN(len(pool))
	if len(pool) > 1 && pool[index].Key == characterState.last[level] {
		index = (index + 1) % len(pool)
	}

	character := pool[index]
	characterState.last[level] = character.Key
	return characterNote{
		Key:   character.Key,
		Name:  character.Name,
		Role:  character.Role,
		Verse: character.Verse,
	}
}
