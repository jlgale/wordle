package wordle

type Diversity struct {
}

func DiversityScale() Diversity {
	return Diversity{}
}

// Weight words based on letter diversity
func (n Diversity) Weights(words []Word) []float64 {
	var scores = make([]float64, len(words))
	for idx, w := range words {
		var letters = w.Letters().Len()
		scores[idx] = float64(letters)
	}
	return scores
}
