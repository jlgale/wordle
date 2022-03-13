package wordle

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMatch(t *testing.T) {
	var matches = []string{
		".....",
		".yG.y",
		"ggggg",
		"GGGGG",
	}
	for _, s := range matches {
		m, err := ParseMatch(s)
		assert.NoError(t, err)
		assert.Equal(t, strings.Replace(s, "g", "G", -1), m.String())
	}
}
