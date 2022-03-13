package wordle

import (
	"fmt"
	"strings"
)

const GuessLimit = 6

type Game struct {
	// Guesses made by the player, up to GuessLimit.
	Guesses []Guess
	// All words that can be played.
	words []Word
	// Words removed from play.
	removed []Word
	// A cache of possibleAnswers answers, deduced from words and Guesses.
	possibleAnswers []Word
}

func NewGame(words, used []Word) Game {
	return Game{
		Guesses:         make([]Guess, 0, GuessLimit),
		words:           words,
		removed:         nil,
		possibleAnswers: words,
	}
}

// Guess at the answer
//
// TODO: the calculation of possible answers from the Match is not
// correct because this implementation assumes that Match yellow, "y",
// is set for _every_ matching instance of a character, when in fact
// it's only set for the first _n_ characters where _n_ is the number
// of instances of that character in the answer. For example, if I
// guess "geese" and the answer is "embed", I'll get a match of
// ".yy...", not ".yy.y" as this code expects.
func (game Game) Guess(word Word, match Match) Game {
	var g = Guess{word, match}
	game.Guesses = append(game.Guesses, g)
	game.possibleAnswers = g.FilterPossible(game.possibleAnswers)
	return game
}

// Don't try Guess the given word. Useful if the official
// game doesn't like a word that we choose.
func (game *Game) RemoveWord(removed Word) {
	game.removed = append(game.removed, removed)
	var filtered []Word
	for _, w := range game.possibleAnswers {
		if w == removed {
			continue
		}
		filtered = append(filtered, w)
	}
	game.possibleAnswers = filtered
}

func (game Game) PossibleAnswers() []Word {
	return game.possibleAnswers
}

func (game Game) Over() bool {
	if len(game.Guesses) >= GuessLimit {
		return true
	}
	for _, g := range game.Guesses {
		if g.Match.Won() {
			return true
		}
	}
	return false
}

func (game Game) Won() bool {
	for _, g := range game.Guesses {
		if g.Match.Won() {
			return true
		}
	}
	return false
}

func (game Game) String() string {
	var b strings.Builder
	for idx, g := range game.Guesses {
		if idx > 0 {
			b.WriteByte('\n')
		}
		fmt.Fprintf(&b, "%d: %s %s", idx+1, g.Word, g.Match)
	}
	return b.String()
}
