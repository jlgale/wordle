package wordle

// Guess is the guess at an answer, along with the Match corresponding to that guess.
type Guess struct {
	Word  Word
	Match Match
}

func (g Guess) MustInclude() (must Letters, mustNot Letters) {
	for idx, c := range g.Word {
		if g.Match[idx] != Grey {
			must = must.AddChar(c)
		} else {
			mustNot = mustNot.AddChar(c)
		}
	}
	// A guess with repeated letters will match Yellow only
	// as many times as the letter appears in the final answer.
	// So we fixup mustNot, removing any letters that we know
	// must be present.
	mustNot = mustNot.Remove(must)
	return
}

// Allows returns true if what we know from this Guess's Green squares
// allows the given answer.
func (g Guess) GreenAllows(answer Word) bool {
	for idx, c := range answer {
		if g.Match[idx] == Green {
			if g.Word[idx] != c {
				return false
			}
		} else {
			if g.Word[idx] == c {
				return false
			}
		}
	}
	return true
}
