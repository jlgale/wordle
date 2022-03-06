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

func play(game *wordle.Game, strategy wordle.Strategy, answer wordle.Word) {
	for !game.Over() {
		guess := strategy.Guess(game)
		match := guess.Match(answer)
		*game = game.Guess(guess, match)
	}
}

func main() {
	var words []wordle.Word
	var strategy wordle.Strategy
	var log zerolog.Logger = zerolog.New(os.Stderr).Level(zerolog.InfoLevel)
	var rng *rand.Rand
	var playStrategy string
	var wordFilePath string
	var randomSeed int64
	var debugLogging bool
	var pow float64
	var scale string
	var fallbackStrategy string

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
			var scalefn func(s wordle.Scale) wordle.Strategy
			switch strings.ToLower(scale) {
			case "random":
				scalefn = func(s wordle.Scale) wordle.Strategy {
					return wordle.WeightedStrategy(rng, s, pow)
				}
			case "top":
				scalefn = func(s wordle.Scale) wordle.Strategy {
					return wordle.TopStrategy(rng, s)
				}
			default:
				return fmt.Errorf("Unrecognized scale function: %s", scale)
			}

			var mkStrategy = func(name string) (wordle.Strategy, error) {
				switch strings.ToLower(name) {
				case "common":
					return scalefn(wordle.CommonScale()), nil
				case "diversity":
					return scalefn(wordle.DiversityScale()), nil
				case "naive":
					return wordle.NaiveStrategy(rng, &log), nil
				case "selective":
					return scalefn(wordle.SelectiveScale(&log)), nil
				default:
					return nil, fmt.Errorf("Unrecognized fallback strategy: %s", name)
				}
			}
			fallback, err := mkStrategy(fallbackStrategy)
			if err != nil {
				return err
			}

			switch strings.ToLower(playStrategy) {
			case "filtering":
				strategy = wordle.FilteringStrategy(rng, &log, fallback)
			default:
				strategy, err = mkStrategy(playStrategy)
				if err != nil {
					return err
				}
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
	root.PersistentFlags().StringVar(&scale, "scale", "random",
		"Choose among weighted words")
	root.PersistentFlags().Float64Var(&pow, "pow", 1.0,
		"Scale weighted strategy by this exponent")
	root.PersistentFlags().StringVar(&fallbackStrategy, "fallback", "diversity",
		"Fallback strategy when a simpler strategy is needed")
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
			switch len(game.Guesses) {
			case 1:
				fmt.Println("Genius")
			case 2:
				fmt.Println("Magnificent")
			case 3:
				fmt.Println("Impressive")
			case 4:
				fmt.Println("Splendid")
			case 5:
				fmt.Println("Great")
			case 6:
				fmt.Println("Phew")
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
