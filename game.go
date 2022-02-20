package wordle

import (
	"fmt"
	"strings"
)

const GuessLimit = 6

type Game struct {
	Guesses []Guess
	words   []Word
	used    []Word
	removed []Word
}

func NewGame(words, used []Word) Game {
	return Game{nil, words, used, nil}
}

func (w *Game) AddGuess(word Word, match Match) {
	w.Guesses = append(w.Guesses, Guess{word, match})
}

func (w *Game) RemoveWord(word Word) {
	w.removed = append(w.removed, word)
}

func (wdl Game) Possible() []Word {
	var mustInclude Letters
	var mustNotInclude Letters
	var mustBe MustBe
	var mustNotBe MustNotBe
	for _, g := range wdl.Guesses {
		mustInclude = mustInclude.Add(g.MustInclude())
		mustNotInclude = mustNotInclude.Add(g.MustNotInclude())
		mustNotBe = mustNotBe.Add(g.MustNotBe())
		mustBe = mustBe.Add(g.MustBe())
	}
	var removed = map[Word]bool{}
	for _, w := range wdl.removed {
		removed[w] = true
	}
	var possible []Word
	for _, w := range wdl.words {
		if removed[w] {
			continue
		}
		if !mustBe.Match(w) {
			continue
		}
		if !mustNotBe.Match(w) {
			continue
		}
		var l = w.Letters()
		if !mustInclude.Remove(l).Empty() {
			continue
		}
		if !mustNotInclude.Intersect(l).Empty() {
			continue
		}
		possible = append(possible, w)
	}
	return possible
}

func (wdl Game) Over() bool {
	if len(wdl.Guesses) >= GuessLimit {
		return true
	}
	for _, g := range wdl.Guesses {
		if g.Match.Won() {
			return true
		}
	}
	return false
}

func (wdl Game) Won() bool {
	for _, g := range wdl.Guesses {
		if g.Match.Won() {
			return true
		}
	}
	return false
}

func (wdl Game) String() string {
	var b strings.Builder
	for idx, g := range wdl.Guesses {
		if idx > 0 {
			b.WriteByte('\n')
		}
		fmt.Fprintf(&b, "%d: %s %s", idx+1, g.Word, g.Match)
	}
	return b.String()
}

type Strategy interface {
	Guess(w *Game) Word
}
