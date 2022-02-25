package wordle

import (
	"math/rand"
	"sort"
)

type Scale interface {
	Weights(words []Word) []float64
}

// Common is a strategy to choose words with the most common letters
// among the possible words. The hypothesis is that this will get more
// "yellow" squares, which is useful.
//
// In practice this is a losing strategy.
type Weighted struct {
	rng   *rand.Rand
	scale Scale
}

func WeightedStrategy(rng *rand.Rand, scale Scale) Weighted {
	return Weighted{rng, scale}
}

func (x Weighted) Guess(game *Game) Word {
	var possible = game.Possible()
	var weights = x.scale.Weights(possible)
	var offset = make([]float64, len(weights))
	var total = 0.0
	for idx, weight := range weights {
		total += weight
		offset[idx] = total
	}
	var choice = x.rng.Float64() * total
	var idx = sort.Search(len(offset), func(i int) bool {
		return offset[i] >= choice
	})
	// if len(offset) < 5 {
	// 	fmt.Printf("offset=%v\n", offset)
	// }
	// fmt.Printf("choice=%f, total=%f, idx=%d, n=%d, nn=%d\n", choice, total, idx, len(weights), len(possible))
	return possible[idx]
}
