package wordle

type UniqueLettersScoring struct {
}

func NewUniqueLettersScoring() UniqueLettersScoring {
	return UniqueLettersScoring{}
}

// Weight words based on letter diversity
func (n UniqueLettersScoring) Weights(words []Word) []float64 {
	var scores = make([]float64, len(words))
	for idx, w := range words {
		var letters = w.Letters().Len()
		scores[idx] = float64(letters)
	}
	return scores
}
