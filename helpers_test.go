package wordle

import (
	"bufio"
	"os"
	"strings"
)

// global list of words used in testing
var words []Word = loadTestWords()

func mkw(s string) Word {
	w, err := ParseWord(s)
	if err != nil {
		panic(err)
	}
	return w
}

func mkm(s string) Match {
	m, err := ParseMatch(s)
	if err != nil {
		panic(err)
	}
	return m
}

func mkl(s string) Letters {
	return NewLetters([]byte(s))
}

func loadTestWords() (words []Word) {
	file, err := os.Open("words")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.IndexByte(line, '#'); idx >= 0 {
			// strip "leading #" style comments
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		w, err := ParseWord(line)
		if err != nil {
			panic(err)
		}
		words = append(words, w)
	}
	return
}
