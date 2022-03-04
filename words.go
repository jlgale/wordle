package wordle

import (
	"fmt"
	"strings"
)

const WordLen = 5

type Word [WordLen]byte

func ParseWord(s string) (w Word, err error) {
	s = strings.ToLower(s)
	if len(s) != WordLen {
		err = fmt.Errorf("Words must be %d letters.", WordLen)
		return
	}
	for i, c := range s {
		if c < 'a' || c > 'z' {
			err = fmt.Errorf("Letter %c not allowed", c)
			return
		}
		w[i] = byte(c)
	}
	return
}

func (w Word) Letters() Letters {
	return NewLetters(w[:])
}

func (w Word) contains(c byte) bool {
	found := false
	for i := 0; i < len(w); i++ {
		if w[i] == c {
			found = true
		}
	}
	return found
}

func (w Word) String() string {
	return string(w[:])
}

func (guess Word) Match(actual Word) Match {
	var lc = actual.LetterCounts()
	var m Match
	for i, g := range guess {
		if lc.Remove(g) {
			m[i] = Yellow
		}
		if g == actual[i] {
			m[i] = Green
		}
	}
	return m
}

func (w Word) LetterCounts() (lc LetterCounts) {
	for _, c := range w {
		lc.Add(c)
	}
	return
}
