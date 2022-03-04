package wordle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLC(t *testing.T) {
	a, _ := ParseWord("abcde")
	b, _ := ParseWord("baedc")
	lca := a.LetterCounts()
	lcb := b.LetterCounts()
	assert.Equal(t, lca, lcb)
	assert.Equal(t, 5, lca.Len())
}

func TestLCAdd(t *testing.T) {
	a, _ := ParseWord("abcde")
	lca := a.LetterCounts()
	assert.Equal(t, 5, lca.Len())
	lca.Remove('f')
	assert.Equal(t, 5, lca.Len())
	lca.Remove('c')
	assert.Equal(t, 4, lca.Len())
}
