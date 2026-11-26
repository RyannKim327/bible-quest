package api

import (
	"math/rand/v2"
	"sync"
)

type shuffledDeck[T any] struct {
	mu        sync.Mutex
	items     []T
	remaining []int
	last      int
}

func newShuffledDeck[T any](items []T) *shuffledDeck[T] {
	deck := &shuffledDeck[T]{items: items, last: -1}
	deck.refill()
	return deck
}

func (deck *shuffledDeck[T]) refill() {
	deck.remaining = make([]int, 0, len(deck.items))
	for i := range deck.items {
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

func (deck *shuffledDeck[T]) next() (T, bool) {
	var empty T

	deck.mu.Lock()
	defer deck.mu.Unlock()

	if len(deck.items) == 0 {
		return empty, false
	}
	if len(deck.remaining) == 0 {
		deck.refill()
	}

	top := len(deck.remaining) - 1
	position := deck.remaining[top]
	deck.remaining = deck.remaining[:top]
	deck.last = position
	return deck.items[position], true
}

type difficultyDecks[T any] struct {
	mu           sync.Mutex
	items        map[difficulty][]T
	byDifficulty map[difficulty]*shuffledDeck[T]
}

func newDifficultyDecks[T any](items map[difficulty][]T) *difficultyDecks[T] {
	return &difficultyDecks[T]{
		items:        items,
		byDifficulty: make(map[difficulty]*shuffledDeck[T]),
	}
}

func (decks *difficultyDecks[T]) next(level difficulty) (T, bool) {
	decks.mu.Lock()
	deck, ok := decks.byDifficulty[level]
	if !ok {
		deck = newShuffledDeck(decks.items[level])
		decks.byDifficulty[level] = deck
	}
	decks.mu.Unlock()

	return deck.next()
}
