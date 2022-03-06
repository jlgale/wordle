package wordle

import (
	"math"
	"math/rand"
	"sort"
)

type Scale interface {
	Weights(words []Word) []float64
}

// Weighted is a strategy to choose words randomly, weighted by
// the score from a given scoring function.
type Weighted struct {
	rng   *rand.Rand
	scale Scale
	pow   float64
}

func WeightedStrategy(rng *rand.Rand, scale Scale, pow float64) Weighted {
	return Weighted{rng, scale, pow}
}

func (x Weighted) Guess(game *Game) Word {
	var possible = game.PossibleAnswers()
	var weights = x.scale.Weights(possible)
	var offset = make([]float64, len(weights))
	var total = 0.0
	for idx, weight := range weights {
		total += math.Pow(weight, x.pow)
		offset[idx] = total
	}
	var choice = x.rng.Float64() * total
	var idx = sort.Search(len(offset), func(i int) bool {
		return offset[i] >= choice
	})
	return possible[idx]
}
