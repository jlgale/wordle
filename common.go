package wordle

import (
	"github.com/rs/zerolog"
)

// Common is a strategy to choose words with the most common letters
// among the possible words. The hypothesis is that this will get more
// "yellow" squares, which is useful.
//
// In practice this is a losing strategy.
type Common struct {
	log *zerolog.Logger
}

func CommonScale(log *zerolog.Logger) Common {
	return Common{log}
}

func (x Common) Weights(words []Word) []float64 {
	var used ['z' - 'a' + 1]int
	for _, word := range words {
		seen := NewLetters(nil)
		for _, c := range word {
			if !seen.Contains(c) {
				used[c-'a'] += 1
				seen = seen.AddChar(c)
			}
		}
	}

	// Assign a score to each word
	scores := make([]float64, len(words))
	for idx, w := range words {
		score := 0
		for _, c := range w {
			score += used[c-'a']
		}
		scores[idx] = float64(score)
	}
	return scores
}
