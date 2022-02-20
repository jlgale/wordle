package main

import (
	"fmt"
	"github.com/jlgale/wordle"
	"github.com/spf13/cobra"
	"math/rand"
	"strings"
)

var (
	wordFilePath string
	randomSeed   int
	repeat       int
	playStrategy string
)

var (
	root = &cobra.Command{
		Use:   "wordle",
		Short: "Play wordle games.",
		Long:  `A utility for playing "wordle" games on the commandline. Useful for exploring playing strategies.`,
	}
)

func init() {
	root.PersistentFlags().StringVar(&wordFilePath, "words", "./words", "Path to accepted word list")
	root.PersistentFlags().IntVar(&randomSeed, "seed", 42, "Random seed")
	root.PersistentFlags().StringVarP(&playStrategy, "strategy", "s", "naive", "Play strategy. One of: naive, diversity")
	interactive := &cobra.Command{
		Use:   "interactive",
		Short: "Interactively guess a wordle answer.",
		RunE: func(cmd *cobra.Command, args []string) error {
			words, err := ReadWordFile(wordFilePath, func(word string, lineno int, err error) error {
				// fmt.Fprintf(os.Stderr, "%s:%d: %s: %v\n", *wordFilePath, lineno, word, err)
				return nil
			})
			if err != nil {
				return err
			}
			rng := rand.New(rand.NewSource(int64(randomSeed)))

			var strategy wordle.Strategy
			switch strings.ToLower(playStrategy) {
			case "naive":
				strategy = wordle.NaiveStrategy(rng)
			case "diversity":
				strategy = wordle.DiversityStrategy(rng)
			default:
				return fmt.Errorf("Unrecognized strategy: %s", playStrategy)
			}
			game := wordle.NewGame(words, nil)
			for !game.Over() {
				guess := strategy.Guess(&game)
				fmt.Println("My guess", guess)
				var matchString string
				for {
					fmt.Print(`describe match (or "again" for a different guess): `)
					fmt.Scanf("%s", &matchString)
					// In case the chosen word is not allowed:
					if strings.ToLower(matchString) == "again" {
						game.RemoveWord(guess)
						break
					}
					match, err := wordle.ParseMatch(matchString)
					if err != nil {
						fmt.Println(err)
						continue
					}
					game.AddGuess(guess, match)
					break
				}
			}
			return nil
		},
	}
	play := &cobra.Command{
		Use:   "play",
		Short: "Play automatically.",
		RunE: func(cmd *cobra.Command, args []string) error {
			words, err := ReadWordFile(wordFilePath, func(word string, lineno int, err error) error {
				return nil
			})
			if err != nil {
				return err
			}
			rng := rand.New(rand.NewSource(int64(randomSeed)))

			var strategy wordle.Strategy
			switch strings.ToLower(playStrategy) {
			case "naive":
				strategy = wordle.NaiveStrategy(rng)
			case "diversity":
				strategy = wordle.DiversityStrategy(rng)
			default:
				return fmt.Errorf("Unrecognized strategy: %s", playStrategy)
			}
			answers := make([]wordle.Word, len(args))
			for idx, s := range args {
				answer, err := wordle.ParseWord(s)
				if err != nil {
					return fmt.Errorf("%s: %w", s, err)
				}
				answers[idx] = answer
			}
			if repeat == 1 {
				for _, answer := range answers {
					game := wordle.NewGame(words, nil)
					play(&game, strategy, answer)
					fmt.Println(game)
				}
			} else if repeat > 1 {
				var guesses int
				var minGuesses int = 7
				var maxGuesses int = 0
				var wins int = 0
				for i := 0; i < repeat; i++ {
					for _, answer := range answers {
						game := wordle.NewGame(words, nil)
						play(&game, strategy, answer)
						if game.Won() {
							wins += 1
						}
						guesses += len(game.Guesses)
						if len(game.Guesses) > maxGuesses {
							maxGuesses = len(game.Guesses)
						}
						if len(game.Guesses) < minGuesses {
							minGuesses = len(game.Guesses)
						}
					}
				}
				games := repeat * len(answers)
				fmt.Printf("Won %d of %d games (%0.1f%%). Guesses: avg %0.1f, min %d, max %d\n",
					wins, games, float64(wins)/float64(games)*100, float64(guesses)/float64(games),
					minGuesses, maxGuesses)
			}
			return nil
		},
	}
	play.PersistentFlags().IntVarP(&repeat, "repeat", "n", 1, "Play multiple games.")
	root.AddCommand(interactive, play)
}

func play(wdl *wordle.Game, strategy wordle.Strategy, word wordle.Word) {
	for !wdl.Over() {
		guess := strategy.Guess(wdl)
		match := guess.Match(word)
		wdl.AddGuess(guess, match)
	}
}

func main() {
	root.Execute()
}
