package wordle

type LetterCounts [26]byte

func (lc *LetterCounts) Add(c byte) {
	lc[c-'a'] += 1
}

func (lc *LetterCounts) Remove(c byte) bool {
	if lc[c-'a'] > 0 {
		lc[c-'a'] -= 1
		return true
	}
	return false
}

func (lc LetterCounts) Len() (n int) {
	for _, c := range lc {
		n += int(c)
	}
	return
}
