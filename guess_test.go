package wordle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkGreenAllows(b *testing.B) {
	var w = mkw("campy")
	var m = mkm("gy.g.")
	var g = Guess{w, m}
	for n := 0; n < b.N; n++ {
		g.GreenAllows(globalWords[n%4096])
	}
}

func BenchmarkFilterPossible(b *testing.B) {
	var g = Guess{
		Word:  mkw("tilde"),
		Match: mkm("yy..g"),
	}
	for n := 0; n < b.N; n++ {
		g.FilterPossible(globalWords)
	}
}

func TestGreenAllows(t *testing.T) {
	var g = Guess{
		Word:  mkw("campy"),
		Match: mkm("gy.g."),
	}
	assert.True(t, g.GreenAllows(mkw("crypt")))
	assert.False(t, g.GreenAllows(mkw("campi")))
	assert.False(t, g.GreenAllows(mkw("curve")))
}

func TestMustInclude(t *testing.T) {
	{
		var g = Guess{
			Word:  mkw("campy"),
			Match: mkm("gy.g."),
		}
		must, mustNot := g.MustInclude()
		assert.Equal(t, mkl("cap"), must)
		assert.Equal(t, mkl("my"), mustNot)
	}
	{
		var g = Guess{
			Word:  mkw("eagle"),
			Match: mkm("....."),
		}
		must, mustNot := g.MustInclude()
		assert.Equal(t, mkl(""), must)
		assert.Equal(t, mkl("eagle"), mustNot)
	}
}

func TestFilterPossible(t *testing.T) {
	var g = Guess{
		Word:  mkw("tilde"),
		Match: mkm("yy..g"),
	}
	possible := g.FilterPossible(globalWords)
	expect := []Word{
		mkw("axite"),
		mkw("boite"),
		mkw("cutie"),
		mkw("evite"),
		mkw("irate"),
		mkw("quite"),
		mkw("retie"),
		mkw("shite"),
		mkw("skite"),
		mkw("smite"),
		mkw("spite"),
		mkw("stime"),
		mkw("stipe"),
		mkw("stire"),
		mkw("stive"),
		mkw("suite"),
		mkw("unite"),
		mkw("untie"),
		mkw("uptie"),
		mkw("urite"),
		mkw("waite"),
		mkw("white"),
		mkw("write"),
	}
	assert.Equal(t, expect, possible)
}
