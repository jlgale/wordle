package wordle

import (
	"math/rand"

	"github.com/rs/zerolog"
)

// Experimentally this setting gives a >97% win rate. Higher values
// don't help much an things get slow quickly. (n*n)
const threshold = 60

type Filtering struct {
	rng *rand.Rand
	log *zerolog.Logger
}

// Select the word that filters the most from the Possible game words.
// We can only know this based on the actual answer, so we instead
// compute the average removed across all possible answers.
//
// This is expensive, so as a fallback we just pick random possible
// answers when that number is above the threshold above. This is okay
// since random possible words do a pretty good job of filtering when
// there are lots of possibilities.
func FilteringStrategy(rng *rand.Rand, log *zerolog.Logger) Filtering {
	return Filtering{rng, log}
}

func (s *Filtering) sampleWords(words []Word, n int) []Word {
	if len(words) <= n {
		return append([]Word(nil), words...)
	}
	var samples = make([]Word, n)
	var from = 0
	for idx := range samples {
		var remaining = n - idx - 1
		var to = len(words) - remaining
		var choice = s.rng.Intn(to - from)
		samples[idx] = words[choice]
		from = choice + 1
	}
	return samples
}

func (n Filtering) Guess(game *Game) Word {
	var possible = game.PossibleAnswers()
	n.log.Printf("%d possible words", len(possible))
	if len(possible) > threshold {
		idx := n.rng.Intn(len(possible))
		return possible[idx]
	}
	if len(possible) == 1 {
		return possible[0]
	}

	// Our candidate words to play are all possible answers plus a
	// random sample of possible words. Among these we'll choose the
	// one that filters the best, on average.
	var candidates = n.sampleWords(game.words, threshold)
	candidates = append(candidates, possible...)

	var choice Word
	var score = -1
	for _, candidate := range candidates {
		var remaining = 0
		for _, answer := range possible {
			if candidate == answer {
				continue
			}
			future := game.Guess(candidate, candidate.Match(answer))
			remaining += len(future.PossibleAnswers())
		}
		if score < 0 || remaining < score {
			choice = candidate
			score = remaining
		}
	}
	n.log.Printf("%s filtered and avg of %f%% of words\n",
		choice, 1-float64(score)/float64(len(possible)*(len(possible)-1)))
	return choice
}
