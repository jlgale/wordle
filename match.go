package wordle

import (
	"fmt"
	"strings"
)

// Match is the result of guessing the answer:
// each letter is colored one of: Grey when
// the letter is not in the answer; Yellow
// when the letter is in the answer; Green
// when the letter is in the answer, in this
// position.
type Match struct {
	exact byte // mask of green squares
	used  byte // mask of colored squares
}

// Won returns true when every square of the
// Match is Green.
func (m Match) Won() bool {
	return m.exact == 0x1f
}

func (m Match) String() string {
	var b strings.Builder
	for idx := 0; idx < WordLen; idx++ {
		switch {
		case maskSet(m.exact, idx):
			b.WriteByte('G')
		case maskSet(m.used, idx):
			b.WriteByte('y')
		default:
			b.WriteByte('.')
		}
	}
	return b.String()
}

// ParseMatch parses a 5-letter string describing the "match" of a
// guess with the answer.
func ParseMatch(s string) (Match, error) {
	var m Match
	s = strings.TrimSpace(s)
	if len(s) != WordLen {
		return m, fmt.Errorf("Unrecognized match description: %s", s)
	}
	for idx, c := range strings.ToLower(s) {
		switch c {
		case 'y':
			m.SetUsed(idx, false)
		case 'g':
			m.SetUsed(idx, true)
		case '.':
			// pass
		default:
			return m, fmt.Errorf("Unrecognized match description: %s", s)
		}
	}
	return m, nil
}

func (m *Match) SetUsed(idx int, exact bool) {
	m.used |= 1 << idx
	if exact {
		m.exact |= 1 << idx
	}
}

func (m Match) Used(idx int) bool {
	return maskSet(m.used, idx)
}

func maskSet(m byte, idx int) bool {
	return m&(1<<idx) != 0
}
