package wordle

// Guess is the guess at an answer, along with the Match corresponding to that guess.
type Guess struct {
	Word  Word
	Match Match
}

func (g Guess) MustInclude() (must Letters, mustNot Letters) {
	for i := 0; i < WordLen; i++ {
		if g.Match[i] != Grey {
			must = must.AddChar(g.Word[i])
		} else {
			mustNot = mustNot.AddChar(g.Word[i])
		}
	}
	// A guess with repeated letters will match Yellow only
	// as many times as the letter appears in the final answer.
	// So we fixup mustNot, removing any letters that we know
	// must be present.
	mustNot = mustNot.Remove(must)
	return
}

func (g Guess) MustBe() (k MustBe) {
	for i := 0; i < WordLen; i++ {
		if g.Match[i] == Green {
			k[i] = g.Word[i]
		}
	}
	return k
}

func (g Guess) MustNotBe() (m MustNotBe) {
	for idx, c := range g.Word {
		if g.Match[idx] != Green {
			m[idx] = NewLetters([]byte{c})
		}
	}
	return m
}

// MustBe tracks letters which are known to be in a given position in the
// final answer.
type MustBe [WordLen]byte

func (k MustBe) Add(l MustBe) MustBe {
	var kk MustBe = k
	for i := 0; i < WordLen; i++ {
		if l[i] != 0 {
			if k[i] != 0 && k[i] != l[i] {
				panic("Disagreement")
			}
			kk[i] = l[i]
		}
	}
	return kk
}

func (k MustBe) Match(w Word) bool {
	for i := 0; i < WordLen; i++ {
		if k[i] != 0 && k[i] != w[i] {
			return false
		}
	}
	return true
}

// MustNotBe tracks letters which are not in a given position in the
// final answer.
type MustNotBe [WordLen]Letters

func (m MustNotBe) Add(l MustNotBe) MustNotBe {
	var mm MustNotBe = m
	for i := 0; i < WordLen; i++ {
		mm[i] = mm[i].Add(l[i])
	}
	return mm
}

func (m MustNotBe) Match(w Word) bool {
	for idx, c := range w {
		if m[idx].Contains(c) {
			return false
		}
	}
	return true
}
