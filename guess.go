package wordle

// Guess is the guess at an answer, along with the Match corresponding to that guess.
type Guess struct {
	Word  Word
	Match Match
}

func (guess Word) Match(actual Word) Match {
	var m Match
	for i := 0; i < WordLen; i++ {
		if guess[i] == actual[i] {
			m[i] = Green
		} else if actual.contains(guess[i]) {
			m[i] = Yellow
		}
	}
	return m
}

func (g Guess) MustInclude() Letters {
	must := make([]byte, 0, WordLen)
	for i := 0; i < WordLen; i++ {
		if g.Match[i] != Grey {
			must = append(must, g.Word[i])
		}
	}
	return NewLetters(must)
}

func (g Guess) MustNotInclude() Letters {
	mustNot := make([]byte, 0, WordLen)
	for i := 0; i < WordLen; i++ {
		if g.Match[i] == Grey {
			mustNot = append(mustNot, g.Word[i])
		}
	}
	return NewLetters(mustNot)
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
