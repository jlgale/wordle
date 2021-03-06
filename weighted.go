package wordle

import (
	"math"
	"math/rand"
	"sort"
)

// WeightedStrategy is a strategy to choose words randomly, weighted by
// the score from a given scoring function.
type WeightedStrategy struct {
	rng     *rand.Rand
	scoring Scoring
	pow     float64
}

func NewWeightedStrategy(rng *rand.Rand, scoring Scoring, pow float64) *WeightedStrategy {
	return &WeightedStrategy{rng, scoring, pow}
}

func (x *WeightedStrategy) Guess(game *Game) Word {
	var possible = game.PossibleAnswers()
	var weights = x.scoring.Weights(possible)
	var idx = weightedSample(x.rng, x.pow, weights)
	return possible[idx]
}

func weightedSample(rng *rand.Rand, pow float64, weights []float64) int {
	var offset = make([]float64, len(weights))
	var total = 0.0
	for idx, weight := range weights {
		total += math.Pow(weight, pow)
		offset[idx] = total
	}
	var choice = rng.Float64() * total
	return sort.Search(len(offset), func(i int) bool {
		return offset[i] >= choice
	})
}
