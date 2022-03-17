package wordle

// Guess is the guess at an answer, along with the Match corresponding to that guess.
type Guess struct {
	Word  Word
	Match Match
}

func (g Guess) MustInclude() (must Letters, mustNot Letters) {
	for idx, c := range g.Word {
		if g.Match.Used(idx) {
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
	var match byte
	for idx := 0; idx < WordLen; idx++ {
		if g.Word[idx] == answer[idx] {
			match |= 1 << idx
		}
	}
	return (g.Match.exact ^ match) == 0
}

// Filter a given list of words to include only those that are
// possible answers given this Guess.
func (g Guess) FilterPossible(words []Word) (possibleAnswers []Word) {
	var mustInclude, mustNotInclude = g.MustInclude()
	for _, w := range words {
		var l = w.Letters()
		if !mustInclude.Remove(l).Empty() {
			continue
		}
		if !mustNotInclude.Intersect(l).Empty() {
			continue
		}
		if !g.GreenAllows(w) {
			continue
		}
		possibleAnswers = append(possibleAnswers, w)
	}
	return
}
