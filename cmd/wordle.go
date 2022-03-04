package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jlgale/wordle"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	words    []wordle.Word
	strategy wordle.Strategy
	log      zerolog.Logger = zerolog.New(os.Stderr).Level(zerolog.InfoLevel)
	rng      *rand.Rand
)

func play(wdl *wordle.Game, strategy wordle.Strategy, word wordle.Word) {
	for !wdl.Over() {
		guess := strategy.Guess(wdl)
		match := guess.Match(word)
		*wdl = wdl.Guess(guess, match)
	}
}

func main() {
	var playStrategy string
	var wordFilePath string
	var randomSeed int64
	var debugLogging bool

	root := &cobra.Command{
		Use:   "wordle",
		Short: "Play wordle games.",
		Long: (`A utility for playing "wordle" games on the commandline. ` +
			`Useful for exploring playing strategies.`),

		// Setup common state
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if debugLogging {
				log = log.Level(zerolog.DebugLevel)
			}
			var err error
			words, err = readWordFile(wordFilePath, func(word string, lineno int, err error) error {
				log.Printf("%s:%d: %s: %v\n", wordFilePath, lineno, word, err)
				return nil
			})
			if err != nil {
				return err
			}
			log.Printf("%s: loaded %d words", wordFilePath, len(words))
			if randomSeed == 0 {
				randomSeed = time.Now().UnixNano()
			}
			rng = rand.New(rand.NewSource(randomSeed))

			switch strings.ToLower(playStrategy) {
			case "common":
				strategy = wordle.WeightedStrategy(rng,
					wordle.CommonScale(&log))
			case "diversity":
				strategy = wordle.WeightedStrategy(rng, wordle.DiversityScale())
			case "filtering":
				strategy = wordle.FilteringStrategy(rng, &log, wordle.WeightedStrategy(rng, wordle.DiversityScale()))
			case "naive":
				strategy = wordle.NaiveStrategy(rng, &log)
			case "selective":
				strategy = wordle.WeightedStrategy(rng, wordle.SelectiveScale(&log))
			default:
				return fmt.Errorf("Unrecognized strategy: %s", playStrategy)
			}
			return nil
		},
	}
	root.PersistentFlags().StringVar(&wordFilePath, "words", "./words",
		"Path to accepted word list")
	root.PersistentFlags().Int64Var(&randomSeed, "seed", 0, "Random seed")
	root.PersistentFlags().StringVarP(&playStrategy, "strategy", "s", "filtering",
		"Play strategy. One of: common, diversity, filtering, naive, selective")
	root.PersistentFlags().BoolVarP(&debugLogging, "debug", "d", false,
		"Enable debug logging")
	interact := &cobra.Command{
		Use:   "interact",
		Short: "Interactively guess a wordle answer.",
		RunE: func(cmd *cobra.Command, args []string) error {
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
					game = game.Guess(guess, match)
					break
				}
			}
			return nil
		},
	}
	var repeat int
	var answerFilePath string
	play := &cobra.Command{
		Use:   "play",
		Short: "Play automatically.",
		RunE: func(cmd *cobra.Command, args []string) error {
			answers := make([]wordle.Word, len(args))
			for idx, s := range args {
				answer, err := wordle.ParseWord(s)
				if err != nil {
					return fmt.Errorf("%s: %w", s, err)
				}
				answers[idx] = answer
			}
			if answerFilePath != "" {
				loaded, err := readWordFile(answerFilePath, func(word string, lineno int, err error) error {
					log.Printf("%s:%d: %s: %v\n", answerFilePath, lineno, word, err)
					return nil
				})
				if err != nil {
					return err
				}
				answers = append(answers, loaded...)
			}
			if repeat == 0 {
				for _, answer := range answers {
					game := wordle.NewGame(words, nil)
					play(&game, strategy, answer)
					fmt.Println(game)
					if !game.Won() {
						fmt.Println("The answer was:", answer)
					}
				}
			} else if repeat > 0 {
				var guesses int
				var minGuesses int = 7
				var maxGuesses int = 0
				var wins int = 0
				for i := 0; i < repeat; i++ {
					for _, answer := range answers {
						log.Debug().Stringer("answer", answer).Msg("New Game")
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
	play.Flags().IntVarP(&repeat, "repeat", "n", 0, "Play multiple games.")
	play.Flags().StringVarP(&answerFilePath, "answers", "a", "", "Load answers from a file.")
	root.AddCommand(interact, play)
	root.Execute()
}
