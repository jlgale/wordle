package wordle

import (
	"math/rand"
)

type FilteringStrategy struct {
	rng       *rand.Rand
	log       Logger
	fallback  Strategy
	threshold int
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
func NewFilteringStrategy(rng *rand.Rand, log Logger, fallback Strategy, threshold int) *FilteringStrategy {
	return &FilteringStrategy{rng, log, fallback, threshold}
}

func (n FilteringStrategy) Guess(game *Game) Word {
	var possible = game.PossibleAnswers()
	if len(possible) > n.threshold {
		return n.fallback.Guess(game)
	}
	if len(possible) == 1 {
		return possible[0]
	}

	// Our candidate words to play are all possible answers plus a
	// random sample of possible words. Among these we'll choose the
	// one that filters the best, on average.
	var candidates = sample(n.rng, game.words, n.threshold)
	candidates = append(candidates, possible...)

	var choice Word
	var choiceRemaining = -1
	for _, candidate := range candidates {
		var remaining = 0
		for _, answer := range possible {
			if candidate == answer {
				continue
			}
			future := game.Guess(candidate, candidate.Match(answer))
			remaining += len(future.PossibleAnswers())
		}
		if choiceRemaining < 0 || remaining < choiceRemaining {
			choice = candidate
			choiceRemaining = remaining
		}
	}
	n.log.Printf("%s filtered an avg of %f%% of words\n",
		choice, 100.0*(1-float64(choiceRemaining)/float64(len(possible)*(len(possible)-1))))
	return choice
}

// sample returns an array of n words chosen randomly, without
// replacement, from the given array.
func sample[T any](rng *rand.Rand, arr []T, n int) []T {
	if len(arr) <= n {
		return append([]T(nil), arr...)
	}
	var samples = make([]T, n)
	var from = 0
	for idx := range samples {
		var remaining = n - idx - 1
		var to = len(arr) - remaining
		var choice = rng.Intn(to - from)
		samples[idx] = arr[choice]
		from = choice + 1
	}
	return samples
}
