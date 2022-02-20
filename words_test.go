package wordle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func ParseWordTest(t *testing.T) {
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
