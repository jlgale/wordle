package wordle

import (
	"math/rand"
)

type Naive struct {
	rng *rand.Rand
}

func NaiveStrategy(rng *rand.Rand) Naive {
	return Naive{rng}
}

func (n Naive) Guess(game *Game) Word {
	var possible = game.PossibleAnswers()
	var idx = n.rng.Intn(len(possible))
	return possible[idx]
}
