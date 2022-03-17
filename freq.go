package wordle

type Freq struct {
	weights       map[Word]float64
	defaultWeight float64
}

// Weight more common words higher.
func FreqScale(weights map[Word]float64, defaultWeight float64) Freq {
	return Freq{weights, defaultWeight}
}

func (f Freq) Weights(words []Word) []float64 {
	var scores = make([]float64, len(words))
	for idx, w := range words {
		score, ok := f.weights[w]
		if ok {
			scores[idx] = score
		} else {
			scores[idx] = f.defaultWeight
		}
	}
	return scores
}
