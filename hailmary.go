package wordle

type HailMary struct {
	normal   Strategy
	hailmary Strategy
}

// A meta strategy which uses a normal strategy for the first
// five guesses and a "hailmary" strategy for the final guess.
func NewHailMary(normal, hailmary Strategy) HailMary {
	return HailMary{normal, hailmary}
}

func (h HailMary) Guess(game *Game) Word {
	if len(game.Guesses) == GuessLimit-1 {
		return h.hailmary.Guess(game)
	}
	return h.normal.Guess(game)
}
