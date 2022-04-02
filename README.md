# Wordle Play

A little utility to "play" wordle games. Given a dictionary of
possible words, it can try to guess at the answer using various
strategies.

Requires the go 1.17 toolchain. Run `go run ./cmd/...` to run the
utility.

```
A utility for playing "wordle" games on the commandline. Useful for exploring playing strategies.

Usage:
  wordle [command]

Available Commands:
  help        Help about any command
  interact    Interactively guess a wordle answer.
  play        Play automatically with the given answer.

Flags:
  -d, --debug                     Enable debug logging
      --exp float                 Scale weighted strategy by this exponent (default 1)
      --fallback string           Fallback strategy when a simpler strategy is needed (default "freq")
      --fallback-threshold int    Threshold where the fallback strategy is used (default 150)
      --hail-mary string          Choose a different strategy for the final guess. (default "freq")
  -h, --help                      help for wordle
  -o, --open stringArray          Force an opening sequence of guesses
      --score string              Choose among weighted words. One of: random, top (default "random")
      --seed int                  Random seed
  -s, --strategy string           Play strategy. One of: common, diversity, filtering, naive, selective (default "filtering")
      --word-frequencies string   Word frequency scores. (default "./word_freq.csv")
      --words string              Path to accepted word list (default "./words")

Use "wordle [command] --help" for more information about a command.
```
