package wordle

import (
	"math/rand"
)

// Top is a strategy to choose word the top words for a given score.
// If multiple words have the same score, choose randomly among them.
type Top struct {
	rng   *rand.Rand
	scale Scale
}

func TopStrategy(rng *rand.Rand, scale Scale) Top {
	return Top{rng, scale}
}

func (x Top) Guess(game *Game) Word {
	var possible = game.PossibleAnswers()
	var weights = x.scale.Weights(possible)
	var top = weights[0]
	var topIdx = 0
	var n = 1
	for idx, w := range weights[1:] {
		if w > top {
			top = w
			topIdx = idx
			n = 1
		} else if w == top {
			n++
		}
	}
	if n == 1 {
		return possible[topIdx]
	}
	var choices []Word
	for idx, w := range weights[topIdx:] {
		if w == top {
			choices = append(choices, possible[idx])
		}
	}
	var choice = x.rng.Intn(len(choices))
	return choices[choice]
}
