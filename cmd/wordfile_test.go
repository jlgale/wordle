package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadWordfile(t *testing.T) {
	discarded := 0
	words, err := readWordFile("../test_answers", func(word string, lineno int, err error) error {
		assert.Equal(t, "invalid", word)
		assert.Equal(t, 9, lineno)
		discarded += 1
		return nil
	})
	assert.Nil(t, err)
	assert.Len(t, words, 20)
	assert.Equal(t, 1, discarded)
}
