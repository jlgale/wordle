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
