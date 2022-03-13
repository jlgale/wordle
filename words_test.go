package wordle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseWord(t *testing.T) {
	_, err := ParseWord("yuk")
	assert.Error(t, err)
	_, err = ParseWord("ab de")
	assert.Error(t, err)
	w1, err := ParseWord("hmble")
	assert.NoError(t, err, "parse failed")
	w2, err := ParseWord("HMBLE")
	assert.NoError(t, err, "parse failed")
	assert.Equal(t, w1, w2)
}

func TestWordMatch(t *testing.T) {
	a := mkw("cloth")
	assert.Equal(t, mkm("..y.."), mkw("petar").Match(a))
	assert.Equal(t, mkm("....y"), mkw("quint").Match(a))
	assert.Equal(t, mkm(".ygy."), mkw("stock").Match(a))
	assert.Equal(t, mkm("ggggg"), mkw("cloth").Match(a))

	assert.Equal(t, mkm("..y.."), mkw("fuzzy").Match(mkw("zilch")))
	assert.Equal(t, mkm("....g"), mkw("eagle").Match(mkw("wince")))
}
