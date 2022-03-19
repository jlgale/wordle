package wordle

type ScoringCache struct {
	inner   Scoring
	words   []Word
	weights []float64
}

func NewScoringCache(scoring Scoring, words []Word) *ScoringCache {
	return &ScoringCache{scoring, words, scoring.Weights(words)}
}

func (c *ScoringCache) Weights(words []Word) []float64 {
	if identical(c.words, words) {
		return c.weights
	}
	return c.inner.Weights(words)
}

// identical tests if two slices are in fact pointing to the same
// array of values.
func identical(s1, s2 []Word) bool {
	switch {
	case len(s1) != len(s2):
		return false
	case len(s1) == 0:
		return true
	default:
		return &s1[0] == &s2[0]
	}
}
