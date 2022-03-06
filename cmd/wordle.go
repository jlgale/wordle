package main

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"sort"
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

	// Setup common state
	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
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

		if debugLogging {
			// Wrap a logger around the scale function
			innerScaleFn := scalefn
			scalefn = func(s wordle.Scale) wordle.Strategy {
				return innerScaleFn(&loggingScale{s, &log})
			}
		}

		var mkStrategy = func(name string) (strategy wordle.Strategy, err error) {
			switch strings.ToLower(name) {
			case "common":
				strategy = scalefn(wordle.CommonScale())
			case "diversity":
				strategy = scalefn(wordle.DiversityScale())
			case "naive":
				strategy = wordle.NaiveStrategy(rng)
			case "selective":
				strategy = scalefn(wordle.SelectiveScale())
			default:
				return nil, fmt.Errorf("Unrecognized fallback strategy: %s", name)
			}
			if debugLogging {
				strategy = &loggingStrategy{strategy, &log}
			}
			return
		}
		fallback, err := mkStrategy(fallbackStrategy)
		if err != nil {
			return err
		}

		switch strings.ToLower(playStrategy) {
		case "filtering":
			strategy = wordle.FilteringStrategy(rng, &log, fallback)
			if debugLogging {
				strategy = &loggingStrategy{strategy, &log}
			}
		default:
			strategy, err = mkStrategy(playStrategy)
			if err != nil {
				return err
			}
		}
		return nil
	}

	interact := &cobra.Command{
		Use:   "interact",
		Short: "Interactively guess a wordle answer.",
		RunE: func(cmd *cobra.Command, args []string) error {
			game := wordle.NewGame(words, nil)
			for !game.Over() {
				guess := strategy.Guess(&game)
			force:
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
					if strings.ToLower(matchString) == "force" {
						fmt.Print(`give me a word to play: `)
						var guessString string
						fmt.Scanf("%s", &guessString)
						var err error
						guess, err = wordle.ParseWord(guessString)
						if err != nil {
							fmt.Println(err)
							continue
						}
						goto force
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

type loggingScale struct {
	inner wordle.Scale
	log   wordle.Logger
}

func (s *loggingScale) Weights(words []wordle.Word) []float64 {
	inner := reflect.TypeOf(s.inner)
	weights := s.inner.Weights(words)
	index := make([]int, len(words))
	for idx := range words {
		index[idx] = idx
	}
	sort.Slice(index, func(i, j int) bool {
		return weights[i] > weights[j]
	})
	if len(words) > 0 {
		s.log.Printf("%s scored %d words [%f,%f], top:",
			inner.Name(), len(words),
			weights[index[len(index)-1]],
			weights[index[0]])
		for idx := 0; idx < 5 && idx < len(index); idx++ {
			s.log.Printf(" %d: %f %s", idx+1, weights[index[idx]], words[index[idx]])
		}
	}
	return weights
}

type loggingStrategy struct {
	inner wordle.Strategy
	log   wordle.Logger
}

func (s *loggingStrategy) Guess(game *wordle.Game) wordle.Word {
	inner := reflect.TypeOf(s.inner)
	w := s.inner.Guess(game)
	s.log.
		Printf("%s chose %q for guess %d, %d possible answers remaining",
			inner.Name(), w, len(game.Guesses)+1, len(game.PossibleAnswers()))
	return w
}
