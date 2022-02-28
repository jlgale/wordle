package wordle

import (
	"math/rand"

	"github.com/rs/zerolog"
)

type Naive struct {
	rng *rand.Rand
	log *zerolog.Logger
}

func NaiveStrategy(rng *rand.Rand, log *zerolog.Logger) Naive {
	return Naive{rng, log}
}

func (n Naive) Guess(game *Game) Word {
	var possible = game.PossibleAnswers()
	n.log.Printf("choosing from %d possible words", len(possible))
	var idx = n.rng.Intn(len(possible))
	return possible[idx]
}
