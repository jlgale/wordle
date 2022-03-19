package wordle

// Strategy is used to choose the next word to play in the given Game.
type Strategy interface {
	Guess(w *Game) Word
}

// Scoring assigns a score, or "weight", to each word in the given array.
// The weights can be independent or dependent on the other words in the
// array.
type Scoring interface {
	Weights(words []Word) []float64
}

// Generic debug logger used by some strategies. Compatible with the zerolog package.
type Logger interface {
	Printf(template string, args ...interface{})
}
