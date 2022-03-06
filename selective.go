package wordle

type Selective struct {
	log Logger
}

func SelectiveScale(log Logger) Selective {
	return Selective{log}
}

func (x Selective) Weights(words []Word) []float64 {
	// Find how often a letter is at a position in the set of possible words
	var found [WordLen]['z' - 'a' + 1]int
	for _, word := range words {
		for idx, c := range word {
			found[idx][c-'a'] += 1
		}
	}

	// Assign a score to each word
	scores := make([]float64, len(words))
	for idx, w := range words {
		score := 0
		for jdx, c := range w {
			score += found[jdx][c-'a']
		}
		scores[idx] = float64(score)
	}
	return scores
}
