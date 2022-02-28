package wordle

import "math/rand"

type Diversity struct {
	rng *rand.Rand
}

func DiversityStrategy(rng *rand.Rand) Diversity {
	return Diversity{rng}
}

func (n Diversity) Guess(game *Game) Word {
	var possible = game.PossibleAnswers()
	var choices = []Word{possible[0]}
	var score = choices[0].Letters().Len()
	for _, w := range possible[1:] {
		wscore := w.Letters().Len()
		if wscore > score {
			choices = []Word{w}
			score = wscore
		} else {
			choices = append(choices, w)
		}
	}
	idx := n.rng.Intn(len(choices))
	return choices[idx]
}
