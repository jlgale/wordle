package wordle

// Play a fixed opening sequence before continuing with a follow-on
// strategy.
type Fixed struct {
	open     []Word
	followOn Strategy
}

func FixedStrategy(open []Word, followOn Strategy) Fixed {
	return Fixed{open, followOn}
}

func (f Fixed) Guess(game *Game) Word {
	idx := len(game.Guesses)
	if idx < len(f.open) {
		return f.open[idx]
	}
	return f.followOn.Guess(game)
}
