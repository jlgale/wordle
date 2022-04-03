package wordle

import (
	"math/rand"
)

type FilteringStrategy struct {
	rng       *rand.Rand
	log       Logger
	// Fallback strategy when there are too many choices
	fallback  Strategy
	// Use the fallback strategy when this many words remain
	threshold int
	// If several words are equally filtering, use this scoring to
	// weight them and then choose randomly.
	tiebreaker Scoring
}

// Select the word that filters the most from the Possible game words.
// We can only know this based on the actual answer, so we estimate this
// by finding the average removed across all possible answers.
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
func NewFilteringStrategy(rng *rand.Rand, log Logger, fallback Strategy, threshold int, tiebreaker Scoring) *FilteringStrategy {
	return &FilteringStrategy{rng, log, fallback, threshold, tiebreaker}
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
	// random sample of impossible answers. Among these we'll
	// choose the one that filters the best, on average.
	var candidates = sample(n.rng, game.words, n.threshold)
	candidates = append(candidates, possible...)

	var choices []Word
	var choiceRemaining = -1
	for _, candidate := range candidates {
		var remaining = 0
		for _, answer := range possible {
			if candidate == answer {
				continue
			}
			future := game.Guess(candidate, candidate.Match(answer))
			remaining += len(future.PossibleAnswers())
			if choiceRemaining >= 0 && remaining > choiceRemaining {
				break // bail early if this isn't going to work
			}
		}
		if choiceRemaining < 0 || remaining < choiceRemaining {
			choices = choices[:0] // truncate
			choices = append(choices, candidate)
			choiceRemaining = remaining
		} else if remaining == choiceRemaining {
			choices = append(choices, candidate)
		}
	}
	var idx int;
	if len(choices) > 1 {
		weights := n.tiebreaker.Weights(choices)
		idx = weightedSample(n.rng, 2.0, weights)
	} else {
		idx = 0
	}
	choice := choices[idx];
	n.log.Printf("%s filtered an avg of %f%% of words, chosen from %d choices\n",
		choice,
		100.0*(1-float64(choiceRemaining)/float64(len(possible)*(len(possible)-1))),
		len(choices),
	)
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
