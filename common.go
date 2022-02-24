package wordle

import (
	"math/rand"

	"github.com/rs/zerolog"
)

// Common is a strategy to choose words with the most common letters
// among the possible words. The hypothesis is that this will get more
// "yellow" squares, which is useful.
//
// In practice this is a losing strategy.
type Common struct {
	rng *rand.Rand
	log *zerolog.Logger
}

func CommonStrategy(rng *rand.Rand, log *zerolog.Logger) Common {
	return Common{rng, log}
}

func (x Common) Guess(game *Game) Word {
	var possible = game.Possible()
	var used ['z' - 'a' + 1]int
	for _, word := range possible {
		seen := NewLetters(nil)
		for _, c := range word {
			if !seen.Contains(c) {
				used[c-'a'] += 1
				seen = seen.AddChar(c)
			}
		}
	}

	// Assign a score to each word
	scores := make([]int, len(possible))
	for idx, w := range possible {
		score := 0
		for _, c := range w {
			score += used[c-'a']
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
