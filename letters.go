package wordle

import (
	"fmt"
	"math/bits"
	"strings"
)

// Letters is a set of wordle letters
type Letters uint32

// NewLetters constructs a new Letters set, initialized with the given characters
func NewLetters(s []byte) (l Letters) {
	for _, c := range s {
		l |= letterMask(c)
	}
	return l
}

func letterMask(c byte) Letters {
	if c < 'a' || c > 'z' {
		panic(fmt.Sprintf("bad letter: %c (%d)", c, c))
	}
	idx := c - 'a'
	return 1 << idx
}

func (a Letters) Add(b Letters) Letters {
	return a | b
}

func (a Letters) AddChar(c byte) Letters {
	return a | letterMask(c)
}

func (a Letters) Intersect(b Letters) Letters {
	return a & b
}

// Remove removes b from a
func (a Letters) Remove(b Letters) Letters {
	return a & ^b
}

// String displays Letters in a regex-like form.
func (c Letters) String() string {
	if c == 0 {
		return "[]"
	}
	var b strings.Builder
	b.WriteByte('[')
	run := 0
	endRun := func(letter byte) {
		if run == 2 {
			b.WriteByte(letter)
		} else if run > 2 {
			b.WriteByte('-')
			b.WriteByte(letter)
		}
		run = 0
	}
	for idx := 0; idx < 26; idx++ {
		letter := 'a' + byte(idx)
		if c.Contains(letter) {
			if run == 0 {
				b.WriteByte(letter)
			}
			run += 1
			continue
		}
		endRun(letter - 1)
	}
	endRun('z')
	b.WriteByte(']')
	return b.String()
}

func (l Letters) Empty() bool {
	return l == 0
}

func (l Letters) Len() int {
	return bits.OnesCount32(uint32(l))
}

func (a Letters) Contains(c byte) bool {
	return letterMask(c)&a != 0
}
