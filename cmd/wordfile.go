package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/jlgale/wordle"
)

// readWordFile loads the contents of the file at the given filename
// as a list of Wordle words, one per line. The format is similar to
// a unix "dict" file.
//
// Comments (beginning with #) are ignored
// Non-conforming words are passed to the given onError handler. If that handler
// returns an error, ReadWordFile stops and returns it.
func readWordFile(filename string, onError func(word string, lineno int, err error) error) ([]wordle.Word, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []wordle.Word
	var seen = map[wordle.Word]bool{}
	scanner := bufio.NewScanner(file)
	lineno := 0
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.IndexByte(line, '#'); idx >= 0 {
			// strip "leading #" style comments
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		lineno += 1
		if line == "" {
			continue
		}
		w, err := wordle.ParseWord(line)
		if err != nil {
			if err := onError(line, lineno, err); err != nil {
				return words, err
			}
			continue
		}
		if seen[w] {
			continue // ignore duplicate words
		}
		seen[w] = true
		words = append(words, w)
	}
	return words, nil
}
