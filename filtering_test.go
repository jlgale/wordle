package wordle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilteringPlay(t *testing.T) {
	rng := mkRand(1)
	fallback := WeightedStrategy(rng, DiversityScale(), 1)
	strategy := FilteringStrategy(rng, globalLog, fallback, 60)
	game := NewGame(globalWords, nil)
	answer, err := ParseWord("cigar")
	assert.Nil(t, err)
	for !game.Over() {
		guess := strategy.Guess(&game)
		match := guess.Match(answer)
		game = game.Guess(guess, match)
	}
	assert.True(t, game.Won())
	assert.Len(t, game.Guesses, 4) // arbitrary, but detect if something changes
}
