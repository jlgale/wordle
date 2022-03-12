package wordle

import (
	"math/rand"
)

// When our number of possible answers is > than threshold, use a
// fallback strategy instead.
//
// Experimentally this setting gives a >97% win rate. Higher values
// don't help much an things get slow quickly. (n*n)
const threshold = 60

type Filtering struct {
	rng      *rand.Rand
	log      Logger
	fallback Strategy
}

// Select the word that filters the most from the Possible game words.
// We can only know this based on the actual answer, so we instead
// compute the average removed across all possible answers.
//
// This is expensive, so we need another strategy for the early game
// when there are many possible answers.
//
// This strategy helps the endgame where most letters are known, but
// the number of remaining possible words is higher than the number of
// remaining guesses. For example, if we play "arbas" and match
// "gg.gg", we have five possible words to try: "areas", "arias",
// "arnas", "arpas" and "arras" Trying them 1 by 1, we might run out
// of guesses. Guessing a word that can't be the answer, but that
// tests the uncommon letters between the words, can filter out
// possibilities more quickly. In the above case, "ferny" (where E, R
// and N) are in 3 of the 5 possible answers, guarantees finding the
// solution in 1 or 2 additional guesses.
func FilteringStrategy(rng *rand.Rand, log Logger, fallback Strategy) Filtering {
	return Filtering{rng, log, fallback}
}

func (n Filtering) Guess(game *Game) Word {
	var possible = game.PossibleAnswers()
	if len(possible) > threshold {
		return n.fallback.Guess(game)
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
	n.log.Printf("%s filtered an avg of %f%% of words\n",
		choice, 100.0*(1-float64(score)/float64(len(possible)*(len(possible)-1))))
	return choice
}

// sampleWords returns an array of n words chosen randomly, without
// replacement, from the given array.
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
