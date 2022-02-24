package wordle

import (
	"math/rand"

	"github.com/rs/zerolog"
)

type Selective struct {
	rng *rand.Rand
	log *zerolog.Logger
}

func SelectiveStrategy(rng *rand.Rand, log *zerolog.Logger) Selective {
	return Selective{rng, log}
}

func (x Selective) Guess(game *Game) Word {
	var possible = game.Possible()
	// Find how often a letter is at a position in the set of possible words
	var found [WordLen]['z' - 'a' + 1]int
	for _, word := range possible {
		for idx, c := range word {
			found[idx][c-'a'] += 1
		}
	}

	// Assign a score to each word
	scores := make([]int, len(possible))
	for idx, w := range possible {
		score := 0
		for jdx, c := range w {
			score += found[jdx][c-'a']
		}
		scores[idx] = score
	}

	// Find the best score
	var topScore = -1
	var fromIdx int
	var nChoices int
	for idx, score := range scores {
		switch {
		case topScore == -1 || score > topScore:
			topScore = score
			fromIdx = idx
			nChoices = 1
		case score == topScore:
			nChoices++
		}
	}

	x.log.Printf("%d Choices with score %d from %d possible words",
		nChoices, topScore, len(possible))
	if nChoices == 1 {
		return possible[fromIdx]
	}

	// If we have multiple choices, chose 1 randomly
	var choice = x.rng.Intn(nChoices)
	for idx, score := range scores[fromIdx:] {
		if score == topScore {
			if choice == 0 {
				return possible[idx]
			}
			choice--
		}
	}
	panic("unreachable")
}
