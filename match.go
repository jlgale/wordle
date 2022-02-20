package wordle

import (
	"fmt"
	"strings"
)

type Square byte

const (
	Grey Square = iota
	Yellow
	Green
)

// Match is the result of guessing the answer:
// each letter is colored one of: Grey when
// the letter is not in the answer; Yellow
// when the letter is in the answer; Green
// when the letter is in the answer, in this
// position.
type Match [WordLen]Square

// Won returns true when every square of the
// Match is Green.
func (m Match) Won() bool {
	for _, sq := range m {
		if sq != Green {
			return false
		}
	}
	return true
}

func (m Match) String() string {
	var b strings.Builder
	for _, sq := range m {
		switch sq {
		case Yellow:
			b.WriteByte('y')
		case Green:
			b.WriteByte('G')
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
			m[idx] = Yellow
		case 'g':
			m[idx] = Green
		case '.':
			m[idx] = Grey
		default:
			return m, fmt.Errorf("Unrecognized match description: %s", s)
		}
	}
	return m, nil
}
