package wordle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var all Letters = mkl("abcdefghijklmnopqrstuvwxyz")

func TestAdd(t *testing.T) {
	a := mkl("abc")
	b := mkl("xyz")
	aa := mkl("abcd")
	e := mkl("")
	assert.Equal(t, mkl("abcxyz"), a.Add(b))
	assert.Equal(t, mkl(""), e.Add(e))
	assert.Equal(t, mkl("abc"), a.Add(a))
	assert.Equal(t, mkl("abcd"), a.Add(aa))
}

func TestRemove(t *testing.T) {
	a := mkl("abcd")
	e := mkl("")
	assert.Equal(t, a, a.Remove(e))
	assert.Equal(t, e, a.Remove(a))
	assert.Equal(t, mkl("ad"), a.Remove(mkl("bc")))
}

func TestIntersect(t *testing.T) {
	a := mkl("abcd")
	b := mkl("efg")
	assert.Equal(t, a, a.Intersect(a))
	assert.Equal(t, mkl(""), a.Intersect(b))
}

func TestString(t *testing.T) {
	assert.Equal(t, "[a-d]", mkl("abcd").String())
	assert.Equal(t, "[acd]", mkl("acd").String())
	assert.Equal(t, "[]", mkl("").String())
	assert.Equal(t, "[a-z]", all.String())
	assert.Equal(t, "[a-cp-t]", mkl("abc").Add(mkl("p")).Add(mkl("qrst")).String())
}

func TestContains(t *testing.T) {
	assert.Equal(t, true, mkl("tuv").Contains('t'))
	assert.Equal(t, true, mkl("tuv").Contains('u'))
	assert.Equal(t, false, mkl("abcefg").Contains('d'))
	assert.Equal(t, false, mkl("").Contains('t'))

	assert.Equal(t, true, mkl("").Empty())
	assert.Equal(t, false, mkl("a").Empty())
}
