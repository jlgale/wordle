package main

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"runtime/pprof"
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
	root := &cobra.Command{
		Use:   "wordle",
		Short: "Play wordle games.",
		Long: (`A utility for playing "wordle" games on the commandline. ` +
			`Useful for exploring playing strategies.`),
	}
	wordsOpt := root.PersistentFlags().String("words", "./words",
		"Path to accepted word list")
	seedOpt := root.PersistentFlags().Int64("seed", 0, "Random seed")
	strategyOpt := root.PersistentFlags().StringP("strategy", "s", "filtering",
		"Play strategy. One of: common, diversity, filtering, naive, selective")
	debugOpt := root.PersistentFlags().BoolP("debug", "d", false,
		"Enable debug logging")
	scaleOpt := root.PersistentFlags().String("scale", "random",
		"Choose among weighted words. One of: random, top")
	expOpt := root.PersistentFlags().Float64("exp", 1.0,
		"Scale weighted strategy by this exponent")
	fallbackOpt := root.PersistentFlags().String("fallback", "diversity",
		"Fallback strategy when a simpler strategy is needed")

	// Setup common state
	var words []wordle.Word
	var strategy wordle.Strategy
	var log zerolog.Logger = zerolog.New(os.Stderr).Level(zerolog.InfoLevel)
	var rng *rand.Rand
	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if *debugOpt {
			log = log.Level(zerolog.DebugLevel)
		}
		var err error
		words, err = readWordFile(*wordsOpt, func(word string, lineno int, err error) error {
			log.Printf("%s:%d: %s: %v\n", *wordsOpt, lineno, word, err)
			return nil
		})
		if err != nil {
			return err
		}
		log.Printf("%s: loaded %d words", *wordsOpt, len(words))
		if *seedOpt == 0 {
			*seedOpt = time.Now().UnixNano()
		}
		rng = rand.New(rand.NewSource(*seedOpt))
		var scalefn func(s wordle.Scale) wordle.Strategy
		switch strings.ToLower(*scaleOpt) {
		case "random":
			scalefn = func(s wordle.Scale) wordle.Strategy {
				return wordle.WeightedStrategy(rng, s, *expOpt)
			}
		case "top":
			scalefn = func(s wordle.Scale) wordle.Strategy {
				return wordle.TopStrategy(rng, s)
			}
		default:
			return fmt.Errorf("Unrecognized scale function: %s", *scaleOpt)
		}

		if *debugOpt {
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
			if *debugOpt {
				strategy = &loggingStrategy{strategy, &log}
			}
			return
		}
		fallback, err := mkStrategy(*fallbackOpt)
		if err != nil {
			return err
		}

		switch strings.ToLower(*strategyOpt) {
		case "filtering":
			strategy = wordle.FilteringStrategy(rng, &log, fallback)
			if *debugOpt {
				strategy = &loggingStrategy{strategy, &log}
			}
		default:
			strategy, err = mkStrategy(*strategyOpt)
			if err != nil {
				return err
			}
		}

		return nil
	}

	interactCmd := &cobra.Command{Use: "interact", Short: "Interactively guess a wordle answer."}
	interactCmd.RunE = func(cmd *cobra.Command, args []string) error {
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
	}

	playCmd := &cobra.Command{Use: "play", Short: "Play automatically."}
	repeatOpt := playCmd.Flags().IntP("repeat", "n", 0, "Play multiple games.")
	answersOpt := playCmd.Flags().StringP("answers", "a", "", "Load answers from a file.")
	cpuProfileOpt := root.PersistentFlags().String("cpu-profile", "",
		"Profile CPU usage and write the given file")
	memProfileOpt := root.PersistentFlags().String("mem-profile", "",
		"Profile memory usage and write the given file")
	playCmd.RunE = func(cmd *cobra.Command, args []string) error {
		answers := make([]wordle.Word, len(args))
		for idx, s := range args {
			answer, err := wordle.ParseWord(s)
			if err != nil {
				return fmt.Errorf("%s: %w", s, err)
			}
			answers[idx] = answer
		}
		if *answersOpt != "" {
			loaded, err := readWordFile(*answersOpt, func(word string, lineno int, err error) error {
				log.Printf("%s:%d: %s: %v\n", *answersOpt, lineno, word, err)
				return nil
			})
			if err != nil {
				return err
			}
			answers = append(answers, loaded...)
		}
		if *repeatOpt == 0 {
			for _, answer := range answers {
				game := wordle.NewGame(words, nil)
				play(&game, strategy, answer)
				fmt.Println(game)
				if !game.Won() {
					fmt.Println("The answer was:", answer)
				}
			}
		} else if *repeatOpt > 0 {
			if *cpuProfileOpt != "" {
				f, err := os.Create(*cpuProfileOpt)
				if err != nil {
					return err
				}
				defer f.Close()
				pprof.StartCPUProfile(f)
				defer pprof.StopCPUProfile()
			}
			var guesses int
			var minGuesses int = 7
			var maxGuesses int = 0
			var wins int = 0
			for i := 0; i < *repeatOpt; i++ {
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
			games := *repeatOpt * len(answers)
			fmt.Printf("Won %d of %d games (%0.1f%%). Guesses: avg %0.1f, min %d, max %d\n",
				wins, games, float64(wins)/float64(games)*100, float64(guesses)/float64(games),
				minGuesses, maxGuesses)
			if *memProfileOpt != "" {
				f, err := os.Create(*memProfileOpt)
				if err != nil {
					return err
				}
				defer f.Close()
				runtime.GC()
				if err := pprof.WriteHeapProfile(f); err != nil {
					return err
				}
			}
		}
		return nil
	}
	root.AddCommand(interactCmd, playCmd)
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
